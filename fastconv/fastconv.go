package fastconv

import (
	"unsafe"
)

// Fast conversion of a byte slice to a string.
// The bytes argument can be nil.
// Must not be used if the byte slice could be mutated.
func String[B ~[]byte](bytes B) string {
	return unsafe.String(unsafe.SliceData(bytes), len(bytes))
}

// Fast conversion of a string to a byte slice.
// Must not be used if the resulting byte slice could be mutated.
func Bytes[S ~string](str S) []byte {
	return unsafe.Slice(unsafe.StringData(string(str)), len(str))
}

// Returns a type pointer.
func TypePointer(a any) uintptr {
	type emptyInterface struct {
		typ unsafe.Pointer
		ptr unsafe.Pointer
	}
	iface := (*emptyInterface)(unsafe.Pointer(&a))
	return uintptr(iface.typ)
}

// Fast conversion of a slice of one type to a slice of another type
// without copying. Use only if type I can be converted to type O.
func Slice[O, I any](input []I) (output []O) {
	outputPtr := (*O)(unsafe.Pointer(unsafe.SliceData(input)))
	return unsafe.Slice(outputPtr, len(input))
}
