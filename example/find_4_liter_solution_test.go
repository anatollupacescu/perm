package diehard

import (
	"testing"

	"github.com/anatollupacescu/perm"
	c "github.com/anatollupacescu/perm/collector"
)

func TestFindStepSequence(t *testing.T) {
	var input []act
	input = append(input,
		act{name: "fill-up-3", cb: func(ctx *context) { ctx.jug_3 = 3 }},
		act{name: "empty-3", cb: func(ctx *context) { ctx.jug_3 = 0 }},
		act{name: "fill-up-5", cb: func(ctx *context) { ctx.jug_5 = 5 }},
		act{name: "empty-5", cb: func(ctx *context) { ctx.jug_5 = 0 }},
		act{name: "3-to-5", cb: transfer_3_to_5},
		act{name: "5-to-3", cb: transfer_5_to_3})

	t.Run("one solution with 6 steps", func(t *testing.T) {
		cl := c.New[context, act](6)

		var (
			solutions []act
			got       int
		)
		in := c.WithFn(cl, func(ctx *context, acc []act, a act) bool {
			if ctx.jug_5 == 4 {
				nacc := make([]act, len(acc))
				copy(nacc, acc)
				solutions = append(nacc, a)
				got++
				return true
			}
			return false
		})

		ml := c.WithMinLen(in, 5)

		ap := c.WithFn(ml, func(ctx *context, _ []act, a act) bool { a.cb(ctx); return false })

		sk := c.WithSkipRules(ap, initialStepIsFillingUp,
			emptyAnEmptyJug, fillUpAfullJug, transferFromEmpty,
			transferToFull, repetitions, redundant)

		want := 1

		perm.New[context, act](input).Perm(sk.Collect)

		if got != want {
			t.Fatalf("want %d solutions, got %d", want, got)
		}

		steps := 6
		if len(solutions) != steps {
			t.Fatalf("want %d steps, got %d", steps, len(solutions))
		}

		for _, s := range solutions {
			t.Log(s.name)
		}
	})

	t.Run("six solutions with 8 steps", func(t *testing.T) {
		cl := c.New[context, act](8)

		var got int
		in := c.WithFn(cl, func(ctx *context, acc []act, a act) bool {
			if ctx.jug_5 == 4 {
				got++
				return true
			}
			return false
		})

		ml := c.WithMinLen(in, 6)

		ap := c.WithFn(ml, func(ctx *context, acc []act, a act) (_ bool) {
			a.cb(ctx)
			return
		})

		sk := c.WithSkipRules(ap, initialStepIsFillingUp,
			emptyAnEmptyJug, fillUpAfullJug, transferFromEmpty,
			transferToFull, repetitions, redundant)

		want := 6

		perm.New[context, act](input).Perm(sk.Collect)

		if got != want {
			t.Fatalf("want %d solutions, got %d", want, got)
		}
	})
}
