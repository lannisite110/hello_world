package main

import (
	"fmt"
	"os"
)

func ifDemo() {
	fmt.Println("===if/else示例===")
	age := 20
	if age >= 18 {
		fmt.Println("已成年")
	} else {
		fmt.Println("未成年")
	}

	//支持初始化语句
	if num := 10; num > 5 {
		fmt.Printf("num(%d)大于5 \n", num)
	}

	//多条件判断
	scores := []int{95, 85, 70, 50}
	for index, score := range scores {
		fmt.Printf("索引%d:%d \n", index, score)
		if score >= 90 {
			fmt.Printf("分数%d:优秀\n", score)
		} else if score >= 80 {
			fmt.Printf("分数%d:良好\n", score)
		} else if score >= 60 {
			fmt.Printf("分数%d:及格\n", score)
		} else {
			fmt.Printf("分数%d:不及格\n", score)
		}
	}
}

func swicthDemo() {
	fmt.Printf("\n ===switch示例===")
	days := []int{1, 3, 5, 7, 10}
	for index, day := range days {
		fmt.Printf("索引%d:%d\n", index, day)
		switch day {
		case 1:
			fmt.Println("Monday")
		case 2:
			fmt.Println("Tuesday")
		case 3:
			fmt.Println("Wednesday")
		case 4:
			fmt.Println("Thursday")
		case 5:
			fmt.Println("Friday")
		case 6, 7:
			fmt.Println("周末")
		default:
			fmt.Printf("Day %d:未知\n", day)
		}
	}

	score := 85
	switch {
	case score >= 90:
		fmt.Println("等级：A")
	case score >= 80:
		fmt.Println("等级：B")
	case score >= 70:
		fmt.Println("等级:C")
	case score >= 60:
		fmt.Println("等级:D")
	default:
		fmt.Println("等级:E")
	}
}

func switchFallthroughDemo() {
	fmt.Println("\n===switch fallthrough 示例 ===")
	//1 默认行为：不穿透
	fmt.Println("1. 默认行为：不穿透")
	num := 1

	switch num {
	case 1:
		fmt.Printf("匹配到case1")
	case 2:
		fmt.Printf("匹配到case2")
	default:
		fmt.Println("默认分支")
	}
	//2. 即使添加break也是多余的
	fmt.Println("\n 2.显式break（虽然多余但是可以增加可读性）")
	num = 1
	switch num {
	case 1:
		fmt.Println("匹配到case1")
		break
	case 2:
		fmt.Println("匹配到case2")
	default:
		fmt.Println("默认分支")
	}
	//4 fallthrough 连续穿透
	fmt.Println("\n 4. fallthrough连续穿透")
	num = 1
	switch num {
	case 1:
		fmt.Println("匹配到case1")
		fallthrough
	case 2:
		fmt.Println("匹配到case2(因为fallthrough)")
		fallthrough
	case 3:
		fmt.Println("执行case3(因为fallthrough)")
	default:
		fmt.Println("默认分支")
	}

	//5 实际应用场景
	fmt.Println("\n 5. 实际应用场景：分级判断(显示所有满足的等级)")
	scores := []int{95, 85, 75, 65}
	for _, score := range scores {
		fmt.Printf("分数%d:", score)
		first := true
		switch {
		case score >= 90:
			fmt.Print("优秀")
			first = false
			fallthrough
		case score >= 100:
			if !first {
				fmt.Print(",良好")
			} else {
				fmt.Print("良好")
			}
			first = false
			fallthrough
		case score >= 60:
			if !first {
				fmt.Print(",及格")
			} else {
				fmt.Print("及格")
			}
		default:
			fmt.Print("不及格")
		}
		fmt.Println()
	}

	//6 对比：不使用fallthrough的正常分级判断
	fmt.Println("\n 6.对比：不使用fallthrough的正常分级判断")
	scores2 := []int{95, 85, 75, 55}
	for _, score := range scores2 {
		fmt.Printf(" 分数 %d:", score)
		switch {
		case score >= 90:
			fmt.Println("优秀")
		case score >= 80:
			fmt.Println("良好")
		case score >= 60:
			fmt.Println("及格")
		default:
			fmt.Println("不及格")
		}
	}
}

