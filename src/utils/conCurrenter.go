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
		timeout:     30 * time.Second, // 默认超时
	}
}

func NewConCurrenterWithTimeout[T any](concurrency int, timeout time.Duration) *ConCurrenter[T] {
	return &ConCurrenter[T]{
		concurrency: concurrency,
		wg:          &sync.WaitGroup{},
		timeout:     timeout,
	}
}

// 优化后的运行方法
func (c *ConCurrenter[T]) Run(params []T, f func(T) error) error {
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
		go c.worker(ctx, taskCh, errCh, f)
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

	// 收集错误
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
func (c *ConCurrenter[T]) worker(ctx context.Context, taskCh <-chan T, errCh chan<- error, f func(T) error) {
	defer c.wg.Done()

	for {
		select {
		case task, ok := <-taskCh:
			if !ok {
				return // 通道已关闭
			}

			// 执行任务
			if err := f(task); err != nil {
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

// 带结果收集的版本
func (c *ConCurrenter[T]) RunWithResults(params []T, f func(T) (interface{}, error)) ([]interface{}, error) {
	if len(params) == 0 {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	taskCh := make(chan T, c.concurrency)
	resultCh := make(chan interface{}, len(params))
	errCh := make(chan error, len(params))

	// 启动工作协程
	for i := 0; i < c.concurrency && i < len(params); i++ {
		c.wg.Add(1)
		go c.workerWithResult(ctx, taskCh, resultCh, errCh, f)
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

	// 等待完成
	go func() {
		c.wg.Wait()
		close(resultCh)
		close(errCh)
	}()

	// 收集结果和错误
	var results []interface{}
	var errors []error

	for {
		select {
		case result, ok := <-resultCh:
			if !ok {
				goto checkErrors
			}
			results = append(results, result)
		case err, ok := <-errCh:
			if !ok {
				goto checkErrors
			}
			if err != nil {
				errors = append(errors, err)
			}
		}
	}

checkErrors:
	if len(errors) > 0 {
		return results, fmt.Errorf("执行过程中发生 %d 个错误", len(errors))
	}

	return results, nil
}

func (c *ConCurrenter[T]) workerWithResult(ctx context.Context, taskCh <-chan T, resultCh chan<- interface{}, errCh chan<- error, f func(T) (interface{}, error)) {
	defer c.wg.Done()

	for {
		select {
		case task, ok := <-taskCh:
			if !ok {
				return
			}

			result, err := f(task)
			if err != nil {
				select {
				case errCh <- err:
				case <-ctx.Done():
					return
				}
			} else {
				select {
				case resultCh <- result:
				case <-ctx.Done():
					return
				}
			}

		case <-ctx.Done():
			return
		}
	}
}
