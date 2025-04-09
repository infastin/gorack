package shot

import (
	"context"
	"sync/atomic"
)

// Holds the error returned from the goroutine.
type E struct {
	err atomic.Value
}

// Starts a goroutine and returns E, which stores the error returned from the goroutine.
// Does not provide any way to determite when goroutine exited.
func GoErr(g func() error) *E {
	e := &E{err: atomic.Value{}}
	go func() {
		if err := g(); err != nil {
			e.err.Store(err)
		}
	}()
	return e
}

// Returns the error returned from the goroutine.
// If goroutine hasn't yet exited, returns nil.
// If goroutine exited without error, also returns nil.
func (e *E) Err() error {
	err, ok := e.err.Load().(error)
	if !ok {
		return nil
	}
	return err
}

// Allows to control state of a goroutine
// and get the error returned from the goroutine.
type G struct {
	s   One
	err atomic.Value
}

// Starts a goroutine and creates One, which is passed to the goroutine
// to control its state.
//
// Returns G, which can be used to control the goroutine
// and get the error returned from the goroutine.
//
// NOTE: Goroutine must make use of (*One).Start for G to function properly.
func Go(ctx context.Context, g func(state *One) error) *G {
	state := &G{s: NewOne(ctx), err: atomic.Value{}}
	go func() {
		if err := g(&state.s); err != nil {
			state.err.Store(err)
		}
	}()
	return state
}

// Closes One passed to the goroutine.
func (g *G) Close(ctx context.Context) error {
	return g.s.Close(ctx)
}

// Returns channel that is closed when goroutine is closed.
func (g *G) Done() <-chan struct{} {
	return g.s.Done()
}

// Returns the error returned from the goroutine.
// If goroutine hasn't yet exited, returns nil.
// If goroutine exited without error, also returns nil.
func (g *G) Err() error {
	err, ok := g.err.Load().(error)
	if !ok {
		return nil
	}
	return err
}
