package shot_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/infastin/gorack/shot"
)

func TestXor_Start_Close(t *testing.T) {
	const numGoroutines = 5
	state := shot.NewXor(context.Background())

	var count atomic.Int64
	count.Add(numGoroutines)

	waitGo(numGoroutines, func(wg *sync.WaitGroup, i int) {
		t.Logf("starting %d'th goroutine", i)
		wg.Done()

		stop, err := state.Start(context.Background())
		if err != nil {
			t.Errorf("%d'th goroutine: Start(): unexpected error: %s", i, err.Error())
			return
		}
		defer stop()

		t.Logf("%d'th goroutine started", i)

		<-state.Context().Done()
		t.Logf("stopping %d'th goroutine", i)
		count.Add(-1)
	})

	var fails int
	for {
		cnt := count.Load()
		if cnt == 1 {
			break
		}
		fails++
		if fails == 100 {
			t.Errorf("only 1 goroutine should be running, got %d", cnt)
			break
		}
		time.Sleep(time.Millisecond)
	}

	if !shouldBe(t, &state, shot.StateRunning) {
		return
	}

	if err := state.Close(context.Background()); err != nil {
		t.Errorf("Close(): unexpected error: %s", err.Error())
		return
	}

	if !shouldBe(t, &state, shot.StateClosed) {
		return
	}

	waitGo(1, func(wg *sync.WaitGroup, _ int) {
		_, err := state.Start(context.Background())
		if err == nil {
			t.Error("Start(): expected an error")
		} else {
			t.Logf("Start(): got expected error: %s", err.Error())
		}
		wg.Done()
	})
}

func TestXor_Start_stop(t *testing.T) {
	const numGoroutines = 5
	state := shot.NewXor(context.Background())

	for i := range numGoroutines {
		waitGo(1, func(wg *sync.WaitGroup, _ int) {
			defer wg.Done()
			stop, err := state.Start(context.Background())
			if err != nil {
				t.Errorf("%d'th goroutine: Start(): unexpected error: %s", i, err.Error())
				return
			}
			stop()
		})
	}

	if !shouldBe(t, &state, shot.StateStopped) {
		return
	}

	if err := state.Close(context.Background()); err != nil {
		t.Errorf("Close(): unexpected error: %s", err.Error())
		return
	}

	if !shouldBe(t, &state, shot.StateClosed) {
		return
	}

	waitGo(1, func(wg *sync.WaitGroup, _ int) {
		_, err := state.Start(context.Background())
		if err == nil {
			t.Error("Start(): expected an error")
		} else {
			t.Logf("Start(): got expected error: %s", err.Error())
		}
		wg.Done()
	})
}

func TestXor_Start_Stop(t *testing.T) {
	state := shot.NewXor(context.Background())

	waitGo(1, func(wg *sync.WaitGroup, _ int) {
		stop, err := state.Start(context.Background())
		wg.Done()
		if err != nil {
			t.Errorf("Start(): unexpected error: %s", err.Error())
			return
		}
		<-state.Context().Done()
		stop()
	})

	if !shouldBe(t, &state, shot.StateRunning) {
		return
	}

	if err := state.Stop(context.Background()); err != nil {
		t.Errorf("Stop(): unexpected error: %s", err.Error())
		return
	}

	if !shouldBe(t, &state, shot.StateStopped) {
		return
	}

	waitGo(1, func(wg *sync.WaitGroup, _ int) {
		stop, err := state.Start(context.Background())
		wg.Done()
		if err != nil {
			t.Errorf("Start(): unexpected error: %s", err.Error())
			return
		}
		<-state.Context().Done()
		stop()
	})

	if !shouldBe(t, &state, shot.StateRunning) {
		return
	}

	if err := state.Close(context.Background()); err != nil {
		t.Errorf("Close(): unexpected error: %s", err.Error())
		return
	}

	if !shouldBe(t, &state, shot.StateClosed) {
		return
	}

	waitGo(1, func(wg *sync.WaitGroup, _ int) {
		_, err := state.Start(context.Background())
		if err == nil {
			t.Error("Start(): expected an error")
		} else {
			t.Logf("Start(): got expected error: %s", err.Error())
		}
		wg.Done()
	})

	if !shouldBe(t, &state, shot.StateClosed) {
		return
	}

	if err := state.Stop(context.Background()); err == nil {
		t.Error("Stop(): expected and error")
		return
	} else {
		t.Logf("Stop(): got expected error: %s", err.Error())
	}

	if !shouldBe(t, &state, shot.StateClosed) {
		return
	}
}
