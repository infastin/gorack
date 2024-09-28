package shot_test

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/infastin/go-rack/shot"
)

type ManyExample struct {
	counter int
	output  chan string
	state   shot.Many
}

func NewManyExample(ctx context.Context) *ManyExample {
	return &ManyExample{
		counter: 0,
		output:  make(chan string, 1),
		state:   shot.NewMany(ctx),
	}
}

func (e *ManyExample) Run() error {
	stop, err := e.state.Start()
	if err != nil {
		return err
	}
	defer stop()

	for {
		select {
		case <-e.state.Context().Done():
			return e.state.Context().Err()
		case <-time.After(time.Second):
			switch e.counter++; {
			case e.counter%3 == 0 && e.counter%5 == 0:
				e.output <- "FizzBuzz"
			case e.counter%3 == 0:
				e.output <- "Fizz"
			case e.counter%5 == 0:
				e.output <- "Fizz"
			default:
				e.output <- strconv.Itoa(e.counter)
			}
			return nil
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

func (e *ManyExample) Output() <-chan string {
	return e.output
}

func ExampleMany() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	many := NewManyExample(ctx)
	defer many.Close()

	go many.Run()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-many.Output():
			fmt.Println(msg)
			go many.Run()
		}
	}
}
