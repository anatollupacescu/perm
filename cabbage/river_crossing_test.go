package river

import (
	"slices"
	"strings"
	"testing"

	"github.com/anatollupacescu/perm"
)

type context struct {
	Farmer, Wolf, Goat, Cabbage bool // true=crossed
}

type act struct {
	name   string
	Mutate func(*context)
}

func TestFindBoatConfiguration(t *testing.T) {
	var input = []act{
		{name: "wolf-across", Mutate: func(ctx *context) { ctx.Wolf = true; ctx.Farmer = true }},
		{name: "wolf-back", Mutate: func(ctx *context) { ctx.Wolf = false; ctx.Farmer = false }},
		{name: "goat-across", Mutate: func(ctx *context) { ctx.Goat = true; ctx.Farmer = true }},
		{name: "goat-back", Mutate: func(ctx *context) { ctx.Goat = false; ctx.Farmer = false }},
		{name: "cabbage-across", Mutate: func(ctx *context) { ctx.Cabbage = true; ctx.Farmer = true }},
		{name: "cabbage-back", Mutate: func(ctx *context) { ctx.Cabbage = false; ctx.Farmer = false }},
		{name: "farmer-across", Mutate: func(ctx *context) { ctx.Farmer = true }},
		{name: "farmer-back", Mutate: func(ctx *context) { ctx.Farmer = false }},
	}

	var solutions [][]string

	var totalChecked, foundPos int
	sink := perm.CollectSize(func(acc []act) {
		totalChecked++

		var ctx context

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

	perm.OfSizeWithSkip(7, skip, sink, input...)

	if len(solutions) != 2 {
		t.Fatalf("wanted two solution, got %d", len(solutions))
	}

	t.Logf("got solution at %d/%d", foundPos, totalChecked)

	for _, r := range solutions {
		t.Log(r)
	}
}

func skip(acc []act, current act) bool {

	// very first step
	if len(acc) == 0 {
		if strings.Contains(current.name, "-back") {
			return true
		}

		var ctx context
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

	var ctx context

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

func atState(acc []act, current act, s []string) bool {
	var aa []string
	for _, a := range acc {
		aa = append(aa, a.name)
	}

	aa = append(aa, current.name)

	eq := slices.Equal(aa, s)

	return eq
}

func TestAtState(t *testing.T) {

	t.Run("left", func(t *testing.T) {
		acc := []act{}
		b := atState(acc, act{name: "goat-across"}, []string{"goat-accross"})
		if b == true {
			t.Fail()
		}
	})

	t.Run("right", func(t *testing.T) {
		acc := []act{{name: "goat-across"}}
		b := atState(acc, act{name: "goat-across"}, []string{})
		if b == true {
			t.Fail()
		}
	})

	t.Run("match 1", func(t *testing.T) {
		acc := []act{{name: "goat-across"}}
		b := atState(acc, act{name: "wolf-across"}, []string{"goat-across", "wolf-across", "random"})
		if b == false {
			t.Fail()
		}
	})

	t.Run("match 2", func(t *testing.T) {
		acc := []act{
			{name: "goat-across"},
			{name: "farmer-back"},
			{name: "wolf-across"},
			{name: "goat-back"},
			{name: "cabbage-across"},
			{name: "farmer-back"},
			// {name: "goat-across"},
		}
		b := atState(acc, act{name: "goat-across"}, []string{"goat-across", "farmer-back", "wolf-across", "goat-back", "cabbage-across", "farmer-back", "goat-across"})
		if b == false {
			t.Fail()
		}
	})
}
