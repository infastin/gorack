package shot

import (
	"context"
	"sync/atomic"
)

type One struct {
	done   chan struct{}
	state  atomic.Int32
	ctx    context.Context
	cancel context.CancelFunc
}

func NewOne(ctx context.Context) One {
	ctx, cancel := context.WithCancel(ctx)
	return One{
		done:   make(chan struct{}),
		state:  atomic.Int32{},
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *One) Start() (stop func(), err error) {
	var state State

	switch state = State(s.state.Load()); state {
	case StateRunning:
		return nil, ErrRunning
	case StateClosed:
		return nil, ErrClosed
	}

	if !s.state.CompareAndSwap(int32(state), int32(StateRunning)) {
		switch state = State(s.state.Load()); state {
		case StateRunning:
			return nil, ErrRunning
		case StateClosed:
			return nil, ErrClosed
		}
	}

	return func() {
		s.state.Store(int32(StateClosed))
		close(s.done)
	}, nil
}

func (s *One) Close(ctx context.Context) error {
	state := State(s.state.Load())
	if state == StateCreated && s.state.CompareAndSwap(int32(state), int32(StateClosed)) {
		close(s.done)
	}

	s.cancel()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.done:
	}

	return nil
}

func (s *One) Context() context.Context {
	return s.ctx
}

func (s *One) Done() <-chan struct{} {
	return s.done
}

func (s *One) State() State {
	return State(s.state.Load())
}
