package perm

// fixed size, no progression
func OfSize[T any](size int, sink func([]T), in ...T) {
	size--
	var perm func(acc []T, sink func([]T))
	perm = func(acc []T, sink func([]T)) {
		for _, v := range in {
			if len(acc) < size {
				perm(append(acc, v), sink)
				continue
			}
			sink(append(acc, v))
		}
	}

	perm(nil, sink)
}

func Of[T any](maxSize int, sink func([]T, T) bool, in ...T) {
	var perm func(acc []T, sink func([]T, T) bool)

	perm = func(acc []T, sink func([]T, T) bool) {
		for _, v := range in {
			if sink(acc, v) {
				continue
			}
			if len(acc) < maxSize-1 {
				perm(append(acc, v), sink)
			}
		}
	}

	perm(nil, sink)
}

// context holds previously computed values
func OfCtx[X, T any](maxSize int, sink func(*X, []T, T) bool, in ...T) {
	var (
		perm func(X, []T)
		ctx  X
	)

	perm = func(ctx X, acc []T) {
		for _, v := range in {
			ctx := ctx
			if sink(&ctx, acc, v) {
				continue
			}
			if len(acc) < maxSize-1 {
				perm(ctx, append(acc, v))
			}
		}
	}

	perm(ctx, nil)
}
