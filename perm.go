package perm

import (
	"iter"
	"math/bits"
)

func Of[T any](maxSize int, sink func([]T, T) bool, in ...T) {
	var perm func(acc []T, sink func([]T, T) bool)

	perm = func(acc []T, sink func([]T, T) bool) {
		for _, v := range in {
			if sink(acc, v) {
				continue
			}
			if len(acc) < maxSize-1 {
				perm(append(acc, v), sink)
			}
		}
	}

	perm(nil, sink)
}

// context holds previously computed values
func OfCtx[X, T any](maxSize int, sink func(*X, []T, T) bool, in ...T) {
	var (
		perm func(X, []T)
		ctx  X
	)

	perm = func(ctx X, acc []T) {
		for _, v := range in {
			ctx := ctx
			if sink(&ctx, acc, v) {
				continue
			}
			if len(acc) < maxSize-1 {
				perm(ctx, append(acc, v))
			}
		}
	}

	perm(ctx, nil)
}

func CombOf[T any](minSize, maxSize int, in ...T) iter.Seq[[]T] {
	if len(in) > 64 {
		in = in[:64]
	}
	if minSize < 0 {
		minSize = 0
	}
	if maxSize > 64 {
		maxSize = 64
	}
	if maxSize > len(in) {
		maxSize = len(in)
	}

	n := len(in)
	total := uint64(1) << n // safe because n <= 64

	return func(yield func([]T) bool) {
		for w := range total {
			// Count set bits to determine combination size upfront
			size := bits.OnesCount64(w)

			if size < minSize || size > maxSize {
				continue
			}

			picked := make([]T, 0, size)
			for i := range n {
				if w&(1<<i) != 0 {
					picked = append(picked, in[i])
				}
			}

			if !yield(picked) {
				return
			}
		}
	}
}
