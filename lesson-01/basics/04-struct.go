package main

import "fmt"

// 定义结构体
type Person struct {
	Name string
	Age  int
}

// 值接收者方法
func (p Person) GetInfo() string {
	return fmt.Sprintf("%s is %d years old", p.Name, p.Age)
}

func (p Person) IncrementAgeWrong() {
	p.Age++
}

func (p *Person) IncrementAge() {
	p.Age++
}

// 指针接收者方法
func (p *Person) ChangeName(name string) {
	p.Name = name
}

// 嵌入结构体
type Employee struct {
	Person //匿名字段
	ID     string
}

func (e Employee) GetEmployeeInfo() string {
	return fmt.Sprintf("ID:%s,%s", e.ID, e.GetInfo())
}

// 接口示例
type Speaker interface {
	Speak() string
}

func (p Person) Speak() string {
	return fmt.Sprintf("hi,i'm %s", p.Name)
}

// 多态示例
func introduce(s Speaker) {
	fmt.Println(s.Speak())
}

func main() {
	fmt.Println("===结构体示例===")
	//初始化结构体
	p1 := Person{
		Name: "alice",
		Age:  25,
	}
	fmt.Println("p1:", p1)
	//简短初始化
	p2 := Person{"bob", 22}
	fmt.Println("p2:", p2)
	//部分初始化
	p3 := Person{Name: "charlie"}
	fmt.Println("p3:", p3)

	//使用值接收者方法
	fmt.Println("\n值接收者方法:")
	fmt.Println(p1.GetInfo())
	fmt.Println(p2.GetInfo())

	//使用指针接收者方法
	fmt.Println("\n指针接收者方法")
	fmt.Printf("年龄:%d \n", p1.Age)
	p1.IncrementAge()
	fmt.Printf("增加年龄后：%d \n", p1.Age)

	//尝试值接收者修改
	p1.IncrementAgeWrong()
	fmt.Printf("使用值接收者后：%d \n", p1.Age)
	//修改name
	p1.ChangeName("Alice Smith")
	fmt.Println("修改姓名后：", p1.GetInfo())

	//嵌入结构体
	fmt.Println("\n===嵌入结构体===")
	Employee := Employee{
		Person: Person{Name: "david", Age: 28},
		ID:     "E001",
	}
	fmt.Println("employee:", Employee)
	fmt.Println("直接访问嵌入的字段：", Employee.Name, Employee.Age)
	fmt.Println("调用嵌入结构图的方法：", Employee.GetInfo())
	fmt.Println("调用自己的方法：", Employee.GetEmployeeInfo())
}
