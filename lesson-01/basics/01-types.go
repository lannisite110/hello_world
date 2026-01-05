package main

import (
	"fmt"
	"unsafe"
)

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
	var u8 uint8 = 255
	var u16 uint16 = 65535
	var u32 uint32 = 4000000000
	var u64 uint64 = 18440000000000000000
	fmt.Printf("uint8:%d,uint16:%d,uint32:%d,uint64:%d\n", u8, u16, u32, u64)

	//类型别名
	var b byte = 65
	var r rune = '中'
	fmt.Printf("byte:%d(%c),rune:%d(%c)\n", b, b, r, r)

	//显示类型占用的内存大小
	fmt.Printf("\n类型大小：\n")
	fmt.Printf("int8 size:%d bytes\n", unsafe.Sizeof(count))
	fmt.Printf("int16 size: %d bytes\n", unsafe.Sizeof(number))
	fmt.Printf("int32 size: %d bytes \n", unsafe.Sizeof(data))
	fmt.Printf("int64 size: %d bytes \n", unsafe.Sizeof(bigNumber))
	fmt.Printf("bool size: %d bytes \n", unsafe.Sizeof(isTrue))

	//浮点数
	var price float32 = 99.99
	var pi float64 = 3.141559265359
	fmt.Printf("\n float3：%.2f, float54: %.10f \n", price, pi)
	fmt.Printf("float32 size: %d bytes, float64 size:%d bytes \n", unsafe.Sizeof(price), unsafe.Sizeof(pi))

	//字符串
	var name string = "Golang"
	var greeting = "hello world"
	fmt.Printf("string : %s \n", name)
	fmt.Printf("类型推断： %s \n", greeting)

	//字符串操作
	str := "hello 世界"
	fmt.Printf("字符串长度：%d(字节数、n)", len(str))
	firstByte := str[0]
	fmt.Printf("第一个字节：%d('%c')\n", firstByte, firstByte)

	//原始字符串字面量
	raw := `这是
	一个多行
	字符串`
	fmt.Println("\n 原始字符串:")
	fmt.Println(raw)

	//字符（使用rune)
	var char rune = 'A'
	fmt.Printf("rune: %c \n", char)

	//复数类型
	var c1 complex64 = 1 + 2i
	var c2 complex128 = 3.14 + 6.28i
	c3 := complex(5.0, 10.0)
	fmt.Println("\n复数类型:")
	fmt.Printf("complex64:%v, 实部:%.1f， 虚部：%.1f  \n", c1, real(c1), imag(c1))
	fmt.Printf("complex128:%v, 实部:%.1f , 虚部：%.1f \n", c2, real(c2), imag(c2))
	fmt.Printf("使用complex函数：%v \n", c3)
	fmt.Printf("complex64 size:%d byte, complex128 size:%d bytes \n",
		unsafe.Sizeof(c1), unsafe.Sizeof(c2))

	//演示：为什么int8的范围是-128到127？
	fmt.Printf("\n计算： 2的次方=%d\n", 1<<8)
	fmt.Printf("int8 可以表示%d 个不同的值\n", 1<<8)
	fmt.Printf("int8可以表示%d 个不同的值\n", 1<<8)
	fmt.Printf("int8范围：%d 到 %d \n", -1<<7, 1<<7-1)

	//演示int8的边界值
	fmt.Printf("\n int8 边界值演示：")
	var minInt8 int8 = -128
	var maxInt8 int8 = 127
	fmt.Printf("最小值：%d \n", minInt8)
	fmt.Printf("最大值：%d \n", maxInt8)
	fmt.Printf("\n ===类型区分正负的关键演示===\n")
	//这两个值的二进制表示完全相同
	var signed int8 = -128
	var unsigned uint8 = 128
	fmt.Printf("有符号 int8的-128：\n")
	fmt.Printf(" 十进制：%d \n", signed)
	fmt.Printf(" 二进制：%08b \n", uint8(signed))

	fmt.Printf("\n 无符号 uint8的128\n")
	fmt.Printf(" 十进制：%d \n", unsigned)
	fmt.Printf(" 二进制：%08b \n", unsigned)

	fmt.Println("\n关键理解")
	fmt.Println(" 二进制表示相同都是10000000")
	fmt.Println(" 类型决定了如何解读这个二进制")
	fmt.Println(" int8的最高位被看作符号位（1表示负数）")
	fmt.Println(" uint8 的所有位都被看做数值（没有符号位）")
	fmt.Println("  因此编辑器通过类型声明知道应该用那种方式解读")
}

// 数组示例
func arrayDemo() {
	fmt.Println("\n===数组示例===")
	//声明数组
	var arr [5]int
	fmt.Printf("声明但不初始化:%v \n", arr)

	arr = [5]int{1, 2, 3, 4, 5}
	fmt.Printf("初始化：%v \n", arr)
	fmt.Println(arr[0])
	fmt.Println(arr[1])
	fmt.Println(arr[2])

	//自动推断长度
	arr2 := [...]int{1, 2, 3}
	fmt.Printf("自动长度：%v,长度：%d \n", arr2, len(arr2))
}

func sliceDemo() {
	fmt.Println("\n===切片示例===")
	//声明切片
	var slice []int
	fmt.Printf("空切片：%v,nil:%t \n", slice, slice == nil)
	//使用make创建切片
	slice = make([]int, 3, 5)
	fmt.Printf("make创建：%v,长度：%d, 容量：%d \n", slice, len(slice), cap(slice))

	//直接初始化
	slice = []int{1, 2, 3, 4, 5}
	fmt.Printf("初始化：%v \n", slice)

	slice = append(slice, 6)
	fmt.Printf("追加后：%v,容量：%d \n", slice, cap(slice))

	//切片截取
	subSlice := slice[1:3]
	fmt.Printf("切片[1:3]: %v \n", subSlice)

	//切片共享底层数据
	subSlice[0] = 999
	fmt.Printf("修改subSlice后：%v \n", slice)
}

// 映射示例
func mapDemo() {
	fmt.Println("\n===映射示例===")

	//声明映射
	var m map[string]int
	fmt.Printf("空映射：%v,nil:%t \n", m, m == nil)

	//使用make创建映射
	m = make(map[string]int)
	fmt.Printf("make创建：%v \n", m)
	m["apple"] = 5
	m["banana"] = 10
	fmt.Printf("添加后：%v \n", m)

	//直接初始化
	m2 := map[string]int{
		"apple":  10,
		"banana": 5,
		"orange": 8,
	}
	fmt.Printf("直接初始化：%v \n", m2)
	//读取值
	value := m2["apple"]
	fmt.Printf("apple的值：%d \n", value)
	value2 := m2["orage"]
	fmt.Printf("orange的值：%d \n", value2)

	m2["orange"] = 100
	fmt.Printf("orange的值：%d\n", m2["orange"])
	//检测key是否存在
	value, ok := m2["grape"]
	fmt.Printf("grape存在：%t，值：%d\n", ok, value)

	valueA, okA := m2["apple"]
	fmt.Printf("apple存在：%t,值：%d \n", okA, valueA)

	fmt.Println("遍历映射")
	for key, value := range m2 {
		fmt.Printf("%s:%d\n", key, value)
	}
}

func main() {
	//basicTypesDemo()
	//arrayDemo()
	//sliceDemo()
	mapDemo()
}
