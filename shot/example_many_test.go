package shot_test

import (
	"context"
	"fmt"
	"time"

	"github.com/infastin/gorack/shot"
)

type FibMany struct {
	values [2]int
	idx    int
	output chan int
	state  shot.Many
}

func NewFibMany(ctx context.Context) *FibMany {
	return &FibMany{
		values: [2]int{-1, 1},
		idx:    0,
		output: make(chan int, 1),
		state:  shot.NewMany(ctx),
	}
}

func (e *FibMany) Run() error {
	stop, err := e.state.Start()
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
			sum := e.values[0] + e.values[1]
			e.values[e.idx] = sum
			e.idx = (e.idx + 1) % 2
			e.output <- sum
			if sum == 13 || sum == 89 {
				return nil
			}
			if sum == 144 {
				ticker.Stop()
			}
		}
	}
}

func (e *FibMany) Close() error {
	if err := e.state.Close(context.Background()); err != nil {
		return err
	}
	close(e.output)
	return nil
}

func (e *FibMany) Output() <-chan int {
	return e.output
}

func (e *FibMany) Done() <-chan struct{} {
	return e.state.Done()
}

func ExampleMany() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	fib := NewFibMany(ctx)
	go fib.Run()
	defer fib.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case <-fib.Done():
			go fib.Run()
		case msg := <-fib.Output():
			fmt.Print(msg, " ")
		}
	}

	// Output: 0 1 1 2 3 5 8 13 21 34 55 89 144
}
