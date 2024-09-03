package perm_test

import (
	"slices"
	"testing"

	"github.com/anatollupacescu/perm"
)

func TestPermOfSize(t *testing.T) {
	sink := func(s []string) {
		if len(s) != 4 {
			t.Fatal("wrong size")
		}
	}

	perm.OfSize(4, sink, "a", "b", "c", "d")
}

func TestPermMoreValsThanSize(t *testing.T) {
	want := [][]string{
		{"a", "a"}, {"a", "b"}, {"a", "c"}, {"a", "d"},
		{"b", "a"}, {"b", "b"}, {"b", "c"}, {"b", "d"},
		{"c", "a"}, {"c", "b"}, {"c", "c"}, {"c", "d"},
		{"d", "a"}, {"d", "b"}, {"d", "c"}, {"d", "d"},
	}

	t.Run("sink", func(t *testing.T) {
		const two = 2

		var res [][]string
		sink := perm.Collect(func(s []string) {
			// could make a copy but because of min len
			// no further appends will happen on this memory location
			res = append(res, s)
		})

		ml := perm.MinLen(two, sink)

		perm.Of(two, ml, "a", "b", "c", "d")

		if len(res) != len(want) {
			t.Fatalf("want result size: %d, got %d", len(want), len(res))
		}

		var index int
		for _, v := range want {
			c := res[index]
			index++
			if !slices.Equal(c, v) {
				t.Fatalf("\nwant: %v\ngot:  %v", v, c)
			}
		}
	})

	type context struct {
		allSameCount bool
	}

	t.Run("perm ctx", func(t *testing.T) {
		var total, index int
		sink := perm.CollectCtx[context](func(ctx *context, got []string) {
			if !slices.Equal(got, want[index]) {
				t.Fatalf("\nwant: %v\ngot:  %v", want, got)
			}
			index++
			if ctx.allSameCount {
				total++
			}
		})

		const two = 2

		sink = perm.FilterCtx(sink, func(ctx *context, acc []string, current string) bool {
			nacc := make([]string, len(acc))
			copy(nacc, acc)
			nacc = append(nacc, current)
			c := slices.Compact(nacc)
			if len(c) == 1 {
				ctx.allSameCount = true
			}
			return false
		})

		sink = perm.MinLenCtx(two, sink)

		perm.OfCtx(two, sink, "a", "b", "c", "d")

		want := 4 // same on the diagonal
		if total != want {
			t.Fatalf("\nwant: %v\ngot:  %v", want, total)
		}
	})
}

func TestPermFewerValsThanSize(t *testing.T) {
	want := [][]string{
		{"1", "1", "1"},
		{"1", "1", "0"},
		{"1", "0", "1"},
		{"1", "0", "0"},
		{"0", "1", "1"},
		{"0", "1", "0"},
		{"0", "0", "1"},
		{"0", "0", "0"},
	}

	t.Run("sink", func(t *testing.T) {
		const size = 3

		var res [][]string
		sink := perm.Collect(func(s []string) {
			res = append(res, s)
		})

		ml := perm.MinLen(size, sink)

		perm.Of(size, ml, "1", "0")

		if len(res) != len(want) {
			t.Fatalf("want result size: %d, got %d", len(want), len(res))
		}

		var index int
		for _, v := range want {
			c := res[index]
			index++
			if !slices.Equal(c, v) {
				t.Fatalf("\nwant: %v\ngot:  %v", v, c)
			}
		}
	})

	t.Run("perm ctx", func(t *testing.T) {
		var index int
		sink := perm.CollectCtx[any](func(_ *any, got []string) {
			if !slices.Equal(got, want[index]) {
				t.Fatalf("\nwant: %v\ngot:  %v", want, got)
			}
			index++
		})

		const three = 3

		sink = perm.MinLenCtx(three, sink)

		perm.OfCtx(three, sink, "1", "0")
	})
}

func TestPerm(t *testing.T) {
	want := [][]string{
		{"a", "a", "a"}, {"a", "a", "b"}, {"a", "a", "c"},
		{"a", "b", "a"}, {"a", "b", "b"}, {"a", "b", "c"},
		{"a", "c", "a"}, {"a", "c", "b"}, {"a", "c", "c"},
		{"b", "a", "a"}, {"b", "a", "b"}, {"b", "a", "c"},
		{"b", "b", "a"}, {"b", "b", "b"}, {"b", "b", "c"},
		{"b", "c", "a"}, {"b", "c", "b"}, {"b", "c", "c"},
		{"c", "a", "a"}, {"c", "a", "b"}, {"c", "a", "c"},
		{"c", "b", "a"}, {"c", "b", "b"}, {"c", "b", "c"},
		{"c", "c", "a"}, {"c", "c", "b"}, {"c", "c", "c"},
	}

	t.Run("sink", func(t *testing.T) {
		const size = 3

		var res [][]string
		c := perm.Collect(func(s []string) {
			res = append(res, s)
		})

		ml := perm.MinLen(size, c)

		perm.Of(size, ml, "a", "b", "c")

		if len(res) != len(want) {
			t.Fatalf("want result size: %d, got %d", len(want), len(res))
		}

		var index int
		for _, v := range want {
			c := res[index]
			index++
			if !slices.Equal(c, v) {
				t.Fatalf("\nwant: %v\ngot:  %v", v, c)
			}
		}
	})
}

func TestMutateCtx(t *testing.T) {
	type context struct {
		count int
	}

	t.Run("sum", func(t *testing.T) {
		var total int
		sink := perm.CollectCtx[context](func(ctx *context, got []int) {
			total += ctx.count
		})

		sink = perm.FilterCtx(sink, func(ctx *context, acc []int, current int) bool {
			ctx.count += current
			return false
		})

		const three = 3
		perm.OfCtx(three, sink, 1, 0)

		want := 17

		if total != want {
			t.Fatalf("\nwant: %v\ngot:  %v", want, total)
		}
	})

	t.Run("sum filter", func(t *testing.T) {
		var total int
		sum := perm.CollectCtx[context](func(ctx *context, got []int) {
			total += ctx.count
		})

		filter := perm.FilterCtx(sum, func(ctx *context, acc []int, current int) bool {
			ctx.count += current
			return false
		})

		const three = 3

		sink := perm.MinLenCtx(three, filter)

		perm.OfCtx(three, sink, 1, 0)

		want := 4

		if total != want {
			t.Fatalf("\nwant: %v\ngot:  %v", want, total)
		}
	})
}
