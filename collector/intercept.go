package collector

type Sink[X, T any] struct {
	collector[X, T]
	sink func(acc []T)
}

func WithSink[X, T any](c collector[X, T], sink func(acc []T)) *Sink[X, T] {
	return &Sink[X, T]{
		collector: c,
		sink:      sink,
	}
}

func (ml *Sink[X, T]) Collect(ctx *X, acc []T, current T) (doSkip, done bool) {
	nacc := make([]T, len(acc))
	copy(nacc, acc)
	ml.sink(append(nacc, current))

	return ml.collector.Collect(ctx, acc, current)
}
