package constraints

import "cmp"

type Int interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~int
}

type Signed = Int

type Uint interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint | ~uintptr
}

type Unsigned = Uint

type Float interface {
	~float32 | ~float64
}

type Complex interface {
	~complex64 | ~complex128
}

type Number interface {
	Int | Uint | Float | Complex
}

type Ordered = cmp.Ordered
