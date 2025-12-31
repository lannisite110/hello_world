package main

import "fmt"

func basicTypesDemo() {
	fmt.Println("===基础类型示例===")

	//bool
	var isTrue bool = true
	fmt.Printf("bool:%t\n", isTrue)
	//整数类型
	var age int = 25
	fmt.Printf("int:%d\n", age)

	//不同长度的整数类型
	var count int8 = 127
	var number int16 = 32767
	var data int32 = 2000000000
	var bigNumber int64 = 9000000000000000000
	fmt.Printf("int8:%d,int16:%d,int32:%d,int64:%d\n", count, number, data, bigNumber)

	//无符号整数

}

func main() {
	basicTypesDemo()
}
