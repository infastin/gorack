package shot

import (
	"context"
	"sync/atomic"
)

// E holds the error returned from the goroutine.
type E struct {
	err atomic.Value
}

// GoErr starts a goroutine and returns E,
// which stores the error returned from the goroutine.
// Does not provide any way to determite when goroutine exited.
func GoErr(f func() error) *E {
	e := &E{err: atomic.Value{}}
	go func() {
		if err := f(); err != nil {
			e.err.Store(err)
		}
	}()
	return e
}

// Err returns the error returned from the goroutine.
// If goroutine hasn't yet exited, returns nil.
// If goroutine exited without error, also returns nil.
func (e *E) Err() error {
	err, ok := e.err.Load().(error)
	if !ok {
		return nil
	}
	return err
}

// G allows to control state of a goroutine
// and get the error returned from the goroutine.
type G struct {
	s   One
	err atomic.Value
}

// Go starts a goroutine and creates One,
// which is passed to the goroutine to control its state.
//
// Returns G, which can be used to control the goroutine
// and get the error returned from the goroutine.
//
// NOTE: Goroutine must make use of (*One).Start for G to function properly.
func Go(ctx context.Context, f func(state *One) error) *G {
	g := &G{s: NewOne(ctx), err: atomic.Value{}}
	go func() {
		if err := f(&g.s); err != nil {
			g.err.Store(err)
		}
	}()
	return g
}

// GoCtx starts a goroutine and creates One, whose Context
// is passed to the goroutine to control its state.
//
// Returns G, which can be used to control the goroutine
// and get the error returned from the goroutine.
//
// NOTE: Compared to Go function, this one calls (*One).Start before f is called
// and calls stop function returned from (*One).Start when f returns.
func GoCtx(ctx context.Context, f func(ctx context.Context) error) *G {
	g := &G{s: NewOne(ctx), err: atomic.Value{}}
	go func() {
		stop, err := g.s.Start()
		if err != nil {
			g.err.Store(err)
			return
		}
		if err := f(g.s.Context()); err != nil {
			g.err.Store(err)
		}
		stop()
	}()
	return g
}

// Close closes One passed to the goroutine.
func (g *G) Close(ctx context.Context) error {
	return g.s.Close(ctx)
}

// Done returns channel that is closed when the goroutine has exited.
func (g *G) Done() <-chan struct{} {
	return g.s.Done()
}

// Err returns the error returned from the goroutine.
// If goroutine hasn't yet exited, returns nil.
// If goroutine exited without error, also returns nil.
func (g *G) Err() error {
	err, ok := g.err.Load().(error)
	if !ok {
		return nil
	}
	return err
}
