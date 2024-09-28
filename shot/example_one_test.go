package shot_test

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/infastin/go-rack/shot"
)

type OneExample struct {
	ticker *time.Ticker
	output chan string
	state  shot.One
}

func NewOneExample(ctx context.Context) *OneExample {
	return &OneExample{
		ticker: time.NewTicker(time.Second),
		output: make(chan string, 1),
		state:  shot.NewOne(ctx),
	}
}

func (e *OneExample) Run() error {
	stop, err := e.state.Start()
	if err != nil {
		return err
	}
	defer stop()

	counter := 0

	for {
		select {
		case <-e.state.Context().Done():
			return e.state.Context().Err()
		case <-e.ticker.C:
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
		}
	}
}

func (e *OneExample) Close() error {
	if err := e.state.Close(context.Background()); err != nil {
		return err
	}
	close(e.output)
	return nil
}

func (e *OneExample) Output() <-chan string {
	return e.output
}

func ExampleOne() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	one := NewOneExample(ctx)
	defer one.Close()

	go one.Run()

	for msg := range one.Output() {
		fmt.Println(msg)
	}
}
