package perm

func MinLen[T any](size int, delegate func(in []T, c T) bool) func(in []T, c T) bool {
	return func(in []T, c T) bool {
		if len(in) < size-1 {
			return false
		}
		return delegate(in, c)
	}
}

func Collect[T any](sink func([]T)) func(in []T, c T) bool {
	return func(in []T, c T) bool {
		nacc := make([]T, 0, len(in)+1)
		nacc = append(nacc, in...)
		sink(append(nacc, c))
		return false
	}
}

func MinLenCtx[X, T any](size int, delegate func(*X, []T, T) bool) func(*X, []T, T) bool {
	return func(ctx *X, in []T, c T) bool {
		if len(in) < size-1 {
			return false
		}
		return delegate(ctx, in, c)
	}
}

func CollectCtx[X, T any](sink func(*X, []T)) func(*X, []T, T) bool {
	return func(ctx *X, in []T, c T) bool {
		nacc := make([]T, 0, len(in)+1)
		nacc = append(nacc, in...)
		sink(ctx, append(nacc, c))
		return false
	}
}

func FilterCtx[X, T any](delegate, skip func(*X, []T, T) bool) func(*X, []T, T) bool {
	return func(ctx *X, in []T, c T) bool {
		if skip(ctx, in, c) {
			return true
		}
		return delegate(ctx, in, c)
	}
}

func Filter[T any](delegate, skip func([]T, T) bool) func([]T, T) bool {
	return func(in []T, c T) bool {
		if skip(in, c) {
			return true
		}
		return delegate(in, c)
	}
}
