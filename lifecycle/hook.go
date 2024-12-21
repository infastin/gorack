package lifecycle

import (
	"context"
	"os"
	"os/signal"
)

type hook struct {
	onStart func(context.Context, context.CancelCauseFunc) error
	onStop  func(context.Context) error
}

type Hook struct {
	Name    string
	OnStart func(context.Context, context.CancelCauseFunc) error
	OnStop  func(context.Context) error
}

type Actor struct {
	Name     string
	Run      func(context.Context) error
	Shutdown func(context.Context) error
}

type SignalError struct {
	Signal os.Signal
}

func (e *SignalError) Error() string {
	return "received signal: " + e.Signal.String()
}

func Signal(signals ...os.Signal) Hook {
	return Hook{
		Name: "signal handler",
		OnStart: func(ctx context.Context, cancel context.CancelCauseFunc) error {
			go func() {
				sigCh := make(chan os.Signal, 1)
				signal.Notify(sigCh, signals...)

				select {
				case <-ctx.Done():
					return
				case sig := <-sigCh:
					cancel(&SignalError{Signal: sig})
				}
			}()
			return nil
		},
	}
}
