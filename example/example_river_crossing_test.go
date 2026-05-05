package example

import (
	"strings"
	"testing"

	"github.com/anatollupacescu/perm"
)

type riverCtx struct {
	Farmer, Wolf, Goat, Cabbage bool // true=crossed
}

type riverAct struct {
	name   string
	mutate func(*riverCtx)
}

func (r riverAct) Mutate(ctx *riverCtx) {
	r.mutate(ctx)
}

var riverActions = []riverAct{
	{name: "wolf-across", mutate: func(ctx *riverCtx) { ctx.Wolf = true; ctx.Farmer = true }},
	{name: "wolf-back", mutate: func(ctx *riverCtx) { ctx.Wolf = false; ctx.Farmer = false }},
	{name: "goat-across", mutate: func(ctx *riverCtx) { ctx.Goat = true; ctx.Farmer = true }},
	{name: "goat-back", mutate: func(ctx *riverCtx) { ctx.Goat = false; ctx.Farmer = false }},
	{name: "cabbage-across", mutate: func(ctx *riverCtx) { ctx.Cabbage = true; ctx.Farmer = true }},
	{name: "cabbage-back", mutate: func(ctx *riverCtx) { ctx.Cabbage = false; ctx.Farmer = false }},
	{name: "farmer-across", mutate: func(ctx *riverCtx) { ctx.Farmer = true }},
	{name: "farmer-back", mutate: func(ctx *riverCtx) { ctx.Farmer = false }},
}

func TestRiverCrossingCtx(t *testing.T) {
	var (
		count     int
		solutions [][]string
	)

	sink := perm.CollectCtx(func(ctx *riverCtx, acc []riverAct) {
		count++
		if ctx.Cabbage && ctx.Wolf && ctx.Goat && ctx.Farmer {
			var names []string
			for _, act := range acc {
				names = append(names, act.name)
			}

			solutions = append(solutions, names)
		}
	})

	mutate := perm.MutateCtx(sink)

	minLen := perm.MinLenCtx(6, mutate)

	filters := perm.FilterCtx(minLen, skipRiverCtx)

	perm.OfCtx(7, filters, riverActions...)

	t.Log(count)

	if len(solutions) != 2 {
		t.Fatal("wanted 2 solutions, got", len(solutions))
	}
}

func skipRiverCtx(_ *riverCtx, acc []riverAct, tail riverAct) bool {
	return skipRiver(acc, tail)
}

func TestRiverCrossing(t *testing.T) {

	var solutions [][]string

	var totalChecked, foundPos int
	sink := perm.Collect(func(acc []riverAct) {
		totalChecked++

		var ctx riverCtx

		for _, a := range acc {
			a.Mutate(&ctx)
		}

		if ctx.Cabbage && ctx.Wolf && ctx.Goat && ctx.Farmer {
			foundPos = totalChecked

			var names []string
			for _, act := range acc {
				names = append(names, act.name)
			}

			solutions = append(solutions, names)
		}
	})

	filter := func(in []riverAct, c riverAct) bool {
		if skipRiver(in, c) {
			return true
		}

		sink(in, c)

		return false
	}

	perm.Of(7, filter, riverActions...)

	if len(solutions) != 2 {
		t.Fatalf("wanted two solution, got %d", len(solutions))
	}

	t.Logf("got solution at %d/%d", foundPos, totalChecked)

	for _, r := range solutions {
		t.Log(r)
	}
}

func skipRiver(acc []riverAct, current riverAct) bool {

	// very first step
	if len(acc) == 0 {
		if strings.Contains(current.name, "-back") {
			return true
		}

		var ctx riverCtx
		current.Mutate(&ctx)

		if ctx.Wolf && ctx.Goat && !ctx.Farmer {
			return true // wolf eats goat
		}
		if ctx.Cabbage && ctx.Goat && !ctx.Farmer {
			return true // goat eats cabbage
		}
		if !ctx.Wolf && !ctx.Goat && ctx.Farmer {
			return true // wolf eats goat
		}
		if !ctx.Cabbage && !ctx.Goat && ctx.Farmer {
			return true // goat eats cabbage
		}
	}

	if len(acc) > 0 {
		// redundant
		last := acc[len(acc)-1]
		if last.name == current.name {
			return true
		}
	}

	var ctx riverCtx

	for _, a := range acc {
		a.Mutate(&ctx)
	}

	if ctx.Wolf && current.name == "wolf-across" {
		return true
	}
	if !ctx.Wolf && current.name == "wolf-back" {
		return true
	}
	if ctx.Cabbage && current.name == "cabbage-across" {
		return true
	}
	if !ctx.Cabbage && current.name == "cabbage-back" {
		return true
	}
	if ctx.Farmer && current.name == "farmer-across" {
		return true
	}
	if !ctx.Farmer && current.name == "farmer-back" {
		return true
	}
	if ctx.Goat && current.name == "goat-across" {
		return true
	}
	if !ctx.Goat && current.name == "goat-back" {
		return true
	}

	if strings.Contains(current.name, "-across") && ctx.Farmer {
		return true
	}
	if strings.Contains(current.name, "-back") && !ctx.Farmer {
		return true
	}

	current.Mutate(&ctx)

	if ctx.Wolf && ctx.Goat && !ctx.Farmer {
		return true // wolf eats goat
	}
	if ctx.Cabbage && ctx.Goat && !ctx.Farmer {
		return true // goat eats cabbage
	}
	if !ctx.Wolf && !ctx.Goat && ctx.Farmer {
		return true // wolf eats goat
	}
	if !ctx.Cabbage && !ctx.Goat && ctx.Farmer {
		return true // goat eats cabbage
	}

	return false
}
