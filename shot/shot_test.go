package shot_test

import (
	"sync"
	"testing"

	"github.com/infastin/gorack/shot"
)

type resource interface {
	State() shot.State
	Done() <-chan struct{}
}

func waitGo(num int, f func(wg *sync.WaitGroup, i int)) {
	var wg sync.WaitGroup
	wg.Add(num)
	for i := range num {
		go f(&wg, i)
	}
	wg.Wait()
}

func shouldBe(t *testing.T, r resource, expected shot.State) bool {
	t.Helper()

	if state := r.State(); state != expected {
		t.Errorf("State(): expected=%s got=%s", expected, state)
		return false
	}

	select {
	case <-r.Done():
		if expected == shot.StateCreated || expected == shot.StateRunning {
			t.Error("Done(): should not be closed")
			return false
		} else {
			t.Log("Done(): is closed")
		}
	default:
		if expected == shot.StateCreated || expected == shot.StateRunning {
			t.Log("Done(): is not closed")
		} else {
			t.Error("Done(): should be closed")
			return false

		}
	}

	return true
}
