package main

import (
	"strings"  // 补充导入 strings 包
	"testing"
	"math/big" // 补充导入 big 包

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// 预解析 ABI（测试用）
var testABI abi.ABI

func init() {
	// 初始化测试用的 ABI
	var err error
	testABI, err = abi.JSON(strings.NewReader(erc20ABIJSON))
	if err != nil {
		panic("failed to parse test ABI: " + err.Error())
	}
}

// TestParseLogEvent_KnownEvent 测试解析已知的 Approval 事件
func TestParseLogEvent_KnownEvent(t *testing.T) {
	// 1. 构造模拟的 Approval 事件日志
	// 计算 Approval 事件的签名哈希
	approvalSig := "Approval(address,address,uint256)"
	approvalTopic0 := crypto.Keccak256Hash([]byte(approvalSig))

	// 构造 indexed 参数（from、to、value）
	fromAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	toAddr := common.HexToAddress("0x0987654321098765432109876543210987654321")
	value := common.BigToHash(big.NewInt(1000000000000000000)) // 1 ETH

	// 构造日志
	mockLog := &types.Log{
		Address:     common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdef"),
		BlockNumber: 20000000,
		TxHash:      common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111"),
		Index:       0,
		Topics: []common.Hash{
			approvalTopic0,    // Topic0: 事件签名
			common.BytesToHash(fromAddr.Bytes()), // Topic1: from
			common.BytesToHash(toAddr.Bytes()),   // Topic2: to
			value,              // Topic3: value
		},
		Data: []byte{}, // 无 non-indexed 参数
	}

	// 2. 执行解析函数（不会 panic 即表示逻辑正常）
	parseLogEvent(mockLog, testABI)

	// 3. 验证核心逻辑（这里主要验证无报错，如需更细粒度验证可捕获输出）
	// 检查事件名称是否能被识别（需修改 parseLogEvent 返回事件名，或捕获 stdout）
}

// TestParseLogEvent_UnknownEvent 测试解析未知事件
func TestParseLogEvent_UnknownEvent(t *testing.T) {
	// 构造未知事件的日志（Topic0 不是 Approval 的签名）
	mockLog := &types.Log{
		Topics: []common.Hash{
			crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)")), // 未知事件
		},
		BlockNumber: 20000000,
		TxHash:      common.HexToHash("0x2222222222222222222222222222222222222222222222222222222222222222"),
	}

	// 执行解析函数，验证不会 panic
	parseLogEvent(mockLog, testABI)
}

// TestParseLogEvent_EmptyTopics 测试空 Topics 的日志
func TestParseLogEvent_EmptyTopics(t *testing.T) {
	mockLog := &types.Log{
		Topics: []common.Hash{}, // 空 Topics
	}

	// 执行解析函数，验证直接返回且无报错
	parseLogEvent(mockLog, testABI)
}