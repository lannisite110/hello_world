package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const erc20ABIJSON = `[
   {
     "anonymous":false,
	 "inputs":[
	 {"indexed":true,"name":"from","type":"address"},
	 {"indexed":true,"name":"to","type":"address"},
	 {"indexed":false,"name":"value","type":"uint256"}
	 ],
	 "name":"Transfer",
	 "type":"event"
   }
]`

type TransferEvent struct {
	BlockNumber uint64    `json:"block_number"`
	TxHash      string    `json:"tx_hash"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Value       string    `json:"value"`
	Timestamp   time.Time `json:"timestamp"`
}

type EventStore struct {
	mu     sync.RWMutex
	events []TransferEvent
	limit  int
}

func NewEventStore(limit int) *EventStore {
	return &EventStore{
		events: make([]TransferEvent, 0, limit),
		limit:  limit,
	}
}

func (s *EventStore) Add(e TransferEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.events) >= s.limit {
		s.events = s.events[1:]
	}
	s.events = append(s.events, e)
}

func (s *EventStore) List() []TransferEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]TransferEvent, len(s.events))
	copy(out, s.events)
	return out
}

func main() {
	// 优先读 WS，再读 RPC（适配 Sepolia）
	rpcURL := os.Getenv("ETH_WS_URL")
	if rpcURL == "" {
		rpcURL = os.Getenv("ETH_RPC_URL") // 修正笔误：PRC → RPC
	}

	if rpcURL == "" {
		log.Fatalf("ETH_WS_URL or ETH_RPC_URL must be set")
	}

	contractHex := os.Getenv("ERC20_CONTRACT")
	if contractHex == "" {
		log.Fatal("ERC20_CONTRACT env is not set")
	}
	contractAddr := common.HexToAddress(contractHex)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 连接 Sepolia 节点（增加超时）
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		log.Fatalf("failed to connect to Ethereum node:%v", err)
	}
	defer client.Close()

	// 解析 ABI
	parsedABI, err := abi.JSON(strings.NewReader(erc20ABIJSON))
	if err != nil {
		log.Fatalf("failed to parse ABI:%v", err)
	}

	// 初始化事件缓存（最多 100 条）
	store := NewEventStore(100)

	// 启动 HTTP 轮询（核心！替换原订阅逻辑）
	go pollTransactionEvents(ctx, client, parsedABI, contractAddr, store)

	// HTTP 接口
	mux := http.NewServeMux()
	mux.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		events := store.List()
		_ = json.NewEncoder(w).Encode(events)
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// 启动 HTTP 服务
	go func() {
		log.Printf("HTTP server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("http server error :%v", err)
		}
	}()

	// 优雅退出
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	fmt.Printf("received signal %s, shutting down ...\n", sig.String())

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	_ = server.Shutdown(shutdownCtx)
	cancel()
}

// 轮询版本：定时查询 Sepolia 区块日志（适配 HTTP 节点）
func pollTransactionEvents(ctx context.Context, client *ethclient.Client, parsedABI abi.ABI, contract common.Address, store *EventStore) {
	log.Printf("start polling Transfer events of %s (Sepolia HTTP mode)", contract.Hex())

	// 初始查询最新区块（避免漏查历史事件）
	latestBlock, err := client.BlockNumber(ctx)
	if err != nil {
		log.Printf("failed to get initial block: %v", err)
		latestBlock = 0
	}
	lastBlockNumber := latestBlock

	// Sepolia 区块出块约 12 秒，轮询间隔设为 15 秒
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 1. 获取最新区块号
			currentBlock, err := client.BlockNumber(ctx)
			if err != nil {
				log.Printf("failed to get latest block: %v", err)
				continue
			}

			// 2. 无新区块则跳过
			if currentBlock <= lastBlockNumber {
				log.Printf("no new blocks (last: %d, current: %d)", lastBlockNumber, currentBlock)
				continue
			}

			log.Printf("scanning blocks from %d to %d", lastBlockNumber+1, currentBlock)

			// 3. 构建日志查询（指定合约+事件）
			query := ethereum.FilterQuery{
				Addresses: []common.Address{contract},
				FromBlock: new(big.Int).SetUint64(lastBlockNumber + 1),
				ToBlock:   new(big.Int).SetUint64(currentBlock),
				Topics: [][]common.Hash{
					{parsedABI.Events["Transfer"].ID}, // 只查 Transfer 事件（过滤无关日志）
				},
			}

			// 4. 查询日志
			logs, err := client.FilterLogs(ctx, query)
			if err != nil {
				log.Printf("failed to filter logs: %v", err)
				continue
			}

			// 5. 解析并保存事件
			for _, vLog := range logs {
				var event struct {
					From  common.Address
					To    common.Address
					Value *big.Int
				}

				// 解码非 indexed 参数
				if err := parsedABI.UnpackIntoInterface(&event, "Transfer", vLog.Data); err != nil {
					log.Printf("failed to unpack log data: %v", err)
					continue
				}

				// 解码 indexed 地址（Topics[1]=from, Topics[2]=to）
				if len(vLog.Topics) >= 3 {
					event.From = common.BytesToAddress(vLog.Topics[1].Bytes())
					event.To = common.BytesToAddress(vLog.Topics[2].Bytes())
				}

				// 添加到缓存
				transferEvent := TransferEvent{
					BlockNumber: vLog.BlockNumber,
					TxHash:      vLog.TxHash.Hex(),
					From:        event.From.Hex(),
					To:          event.To.Hex(),
					Value:       event.Value.String(),
					Timestamp:   time.Now(),
				}
				store.Add(transferEvent)

				// 打印日志（方便调试）
				log.Printf("captured Transfer event:")
				log.Printf("  Block: %d, TxHash: %s", vLog.BlockNumber, vLog.TxHash.Hex())
				log.Printf("  From: %s, To: %s, Value: %s", event.From.Hex(), event.To.Hex(), event.Value.String())
			}

			// 6. 更新最后查询的区块号
			lastBlockNumber = currentBlock

		case <-ctx.Done():
			log.Println("context cancelled, stop polling")
			return
		}
	}
}
