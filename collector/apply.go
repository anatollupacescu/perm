package collector

type Apply[X, T any] struct {
	collector[X, T]
	apply func(*X, T)
}

func WithApply[X, T any](c collector[X, T], apply func(*X, T)) *Apply[X, T] {
	return &Apply[X, T]{
		collector: c,
		apply:     apply,
	}
}

func (ap *Apply[X, T]) Collect(ctx *X, acc []T, current T) (doSkip, done bool) {
	ap.apply(ctx, current)
	return ap.collector.Collect(ctx, acc, current)
}
