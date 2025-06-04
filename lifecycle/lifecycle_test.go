package lifecycle_test

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/infastin/gorack/lifecycle"
)

func FuzzLifecycle_startContext(f *testing.F) {
	for n := 3; n <= 6; n++ {
		f.Add(n)
	}
	f.Fuzz(func(t *testing.T, n int) {
		lc := lifecycle.New()

		ctxs := make([]context.Context, n)
		for i := range n {
			lc.Append(lifecycle.Hook{
				OnStart: func(ctx context.Context, ccf context.CancelCauseFunc) error {
					ctxs[i] = ctx
					return nil
				},
			})
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		if err := lc.Run(ctx); err == nil {
			t.Error("expected an error")
			return
		}

		for i := 1; i < n; i++ {
			if ctxs[i-1] == ctxs[i] {
				t.Errorf("expected %d'th context to not be equal to %d'th context", i-1, i)
			}
		}
	})
}

func TestLifecycle_cancel(t *testing.T) {
	lc := lifecycle.New()

	canExit := make(chan struct{})
	canceled := make(chan struct{})

	lc.Append(lifecycle.Hook{
		OnStart: func(ctx context.Context, ccf context.CancelCauseFunc) error {
			go func() {
				ccf(nil)
				close(canceled)
			}()
			return nil
		},
	})

	lc.Go(lifecycle.Actor{
		Run: func(ctx context.Context) error {
			<-canceled
			select {
			case <-ctx.Done():
				t.Error("expected context to not be closed")
			case <-time.After(50 * time.Millisecond):
			}
			close(canExit)
			return nil
		},
	})

	lc.Append(lifecycle.Hook{
		OnStop: func(ctx context.Context) error {
			<-canExit
			return nil
		},
	})

	if err := lc.Run(context.Background()); err == nil {
		t.Error("expected an error")
		return
	}
}

func FuzzLifecycle_stopOrder(f *testing.F) {
	for n := 3; n <= 6; n++ {
		f.Add(n)
	}
	f.Fuzz(func(t *testing.T, n int) {
		lc := lifecycle.New()

		expected := make([]int, n)
		got := make([]int, 0, n)

		for i := range n {
			expected[i] = n - (i + 1)
			lc.Append(lifecycle.Hook{
				OnStop: func(ctx context.Context) error {
					got = append(got, i)
					return nil
				},
			})
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		if err := lc.Run(ctx); err == nil {
			t.Error("expected an error")
			return
		}

		if !slices.Equal(expected, got) {
			t.Errorf("slices must be equal: expected=%v got=%v", expected, got)
		}
	})
}
