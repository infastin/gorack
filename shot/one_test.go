package shot_test

import (
	"context"
	"sync"
	"testing"

	"github.com/infastin/gorack/shot"
)

func TestOne_Start_Close(t *testing.T) {
	state := shot.NewOne(context.Background())

	waitGo(1, func(wg *sync.WaitGroup, _ int) {
		stop, err := state.Start()
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

	waitGo(1, func(wg *sync.WaitGroup, _ int) {
		_, err := state.Start()
		if err == nil {
			t.Error("Start(): expected an error")
		} else {
			t.Logf("Start(): got expected error: %s", err.Error())
		}
		wg.Done()
	})

	if err := state.Close(context.Background()); err != nil {
		t.Errorf("Close(): unexpected error: %s", err.Error())
		return
	}

	if !shouldBe(t, &state, shot.StateClosed) {
		return
	}

	waitGo(1, func(wg *sync.WaitGroup, _ int) {
		_, err := state.Start()
		if err == nil {
			t.Error("Start(): expected an error")
		} else {
			t.Logf("Start(): got expected error: %s", err.Error())
		}
		wg.Done()
	})
}

func TestOne_Start_stop(t *testing.T) {
	state := shot.NewOne(context.Background())

	waitGo(1, func(wg *sync.WaitGroup, _ int) {
		defer wg.Done()
		stop, err := state.Start()
		if err != nil {
			t.Errorf("Start(): unexpected error: %s", err.Error())
			return
		}
		stop()
	})

	if !shouldBe(t, &state, shot.StateClosed) {
		return
	}

	waitGo(1, func(wg *sync.WaitGroup, _ int) {
		_, err := state.Start()
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

	if err := state.Close(context.Background()); err != nil {
		t.Errorf("Close(): unexpected error: %s", err.Error())
		return
	}

	if !shouldBe(t, &state, shot.StateClosed) {
		return
	}
}
