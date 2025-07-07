package shot

import (
	"context"
	"sync"
)

// Many allows to control the lifecycle of a resource
// so the resource can have only one instance running.
type Many struct {
	done      chan struct{}
	state     State
	ctx       context.Context
	cancel    context.CancelFunc
	mu        sync.RWMutex
	parentCtx context.Context
}

// NewMany creates Many with the given parent context.
func NewMany(parent context.Context) Many {
	ctx, cancel := context.WithCancel(parent)
	return Many{
		done:      make(chan struct{}),
		state:     StateCreated,
		ctx:       ctx,
		cancel:    cancel,
		mu:        sync.RWMutex{},
		parentCtx: parent,
	}
}

// Start puts the resource in the Running state
// and returns stop function that is used
// to signal that the resource has exited.
//
// If resource hasn't been closed, calling the returned stop function
// will put the resource in the Stopped state, which allows for the resource
// to be started again with Start method.
//
// The returned stop function must be called
// for the resource to properly exit.
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
		m.ctx, m.cancel = context.WithCancel(m.parentCtx)
		m.done = make(chan struct{})
	}

	m.state = StateRunning

	return m.onExit, nil
}

func (m *Many) onExit() {
	m.mu.Lock()
	if m.ctx.Err() != nil {
		m.state = StateClosed
	} else {
		m.state = StateStopped
		m.cancel()
	}
	close(m.done)
	m.mu.Unlock()
}

// Close transitions the resource into the Closed state,
// cancels the context returned from Context method
// and prevents the resource from being started at all.
//
// If the resource is in the Running state, waits for the resource to exit.
//
// Context passed to this method can be canceled to pass control back to the caller
// if the resource takes too much time to exit.
// Canceling the context passed to this method doesn't affect the resource in any way.
func (m *Many) Close(ctx context.Context) error {
	done := m.close()
	if done == nil {
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
	}

	return nil
}

func (m *Many) close() (done <-chan struct{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch m.state {
	case StateCreated:
		m.state = StateClosed
		m.cancel()
		close(m.done)
		return nil
	case StateStopped:
		m.state = StateClosed
		m.cancel()
		return nil
	case StateClosed:
		return nil
	}

	m.cancel()
	return m.done
}

// Context returns the context of the resource.
//
// Context is cancelled when the resource exits,
// is closed by Close method, or when the parent context is cancelled.
//
// Resource must use the context returned from this method
// to exit when the context has been cancelled.
func (m *Many) Context() context.Context {
	return m.ctx
}

// Done returns the channel that is closed
// when the resource exits or is closed.
func (m *Many) Done() <-chan struct{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.done
}

// State returns current state of the resource.
func (m *Many) State() State {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state
}
