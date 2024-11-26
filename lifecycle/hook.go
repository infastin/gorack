package lifecycle

import (
	"context"
	"os"
	"os/signal"
)

type Hook struct {
	OnStart func(context.Context, context.CancelCauseFunc) error
	OnStop  func(context.Context) error
}

type Go struct {
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
