package shot

import (
	"context"
	"sync"
)

type Xor struct {
	mu         sync.RWMutex
	done       chan struct{}
	state      State
	ctx        context.Context
	cancel     context.CancelFunc
	stopCtx    context.Context
	stopCancel context.CancelFunc
}

func NewXor(ctx context.Context) Xor {
	ctx, cancel := context.WithCancel(ctx)
	return Xor{
		mu:         sync.RWMutex{},
		done:       make(chan struct{}),
		state:      StateCreated,
		ctx:        ctx,
		cancel:     cancel,
		stopCtx:    ctx,
		stopCancel: cancel,
	}
}

func (x *Xor) Start(ctx context.Context) (stop func(), err error) {
	x.mu.Lock()
	defer x.mu.Unlock()

retry:
	switch x.state {
	case StateRunning:
		x.mu.Unlock()
		x.stopCancel()
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-x.done:
			x.mu.Lock()
			goto retry
		}
	case StateClosed:
		return nil, ErrClosed
	case StateStopped:
		x.done = make(chan struct{})
	}

	x.stopCtx, x.stopCancel = context.WithCancel(x.ctx)
	x.state = StateRunning

	return func() {
		x.mu.Lock()
		if x.ctx.Err() != nil {
			x.state = StateClosed
		} else {
			x.state = StateStopped
		}
		x.mu.Unlock()
		close(x.done)
	}, nil
}

func (x *Xor) Stop(ctx context.Context) error {
	x.mu.Lock()
	if x.state == StateCreated {
		x.state = StateStopped
		close(x.done)
	}
	x.mu.Unlock()

	x.stopCancel()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-x.done:
	}

	return nil
}

func (x *Xor) Close(ctx context.Context) error {
	x.mu.Lock()
	if x.state == StateCreated {
		x.state = StateClosed
		close(x.done)
	}
	x.mu.Unlock()

	x.cancel()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-x.done:
	}

	return nil
}

func (x *Xor) Context() context.Context {
	x.mu.RLock()
	defer x.mu.RUnlock()
	return x.stopCtx
}

func (x *Xor) Done() <-chan struct{} {
	x.mu.RLock()
	defer x.mu.RUnlock()
	return x.done
}

func (x *Xor) State() State {
	x.mu.RLock()
	defer x.mu.RUnlock()
	return x.state
}