// for循环示例
func forloopDmeo() {
	fmt.Println("\n ===for循环示例===")
	//基本for循环
	fmt.Println("基本for循环")
	for i := 0; i < 5; i++ {
		fmt.Printf("%d", i)
	}
	fmt.Println()

	//类似while循环
	fmt.Println("类似while循环：")
	i := 0
	for i < 5 {
		fmt.Printf(" %d", i)
		i++
	}
	fmt.Println()
	//死循环需要break
	fmt.Println("循环直到满足条件")
	i = 0
	for {
		if i >= 5 {
			break
		}
		fmt.Printf("%d", i)
		i++
	}
	fmt.Println()

	//range遍历数组
	fmt.Println("range遍历数组：")
	arr := []int{10, 20, 30, 40, 50}
	for index, value := range arr {
		fmt.Printf("索引[%d]=%d\n", index, value)
	}
	//只遍历值
	fmt.Println("只遍历值：")
	for _, value := range arr {
		fmt.Printf(" %d ", value)
	}
	fmt.Println()

	//遍历字符串
	fmt.Println("遍历字符串")
	str := "hello 世界"
	for i, char := range str {
		fmt.Printf("[%d]%c \n", i, char)
	}
}

func deferDemo() {
	fmt.Println("\n ===defer示例===")
	fmt.Println("\n 7. defer用于资源清理")
	if err := readFile("02-control-flow.go"); err != nil {
		fmt.Printf("错误：%v \n", err)
	}
	if err := readFile("noneeixtent.txt"); err != nil {
		fmt.Printf("错误：%v \n", err)
	}
}

func readFile(filename string) error {
	defer func() {
		fmt.Printf("清理资源：%s\n", filename)
	}()

	fmt.Printf("准备打开：%s \n", filename)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	defer fmt.Printf("关闭文件：%s \n", filename)
	fmt.Printf("成功打开：%s \n", filename)
	return nil
}

// 演示defer在return之后执行
func returnWithDefer() int {
	defer fmt.Println("defer:在return之后执行")
	fmt.Println("return:先准备返回值")
	return 42
}

// 演示defer捕获变量的时机
func deferValueDemo() {
	i := 0
	defer fmt.Println("defer1:i=", i)
	i++
	defer fmt.Println("defer2 i=", i)
	i++
	fmt.Println("函数内：i=", i)

}

// 演示defer闭包捕获最终值
func deferClosureDemo() {
	i := 0
	i++
	i++
	defer func() {
		fmt.Println("defer闭包：i=", i)
	}()
	fmt.Println("函数内：i=", i)
}

// 演示defer在panic之后也会执行
func deferWithPanic() {
	defer func() {
		fmt.Println("清理工作：defer在panic之后执行")
	}()
	fmt.Println("开始执行")
	fmt.Println("即将panic...")
	fmt.Println("(演示结束，实际Panic会执行defer)")
}

// painc和recover示例
func panicRecoverDemo() {
	fmt.Println("\n ===panic和recover示例===")
	fmt.Println("触发Panic但使用Recover捕获")
	safeOperation()

	fmt.Println("\n 正常执行Panic")
	riskyOperation()
}

func safeOperation() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获panic：%v \n", r)
			fmt.Println("程序继续执行")
		}
	}()
	fmt.Println("即将执行Panic")
	panic("发生错误")
	fmt.Println("这行不会执行")
}

func riskyOperation() {
	fmt.Println("这会导致程序崩溃")
	panic("致命错误")
	fmt.Println("这行不会执行")
}
func main() {
	//ifDemo()
	//swicthDemo()
	//switchFallthroughDemo()
	//forloopDmeo()
	//deferDemo()
	//returnWithDefer()
	//deferValueDemo()
	//deferClosureDemo()
	//deferWithPanic()
	panicRecoverDemo()
}
