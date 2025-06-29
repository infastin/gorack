package errors

import "sync/atomic"

var numFrames atomic.Int64

// SetTrace sets the number of frames a stack trace will contain.
// If it's a zero then the stack trace won't be included in errors (default).
// The number of frames can't be negative.
func SetTrace(n int64) {
	if n < 0 {
		panic("number of frame can't be negative")
	}
	numFrames.Store(n)
}

func Trace() bool {
	return numFrames.Load() != 0
}
