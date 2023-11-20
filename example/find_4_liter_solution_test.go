package diehard

import (
	"slices"
	"testing"

	"github.com/anatollupacescu/perm"
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

	e := new(perm.Collector[context, act])

	e.AddSkipRule(initialStepIsFillingUp, // no point in starting with anything else
		emptyAnEmptyJug, fillUpAfullJug, transferFromEmpty,
		transferToFull, repetitions, redundant)

	e.AddInvariant("match 4 liters", func(ctx *context) bool {
		return ctx.jug_5 == 4 // solution: we've got 4 liters in the 5 liter jug
	})

	want := 6
	e.WantSolutions(want)

	perm.New[context, act](input).Perm(e.Collect, 8)

	solutions := e.Solutions()

	if len(solutions) != want {
		t.Fatalf("want %d solutions, got %d", want, len(solutions))
	}

	slices.SortFunc(solutions, func(a, b perm.Solution[act]) int {
		return len(a.Steps) - len(b.Steps)
	})

	for _, s := range solutions[0].Steps {
		t.Log(s.name)
	}
}
