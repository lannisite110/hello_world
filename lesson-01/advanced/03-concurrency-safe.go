package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// 使用Mutex的线程安全计数器
type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func NewSafeCounter() *SafeCounter {
	return &SafeCounter{}
}
func (sc *SafeCounter) Increment(m int) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.count++
	fmt.Printf("goroutine %d incremented counter %d\n", m, sc.count)
}
func (sc *SafeCounter) GetCount() int {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.count
}

func mutexDemo() {
	fmt.Println("mutex示例")
	counter := NewSafeCounter()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			counter.Increment(3)
		}(i)
	}
	wg.Wait()
	fmt.Printf("最终计数：%d (期望：100)\n", counter.GetCount())
}

// 使用RWMutex
type SafeMap struct {
	mu   sync.RWMutex
	data map[string]int
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		data: make(map[string]int),
	}
}

func (sm *SafeMap) Read(key string) int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.data[key]
}

func (sm *SafeMap) Write(key string, value int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[key] = value
}

func (sm *SafeMap) WriteVerbose(label, key string, value int, hold time.Duration) {
	fmt.Printf("[%s]准备写入 %q \n", label, key)
	sm.mu.Lock()
	fmt.Printf("[%s]获得写锁，开始写入%q = %d \n", label, key, value)
	if hold > 0 {
		time.Sleep(hold)
	}
	sm.data[key] = value
	fmt.Printf("[%s]写入完成，释放写锁", label)
	sm.mu.Unlock()
}

func (sm *SafeMap) ReadVerbose(label, key string, hold time.Duration) int {
	fmt.Printf("[%s]等待读锁以读取%q \n", label, key)
	sm.mu.RLock()
	fmt.Printf("[%s]获得读锁，读取%q \n", label, key)
	if hold > 0 {
		time.Sleep(hold)
	}
	value, ok := sm.data[key]
	if ok {
		fmt.Printf("[%s]读取结果 %q = %d \n", label, key, value)
	} else {
		fmt.Printf("[%s]%q 尚未写入 \n", label, key)
		value = -1
	}
	fmt.Printf("[%s]释放读锁 %q \n", label, key)
	sm.mu.RUnlock()
	return value
}

func (sm *SafeMap) GetAll() map[string]int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	result := make(map[string]int)
	for k, v := range sm.data {
		result[k] = v
	}
	return result
}

func rwmutexDemo() {
	fmt.Println("RWmutex示例")
	sm := NewSafeMap()
	var wg sync.WaitGroup
	fmt.Println("场景1：写独占锁，读操作需要等待写结束")
	wg.Add(1)
	go func() {
		defer wg.Done()
		sm.WriteVerbose("Writer#1", "shared", 42, 3000*time.Millisecond)
	}()
	time.Sleep(50 * time.Millisecond)

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sm.ReadVerbose(fmt.Sprintf("Reader#%d", id), "shared", 1500*time.Millisecond)
		}(i)
	}
	fmt.Println("场景2：只有读操作可以并发执行")
	for i := 4; i < 6; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sm.ReadVerbose(fmt.Sprintf("Reader#%d", id), "shared", 100*time.Millisecond)
		}(i)
	}
	wg.Wait()
	fmt.Println("\n最终数据")
	fmt.Println()
}

// context超时控制
func contextTimeoutDemo() {
	fmt.Println("context超时示例")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	ch := make(chan string)

	go func() {
		time.Sleep(2000 * time.Millisecond)
		ch <- "result"
	}()
	select {
	case result := <-ch:
		fmt.Println("收到：", result)
	case <-ctx.Done():
		fmt.Println("超时：", ctx.Err())
	}
	fmt.Println()
}

// Context取消控制
func contextCancelDemo() {
	fmt.Println("Context取消示例")
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				fmt.Println("Goroutine被取消：", ctx.Err())
				return
			default:
				fmt.Printf("工作中：%d...\n", i)
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
	time.Sleep(3 * time.Second)
	fmt.Println("取消所有goroutine")
	cancel()
	time.Sleep(500 * time.Millisecond)
	fmt.Println()
}

type SafeLogger struct {
	mu      sync.RWMutex
	logs    []string
	maxSize int
}

func NewSafeLogger(maxSize int) *SafeLogger {
	return &SafeLogger{
		logs:    make([]string, 0),
		maxSize: maxSize,
	}
}

func (s1 *SafeLogger) Log(message string) {
	s1.mu.Lock()
	defer s1.mu.Unlock()
	s1.logs = append(s1.logs, message)
	if len(s1.logs) > s1.maxSize {
		s1.logs = s1.logs[1:]
	}
}

func (s1 *SafeLogger) GetAllLogs() []string {
	s1.mu.Lock()
	defer s1.mu.Unlock()

	result := make([]string, len(s1.logs))
	copy(result, s1.logs)
	return result
}

func loggerDemo() {
	fmt.Println("并发安全日志示例")
	logger := NewSafeLogger(10)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			logger.Log(fmt.Sprintf("Log message %d from goroutine %d", id, id))
		}(i)
	}
	wg.Wait()
	//读取所有日志
	logs := logger.GetAllLogs()
	fmt.Printf("总共：%d 条日志(只保留最后10条)：\n", len(logs))
	for _, log := range logs {
		fmt.Println("", log)
	}
	fmt.Println()
}
func main() {
	//mutexDemo()
	//rwmutexDemo()
	//contextTimeoutDemo()
	//contextCancelDemo()
	loggerDemo()
}
