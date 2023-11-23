package collector

type MinLenCollector[X, T any] struct {
	collector[X, T]
	minSize int
}

func WithMinLen[X, T any](c collector[X, T], minSize int) *MinLenCollector[X, T] {
	return &MinLenCollector[X, T]{
		collector: c,
		minSize:   minSize - 1,
	}
}

func (ml *MinLenCollector[X, T]) Collect(ctx *X, acc []T, current T) (doSkip, done bool) {
	if len(acc) < ml.minSize {
		return false, false
	}

	return ml.collector.Collect(ctx, acc, current)
}
