package collector

type Fn[X, T any] struct {
	collector[X, T]
	apply func(*X, []T, T) bool
}

func WithFn[X, T any](c collector[X, T], apply func(*X, []T, T) bool) *Fn[X, T] {
	return &Fn[X, T]{
		collector: c,
		apply:     apply,
	}
}

func (fn *Fn[X, T]) Collect(ctx *X, acc []T, current T) (doSkip, done bool) {
	if fn.apply(ctx, acc, current) {
		return true, false
	}
	return fn.collector.Collect(ctx, acc, current)
}
