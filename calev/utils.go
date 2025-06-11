package calev

type uintSet interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64
}

func rangeToSet[T uintSet](low, high, step int) T {
	var value T
	switch step {
	case 0:
		value = 1 << low
	case 1:
		value = (1<<(high+1) - 1) &^ (1<<low - 1)
	default:
		for i := low; i <= high; i += step {
			value |= 1 << i
		}
	}
	return value
}

func rangeToSetReverse[T uintSet](low, high, step, max int) T {
	var value T
	switch step {
	case 0:
		value = 1 << (max - low)
	case 1:
		value = (1<<(max-low+1) - 1) &^ (1<<(max-high) - 1)
	default:
		for i := low; i <= high; i += step {
			value |= 1 << (max - i)
		}
	}
	return value
}
