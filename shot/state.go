package shot

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
