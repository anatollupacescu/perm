package collector

type Collector[X, T any] struct {
	size int
}

func New[X any, T any](size int) *Collector[X, T] {
	return &Collector[X, T]{size: size - 1}
}

func (e *Collector[X, T]) Collect(ctx *X, acc []T, current T) (doSkip, done bool) {
	if len(acc) >= e.size {
		doSkip = true
	}

	return
}
