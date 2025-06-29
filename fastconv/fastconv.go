package fastconv

import (
	"fmt"
	"reflect"
	"unsafe"
)

// String converts a byte slice into a string without extra allocations.
// The bytes argument can be nil.
// Must not be used if the byte slice could be mutated.
func String[B ~[]byte](b B) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// Bytes converts a string into a byte slice without extra allocations.
// Must not be used if the resulting byte slice could be mutated.
func Bytes[S ~string](s S) []byte {
	return unsafe.Slice(unsafe.StringData(string(s)), len(s))
}

// TypePointer returns a type pointer.
func TypePointer(a any) uintptr {
	type emptyInterface struct {
		typ unsafe.Pointer
		ptr unsafe.Pointer
	}
	iface := (*emptyInterface)(unsafe.Pointer(&a))
	return uintptr(iface.typ)
}

// Slice converts a slice of one type to a slice of another type without extra allocations.
// Use only when type In can be converted to type Out.
func Slice[Out, In any, InSlice ~[]In](input InSlice) (output []Out) {
	outputPtr := (*Out)(unsafe.Pointer(unsafe.SliceData(input)))
	return unsafe.Slice(outputPtr, len(input))
}

// SliceSafe converts a slice of one type to a slice of another type without extra allocations.
// Will panic if type In can't be converted to type Out.
func SliceSafe[Out, In any, InSlice ~[]In](input InSlice) (output []Out) {
	if itp, otp := reflect.TypeFor[In](), reflect.TypeFor[Out](); itp.ConvertibleTo(otp) {
		panic(fmt.Sprintf("fastconv: %T is not convertible to %T", itp, otp))
	}
	outputPtr := (*Out)(unsafe.Pointer(unsafe.SliceData(input)))
	return unsafe.Slice(outputPtr, len(input))
}
