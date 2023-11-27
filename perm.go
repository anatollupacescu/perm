package perm

type node[X, T any] struct {
	v    T
	refs []*node[X, T]
}

type collector[X, T any] func(ctx *X, acc []T, current T) (bool, bool)

func New[X, T any](in []T) *node[X, T] {
	var refs = make([]*node[X, T], 0, len(in))
	for _, v := range in {
		refs = append(refs, &node[X, T]{v: v})
	}

	// recursive struct
	for _, r := range refs {
		r.refs = refs
	}

	return &node[X, T]{refs: refs}
}

func (n *node[X, T]) Perm(ev collector[X, T]) {
	var ctx X
	n.perm(ctx, ev, nil)
}

func (n *node[X, T]) perm(ctx X, collect collector[X, T], acc []T) {
	for _, v := range n.refs {
		ctx := ctx
		skip, done := collect(&ctx, acc, v.v)
		if skip {
			continue
		}
		if done {
			break
		}
		v.perm(ctx, collect, append(acc, v.v))
	}
}

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
