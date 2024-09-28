package shot

import (
	"context"
	"sync"
)

type Many struct {
	mu     sync.RWMutex
	done   chan struct{}
	state  State
	ctx    context.Context
	cancel context.CancelFunc
}

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
		m.mu.Unlock()
		close(m.done)
	}, nil
}

func (m *Many) Close(ctx context.Context) error {
	m.mu.Lock()
	if m.state == StateCreated {
		m.state = StateClosed
		close(m.done)
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

func (m *Many) Context() context.Context {
	return m.ctx
}

func (m *Many) Done() <-chan struct{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.done
}

func (m *Many) State() State {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state
}
