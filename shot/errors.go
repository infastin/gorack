package shot

import "errors"

var (
	// Tried to call Start on resource in the the Running state.
	ErrRunning = errors.New("already running")
	// Tried to call Start on the resource in the Closed state.
	ErrClosed = errors.New("closed")
)
