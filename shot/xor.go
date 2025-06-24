package shot

import (
	"context"
	"sync"
)

// Xor allows to control the lifecycle of a resource
// so the resource can have only one instance running,
// while also allowing for the resource to be replaced.
type Xor struct {
	done      chan struct{}
	state     State
	ctx       context.Context
	cancel    context.CancelFunc
	quitting  bool
	mu        sync.RWMutex
	parentCtx context.Context
}

// NewXor creates Xor with the given parent context.
func NewXor(parent context.Context) Xor {
	ctx, cancel := context.WithCancel(parent)
	return Xor{
		done:      make(chan struct{}),
		state:     StateCreated,
		ctx:       ctx,
		cancel:    cancel,
		quitting:  false,
		mu:        sync.RWMutex{},
		parentCtx: parent,
	}
}

// Start puts the resource in the Running state and returns stop function
// that is used to signal that resource has exited.
//
// If the resource is already running, cancels the context returned from Context method
// and waits for the resource to exit, after which the resource is started again.
//
// If the resource hasn't been closed, calling the returned stop function
// will put the resource in the Stopped state, which allows for the resource
// to be started again with Start method.
//
// The returned stop function must be called
// for the resource to properly exit.
//
// Returns an error if the resource has been closed.
//
// Context passed to this method can be canceled to pass control back to the caller
// if the resource takes too much time to restart.
// In such case, the resource will eventually close, but will not start again.
func (x *Xor) Start(ctx context.Context) (stop func(), err error) {
	x.mu.Lock()
retry:
	switch x.state {
	case StateRunning:
		x.cancel()
		x.mu.Unlock()
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-x.done:
			x.mu.Lock()
			goto retry
		}
	case StateClosed:
		x.mu.Unlock()
		return nil, ErrClosed
	case StateStopped:
		x.ctx, x.cancel = context.WithCancel(x.parentCtx)
		x.done = make(chan struct{})
	}
	x.state = StateRunning
	x.mu.Unlock()
	return x.onExit, nil
}

func (x *Xor) onExit() {
	x.mu.Lock()
	if x.quitting || x.parentCtx.Err() != nil {
		x.state = StateClosed
	} else {
		x.state = StateStopped
	}
	if x.ctx.Err() == nil {
		x.cancel()
	}
	close(x.done)
	x.mu.Unlock()
}

// Stop transitions the resource into Stopped state
// and cancels the context returned from Context method.
//
// If the resource is in the Closed state, returns ErrClosed.
//
// If the resource is in the Running state, waits for the resource to exit.
//
// Context passed to this method can be canceled to pass control back to the caller
// if the resource takes too much time to exit.
// Canceling the context passed to this method doesn't affect the resource in any way.
func (x *Xor) Stop(ctx context.Context) error {
	if nowait, err := x.stop(); nowait {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-x.done:
	}

	return nil
}

func (x *Xor) stop() (bool, error) {
	x.mu.Lock()
	defer x.mu.Unlock()

	switch x.state {
	case StateCreated:
		x.state = StateStopped
		x.cancel()
		close(x.done)
		return true, nil
	case StateStopped:
		return true, nil
	case StateClosed:
		return true, ErrClosed
	}

	x.cancel()
	return false, nil
}

// Close transitions the resource into Closed state,
// cancels the context returned from Context method,
// and prevents the resource from being started at all.
//
// If the resource is in the Running state, waits for resource to exit.
//
// Context passed to this method can be canceled to pass control back to the caller
// if the resource takes too much time to exit.
// Canceling the context passed to this method doesn't affect the resource in any way.
func (x *Xor) Close(ctx context.Context) error {
	if x.close() {
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-x.done:
	}

	return nil
}

func (x *Xor) close() bool {
	x.mu.Lock()
	defer x.mu.Unlock()

	switch x.state {
	case StateCreated:
		x.state = StateClosed
		x.cancel()
		close(x.done)
		return true
	case StateStopped:
		x.state = StateClosed
		x.cancel()
		return true
	case StateClosed:
		return true
	}

	x.quitting = true
	x.cancel()

	return false
}

// Context returns the context of the resource.
//
// Context is cancelled when the resource is closed by Close method,
// is stopped by Stop method, or when the parent context is cancelled.
//
// Resource must use the context returned from this method
// to exit when the context has been cancelled.
func (x *Xor) Context() context.Context {
	x.mu.RLock()
	defer x.mu.RUnlock()
	return x.ctx
}

// Done returns the channel that is closed
// when the resource exits or is closed.
func (x *Xor) Done() <-chan struct{} {
	x.mu.RLock()
	defer x.mu.RUnlock()
	return x.done
}

// State returns current state of the resource.
func (x *Xor) State() State {
	x.mu.RLock()
	defer x.mu.RUnlock()
	return x.state
}
