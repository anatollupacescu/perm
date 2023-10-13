package perm

import (
	"slices"
	"testing"
)

func TestPermMore(t *testing.T) {
	want := [][]string{
		{"a", "a"}, {"a", "b"}, {"a", "c"}, {"a", "d"}, {"b", "a"}, {"b", "b"}, {"b", "c"}, {"b", "d"},
		{"c", "a"}, {"c", "b"}, {"c", "c"}, {"c", "d"}, {"d", "a"}, {"d", "b"}, {"d", "c"}, {"d", "d"},
	}

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
}

func TestPermFewer(t *testing.T) {
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
}
