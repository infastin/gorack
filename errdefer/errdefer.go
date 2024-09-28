package errdefer

func Close(err *error, fn func() error) error {
	if *err != nil {
		return fn()
	}
	return nil
}

func CloseNoError(err *error, fn func()) {
	if *err != nil {
		fn()
	}
}
