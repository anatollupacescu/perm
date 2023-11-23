package perm

import (
	"slices"
	"testing"

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

	New[any, string]([]string{"x", "y", "z"}).Perm(sr.Collect)
}

func TestPermMoreValsThanSize(t *testing.T) {
	want := [][]string{
		{"a", "a"}, {"a", "b"}, {"a", "c"}, {"a", "d"}, {"b", "a"}, {"b", "b"}, {"b", "c"}, {"b", "d"},
		{"c", "a"}, {"c", "b"}, {"c", "c"}, {"c", "d"}, {"d", "a"}, {"d", "b"}, {"d", "c"}, {"d", "d"},
	}

	t.Run("sink", func(t *testing.T) {
		tree := New[any, string]([]string{"a", "b", "c", "d"})

		var res [][]string
		sink := func(in []string) {
			res = append(res, in)
		}

		perm(tree, nil, sink, 2)

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

		New[any, string]([]string{"a", "b", "c", "d"}).Perm(ml.Collect)

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
		tree := New[any, string]([]string{"1", "0"})

		var res [][]string
		sink := func(in []string) {
			res = append(res, in)
		}

		perm(tree, nil, sink, 3)

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

		New[any, string]([]string{"1", "0"}).Perm(ml.Collect)

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

func TestPerm(t *testing.T) {
	want := [][]string{
		{"a", "a", "a"},
		{"a", "a", "b"},
		{"a", "a", "c"},

		{"a", "b", "a"},
		{"a", "b", "b"},
		{"a", "b", "c"},

		{"a", "c", "a"},
		{"a", "c", "b"},
		{"a", "c", "c"},

		{"b", "a", "a"},
		{"b", "a", "b"},
		{"b", "a", "c"},

		{"b", "b", "a"},
		{"b", "b", "b"},
		{"b", "b", "c"},

		{"b", "c", "a"},
		{"b", "c", "b"},
		{"b", "c", "c"},

		{"c", "a", "a"},
		{"c", "a", "b"},
		{"c", "a", "c"},

		{"c", "b", "a"},
		{"c", "b", "b"},
		{"c", "b", "c"},

		{"c", "c", "a"},
		{"c", "c", "b"},
		{"c", "c", "c"},
	}

	t.Run("sinc", func(t *testing.T) {
		tree := New[any, string]([]string{"a", "b", "c"})

		var res [][]string
		sink := func(in []string) {
			res = append(res, in)
		}

		perm(tree, nil, sink, 3)

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

		New[any, string]([]string{"a", "b", "c"}).Perm(ml.Collect)

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
