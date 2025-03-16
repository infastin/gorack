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

func ConvertPtr[T, U any](ptr *T, fn func(T) U) *U {
	if ptr == nil {
		return nil
	}
	res := new(U)
	*res = fn(*ptr)
	return res
}

func Deref[T any](ptr *T, def T) T {
	if ptr == nil {
		return def
	}
	return *ptr
}
