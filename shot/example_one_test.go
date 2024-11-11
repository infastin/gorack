package shot_test

import (
	"context"
	"fmt"
	"time"

	"github.com/infastin/gorack/shot"
)

type OneExample struct {
	output chan int
	state  shot.One
}

func NewOneExample(ctx context.Context) *OneExample {
	return &OneExample{
		output: make(chan int, 1),
		state:  shot.NewOne(ctx),
	}
}

func (e *OneExample) Run() error {
	stop, err := e.state.Start()
	if err != nil {
		return err
	}
	defer stop()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	values := [2]int{-1, 1}
	idx := 0

	for {
		select {
		case <-e.state.Context().Done():
			return e.state.Context().Err()
		case <-ticker.C:
			sum := values[0] + values[1]
			values[idx] = sum
			idx = (idx + 1) % 2
			e.output <- sum
			if sum == 144 {
				ticker.Stop()
			}
		}
	}
}

func (e *OneExample) Close() error {
	if err := e.state.Close(context.Background()); err != nil {
		return err
	}
	return nil
}

func (e *OneExample) Output() <-chan int {
	return e.output
}

func ExampleOne() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	one := NewOneExample(ctx)
	defer one.Close()

	go one.Run()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-one.Output():
			fmt.Print(msg, " ")
		}
	}

	// Output: 0 1 1 2 3 5 8 13 21 34 55 89 144
}
