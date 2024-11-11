package shot_test

import (
	"context"
	"fmt"
	"time"

	"github.com/infastin/gorack/shot"
)

type ManyExample struct {
	values [2]int
	idx    int
	output chan int
	state  shot.Many
}

func NewManyExample(ctx context.Context) *ManyExample {
	return &ManyExample{
		values: [2]int{-1, 1},
		idx:    0,
		output: make(chan int, 1),
		state:  shot.NewMany(ctx),
	}
}

func (e *ManyExample) Run() error {
	stop, err := e.state.Start()
	if err != nil {
		return err
	}
	defer stop()

	ticker := time.NewTicker(time.Second)
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

func (e *ManyExample) Close() error {
	if err := e.state.Close(context.Background()); err != nil {
		return err
	}
	close(e.output)
	return nil
}

func (e *ManyExample) Output() <-chan int {
	return e.output
}

func (e *ManyExample) Done() <-chan struct{} {
	return e.state.Done()
}

func ExampleMany() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	many := NewManyExample(ctx)
	defer many.Close()

	go many.Run()

	for {
		select {
		case <-ctx.Done():
			return
		case <-many.Done():
			go many.Run()
		case msg := <-many.Output():
			fmt.Print(msg, " ")
		}
	}

	// Output: 0 1 1 2 3 5 8 13 21 34 55 89 144
}
