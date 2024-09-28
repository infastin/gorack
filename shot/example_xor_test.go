package shot_test

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/infastin/go-rack/shot"
)

type XorExample struct {
	output chan string
	state  shot.Xor
}

func NewXorExample(ctx context.Context) *XorExample {
	return &XorExample{
		output: make(chan string, 1),
		state:  shot.NewXor(ctx),
	}
}

func (e *XorExample) Run(counter int) error {
	stop, err := e.state.Start(context.Background())
	if err != nil {
		return err
	}
	defer stop()

	for {
		select {
		case <-e.state.Context().Done():
			return e.state.Context().Err()
		case <-time.After(time.Second):
			switch counter++; {
			case counter%3 == 0 && counter%5 == 0:
				e.output <- "FizzBuzz"
			case counter%3 == 0:
				e.output <- "Fizz"
			case counter%5 == 0:
				e.output <- "Fizz"
			default:
				e.output <- strconv.Itoa(counter)
			}
			return nil
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

func (e *XorExample) Output() <-chan string {
	return e.output
}

func ExampleXor() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	xor := NewXorExample(ctx)
	defer xor.Close()

	counters := make(chan int, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(1500 * time.Millisecond):
				counters <- rand.Int()
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case counter := <-counters:
			go xor.Run(counter)
		case msg := <-xor.Output():
			fmt.Println(msg)
		}
	}
}
