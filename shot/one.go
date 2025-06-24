package shot

import (
	"context"
	"sync/atomic"
)

// One allows to control the lifecycle of a resource
// so the resource can only be run once
// and can have only one instance running.
type One struct {
	done   chan struct{}
	state  atomic.Int32
	ctx    context.Context
	cancel context.CancelFunc
}

// NewOne creates One with the given parent context.
func NewOne(parent context.Context) One {
	ctx, cancel := context.WithCancel(parent)
	return One{
		done:   make(chan struct{}),
		state:  atomic.Int32{},
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start puts the resource in the Running state
// and returns stop function that is used
// to signal that the resource has exited.
//
// Calling the returned stop function
// will put the resource in the Closed state.
// The returned stop function must be called
// for the resource to properly exit.
//
// Returns an error if the resource has been closed
// or currently running.
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

	return s.onExit, nil
}

func (s *One) onExit() {
	s.state.Store(int32(StateClosed))
	s.cancel()
	close(s.done)
}

// Close transitions the resource into the Closed state,
// and cancels the context returned from Context method.
//
// If resource was in the Created state,
// prevents resource from being started at all.
//
// If resource is in the Running state,
// waits for resource to exit.
//
// Context passed to this method can be canceled to pass control back to the caller
// if resource takes too much time to exit.
// Canceling the context passed to this method doesn't affect the resource in any way.
func (s *One) Close(ctx context.Context) error {
	switch state := State(s.state.Load()); state {
	case StateCreated:
		if s.state.CompareAndSwap(int32(state), int32(StateClosed)) {
			s.cancel()
			close(s.done)
			return nil
		}
	case StateClosed:
		return nil
	}

	s.cancel()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.done:
	}

	return nil
}

// Context returns the context of the resource.
//
// Context is cancelled when the resource exits,
// is closed by Close method, or when the parent context is cancelled.
//
// Resource must use the context returned from this method
// to exit when the context has been cancelled.
func (s *One) Context() context.Context {
	return s.ctx
}

// Done returns the channel that is closed
// when the resource exits or is closed.
func (s *One) Done() <-chan struct{} {
	return s.done
}

// State returns current state of the resource.
func (s *One) State() State {
	return State(s.state.Load())
}
