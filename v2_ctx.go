package perm

type FnCtx[X, T any] func(*X, []T, T) bool

func DevNullCtx[X, T any]() FnCtx[X, T] {
	return func(*X, []T, T) bool { return false }
}

func MutateCtx2[X any, T interface{ Mutate(*X) }](sink FnCtx[X, T]) FnCtx[X, T] {
	return func(ctx *X, in []T, c T) bool {
		c.Mutate(ctx)
		return sink(ctx, in, c)
	}
}

func TakeCtx[X, T any](sink FnCtx[X, T], take ...FnCtx[X, T]) FnCtx[X, T] {
	return func(ctx *X, t1 []T, t2 T) bool {
		for _, t := range take {
			if t(ctx, t1, t2) {
				return sink(ctx, t1, t2)
			}
		}

		return false
	}
}

func DropCtx[X, T any](sink FnCtx[X, T], drop FnCtx[X, T]) FnCtx[X, T] {
	return func(ctx *X, t1 []T, t2 T) bool {
		if drop(ctx, t1, t2) {
			return false
		}

		return sink(ctx, t1, t2)
	}
}

func AndCtx[X, T any](filter ...FnCtx[X, T]) FnCtx[X, T] {
	return func(ctx *X, t1 []T, t2 T) bool {
		for _, f := range filter {
			if !f(ctx, t1, t2) {
				return false
			}
		}

		return true
	}
}

func OrCtx[X, T any](filter ...FnCtx[X, T]) FnCtx[X, T] {
	return func(ctx *X, t1 []T, t2 T) bool {
		var matched bool
		for _, f := range filter {
			if f(ctx, t1, t2) {
				matched = true
				break
			}
		}

		return matched
	}
}
