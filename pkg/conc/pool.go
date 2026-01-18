package conc

import (
	"fmt"
)

type PoolHandler[T any, R any] interface {
	Handle(
		task T,
		result chan<- R,
	)
}
type Worker[T any, R any] interface {
	Stat() string
	ID() int
	IncrementTasks()
	Work(task T, result chan<- R)
	Close() error
}

type Pool[T any, R any] struct {
	workers []Worker[T, R]
	Name    string
	pool    chan Worker[T, R]
}

func NewWorkerPool[T any, R any](workers []Worker[T, R]) *Pool[T, R] {
	return &Pool[T, R]{
		workers: workers,
		pool:    make(chan Worker[T, R], len(workers)),
	}
}

func (p *Pool[T, R]) Create() {
	for _, worker := range p.workers {
		p.pool <- worker
	}
}

func (p *Pool[T, R]) Work(task T, result chan<- R) {
	worker := <-p.pool
	go func(w Worker[T, R], t T, result chan<- R) {
		w.Work(t, result)
		w.IncrementTasks()
		p.pool <- w
	}(worker, task, result)
}

func (p *Pool[T, R]) Wait() error {
	for _, worker := range p.workers {
		if err := worker.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Pool[T, R]) PrintStats() {
	report := fmt.Sprintf("------------------------------- Worker Pool: %s Stats -------------------------------\n", p.Name)
	for _, w := range p.workers {
		report += w.Stat() + "\n"
	}
	fmt.Println(report)
}

type HandlerFunc[T any, R any] func(
	task T,
	result chan<- R,
)

func (f HandlerFunc[T, R]) Handle(
	task T,
	result chan<- R,
) {
	f(task, result)
}
