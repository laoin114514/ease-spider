package utils

import "sync"

type ConCurrenter[T any] struct {
	wg          *sync.WaitGroup
	concurrency int
}

func NewConCurrenter[T any](concurrency int) *ConCurrenter[T] {
	return &ConCurrenter[T]{concurrency: concurrency, wg: &sync.WaitGroup{}}
}

func (c *ConCurrenter[T]) Run(parms []T, f func(T)) {
	ch := make(chan T, c.concurrency-1)

	//构建并发器
	for i := 0; i < c.concurrency; i++ {
		go c.GoroutineWorker(ch, f)
	}

	//发布所有任务
	for _, parm := range parms {
		ch <- parm
	}
	c.wg.Wait()
}

func (c *ConCurrenter[T]) GoroutineWorker(ch <-chan T, f func(T)) {
	for item := range ch {
		c.wg.Add(1)
		f(item)
		c.wg.Done()
	}
}
