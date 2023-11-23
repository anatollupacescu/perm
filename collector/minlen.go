package perm

type MinLenCollector[X, T any] struct {
	*Collector[X, T]
}

func WithMinLen[X, T any](c *Collector[X, T]) *MinLenCollector[X, T] {
	return &MinLenCollector[X, T]{
		Collector: c,
	}
}

func (ml *MinLenCollector[X, T]) Collect(ctx *X, acc []T, current T) (doSkip, done bool) {
	if len(acc) < ml.size {
		return false, false
	}

	return ml.Collector.Collect(ctx, acc, current)
}
