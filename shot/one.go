package shot

import (
	"context"
	"sync/atomic"
)

// Allows to control the lifecycle of a resource
// so the resource can only be run once
// and can have only one instance running.
type One struct {
	done   chan struct{}
	state  atomic.Int32
	ctx    context.Context
	cancel context.CancelFunc
}

// Creates One with the given parent Context.
func NewOne(ctx context.Context) One {
	ctx, cancel := context.WithCancel(ctx)
	return One{
		done:   make(chan struct{}),
		state:  atomic.Int32{},
		ctx:    ctx,
		cancel: cancel,
	}
}

// Puts the resource in the Running state
// and returns stop function that is used
// to signal that resource has exited.
//
// Calling the returned stop function
// will put the resource in the Closed state.
//
// The returned stop function must be called for the
// channel returned from Done method to become closed.
//
// Returns an error if the resource has been closed or currently running.
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

// Transitions the resource into the Closed state,
// and cancels Context returned from Context method.
//
// If resource was in the Created state,
// prevents resource from being started at all.
//
// If resource is in the Running state, waits for resource to call stop function
// returned from Start method (i.e. waits for resource to exit).
//
// Context passed to this function can be canceled to pass control back to the caller
// if resource takes too much time to exit. Canceling Context passed to this function
// doesn't affect the resource in any way.
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

// Returns Context of the resource.
//
// Context will be closed when resource is closed
// or resource has exited.
//
// Resource must use Context returned from this method to detect
// when closing the resource was requested with Close method.
func (s *One) Context() context.Context {
	return s.ctx
}

// Returns the channel that will be closed
// when resource is closed or resource has exited.
func (s *One) Done() <-chan struct{} {
	return s.done
}

// Returns current state of the resource.
func (s *One) State() State {
	return State(s.state.Load())
}
