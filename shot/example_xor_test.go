package shot_test

import (
	"context"
	"fmt"
	"time"

	"github.com/infastin/gorack/shot"
)

type FibXor struct {
	output chan int
	state  shot.Xor
}

func NewFibXor(ctx context.Context) *FibXor {
	return &FibXor{
		output: make(chan int, 1),
		state:  shot.NewXor(ctx),
	}
}

func (e *FibXor) Run(values [2]int, idx int) error {
	stop, err := e.state.Start(context.Background())
	if err != nil {
		return err
	}
	defer stop()

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-e.state.Context().Done():
			return e.state.Context().Err()
		case <-ticker.C:
			sum := values[0] + values[1]
			values[idx] = sum
			idx = (idx + 1) % 2
			e.output <- sum
			if sum == 377 {
				ticker.Stop()
			}
		}
	}
}

func (e *FibXor) Close() error {
	if err := e.state.Close(context.Background()); err != nil {
		return err
	}
	close(e.output)
	return nil
}

func (e *FibXor) Output() <-chan int {
	return e.output
}

func ExampleXor() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	fib := NewFibXor(ctx)
	go fib.Run([2]int{-1, 1}, 0)
	defer fib.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-fib.Output():
			if msg == 13 {
				go fib.Run([2]int{13, 21}, 0)
			}
			if msg == 89 {
				go fib.Run([2]int{89, 144}, 0)
			}
			fmt.Print(msg, " ")
		}
	}

	// Output: 0 1 1 2 3 5 8 13 34 55 89 233 377
}
