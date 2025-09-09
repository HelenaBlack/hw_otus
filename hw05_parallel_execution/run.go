// Package hw05parallelexecution реализует параллельное выполнение задач с ограничением по ошибкам.
package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
)

// ErrErrorsLimitExceeded возвращается, если количество ошибок превысило лимит.
var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

// Task представляет собой функцию, которая выполняет некоторую работу и возвращает ошибку, если что-то пошло не так.
type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n int, m int) error {
	if m <= 0 {
		m = 1
	}

	if len(tasks) == 0 || n <= 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	taskChan := make(chan Task, len(tasks))
	sendTasks(ctx, taskChan, tasks)

	var wg sync.WaitGroup
	var mu sync.Mutex
	errorCount := 0
	completed := 0

	processError := func(err error) bool {
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			errorCount++
			if m > 0 && errorCount >= m {
				cancel()
				return true
			}
		}
		completed++
		return false
	}

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case task, ok := <-taskChan:
					if !ok {
						return
					}
					err := task()
					if processError(err) {
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	wg.Wait()
	mu.Lock()
	defer mu.Unlock()
	if m > 0 && errorCount >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func sendTasks(ctx context.Context, taskChan chan<- Task, tasks []Task) {
	defer close(taskChan)
	for _, task := range tasks {
		select {
		case taskChan <- task:
		case <-ctx.Done():
			return
		}
	}
}
