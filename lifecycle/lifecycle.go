package lifecycle

import (
	"context"
	"errors"
	"time"
)

type config struct {
	stopTimeout time.Duration
}

func defaultConfig() config {
	return config{
		stopTimeout: time.Minute,
	}
}

type Option func(cfg *config)

func WithStopTimeout(stopTimeout time.Duration) Option {
	return func(cfg *config) {
		cfg.stopTimeout = stopTimeout
	}
}

type Lifecycle struct {
	stopTimeout time.Duration
	hooks       []Hook
}

func New(opts ...Option) *Lifecycle {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	return &Lifecycle{
		stopTimeout: cfg.stopTimeout,
		hooks:       make([]Hook, 0),
	}
}

func (l *Lifecycle) Append(hook Hook) {
	l.hooks = append(l.hooks, hook)
}

func (l *Lifecycle) Run(ctx context.Context) error {
	if len(l.hooks) == 0 {
		return nil
	}

	hookCtx, hookCancel := context.WithCancelCause(ctx)

	for _, hook := range l.hooks {
		if hook.OnStart == nil {
			continue
		}

		if err := hook.OnStart(hookCtx, hookCancel); err != nil {
			return err
		}
	}

	<-hookCtx.Done()

	errs := make([]error, 0, len(l.hooks)+1)
	errs = append(errs, context.Cause(hookCtx))

	stopCtx, cancel := context.WithTimeout(ctx, l.stopTimeout)
	defer cancel()

	for i := len(l.hooks) - 1; i >= 0; i-- {
		if l.hooks[i].OnStop == nil {
			continue
		}

		if err := l.hooks[i].OnStop(stopCtx); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
