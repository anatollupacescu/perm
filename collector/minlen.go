package perm

type MinLenCollector[X any, T Transformation[X]] struct {
	*Collector[X, T]
}

func WithMinLen[X any, T Transformation[X]](c *Collector[X, T]) *MinLenCollector[X, T] {
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
