package diehard

import (
	"testing"

	"github.com/anatollupacescu/perm"
)

var input = append([]act{},
	act{name: "fill-up-3", Mutate: func(ctx *context) { ctx.jug_3 = 3 }},
	act{name: "empty-3", Mutate: func(ctx *context) { ctx.jug_3 = 0 }},
	act{name: "fill-up-5", Mutate: func(ctx *context) { ctx.jug_5 = 5 }},
	act{name: "empty-5", Mutate: func(ctx *context) { ctx.jug_5 = 0 }},
	act{name: "3-to-5", Mutate: transfer_3_to_5},
	act{name: "5-to-3", Mutate: transfer_5_to_3})

var skipRulesCtx = perm.RuleSetCtx[context, act]([]func(*context, []act, act) bool{emptyAnEmptyJug, fillUpAfullJug, transferFromEmpty, transferToFull})

func TestPermExpensiveCtx(t *testing.T) {
	var solutions [][]act

	// looking for the sequence of steps that leeds 4 liters of water
	invariant := func(ctx *context, acc []act, current act) bool {
		if ctx.jug_5 == 4 {
			// dont forget to make a copy because this memory location will be mutated
			nacc := make([]act, len(acc))
			copy(nacc, acc)
			solutions = append(solutions, append(nacc, current))
			return true // skip next because all other longer solution will
			// be just this one with extra redundant actions at the end
		}

		return false
	}

	mutate := perm.FilterCtx(invariant, func(ctx *context, _ []act, current act) bool {
		current.Mutate(ctx)
		return false
	})

	skipRules := perm.RuleSetCtx[context, act]([]func(*context, []act, act) bool{repetitionsCtx, redundantCtx})
	filter := perm.FilterCtx(mutate, skipRules.Match)

	// cheap operations next
	filterCtx := perm.FilterCtx(filter, skipRulesCtx.Match)

	// only makes sense to start with filling up
	startState := perm.RuleSetCtx[context, act]([]func(*context, []act, act) bool{initStepIsFillingUpCtx})
	sink := perm.FilterCtx(filterCtx, startState.Match)

	t.Run("no solutions within size 5", func(t *testing.T) {
		solutions = nil
		perm.OfCtx(5, sink, input...)
		if want, got := 0, len(solutions); got != want {
			t.Fatalf("want %d solutions, got %d", want, got)
		}
	})

	t.Run("one solution with 6 steps within size 6", func(t *testing.T) {
		solutions = nil
		perm.OfCtx(6, sink, input...)
		if want, got := 1, len(solutions); got != want {
			t.Fatalf("want %d solutions, got %d", want, got)
		}
		print(t, solutions)
	})

	t.Run("one solutions with 6 steps within size 7", func(t *testing.T) {
		solutions = nil
		perm.OfCtx(7, sink, input...)
		if want, got := 1, len(solutions); got != want {
			t.Fatalf("want %d solutions, got %d", want, got)
		}
		print(t, solutions)
	})

	t.Run("four solutions with 8 steps", func(t *testing.T) {
		solutions = nil
		perm.OfCtx(8, sink, input...)
		if want, got := 4, len(solutions); got != want {
			t.Fatalf("want %d solutions, got %d", want, got)
		}
		print(t, solutions)
	})
}

func build(acc []act, extra ...act) *context {
	var ctx = new(context)
	for _, a := range acc {
		a.Mutate(ctx)
	}
	for _, a := range extra {
		a.Mutate(ctx)
	}
	return ctx
}

func TestPermCheapContext(t *testing.T) {
	var solutions [][]act

	sink := func(acc []act, current act) bool {
		var ctx = build(acc)
		current.Mutate(ctx)

		// invariant
		if ctx.jug_5 == 4 {
			// dont forget to make a copy because this memory location will be 'appended' to
			nacc := make([]act, len(acc))
			copy(nacc, acc)
			solutions = append(solutions, append(nacc, current))
			return true // skip next because all other longer solution will
			// be just this one with extra redundant actions at the end
		}

		return false
	}

	// filter out actions irrelevant for the current context
	sink = perm.Filter(sink, func(acc []act, current act) bool {
		var ctx = build(acc)
		return skipRulesCtx.Match(ctx, nil, current)
	})

	skipRules := perm.RuleSet[act]([]func([]act, act) bool{repetitions, redundant})
	sink = perm.Filter(sink, skipRules.Match)

	startState := perm.RuleSet[act]([]func([]act, act) bool{initStepIsFillingUp})
	sink = perm.Filter(sink, startState.Match)

	t.Run("no solutions within size 5", func(t *testing.T) {
		solutions = nil
		perm.Of(5, sink, input...)
		if want, got := 0, len(solutions); got != want {
			t.Fatalf("want %d solutions, got %d", want, got)
		}
	})

	t.Run("one solution with 6 steps up to size 6", func(t *testing.T) {
		solutions = nil
		perm.Of(6, sink, input...)
		if want, got := 1, len(solutions); got != want {
			t.Fatalf("want %d solutions, got %d", want, got)
		}
		print(t, solutions)
	})

	t.Run("one solutions with 6 steps up to size 7", func(t *testing.T) {
		solutions = nil
		perm.Of(7, sink, input...)
		if want, got := 1, len(solutions); got != want {
			t.Fatalf("want %d solutions, got %d", want, got)
		}
		print(t, solutions)
	})

	t.Run("four solutions with up to 8 steps", func(t *testing.T) {
		solutions = nil
		perm.Of(8, sink, input...)
		if want, got := 4, len(solutions); got != want {
			t.Fatalf("want %d solutions, got %d", want, got)
		}
		print(t, solutions)
	})
}

func print(t *testing.T, solutions [][]act) {
	var doPrint bool //= true
	if doPrint {
		for _, v := range solutions {
			t.Log("len", len(v))
			for _, s := range v {
				t.Log(s.name)
			}
		}
	}
}
