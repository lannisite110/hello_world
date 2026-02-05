package main

import (
	"math"
	"math/big"
	"testing"
)

// TestWeiToEth 测试 Wei 转 ETH 的核心函数
func TestWeiToEth(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name     string
		weiInput *big.Int
		wantEth  *big.Float
	}{
		{
			name:     "1 ETH = 1e18 Wei",
			weiInput: big.NewInt(1e18),
			wantEth:  big.NewFloat(1.0),
		},
		{
			name:     "0 Wei",
			weiInput: big.NewInt(0),
			wantEth:  big.NewFloat(0.0),
		},
		{
			name:     "1.5 ETH = 1.5e18 Wei",
			weiInput: big.NewInt(1500000000000000000),
			wantEth:  big.NewFloat(1.5),
		},
		{
			name:     "大数测试",
			weiInput: new(big.Int).Mul(big.NewInt(1234), big.NewInt(int64(math.Pow10(18)))),
			wantEth:  big.NewFloat(1234.0),
		},
	}

	// 遍历测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := weiToEth(tt.weiInput)
			// 比较 big.Float 的值（精度到 18 位）
			if got.Cmp(tt.wantEth) != 0 {
				t.Errorf("weiToEth() = %s, want %s", got.Text('f', 18), tt.wantEth.Text('f', 18))
			}
		})
	}
}
