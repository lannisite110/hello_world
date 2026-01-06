package main

import "fmt"

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
}

func main() {
	//ifDemo()
	swicthDemo()
}
