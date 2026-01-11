package main

import (
	"fmt"
	"sync"
	"time"
)

// 生产者消费者模式示例
func producerConsumerDemo() {
	fmt.Println("生产者消费者模式")
	//创建带缓冲的channel
	jobs := make(chan int, 5)
	results := make(chan int, 5)

	//启动多个消费者
	numWorkers := 3
	var wg sync.WaitGroup

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				fmt.Printf("worker %d processing job %d \n", id, job)
				time.Sleep(100 * time.Millisecond)
				results <- job * 2
			}
		}(w)
	}

	//生产者
	go func() {
		defer close(jobs)
		for i := 1; i <= 10; i++ {
			fmt.Printf("Producing job %d\n", i)
			jobs <- i
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()
	//读取结果
	fmt.Println("Results:")
	for result := range results {
		fmt.Printf("Result:%d\n", result)
	}
}

// worker pool模式
type Job struct {
	ID      int
	Payload string
}

type WorkerPool struct {
	numWorkers int
	jobs       chan Job
	results    chan Job
	wg         sync.WaitGroup
}

func NewWorkerPool(numWorkers, jobBuffer int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		jobs:       make(chan Job, jobBuffer),
		results:    make(chan Job, jobBuffer),
	}
}
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go func(id int) {
			defer wp.wg.Done()
			for job := range wp.jobs {
				job.Process(id)
				wp.results <- job
			}
		}(i)
	}
}
func (j *Job) Process(workerID int) {
	fmt.Printf("worker %d processing job %d:%s\n", workerID, j.ID, j.Payload)
	time.Sleep(200 * time.Millisecond)
}

func (wp *WorkerPool) AddJob(job Job) {
	wp.jobs <- job
}

func (wp *WorkerPool) Close() {
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
}

func (wp *WorkerPool) GetResults() <-chan Job {
	return wp.results
}

func workerPoolDemo() {
	fmt.Println("worker pool 模式")
	pool := NewWorkerPool(3, 5)
	pool.Start()
	//添加任务
	for i := 1; i <= 10; i++ {
		job := Job{
			ID:      i,
			Payload: fmt.Sprintf("Task %d", i),
		}
		pool.AddJob(job)
	}

	//关闭并收集结果
	go func() {
		for result := range pool.GetResults() {
			fmt.Printf("Completed: Job %d\n", result.ID)
		}
	}()
	time.Sleep(3 * time.Second)
	pool.Close()
}

func main() {
	producerConsumerDemo()
	workerPoolDemo()
}
