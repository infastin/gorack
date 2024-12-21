package lifecycle

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

type config struct {
	stopTimeout time.Duration
	logger      *slog.Logger
}

func defaultConfig() config {
	return config{
		stopTimeout: time.Minute,
		logger:      slog.New(discardHandler{}),
	}
}

type Option func(cfg *config)

func WithStopTimeout(stopTimeout time.Duration) Option {
	return func(cfg *config) {
		cfg.stopTimeout = stopTimeout
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(cfg *config) {
		cfg.logger = logger
	}
}

type Lifecycle struct {
	stopTimeout time.Duration
	logger      *slog.Logger
	hooks       []hook
}

func New(opts ...Option) *Lifecycle {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	return &Lifecycle{
		stopTimeout: cfg.stopTimeout,
		logger:      cfg.logger,
		hooks:       make([]hook, 0),
	}
}

func (l *Lifecycle) Append(h Hook) {
	hook := hook{}

	logger := l.logger
	if h.Name != "" {
		logger = logger.With(slog.String("name", h.Name))
	}

	if h.OnStart != nil {
		hook.onStart = func(ctx context.Context, cancel context.CancelCauseFunc) error {
			logger.Info("running start hook")
			if err := h.OnStart(ctx, cancel); err != nil {
				logger.Error("failed to run start hook", slog.String("error", err.Error()))
				return err
			}
			return nil
		}
	}

	if h.OnStop != nil {
		hook.onStop = func(ctx context.Context) error {
			logger.Info("running stop hook")
			if err := h.OnStop(ctx); err != nil {
				logger.Error("failed to run stop hook", slog.String("error", err.Error()))
				return err
			}
			return nil
		}
	}

	l.hooks = append(l.hooks, hook)
}

func (l *Lifecycle) Go(actor Actor) {
	hook := hook{}

	logger := l.logger
	if actor.Name != "" {
		logger = logger.With(slog.String("name", actor.Name))
	}

	if actor.Run != nil {
		hook.onStart = func(ctx context.Context, cancel context.CancelCauseFunc) error {
			go func() {
				logger.Info("running start hook")
				if err := actor.Run(ctx); err != nil {
					logger.Error("failed to run start hook", slog.String("error", err.Error()))
					cancel(err)
				}
			}()
			return nil
		}
	}

	if actor.Shutdown != nil {
		hook.onStop = func(ctx context.Context) error {
			logger.Info("running stop hook")
			if err := actor.Shutdown(ctx); err != nil {
				logger.Error("failed to run stop hook", slog.String("error", err.Error()))
				return err
			}
			return nil
		}
	}

	l.hooks = append(l.hooks, hook)
}

func (l *Lifecycle) Run(ctx context.Context) error {
	if len(l.hooks) == 0 {
		return nil
	}

	hookCtx, hookCancel := context.WithCancelCause(ctx)
	defer hookCancel(nil)

	k := 0
	for ; k < len(l.hooks); k++ {
		hook := l.hooks[k]
		if hook.onStart == nil {
			continue
		}
		if err := hook.onStart(hookCtx, hookCancel); err != nil {
			hookCancel(err)
			break
		}
	}

	<-hookCtx.Done()

	errs := make([]error, 0, len(l.hooks)+1)
	errs = append(errs, context.Cause(hookCtx))

	stopCtx, cancel := context.WithTimeout(ctx, l.stopTimeout)
	defer cancel()

	for i := k - 1; i >= 0; i-- {
		if l.hooks[i].onStop == nil {
			continue
		}
		if err := l.hooks[i].onStop(stopCtx); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
