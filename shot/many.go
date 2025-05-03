package shot

import (
	"context"
	"sync"
)

// Allows to control the lifecycle of a resource
// so the resource can have only one instance running.
type Many struct {
	mu     sync.RWMutex
	done   chan struct{}
	state  State
	ctx    context.Context
	cancel context.CancelFunc
}

// Creates Many with the given parent Context.
func NewMany(ctx context.Context) Many {
	ctx, cancel := context.WithCancel(ctx)
	return Many{
		mu:     sync.RWMutex{},
		done:   make(chan struct{}),
		state:  StateCreated,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Puts the resource in the Running state
// and returns stop function that is used
// to signal that resource has exited.
//
// If resource hasn't been closed, calling the returned stop function
// will put the resource in the Stopped state, which allows for the resource
// to be started again with Start method.
//
// The returned stop function must be called for the
// channel returned from Done method to become closed.
//
// Returns an error if the resource has been closed or currently running.
func (m *Many) Start() (stop func(), err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch m.state {
	case StateRunning:
		return nil, ErrRunning
	case StateClosed:
		return nil, ErrClosed
	case StateStopped:
		m.done = make(chan struct{})
	}

	m.state = StateRunning

	return func() {
		m.mu.Lock()
		if m.ctx.Err() != nil {
			m.state = StateClosed
		} else {
			m.state = StateStopped
		}
		close(m.done)
		m.mu.Unlock()
	}, nil
}

// Transitions the resource into the Closed state,
// cancels Context returned from Context method
// and prevents resource from being started at all.
//
// If resource is in the Running state, waits for resource to call stop function
// returned from Start method (i.e. waits for resource to exit).
//
// Context passed to this function can be canceled to pass control back to the caller
// if resource takes too much time to exit. Canceling Context passed to this function
// doesn't affect the resource in any way.
func (m *Many) Close(ctx context.Context) error {
	m.mu.Lock()
	switch m.state {
	case StateCreated:
		m.state = StateClosed
		close(m.done)
	case StateStopped:
		m.state = StateClosed
		// done is already closed here
	}
	m.mu.Unlock()

	m.cancel()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-m.done:
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
func (m *Many) Context() context.Context {
	return m.ctx
}

// Returns the channel that will be closed
// when resource is closed or resource has exited.
func (m *Many) Done() <-chan struct{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.done
}

// Returns current state of the resource.
func (m *Many) State() State {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state
}
