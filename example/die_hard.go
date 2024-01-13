package diehard

import "slices"

type context struct {
	jug_3, jug_5 int
}

type act struct {
	name   string
	Mutate func(*context)
}

func transfer_3_to_5(ctx *context) {
	j5 := ctx.jug_5
	ctx.jug_5 += ctx.jug_3
	if ctx.jug_5 > 5 { // overflows back into jug 5
		ctx.jug_3 = ctx.jug_5 - 5
		ctx.jug_5 = 5
		return
	}
	ctx.jug_3 -= ctx.jug_5 - j5
}

func transfer_5_to_3(ctx *context) {
	j3 := ctx.jug_3
	ctx.jug_3 += ctx.jug_5
	if ctx.jug_3 > 3 {
		ctx.jug_5 = ctx.jug_3 - 3
		ctx.jug_3 = 3
		return
	}
	ctx.jug_5 -= ctx.jug_3 - j3
}

// skip rules
func emptyAnEmptyJug(ctx *context, _ []act, current act) bool {
	if current.name == "empty-3" && ctx.jug_3 == 0 {
		return true
	}
	if current.name == "empty-5" && ctx.jug_5 == 0 {
		return true
	}
	return false
}

func fillUpAfullJug(ctx *context, _ []act, current act) bool {
	if current.name == "fill-up-3" && ctx.jug_3 == 3 {
		return true
	}
	if current.name == "fill-up-5" && ctx.jug_5 == 5 {
		return true
	}
	return false
}

func transferFromEmpty(ctx *context, _ []act, current act) bool {
	if current.name == "3-to-5" && ctx.jug_3 == 0 {
		return true
	}
	if current.name == "5-to-3" && ctx.jug_5 == 0 {
		return true
	}
	return false
}

func transferToFull(ctx *context, _ []act, current act) bool {
	if current.name == "3-to-5" && ctx.jug_5 == 5 {
		return true
	}
	if current.name == "5-to-3" && ctx.jug_3 == 3 {
		return true
	}
	return false
}

var acceptable = []string{"fill-up-3", "fill-up-5"}

func initStepIsFillingUpCtx(_ *context, acc []act, current act) bool {
	return initStepIsFillingUp(acc, current)
}

func initStepIsFillingUp(acc []act, current act) bool {
	if len(acc) == 0 {
		validFirstInput := slices.Contains(acceptable, current.name)
		return !validFirstInput
	}
	return false
}

func repetitionsCtx(_ *context, acc []act, current act) bool {
	return repetitions(acc, current)
}

func repetitions(acc []act, current act) bool {
	if len(acc) == 0 {
		return false
	}
	return prev(acc).name == current.name
}

func redundantCtx(_ *context, acc []act, current act) bool {
	return redundant(acc, current)
}

func redundant(acc []act, current act) bool {
	if len(acc) == 0 {
		return false
	}
	seq := []string{prev(acc).name, current.name}
	slices.Sort(seq)
	if slices.Equal(seq, []string{"empty-3", "fill-up-3"}) {
		return true
	}
	if slices.Equal(seq, []string{"empty-5", "fill-up-5"}) {
		return true
	}
	if slices.Equal(seq, []string{"3-to-5", "5-to-3"}) {
		return true
	}
	return false
}

func prev(acc []act) act {
	return acc[len(acc)-1]
}
