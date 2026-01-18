package conc

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

type Ticker struct {
	updateTicker chan time.Duration
	stop         chan struct{}
	closed       atomic.Bool
}

func NewTicker() *Ticker {
	return &Ticker{
		updateTicker: make(chan time.Duration),
		stop:         make(chan struct{}),
		closed:       atomic.Bool{},
	}
}

func (t *Ticker) Start(ctx context.Context, interval time.Duration, handle func()) {
	ticker := time.NewTicker(interval)

	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.stop:
			return
		case <-ticker.C:
			handle()
		case interval := <-t.updateTicker:
			if interval <= 0 {
				continue
			}
			ticker.Reset(interval)
		}
	}
}

func (t *Ticker) UpdateInterval(d time.Duration) error {
	if t.closed.Load() {
		return fmt.Errorf("ticker is closed")
	}
	fmt.Println("Updating ticker interval to:", d)
	go func() {
		t.updateTicker <- d
	}()
	return nil
}

func (t *Ticker) Stop() error {
	if t.closed.Load() {
		return fmt.Errorf("ticker is closed")
	}
	t.closed.Store(true)
	go func() {
		t.stop <- struct{}{}
	}()
	return nil
}
