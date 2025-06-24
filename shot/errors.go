package shot

import "errors"

var (
	ErrRunning = errors.New("already running")
	ErrClosed  = errors.New("closed")
)
