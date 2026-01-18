package exchanger

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"sync"
)

type LiveExchanger[T any] struct {
	Name          string
	Host          string
	Port          string
	receivedTasks int
	wg            *sync.WaitGroup
	cancel        context.CancelFunc
	parse         func(raw string) (T, error)
}

func NewLiveExchanger[T any](name, host, port string, parse func(raw string) (T, error)) (*LiveExchanger[T], error) {
	if host == "" || port == "" {
		return nil, nil
	}
	if parse == nil {
		return nil, fmt.Errorf("parse function is required")
	}
	return &LiveExchanger[T]{
		Name:          name,
		Host:          host,
		Port:          port,
		receivedTasks: 0,
		wg:            &sync.WaitGroup{},
		parse:         parse,
	}, nil
}

func (l *LiveExchanger[T]) Stream(ctx context.Context, out chan<- Task[T], results chan<- Result) {
	ctx, cancel := context.WithCancel(ctx)

	l.cancel = cancel
	defer cancel()

	select {
	case <-ctx.Done():
		l.sendResult(results, ctx.Err())
		return
	default:
	}

	conn, err := net.Dial("tcp", net.JoinHostPort(l.Host, l.Port))
	if err != nil {
		l.sendResult(results, err)
		return
	}

	if err = l.handle(ctx, conn, out); err != nil {
		l.sendResult(results, err)
		return
	}
	l.sendResult(results, nil)
}

func (l *LiveExchanger[T]) Stop() error {
	if l.cancel != nil {
		l.cancel()
		return nil
	}

	return fmt.Errorf("exchanger %s not running", l.Name)
}

func (l *LiveExchanger[T]) handle(ctx context.Context, conn net.Conn, out chan<- Task[T]) error {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		parsed, err := l.parse(scanner.Text())
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return nil
		case out <- WrapTask(l.Name, parsed):
			l.receivedTasks++
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return fmt.Errorf("connection to exchanger %s closed", l.Name)
}

func (l *LiveExchanger[T]) sendResult(results chan<- Result, err error) {
	results <- Result{Name: l.Name, Host: l.Host, Port: l.Port, ReceivedTasks: l.receivedTasks, Err: err}
}
