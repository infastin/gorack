package send

import "context"

func Drop[T any](ch chan<- T, val T) (sent bool) {
	select {
	case ch <- val:
		return true
	default:
		return false
	}
}

func DropCtx[T any](ctx context.Context, ch chan<- T, val T) (sent bool) {
	select {
	case <-ctx.Done():
		return false
	default:
		select {
		case <-ctx.Done():
			return false
		case ch <- val:
			return true
		}
	}
}

func WaitCtx[T any](ctx context.Context, ch chan<- T, val T) (sent bool) {
	select {
	case ch <- val:
		return true
	default:
		select {
		case <-ctx.Done():
			return false
		case ch <- val:
			return true
		}
	}
}

func Replace[T any](ch chan T, val T) (replaced bool) {
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

func ReplaceCb[T any](ch chan T, val T, cb func(old T)) (replaced bool) {
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
