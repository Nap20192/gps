package exchanger

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type Exchanger[T any] interface {
	Stream(ctx context.Context, out chan<- Task[T], results chan<- Result)
	Stop() error

}
type Task[T any] struct {
	Exchanger string
	Data T
}

func WrapTask[T any](from string, data T) Task[T] {
	return Task[T]{
		Exchanger: from,
		Data: data,
	}
}

type Pool[T any] struct {
	MaxCount   int
	Exchangers map[string]Exchanger[T]
	numClients int

	wg     *sync.WaitGroup
	out    chan Task[T]
	result chan Result
	mu     sync.Mutex
}

func NewPool[T any](maxCount int) *Pool[T] {
	pool := &Pool[T]{
		MaxCount:   maxCount,
		Exchangers: make(map[string]Exchanger[T]),
		numClients: 0,
		wg:         &sync.WaitGroup{},
		out:        make(chan Task[T]),
		result:     make(chan Result),
	}
	return pool
}

func (p *Pool[T]) GetConnectedExchangers() map[string]bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	connected := make(map[string]bool)
	for name := range p.Exchangers {
		connected[name] = true
	}
	return connected
}

func (p *Pool[T]) Add(ctx context.Context, name, host, port string, parse func(raw string) (T, error)) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	n := p.numClients + 1
	if n >= p.MaxCount {
		return fmt.Errorf("max exchangers limit reached: %d", p.MaxCount)
	}
	p.numClients = n

	if _, exists := p.Exchangers[name]; exists {
		return fmt.Errorf("exchanger with name %s already exists", name)
	}

	worker, err := NewLiveExchanger(name, host, port, parse)
	if err != nil {
		return err
	}
	p.Exchangers[name] = worker

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		worker.Stream(ctx, p.out, p.result)
		p.mu.Lock()
		defer p.mu.Unlock()
		p.numClients--
		delete(p.Exchangers, name)
	}()

	return nil
}

func (p *Pool[T]) AddTest(ctx context.Context, name string, generate func(name string) T) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	n := p.numClients + 1

	if n >= p.MaxCount {
		return fmt.Errorf("max exchangers limit reached: %d", p.MaxCount)
	}

	p.numClients = n

	if _, exists := p.Exchangers[name]; exists {
		return fmt.Errorf("exchanger with name %s already exists", name)
	}

	worker, err := NewTestExchanger(name, 100*time.Millisecond, generate)
	if err != nil {
		return err
	}

	p.Exchangers[name] = worker

	p.wg.Add(1)

	go func() {
		defer p.wg.Done()
		worker.Stream(ctx, p.out, p.result)
		p.mu.Lock()
		defer p.mu.Unlock()
		p.numClients--
		delete(p.Exchangers, name)
	}()
	return nil
}

func (p *Pool[T]) Remove(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if exchanger, ok := p.Exchangers[name]; ok {
		exchanger.Stop()
		p.numClients--
		delete(p.Exchangers, name)
	}
}

func (p *Pool[T]) StopPool() {
	p.mu.Lock()
	fmt.Println(len(p.Exchangers))
	for n, exchanger := range p.Exchangers {
		slog.Warn("stopping exchanger...", "name", n)
		exchanger.Stop()
	}
	p.mu.Unlock()

	p.wg.Wait()
	close(p.out)
	close(p.result)
}

func (p *Pool[T]) Out() <-chan Task[T] {
	return p.out
}

type Result struct {
	Name          string
	Host          string
	Port          string
	ReceivedTasks int
	Err           error
}

func (p *Pool[T]) Results() <-chan Result {
	return p.result
}
