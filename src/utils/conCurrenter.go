package utils

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type ConCurrenter[T any] struct {
	wg          *sync.WaitGroup
	concurrency int
	timeout     time.Duration
}

func NewConCurrenter[T any](concurrency int) *ConCurrenter[T] {
	return &ConCurrenter[T]{
		concurrency: concurrency,
		wg:          &sync.WaitGroup{},
		timeout:     60 * time.Second, // 默认超时,防止卡死
	}
}

func NewConCurrenterWithTimeout[T any](concurrency int, timeout time.Duration) *ConCurrenter[T] {
	return &ConCurrenter[T]{
		concurrency: concurrency,
		wg:          &sync.WaitGroup{},
		timeout:     timeout, // 超时时间
	}
}

// 并发器主函数
func (c *ConCurrenter[T]) Run(params []T, handler func(T) error) error {
	if len(params) == 0 {
		return nil
	}

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	// 创建任务和错误通道
	taskCh := make(chan T, c.concurrency)
	errCh := make(chan error, len(params))

	// 启动工作协程
	for i := 0; i < c.concurrency && i < len(params); i++ {
		c.wg.Add(1)
		go c.worker(ctx, taskCh, errCh, handler)
	}

	// 发送任务
	go func() {
		defer close(taskCh)
		for _, param := range params {
			select {
			case taskCh <- param:
			case <-ctx.Done():
				return
			}
		}
	}()

	// 等待所有任务完成
	go func() {
		c.wg.Wait()
		close(errCh)
	}()

	// 收集错误，使用错误通道来阻塞，直到所有任务完成
	var errors []error
	for err := range errCh {
		if err != nil {
			errors = append(errors, err)
		}
	}

	// 检查是否有错误
	if len(errors) > 0 {
		return fmt.Errorf("执行过程中发生 %d 个错误", len(errors))
	}

	return nil
}

// 工作协程
func (c *ConCurrenter[T]) worker(ctx context.Context, taskCh <-chan T, errCh chan<- error, handler func(T) error) {
	defer c.wg.Done()

	for {
		select {
		case task, ok := <-taskCh:
			if !ok {
				return // 通道已关闭
			}

			// 执行任务
			if err := handler(task); err != nil {
				select {
				case errCh <- err:
				case <-ctx.Done():
					return
				}
			}

		case <-ctx.Done():
			return // 超时或取消
		}
	}
}
