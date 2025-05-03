package shot

import "fmt"

// State represents state of a resource.
type State int32

const (
	// Resource was just created and not running yet.
	StateCreated State = iota
	// Resource is running.
	StateRunning
	// Resource was stopped and can be started again.
	StateStopped
	// Resource was closed and can't be started again.
	StateClosed
)

func (s State) String() string {
	switch s {
	case StateCreated:
		return "created"
	case StateRunning:
		return "running"
	case StateStopped:
		return "stopped"
	case StateClosed:
		return "closed"
	default:
		panic(fmt.Sprintf("invalid state: %d", s))
	}
}
