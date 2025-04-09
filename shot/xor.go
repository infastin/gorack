package shot

import (
	"context"
	"sync"
)

// Allows to control the lifecycle of a resource
// so the resource can have only one instance running,
// while also allowing for the resource to be replaced.
type Xor struct {
	mu         sync.RWMutex
	done       chan struct{}
	state      State
	ctx        context.Context
	cancel     context.CancelFunc
	stopCtx    context.Context
	stopCancel context.CancelFunc
}

// Creates Xor with the given parent Context.
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

// Puts the resource in the Running state
// and returns stop function that is used
// to signal that resource has exited.
//
// If the resource is already running, cancels Context returned from Context method
// and waits for resource to call stop function returned from Start method (i.e. waits for resource to exit),
// after which resource is started again.
//
// If resource hasn't been closed, calling the returned stop function
// will put the resource in the Stopped state, which allows for the resource
// to be started again with Start method.
//
// The returned stop function must be called for the
// channel returned from Done method to become closed.
//
// Returns an error if the resource has been closed.
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

// Transitions the resource into Stopped state,
// cancels Context returned from Context method.
//
// If resource is in the Running state waits for resource to call stop function returned
// from Start method (i.e. waits for resource to exit).
//
// Context passed to this function can be canceled to pass control back to the caller
// if resource takes too much time to exit. Canceling Context passed to this function
// doesn't affect the resource in any way.
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

// Transitions the resource into Closed state,
// cancels Context returned from Context method,
// and prevents resource from being started at all.
//
// If resource is in the Running state, waits for resource to call stop function
// returned from Start method (i.e. waits for resource to exit).
//
// Context passed to this function can be canceled to pass control back to the caller
// if resource takes too much time to exit. Canceling Context passed to this function
// doesn't affect the resource in any way.
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

// Returns Context of the resource.
//
// Context will be closed when resource is closed or stoped,
// or resource has exited.
//
// Resource must use Context returned from this method to detect
// when closing/stopping the resource was requested
// with Close/Stop methods respectively.
func (x *Xor) Context() context.Context {
	x.mu.RLock()
	defer x.mu.RUnlock()
	return x.stopCtx
}

// Returns the channel that will be closed
// when resource is closed or stopped, or resource has exited.
func (x *Xor) Done() <-chan struct{} {
	x.mu.RLock()
	defer x.mu.RUnlock()
	return x.done
}

// Returns current state of the resource.
func (x *Xor) State() State {
	x.mu.RLock()
	defer x.mu.RUnlock()
	return x.state
}
