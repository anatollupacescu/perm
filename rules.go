package perm

type RuleSetCtx[X, T any] []func(*X, []T, T) bool

func (m RuleSetCtx[X, T]) Match(ctx *X, acc []T, elem T) bool {
	for _, rule := range m {
		if rule(ctx, acc, elem) {
			return true
		}
	}
	return false
}

type RuleSet[T any] []func([]T, T) bool

func (m RuleSet[T]) Match(acc []T, c T) bool {
	for _, rule := range m {
		if rule(acc, c) {
			return true
		}
	}
	return false
}
