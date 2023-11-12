package perm

type node[X, T any] struct {
	v    T
	refs []*node[X, T]
}

type collector[X, T any] func(ctx *X, acc []T, current T) bool

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

func (n *node[X, T]) Perm(ev collector[X, T], size int) {
	var ctx X
	n.perm(ctx, ev, nil, size)
}

func (n *node[X, T]) perm(ctx X, collect collector[X, T], acc []T, size int) {
	var done bool
	if len(acc) == size-1 {
		done = true
	}
	for _, v := range n.refs {
		ctx := ctx
		if skip := collect(&ctx, acc, v.v); skip {
			continue
		}
		if !done {
			v.perm(ctx, collect, append(acc, v.v), size)
		}
	}
}

func perm[X, T any](n *node[X, T], acc []T, sink func([]T), size int) {
	var done bool
	if len(acc) == size-1 {
		done = true
	}
	for _, v := range n.refs {
		if done {
			sink(append(acc, v.v))
			continue
		}
		perm(v, append(acc, v.v), sink, size)
	}
}
