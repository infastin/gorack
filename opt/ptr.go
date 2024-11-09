package opt

func Ptr[T any](value T) *T {
	ptr := new(T)
	*ptr = value
	return ptr
}

func ZeroPtr[T comparable](value T) *T {
	var zero T
	if value == zero {
		return nil
	}
	ptr := new(T)
	*ptr = value
	return ptr
}
