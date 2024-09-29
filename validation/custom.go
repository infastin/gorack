package validation

func Custom[T Validatable](v T) error {
	return v.Validate()
}

func CustomRule[T Validatable]() AnyRule[T] {
	return AnyRuleFunc[T](func(v T) error {
		return v.Validate()
	})
}
