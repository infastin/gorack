package opt

type Opt[T any] interface {
	Set(value T)
	Ptr() *T
	IsZero() bool
	Get() (value T, ok bool)
	Or(value T) T
}

func ConvertOpt[T any, O Opt[T], U any](opt O, f func(T) U) U {
	value, ok := opt.Get()
	if !ok {
		var zero U
		return zero
	}
	return f(value)
}
