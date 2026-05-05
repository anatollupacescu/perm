package perm

import (
	"slices"
	"testing"
)

func TestFilter(t *testing.T) {
	sink := CollectF(func(in []string) {
		if len(in) > 1 {
			target := in[:2]

			if slices.Equal([]string{"a", "a"}, target) {
				t.Fatal("should have not passed filter")
			}

			if slices.Equal([]string{"b", "b"}, target) {
				t.Fatal("should have not passed filter")
			}
		}
	})

	skipAA := func(in []string, tail string) bool {
		return slices.Equal(append(in, tail), []string{"a", "a"})
	}

	skipBB := func(in []string, tail string) bool {
		return slices.Equal(append(in, tail), []string{"b", "b"})
	}

	filter := Filter(sink, skipAA)
	filter = Filter(filter, skipBB)

	Of(3, filter, "a", "b", "c")
}

func TestBeginsWith(t *testing.T) {
	sink := CollectF(func(in []string) {
		if in[0] != "$" {
			t.Fatal("only expecting $ as the first element")
		}
	})

	bw := BeginsWith(sink, func(in string) string { return in }, "$")

	Of(2, bw, "$", "#")
}

func TestMinLen(t *testing.T) {
	t.Run("no min len", func(t *testing.T) {
		var collected int
		sink := CollectF(func([]string) {
			collected++
		})

		Of(2, sink, "alfa", "beta")

		if collected != 6 {
			t.Fatal("want 6 ")
		}
	})

	t.Run("min len enforces length", func(t *testing.T) {
		var collected int
		sink := CollectF(func(in []string) {
			if len(in) != 2 {
				t.Fatal("unexpected size min len")
			}
			collected++
		})

		minLen := MinLen(2, sink)
		Of(2, minLen, "alfa", "beta")

		if collected != 4 {
			t.Fatal("want 4 got", collected)
		}
	})

	t.Run("min len does not limit length", func(t *testing.T) {
		var collected int
		sink := CollectF(func(in []string) {
			if len(in) < 2 {
				t.Fatal("unexpected size min len")
			}
			collected++
		})

		minLen := MinLen(2, sink)
		Of(3, minLen, "alfa", "beta", "gamma")

		if collected != 36 {
			t.Fatal("want 36 got", collected)
		}
	})
}
