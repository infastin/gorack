package shot

type State int32

const (
	StateCreated State = iota
	StateRunning
	StateStopped
	StateClosed
)
