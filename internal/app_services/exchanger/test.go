package exchanger

import (
	"context"
	"fmt"
	"time"
)

type TestExchanger[T any] struct {
	Name     string
	interval time.Duration
	cancel   context.CancelFunc
	generate func(name string) T
}

func NewTestExchanger[T any](name string, interval time.Duration, generate func(name string) T) (*TestExchanger[T], error) {
	if interval <= 0 {
		interval = time.Millisecond * 500
	}
	if generate == nil {
		return nil, fmt.Errorf("generate function is required")
	}
	return &TestExchanger[T]{
		Name:     name,
		interval: interval,
		generate: generate,
	}, nil
}

func (t *TestExchanger[T]) Stream(ctx context.Context, out chan<- Task[T], results chan<- Result) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	t.cancel = cancel
	ticker := time.NewTicker(t.interval)

	for {
		select {
		case <-ticker.C:
			out <- WrapTask(t.Name, t.generate(t.Name))
		case <-ctx.Done():
			results <- Result{Name: t.Name, Err: nil}
			return
		}
	}
}

func (t *TestExchanger[T]) Stop() error {
	t.cancel()
	return nil
}

func generateTestData(name string) string {

	return string("qwewqe")
}
