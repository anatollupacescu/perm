package diehard

import (
	"testing"

	"github.com/anatollupacescu/perm"
)

func matchAnyCtx[T any](ctx *context, elem T, rules ...func(ctx *context, elem T) bool) bool {
	for _, rule := range rules {
		if rule(ctx, elem) {
			return true
		}
	}
	return false
}

func matchAny[T any](acc []T, rules ...func(acc []T) bool) bool {
	for _, rule := range rules {
		if rule(acc) {
			return true
		}
	}
	return false
}

func TestSimplePerm(t *testing.T) {
	var input []act
	input = append(input,
		act{name: "fill-up-3", cb: func(ctx *context) { ctx.jug_3 = 3 }},
		act{name: "empty-3", cb: func(ctx *context) { ctx.jug_3 = 0 }},
		act{name: "fill-up-5", cb: func(ctx *context) { ctx.jug_5 = 5 }},
		act{name: "empty-5", cb: func(ctx *context) { ctx.jug_5 = 0 }},
		act{name: "3-to-5", cb: transfer_3_to_5},
		act{name: "5-to-3", cb: transfer_5_to_3})

	skipRulesCtx := []func(ctx *context, acc act) bool{emptyAnEmptyJug, fillUpAfullJug, transferFromEmpty, transferToFull}
	skipRules := []func(acc []act) bool{initialStepIsFillingUp, repetitions, redundant}

	var solutions [][]act

	sink := func(acc []act) bool {
		if matchAny(acc, skipRules...) {
			return true
		}

		var ctx = new(context)

		for _, a := range acc[:len(acc)-1] {
			a.cb(ctx)
		}

		last := acc[len(acc)-1]

		if matchAnyCtx(ctx, last, skipRulesCtx...) {
			return true
		}

		last.cb(ctx)

		// invariant
		if ctx.jug_5 == 4 {
			// dont forget to make a copy because this memory location will be 'appended' to
			nacc := make([]act, len(acc))
			copy(nacc, acc)
			solutions = append(solutions, nacc)
			return true // skip next because all other longer solution will
			// be just this one with extra redundant actions at the end
		}

		return false
	}

	t.Run("no solutions within size 5", func(t *testing.T) {
		solutions = nil
		perm.Perm(sink, 5, input...)
		if want, got := 0, len(solutions); got != want {
			t.Fatalf("want %d solutions, got %d", want, got)
		}
	})

	t.Run("one solution with 6 steps within size 6", func(t *testing.T) {
		solutions = nil
		perm.Perm(sink, 6, input...)
		if want, got := 1, len(solutions); got != want {
			t.Fatalf("want %d solutions, got %d", want, got)
		}
		print(t, solutions)
	})

	t.Run("one solutions with 6 steps within size 7", func(t *testing.T) {
		solutions = nil
		perm.Perm(sink, 7, input...)
		if want, got := 1, len(solutions); got != want {
			t.Fatalf("want %d solutions, got %d", want, got)
		}
		print(t, solutions)
	})

	t.Run("one solutions with 8 steps", func(t *testing.T) {
		solutions = nil
		perm.Perm(sink, 8, input...)
		if want, got := 4, len(solutions); got != want {
			t.Fatalf("want %d solutions, got %d", want, got)
		}
		print(t, solutions)
	})
}

func print(t *testing.T, solutions [][]act) {
	for _, v := range solutions {
		t.Log("len", len(v))
		for _, s := range v {
			t.Log(s.name)
		}
	}
}
