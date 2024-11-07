package errors

import (
	"fmt"
	"io"
	"iter"
	"path"
	"runtime"
	"strconv"
)

type Frame runtime.Frame

// Formats the frame according to the fmt.Formatter interface.
//
//	%f    source file
//	%d    source line
//	%x    function name
//	%s    equivalent to %x\n\t%f:%d
//	%v    equivalent to %s
func (f *Frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 'f':
		io.WriteString(s, path.Base(f.File))
	case 'd':
		io.WriteString(s, strconv.Itoa(f.Line))
	case 'x':
		io.WriteString(s, f.Func.Name())
	case 's', 'v':
		io.WriteString(s, f.Function)
		io.WriteString(s, "\n\t")
		io.WriteString(s, f.File)
		io.WriteString(s, ":")
		io.WriteString(s, strconv.Itoa(f.Line))
	}
}

// Formats a stacktrace frame as a text string.
// The output is the same as that of fmt.Sprintf("%x %f:%d", f).
func (f *Frame) String() string {
	return fmt.Sprintf("%s %s:%d", f.Func.Name(), f.File, f.Line)
}

// Formats a stacktrace frame as a text string.
// The output is the same as that of fmt.Sprintf("%x %f:%d", f).
func (f Frame) MarshalText() ([]byte, error) {
	return fmt.Appendf(nil, "%s %s:%d", f.Func.Name(), f.File, f.Line), nil
}

type StackTrace []uintptr

// Formats the stack of Frames according to the fmt.Formatter interface.
//
//	%s   prints filename, function, and line number for each Frame in the stack.
//	%v   equivalent to %s
func (st StackTrace) Format(s fmt.State, verb rune) {
	switch verb {
	case 's', 'v':
		for _, frame := range st.Iter() {
			frame.Format(s, verb)
			io.WriteString(s, "\n")
		}
	}
}

func (st StackTrace) Iter() iter.Seq2[int, Frame] {
	return func(yield func(int, Frame) bool) {
		frames := runtime.CallersFrames(st)
		for i := 0; ; i++ {
			f, ok := frames.Next()
			if !ok || !yield(i, Frame(f)) {
				return
			}
		}
	}
}

func callers(skip int) StackTrace {
	n := numFrames.Load()
	if n == 0 {
		return nil
	}
	pcs := make([]uintptr, n)
	n = int64(runtime.Callers(skip, pcs))
	return pcs[:n]
}
