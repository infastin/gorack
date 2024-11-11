package shot_test

import (
	"context"
	"fmt"
	"time"

	"github.com/infastin/gorack/shot"
)

type XorExample struct {
	output chan int
	state  shot.Xor
}

func NewXorExample(ctx context.Context) *XorExample {
	return &XorExample{
		output: make(chan int, 1),
		state:  shot.NewXor(ctx),
	}
}

func (e *XorExample) Run(values [2]int, idx int) error {
	stop, err := e.state.Start(context.Background())
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

func (e *XorExample) Close() error {
	if err := e.state.Close(context.Background()); err != nil {
		return err
	}
	close(e.output)
	return nil
}

func (e *XorExample) Output() <-chan int {
	return e.output
}

func ExampleXor() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	xor := NewXorExample(ctx)
	defer xor.Close()

	go xor.Run([2]int{-1, 1}, 0)

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-xor.Output():
			if msg == 13 {
				go xor.Run([2]int{13, 21}, 0)
			}
			if msg == 89 {
				go xor.Run([2]int{89, 144}, 0)
			}
			fmt.Print(msg, " ")
		}
	}

	// Output: 0 1 1 2 3 5 8 13 34 55 89 233 377
}
