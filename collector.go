package perm

type Transformation[X any] interface {
	Apply(*X)
}

type Collector[X any, T Transformation[X]] struct {
	scount int // want solution count
	halt   bool

	solutions  []Solution[T]
	invariants []invariant[X]

	skipRules []func(ctx *X, acc []T, current T) bool
}

func (c *Collector[X, T]) AddSkipRule(f ...func(ctx *X, acc []T, current T) bool) {
	c.skipRules = append(c.skipRules, f...)
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

func (e *Collector[X, T]) Collect(ctx *X, acc []T, current T) (doSkip bool) {
	if e.halt {
		return true
	}

	for _, skipRule := range e.skipRules {
		if skipRule(ctx, acc, current) { // match
			return true
		}
	}

	current.Apply(ctx)

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
			}
		}
	}

	return e.halt
}
