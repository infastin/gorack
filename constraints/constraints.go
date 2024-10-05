package constraints

type Int interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~int
}

type Uint interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint | ~uintptr
}

type Float interface {
	~float32 | ~float64
}

type Number interface {
	Int | Uint | Float
}

type Ordered interface {
	~string | Int | Uint | Float
}
