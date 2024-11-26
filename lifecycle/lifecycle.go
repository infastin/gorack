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
	hooks       []Hook
}

func New(opts ...Option) *Lifecycle {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	return &Lifecycle{
		stopTimeout: cfg.stopTimeout,
		logger:      cfg.logger,
		hooks:       make([]Hook, 0),
	}
}

func (l *Lifecycle) Append(hook Hook) {
	l.hooks = append(l.hooks, hook)
}

func (l *Lifecycle) Go(g Go) {
	h := Hook{}

	logger := l.logger
	if g.Name != "" {
		logger = logger.With(slog.String("name", g.Name))
	}

	if g.Run != nil {
		h.OnStart = func(ctx context.Context, cancel context.CancelCauseFunc) error {
			go func() {
				logger.Info("running start hook")
				if err := g.Run(ctx); err != nil {
					logger.Error("failed to run start hook", slog.String("error", err.Error()))
					cancel(err)
				}
			}()
			return nil
		}
	}
	if g.Shutdown != nil {
		h.OnStop = func(ctx context.Context) error {
			logger.Info("running stop hook")
			if err := g.Shutdown(ctx); err != nil {
				logger.Error("failed to run stop hook", slog.String("error", err.Error()))
				return err
			}
			return nil
		}
	}

	l.hooks = append(l.hooks, h)
}

func (l *Lifecycle) Run(ctx context.Context) error {
	if len(l.hooks) == 0 {
		return nil
	}

	hookCtx, hookCancel := context.WithCancelCause(ctx)

	k := 0
	for ; k < len(l.hooks); k++ {
		hook := l.hooks[k]
		if hook.OnStart == nil {
			continue
		}
		if err := hook.OnStart(hookCtx, hookCancel); err != nil {
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
		if l.hooks[i].OnStop == nil {
			continue
		}
		if err := l.hooks[i].OnStop(stopCtx); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
