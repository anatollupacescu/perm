package collector

type collector[X, T any] interface {
	Collect(ctx *X, acc []T, current T) (doSkip, done bool)
}

type SkipRuleCollector[X, T any] struct {
	collector[X, T]
	rules []func(ctx *X, acc []T, current T) bool
}

func WithSkipRules[X, T any](c collector[X, T], rules ...func(ctx *X, acc []T, current T) bool) *SkipRuleCollector[X, T] {
	return &SkipRuleCollector[X, T]{
		collector: c,
		rules:     rules,
	}
}

func (sr *SkipRuleCollector[X, T]) Collect(ctx *X, acc []T, current T) (doSkip, done bool) {
	for _, skipRule := range sr.rules {
		if skipRule(ctx, acc, current) { // match
			return true, false
		}
	}

	return sr.collector.Collect(ctx, acc, current)
}
