package example

import (
	"log"
	"testing"

	"github.com/anatollupacescu/perm"
)

const (
	closed = iota
	open
	locked
	unlocked
)

type doorCtx struct {
	door, lock rune
}

type doorAct struct {
	name   string
	Mutate func(*doorCtx)
}

func TestLockOnOpenDoor(t *testing.T) {
	var input = []doorAct{
		{name: "close-door", Mutate: func(ctx *doorCtx) { ctx.door = closed }},
		{name: "open-door", Mutate: func(ctx *doorCtx) { ctx.door = open }},
		{name: "lock", Mutate: func(ctx *doorCtx) { ctx.lock = locked }},
		{name: "unlock", Mutate: func(ctx *doorCtx) { ctx.lock = unlocked }},
	}

	var res [][]string

	sink := perm.CollectF(func(s []doorAct) {
		var names []string
		for _, act := range s {
			names = append(names, act.name)
		}
		res = append(res, names)
	})

	fault := perm.Filter(sink, fault) // solution

	var count int
	counter := perm.Count(fault, &count)

	firstStep := perm.BeginsWith(counter, actionName, "unlock", "open-door")

	peek := perm.Peek(firstStep, actionName, func(s []string) { log.Println(s) })

	repetitions := perm.Filter(peek, duplicates)

	perm.Of(3, repetitions, input...)

	if len(res) != 1 {
		t.Fatal("wanted 1 solution, got", len(res))
	}

	t.Log("total solutions checked", count)

	for _, r := range res {
		t.Log(r)
	}
}

func actionName(d doorAct) string { return d.name }

func duplicates(acc []doorAct, a doorAct) bool {
	if len(acc) == 0 {
		return false
	}

	// redundant
	last := acc[len(acc)-1]
	if last.name == a.name {
		return true
	}

	return false
}

func fault(acc []doorAct, a doorAct) bool {
	var c doorCtx

	for _, a := range acc {
		a.Mutate(&c)
	}

	a.Mutate(&c)

	return !(c.door == open && c.lock == locked)
}
