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
			err := h.OnStart(ctx, cancel)
			if err != nil {
				logger.Error("start hook ran with failure", slog.String("error", err.Error()))
			} else {
				logger.Info("start hook ran successfully")
			}
			return err
		}
	}

	if h.OnStop != nil {
		hook.onStop = func(ctx context.Context) error {
			logger.Info("running stop hook")
			err := h.OnStop(ctx)
			if err != nil {
				logger.Error("stop hook ran with failure", slog.String("error", err.Error()))
			} else {
				logger.Info("stop hook ran successfully")
			}
			return err
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
				err := actor.Run(ctx)
				if err != nil {
					logger.Error("stop hook ran with failure", slog.String("error", err.Error()))
				} else {
					logger.Info("stop hook ran successfully")
				}
				cancel(err)
			}()
			return nil
		}
	}

	if actor.Shutdown != nil {
		hook.onStop = func(ctx context.Context) error {
			logger.Info("running stop hook")
			err := actor.Shutdown(ctx)
			if err != nil {
				logger.Error("stop hook ran with failure", slog.String("error", err.Error()))
			} else {
				logger.Info("stop hook ran successfully")
			}
			return err
		}
	}

	l.hooks = append(l.hooks, hook)
}

func (l *Lifecycle) GoFunc(fn func(ctx context.Context) error) {
	l.Go(Actor{Run: fn})
}

func (l *Lifecycle) Run(ctx context.Context) error {
	if len(l.hooks) == 0 {
		return nil
	}

	noCancelCtx := context.WithoutCancel(ctx)

	hooksCtx, hooksCancel := context.WithCancelCause(ctx)
	defer hooksCancel(nil)

	k := 0
	for ; k < len(l.hooks); k++ {
		hook := &l.hooks[k]
		if hook.onStart == nil {
			continue
		}
		hook.startCtx, hook.cancelStartCtx = context.WithCancel(noCancelCtx)
		if err := hook.onStart(hook.startCtx, hooksCancel); err != nil {
			hooksCancel(err)
			break
		}
	}

	<-hooksCtx.Done()

	errs := make([]error, 0, len(l.hooks)+1)
	errs = append(errs, context.Cause(hooksCtx))

	stopCtx, cancel := context.WithTimeout(noCancelCtx, l.stopTimeout)
	defer cancel()

	for i := k - 1; i >= 0; i-- {
		hook := &l.hooks[i]
		if hook.onStart != nil {
			hook.cancelStartCtx()
		}
		if hook.onStop == nil {
			continue
		}
		if err := hook.onStop(stopCtx); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
