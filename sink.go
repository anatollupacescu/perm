package perm

func CollectF[T any](sink func([]T)) func(in []T, c T) bool {
	return func(in []T, c T) bool {
		nacc := make([]T, 0, len(in)+1)
		nacc = append(nacc, in...)
		sink(append(nacc, c))
		return false
	}
}

func Collect[T any](sink func([]T)) func(in []T, c T) {
	return func(in []T, c T) {
		nacc := make([]T, 0, len(in)+1)
		nacc = append(nacc, in...)
		sink(append(nacc, c))
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

func Filter[T any](delegate, skip func([]T, T) bool) func([]T, T) bool {
	return func(in []T, c T) bool {
		if skip(in, c) {
			return true
		}
		return delegate(in, c)
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

func MutateCtx[X any, T interface{ Mutate(*X) }](delegate func(*X, []T, T) bool) func(*X, []T, T) bool {
	return func(ctx *X, in []T, c T) bool {
		for _, v := range in {
			v.Mutate(ctx)
		}

		c.Mutate(ctx)

		return delegate(ctx, in, c)
	}
}

func BeginsWith[T any](delegate func([]T, T) bool, name func(T) string, values ...string) func([]T, T) bool {
	return func(in []T, c T) bool {

		var slice = append(in, c)

		if len(values) > len(slice) {
			return false
		}

		for i, v := range values {
			if v != name(slice[i]) {
				return false
			}
		}

		return delegate(in, c)
	}
}

func Count[T any](delegate func([]T, T) bool, counter *int) func([]T, T) bool {
	return func(in []T, c T) bool {
		*counter++
		delegate(in, c)
		return false
	}
}

func Peek[T any](delegate func([]T, T) bool, name func(T) string, peek func([]string)) func([]T, T) bool {
	return func(in []T, c T) bool {
		var names []string
		for _, v := range append(in, c) {
			names = append(names, name(v))
		}

		peek(names)

		delegate(in, c)
		return false
	}
}

func MinLen[T any](minSize int, delegate func(in []T, c T) bool) func(in []T, c T) bool {
	return func(in []T, c T) bool {
		if len(in)+1 >= minSize {
			delegate(in, c)
		}
		return false
	}
}

func MinLenCtx[X, T any](minSize int, delegate func(*X, []T, T) bool) func(*X, []T, T) bool {
	return func(ctx *X, in []T, c T) bool {
		if len(in)+1 >= minSize {
			return delegate(ctx, in, c)
		}
		return false
	}
}
