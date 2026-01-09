package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// 基本goroutine
func basicGoroutine() {
	fmt.Println("基本Goroutine")
	// 启动goroutine
	go sayHello("goroutine 1")
	go sayHello("goroutine 2")
	go sayHello("goroutine 3")

	//等待goroutine
	time.Sleep(time.Second)
	fmt.Println()
}

func sayHello(name string) {
	fmt.Printf("Hello form %s \n", name)
}

//使用waitGroup

func waitGroupDemo() {
	fmt.Println("WaitGroup示例")
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Task %d started \n", id)
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Task %d finished\n", id)
		}(i)
	}
	wg.Wait()
	fmt.Println("All tasks completed")
}

// channel通信
func channelDemo() {
	fmt.Println("channel 示例")
	//无缓冲
	ch := make(chan string)
	go func() {
		defer close(ch)
		ch <- "hello"
		ch <- "world"
	}()
	time.Sleep(3000 * time.Millisecond)

	for msg := range ch {
		fmt.Println("Received:", msg)
	}
	fmt.Println()
}

func bufferedChannelDemo() {
	fmt.Println("缓冲channel示例")
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	//ch<-4
	fmt.Println("写入完成")
	fmt.Println("读取：", <-ch)
	fmt.Println("读取：", <-ch)
	fmt.Println("读取：", <-ch)
	fmt.Println()
}

func bufferedChannelDemoNew() {
	fmt.Println("缓冲Channel示例")
	ch := make(chan int, 3)
	defer close(ch)
	//启动读取goroutine
	go func() {
		for v := range ch {
			fmt.Println("routine1读取：", v)
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	go func() {
		for v := range ch {
			fmt.Println("routine2 读取：", v)
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	for i := 1; i < 40; i++ {
		ch <- i
		fmt.Println("写入：", i)
	}
}

// select语句
func selectDemo() {
	fmt.Println("select示例")
	ch1 := make(chan string)
	ch2 := make(chan string)
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "from ch1"
	}()

	go func() {
		time.Sleep(1 * time.Second)
		ch2 <- "from ch2"
	}()
	time.Sleep(1000 * time.Millisecond)

	select {
	case msg1 := <-ch1:
		fmt.Println(msg1)
	case msg2 := <-ch2:
		fmt.Println(msg2)
	}
}

// select with timeout
func timeoutDemo() {
	fmt.Println("select超时示例")
	ch := make(chan string)
	go func() {
		time.Sleep(2 * time.Second)
		ch <- "result"
	}()

	select {
	case msg := <-ch:
		fmt.Println("收到：", msg)
	case <-time.After(2 * time.Second):
		fmt.Println("超时了")
	}
	fmt.Println()
}

// select with default
func nonBlockingDemo() {
	fmt.Println("非阻塞select示例")
	ch := make(chan int)
	select {
	case value := <-ch:
		fmt.Println("收到：", value)
	default:
		fmt.Println("没有值可读(非阻塞)")
	}

	//非阻塞发送
	select {
	case ch <- 42:
		fmt.Println("发送成功")
	default:
		fmt.Println("channel已满，发送失败")
	}
}

// 循环监听多个channel
func loopSelectDemo() {
	fmt.Println("循环select示例")
	ch1 := make(chan int)
	ch2 := make(chan string)
	go func() {
		for i := 0; i < 5; i++ {
			ch1 <- i
			time.Sleep(100 * time.Millisecond)
		}
		close(ch1)
	}()

	go func() {
		for i := 0; i < 3; i++ {
			ch2 <- fmt.Sprintf("msg-%d", i)
			time.Sleep(150 * time.Millisecond)
		}
		close(ch2)
	}()

	for ch1 != nil || ch2 != nil {
		select {
		case val, ok := <-ch1:
			if !ok {
				ch1 = nil
				continue
			}
			fmt.Println("ch1:", val)
		case msg, ok := <-ch2:
			if !ok {
				ch2 = nil
				continue
			}
			fmt.Println("ch2:", msg)
		default:
			fmt.Println("default")
			time.Sleep(50 * time.Millisecond)
		}
	}
	fmt.Println("所有channel已关闭 ")
}

// 使用退出信号关闭goroutine
func quitChannelDemo() {
	fmt.Println("Quit Channel示例")
	jobs := make(chan int)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case job := <-jobs:
				fmt.Printf("处理任务：%d \n", job)
				time.Sleep(500 * time.Millisecond)
			case <-quit:
				fmt.Println("收到退出信号")
				return
			}
		}
	}()

	for i := 0; i < 30; i++ {
		jobs <- i
	}
	fmt.Println("发送任务完成")
	quit <- struct{}{}
	time.Sleep(100 * time.Millisecond)
	fmt.Println()
}

// 关闭channel后的读取
func closedChannelDemo() {
	fmt.Println("关闭Channel示例")
	ch := make(chan int)
	close(ch)
	select {
	case val, ok := <-ch:
		fmt.Printf("val:%d,ok:%v\n", val, ok)
	default:
		fmt.Println("default")
	}
	fmt.Println()
}

// 随机选择示例
func fairnessDemo() {
	fmt.Println("select公平性示例")
	ch1 := make(chan int)
	ch2 := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			ch1 <- i
		}
		close(ch1)
	}()
	go func() {
		for i := 0; i < 5; i++ {
			ch2 <- i * 10
		}
		close(ch2)
	}()
	for ch1 != nil || ch2 != nil {
		select {
		case val, ok := <-ch1:
			if !ok {
				ch1 = nil
				continue
			}
			fmt.Printf("ch1:%d \n", val)
		case val, ok := <-ch2:
			if !ok {
				ch2 = nil
				continue
			}
			fmt.Printf("ch2:%d \n", val)
		}
	}
	fmt.Println()
}

// Context基础示例
func contexWithSelectDemo() {
	fmt.Println("context与select结合示例")
	ch := make(chan int)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	go func() {
		time.Sleep(2 * time.Second)
		ch <- 42
	}()
	select {
	case result := <-ch:
		fmt.Println("收到结果：", result)
	case <-ctx.Done():
		fmt.Println("超时：", ctx.Err())
	}
	fmt.Println()
}

// context取消信号传递
func contextCancelDemo() {
	fmt.Println("context取消信号传递")
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Goroutine被取消：", ctx.Err())
				return
			default:
				fmt.Println("工作中")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
	time.Sleep(1500 * time.Millisecond)
	fmt.Println("发送取消信号")
	cancel()
	time.Sleep(300 * time.Millisecond)
	fmt.Println()
}

// 携程泄露示例
func leakExampleRight() {
	fmt.Println("携程泄露修复示例")
	ch := make(chan int, 1)
	go func() {
		defer close(ch)
		ch <- 1
	}()
	select {
	case val := <-ch:
		fmt.Println("接收到：", val)
	case <-time.After(500 * time.Millisecond):
		fmt.Println("发送阻塞，触发超时保护")
	}
	fmt.Println()
}

func leakExample() {
	ch := make(chan int)
	go func() {
		ch <- 1
	}()

}
func main() {
	//basicGoroutine()
	//waitGroupDemo()
	//channelDemo()
	//bufferedChannelDemo()
	//bufferedChannelDemoNew()
	//timeoutDemo()
	//nonBlockingDemo()
	//loopSelectDemo()
	//closedChannelDemo()
	//fairnessDemo()
	//contexWithSelectDemo()
	//contextCancelDemo()
	//leakExample()
	leakExampleRight()
}
