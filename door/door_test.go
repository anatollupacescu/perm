package door

import (
	"testing"

	"github.com/anatollupacescu/perm"
)

type act struct {
	name   string
	Mutate func(*context)
}

const (
	closed = iota
	open
	locked
	unlocked
)

type context struct {
	door, lock rune
}

func TestLockOnOpenDoor(t *testing.T) {
	var input = []act{
		{name: "close-door", Mutate: func(ctx *context) { ctx.door = closed }},
		{name: "open-door", Mutate: func(ctx *context) { ctx.door = open }},
		{name: "lock", Mutate: func(ctx *context) { ctx.lock = locked }},
		{name: "unlock", Mutate: func(ctx *context) { ctx.lock = unlocked }},
	}

	var res [][]string

	var totalChecked, foundPos int
	sink := perm.Collect(func(s []act) {
		totalChecked++

		if !fault(s) {
			return
		}

		foundPos = totalChecked

		var names []string
		for _, act := range s {
			names = append(names, act.name)
		}
		res = append(res, names)
	})

	sink = perm.MinLen(2, sink)

	sink = perm.Filter(sink, skip)

	perm.Of(2, sink, input...)

	t.Logf("got solution at %d/%d", foundPos, totalChecked)

	for _, r := range res {
		t.Log(r)
	}
}

func skip(acts []act, a act) bool {
	if len(acts) > 0 {
		// reduntant
		last := acts[len(acts)-1]
		if last.name == a.name {
			return true
		}

		// when door is locked we can only unlock it
		var ctx context

		for _, a := range acts {
			a.Mutate(&ctx)
		}

		if ctx.lock == locked && a.name != "unlock" {
			return true
		}
	}

	return false
}

func fault(s []act) bool {
	var c context

	for _, a := range s {
		a.Mutate(&c)
	}

	return c.door == open && c.lock == locked
}
