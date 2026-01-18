package conc

import (
	"sync"
	"testing"
)

type testWorker struct {
	id    int
	count int
	done  chan<- struct{}
}

func (w *testWorker) Stat() string {
	return ""
}

func (w *testWorker) ID() int {
	return w.id
}

func (w *testWorker) IncrementTasks() {
	w.count++
	if w.done != nil {
		w.done <- struct{}{}
	}
}

func (w *testWorker) Work(task int, result chan<- int) {
	result <- task * 2
}

func (w *testWorker) Close() error {
	return nil
}

func TestPoolProcessesTasks(t *testing.T) {
	done := make(chan struct{}, 2)
	workers := []Worker[int, int]{
		&testWorker{id: 1, done: done},
		&testWorker{id: 2, done: done},
	}

	pool := NewWorkerPool[int, int](workers)
	pool.Create()

	results := make(chan int, 2)
	pool.Work(2, results)
	pool.Work(5, results)

	got1 := <-results
	got2 := <-results
	if got1+got2 != 14 {
		t.Fatalf("unexpected results sum: %d", got1+got2)
	}

	for i := 0; i < 2; i++ {
		<-done
	}

	for _, w := range workers {
		tw := w.(*testWorker)
		if tw.count == 0 {
			t.Fatalf("expected worker %d to have processed at least one task", tw.id)
		}
	}
}

func TestPoolHandlesConcurrentTasks(t *testing.T) {
	done := make(chan struct{}, 6)
	workers := []Worker[int, int]{
		&testWorker{id: 1, done: done},
		&testWorker{id: 2, done: done},
		&testWorker{id: 3, done: done},
	}

	pool := NewWorkerPool[int, int](workers)
	pool.Create()

	results := make(chan int, 6)

	var wg sync.WaitGroup

	for i := 1; i <= 6; i++ {
		wg.Add(1)
		go func(task int) {
			defer wg.Done()
			pool.Work(task, results)
		}(i)
	}

	wg.Wait()
	count := 0
	for i := 0; i < 6; i++ {
		<-results
		count++
	}
	if count != 6 {
		t.Fatalf("expected 6 results, got %d", count)
	}

	for i := 0; i < 6; i++ {
		<-done
	}
}
