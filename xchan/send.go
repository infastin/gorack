package xchan

import "context"

func SendContext[T any](ctx context.Context, ch chan<- T, val T) (sent bool) {
	select {
	case <-ctx.Done():
		return false
	case ch <- val:
		return true
	}
}

func SendDrop[T any](ch chan<- T, val T) (dropped bool) {
	select {
	case ch <- val:
		return false
	default:
		return true
	}
}

func SendReplace[T any](ch chan T, val T) (replaced bool) {
	for {
		select {
		case ch <- val:
			return replaced
		default:
			select {
			case <-ch:
				replaced = true
			default:
			}
		}
	}
}

func SendReplaceFunc[T any](ch chan T, val T, cb func(old T)) (replaced bool) {
	for {
		select {
		case ch <- val:
			return replaced
		default:
			select {
			case old := <-ch:
				cb(old)
				replaced = true
			default:
			}
		}
	}
}
