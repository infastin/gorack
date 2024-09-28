package shot

import "errors"

var (
	ErrRunning    = errors.New("already running")
	ErrStopped    = errors.New("stopped")
	ErrNotRunning = errors.New("not running")
	ErrClosed     = errors.New("closed")
)
