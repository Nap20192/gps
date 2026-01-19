package exchanger

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	pool := NewPool[string](5)

	parse := func(raw string) (string, error) {
		return raw, nil
	}

	err := pool.Add(ctx, "exchanger1", "localhost", "40101", parse)
	if err != nil {
		t.Fatal(err)
	}
	err = pool.Add(ctx, "exchanger2", "localhost", "40102", parse)
	if err != nil {
		t.Fatal(err)
	}
	err = pool.Add(ctx, "exchanger3", "localhost", "40103", parse)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		time.Sleep(9 * time.Second)
		fmt.Println("Stopping pool...")
		pool.StopPool()
	}()

	go func() {
		for range pool.Out() {
			fmt.Println("New task in pool out channel")
		}
	}()
	go func() {
		for i := range pool.Results() {
			if i.Err != nil {
				fmt.Printf("Result from %s: %s\n", i.Name, i.Err.Error())
			} else {
				fmt.Printf("Result from %s: success\n", i.Name)
			}
		}
	}()

	time.Sleep(10 * time.Second)
}
