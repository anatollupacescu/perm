package perm

import (
	"slices"
)

type Fn[T any] func([]T, T) bool

func DevNull[T any]() Fn[T] {
	return func([]T, T) bool { return false }
}

func Take[T any](sink Fn[T], take Fn[T]) Fn[T] {
	return func(t1 []T, t2 T) bool {
		if take(t1, t2) {
			return sink(t1, t2)
		}

		return false
	}
}

func Drop[T any](sink Fn[T], drop Fn[T]) Fn[T] {
	return func(t1 []T, t2 T) bool {
		if drop(t1, t2) {
			return false
		}
		return sink(t1, t2)
	}
}

func HasLen[T any](exactSize int) Fn[T] {
	return func(t1 []T, t2 T) bool {
		want := append(t1, t2)
		return len(want) == exactSize
	}
}

func Duplicate[T comparable]() Fn[T] {
	return func(t1 []T, t2 T) bool {
		if len(t1) > 0 {
			index := len(t1) - 1
			return t1[index] == t2
		}
		return false
	}
}

func BeginsWithFn[T comparable](values ...T) Fn[T] {
	return func(t1 []T, t2 T) bool {
		want := append(t1, t2)
		if len(want) < len(values) {
			return false
		}
		target := want[:len(values)]
		return slices.Equal(target, values)
	}
}

func And[T any](filter ...Fn[T]) Fn[T] {
	return func(t1 []T, t2 T) bool {
		for _, f := range filter {
			if !f(t1, t2) {
				return false
			}
		}

		return true
	}
}

func Or[T any](filter ...Fn[T]) Fn[T] {
	return func(t1 []T, t2 T) bool {
		var matched bool
		for _, f := range filter {
			if f(t1, t2) {
				matched = true
				break
			}
		}

		return matched
	}
}
