package collector

import "fmt"

type Collector[X, T any] struct {
	scount int // want solution count
	halt   bool
	size   int

	solutions  []Solution[T]
	invariants []invariant[X]
}

func New[X any, T any](size int, invs ...func(*X) bool) *Collector[X, T] {
	c := &Collector[X, T]{size: size - 1}
	for i, inv := range invs {
		c.invariants = append(c.invariants, invariant[X]{
			name:  fmt.Sprintf("main-%d", i),
			check: inv,
		})
	}
	return c
}

func (c *Collector[X, T]) AddInvariant(name string, check func(*X) bool) {
	c.invariants = append(c.invariants, invariant[X]{
		name:  name,
		check: check,
	})
}

func (c *Collector[X, T]) WantSolutions(count int) {
	c.scount = count
}

func (c *Collector[X, T]) Solutions() []Solution[T] {
	return c.solutions
}

type invariant[X any] struct {
	name  string
	check func(*X) bool
}

type Solution[T any] struct {
	Name  string
	Steps []T
}

func (e *Collector[X, T]) Collect(ctx *X, acc []T, current T) (doSkip, done bool) {
	if e.halt {
		return false, true
	}

	for _, inv := range e.invariants {
		if inv.check(ctx) {
			nacc := make([]T, len(acc))
			copy(nacc, acc)
			s := Solution[T]{
				Name:  inv.name,
				Steps: append(nacc, current),
			}
			e.solutions = append(e.solutions, s)
			doSkip = true
			if len(e.solutions) >= e.scount {
				e.halt = true
				return false, true
			}
		}
	}

	if len(acc) >= e.size {
		return true, false
	}

	return
}
