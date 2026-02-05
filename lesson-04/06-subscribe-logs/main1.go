package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 06-subscribe-logs.go
// 订阅指定合约的日志事件（支持 Transfer + Approval 事件）
const erc20ABIJSON = `[
    {
        "anonymous":false,
        "inputs":[
            {"indexed":true,"name":"owner","type":"address"},
            {"indexed":true,"name":"spender","type":"address"},
            {"indexed":false,"name":"value","type":"uint256"}
        ],
        "name":"Approval",
        "type":"event"
    },
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

func main(){
	contractAddr:=flag.String("contract","","Contract address to subscribe logs from (required)")
    flag.Parse()

	if *contractAddr==""{
		log.Fatalf("missing -contract flag")
	}

	rpcURL:=os.Getenv("ETH_WS_URL")
	if rpcURL==""{
		rpcURL = os.Getenv("ETH_RPC_URL")
	}

	if rpcURL==""{
		log.Fatalf("ETH_WS_URL or ETH_RPC_URL must be set")
	}

	ctx,cancel:=context.WithCancel(context.Background())
    defer cancel()

	client,err:=ethclient.DialContext(ctx,rpcURL)
	if err!=nil{
		log.Fatalf("failed to connect to Ethereum node:%v",err)
	}
	defer client.Close()

	// 解析 ABI
	parsedABI,err:=abi.JSON(strings.NewReader(erc20ABIJSON))
	if err!=nil{
		log.Fatalf("failed to parse ABI:%v",err)
	}
	contract:=common.HexToAddress(*contractAddr)

	query:=ethereum.FilterQuery{
		Addresses:[]common.Address{contract},
	}

	logsCh:=make(chan types.Log)
	sub,err:=client.SubscribeFilterLogs(ctx,query,logsCh)
	if err!=nil{
		log.Fatalf("failed to subscribe logs: %v",err)
	}

	fmt.Printf("Subscribed to logs of contract %s via %s\n",contract.Hex(),rpcURL)
	fmt.Printf("Listening for Transfer/Approval events...\n\n")

	sigCh:=make(chan os.Signal,1)
	signal.Notify(sigCh,syscall.SIGINT,syscall.SIGTERM)
	for{
		select{
		case vLog:=<-logsCh:
			// 解析日志事件
			parseLogEvent(&vLog,parsedABI)
		case err:=<-sub.Err():
			log.Printf("subscription error:%v",err)
			return
		case sig:=<-sigCh:
			fmt.Printf("received signal %s,shutting down...\n",sig.String())
			return
		case <-ctx.Done():
			fmt.Println("context cancelled, exiting")
			return
		}
	}
}

// parseLogEvent 解析 Transfer/Approval 事件
func parseLogEvent(vLog *types.Log, parsedABI abi.ABI){
	if len(vLog.Topics)==0{
		return
	}

	eventTopic:=vLog.Topics[0]
	var eventName string
	var eventSig  abi.Event

	// 匹配事件签名
	for name,event :=range parsedABI.Events{
		eventSigHash:=crypto.Keccak256Hash([]byte(event.Sig))
		if eventSigHash==eventTopic{
			eventName = name
			eventSig = event
			break
		}
	}

	if eventName==""{
		fmt.Printf("[%s] Unknown Event - Block: %d, Tx: %s, Topic[0]:%s \n",
	    time.Now().Format(time.RFC3339),vLog.BlockNumber,vLog.TxHash.Hex(), eventTopic.Hex())
	    return
	}

	// 输出事件基本信息
	fmt.Printf("----------------------------------------\n")
	fmt.Printf("[%s] Event:%s \n",time.Now().Format(time.RFC3339),eventName)
	fmt.Printf("  Block Number  :%d\n",vLog.BlockNumber)
	fmt.Printf("  Tx Hash       :%s\n",vLog.TxHash.Hex())
	fmt.Printf("  Log Index     :%d\n",vLog.Index)
	fmt.Printf("  Contract      :%s\n",vLog.Address.Hex())
	fmt.Printf("  Topics Count  :%d\n",len(vLog.Topics))

	// 解析 indexed 参数（Topics[1..N]）
	fmt.Printf("\n   Indexed Parameters (from Topics):\n")
	indexedParamIndex:=0
	for i,input:=range eventSig.Inputs{
		if !input.Indexed{
			continue
		}
		topicIndex:=1+indexedParamIndex
		indexedParamIndex++

		if topicIndex >=len(vLog.Topics){
			continue
		}
		topic:=vLog.Topics[topicIndex]
		fmt.Printf("    [%d]%s(%s):",i+1,input.Name,input.Type)
		switch input.Type.T {
		case abi.AddressTy:
			addr:=common.BytesToAddress(topic.Bytes())
			fmt.Printf("%s\n",addr.Hex())
		case abi.IntTy,abi.UintTy:
			value:=new(big.Int).SetBytes(topic.Bytes())
			fmt.Printf("%s\n",value.String())
		case abi.BoolTy:
			fmt.Printf("%t\n",topic[31]!=0)
		default:
			fmt.Printf("%s(raw)\n",topic.Hex())
		}
	}

	// 解析非 indexed 参数（Data 字段）
	if len(vLog.Data)>0{
		fmt.Printf("\n Non-Indexed Parameters (from Data):\n")
		values,err:=parsedABI.Unpack(eventName,vLog.Data)
		if err!=nil{
			fmt.Printf("   Error decoding data: %v\n",err)
		}else{
			nonIndexedIdx:=0
			for i,input :=range eventSig.Inputs{
				if !input.Indexed{
					if nonIndexedIdx < len(values){
						value:=values[nonIndexedIdx]
						fmt.Printf("   [%d]%s(%s):",i+1,input.Name,input.Type)
						switch v:=value.(type){
						case *big.Int:
							// 处理代币小数（ERC20 通常 6/18 位小数）
							fmt.Printf("%s\n",v.String())
						case common.Address:
							fmt.Printf("%s\n",v.Hex())
						case []byte:
							fmt.Printf("0x%x\n",v)
						default:
							fmt.Printf("%v\n",v)
						}
						nonIndexedIdx++
					}
				}
			}
		}
	}else{
		fmt.Printf("\n Non-Indexed Parameters: None\n")
	}
	fmt.Printf("--------------------------------------------------\n\n")
}