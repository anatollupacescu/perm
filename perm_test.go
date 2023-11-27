package perm_test

import (
	"slices"
	"testing"

	"github.com/anatollupacescu/perm"
	c "github.com/anatollupacescu/perm/collector"
)

func TestSkipRule(t *testing.T) {
	cl := c.New[any, string](2)

	sk := c.WithSink(cl, func(acc []string) {
		// t.Log(acc)
	})

	ml := c.WithMinLen(sk, 2)

	sr := c.WithSkipRules(ml, func(ctx *any, acc []string, current string) bool {
		return len(acc) == 1 && current == "y"
	})

	perm.New[any, string]([]string{"x", "y", "z"}).Perm(sr.Collect)
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
		c := perm.Collect(func(s []string) {
			res = append(res, s)
		})

		ml := perm.MinLen(two, c)

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

	t.Run("collector", func(t *testing.T) {
		cl := c.New[any, string](2)

		var res [][]string
		sk := c.WithSink(cl, func(acc []string) {
			res = append(res, acc)
		})

		ml := c.WithMinLen(sk, 2)

		perm.New[any, string]([]string{"a", "b", "c", "d"}).Perm(ml.Collect)

		var index int
		for _, v := range want {
			c := res[index]
			index++
			if !slices.Equal(c, v) {
				t.Fatalf("\nwant: %v\ngot:  %v", v, c)
			}
		}
	})

	type context struct{}

	t.Run("perm ctx", func(t *testing.T) {
		var index int
		sink := perm.CollectCtx[context](func(_ *context, got []string) {
			if !slices.Equal(got, want[index]) {
				t.Fatalf("\nwant: %v\ngot:  %v", want, got)
			}
			index++
		})

		const two = 2

		sink = perm.MinLenCtx(two, sink)

		perm.OfCtx(two, sink, "a", "b", "c", "d")
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
		c := perm.Collect(func(s []string) {
			res = append(res, s)
		})

		ml := perm.MinLen(size, c)

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

	t.Run("collector", func(t *testing.T) {
		cl := c.New[any, string](2)

		var res [][]string
		sk := c.WithSink(cl, func(acc []string) {
			res = append(res, acc)
		})

		ml := c.WithMinLen(sk, 3)

		perm.New[any, string]([]string{"1", "0"}).Perm(ml.Collect)

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

	t.Run("sinc", func(t *testing.T) {
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

	t.Run("collector", func(t *testing.T) {
		cl := c.New[any, string](3)

		var res [][]string
		sk := c.WithSink(cl, func(acc []string) {
			res = append(res, acc)
		})

		ml := c.WithMinLen(sk, 3)

		perm.New[any, string]([]string{"a", "b", "c"}).Perm(ml.Collect)

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
		sink := perm.CollectCtx[context](func(ctx *context, got []int) {
			total += ctx.count
		})

		sink = perm.FilterCtx(sink, func(ctx *context, acc []int, current int) bool {
			ctx.count += current
			return false
		})

		const three = 3

		sink = perm.MinLenCtx(three, sink)

		perm.OfCtx(three, sink, 1, 0)

		want := 4

		if total != want {
			t.Fatalf("\nwant: %v\ngot:  %v", want, total)
		}
	})
}
