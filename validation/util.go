package validation

type validatorScope []bool

func (s validatorScope) Ok() bool {
	return len(s) == 0 || s[len(s)-1]
}

func (s validatorScope) Empty() bool {
	return len(s) != 0
}

func (s validatorScope) Set(condition bool) {
	s[len(s)-1] = !s[len(s)-1]
}

func (s validatorScope) Push(condition bool) validatorScope {
	return append(s, condition)
}

func (s validatorScope) Pop() validatorScope {
	return s[:len(s)-1]
}
