package perm_test

import (
	"slices"
	"testing"

	. "github.com/anatollupacescu/perm"
)

type v int
type ctx struct{ total int }

func (v v) Mutate(c *ctx) {
	c.total += int(v)
}

func TestAndOrCtx(t *testing.T) {
	isEven := func(_ *ctx, _ []int, v int) bool {
		return v == 2
	}

	isOne := func(_ *ctx, _ []int, v int) bool {
		return v == 1
	}

	sink := func(c *ctx, in []int, v int) bool {
		t.Log(append(in, v))
		return false
	}

	filter := OrCtx(isEven, isOne)

	take := TakeCtx(sink, filter)

	OfCtx(2, take, 0, 1)
}

func TestAndOr(t *testing.T) {
	isEven := func(_ []int, v int) bool {
		return v == 2
	}

	isOne := func(_ []int, v int) bool {
		return v == 1
	}

	sink := func(in []int, v int) bool {
		t.Log(append(in, v))
		return false
	}

	take := Take(sink, Or(isEven, isOne))

	Of(2, take, 0, 1, 2)
}

func TestV2Ctx(t *testing.T) {
	stat := NewStat()

	var want = [][]v{
		{1, 2},
		{2, 1},
	}

	var step int
	sum3 := func(c *ctx, in []v, ta v) bool {
		if c.total == 3 {
			if !slices.Equal(append(in, ta), want[step]) {
				t.Fatalf("slices not equal, want %v, got %v", want[step], append(in, ta))
			}
			step++
			return true
		}

		return false
	}

	take := TakeStatCtx(DevNullCtx[ctx, v](), stat, "take(sum=3)", sum3)

	mut := MutateCtx2(take)

	OfCtx(2, mut, 1, 2)

	t.Log("test stats:")
	stat.Print(log(t))
}

func TestMatchTail(t *testing.T) {
	stat := NewStat()

	var want = [][]string{
		{"1", "b"},
		{"b"},
		{"b", "b"},
	}

	var step int
	assert := func(in []string, ta string) bool {
		if !slices.Equal(append(in, ta), want[step]) {
			t.Fatalf("slices not equal, want %v, got %v at step %d", want[step], append(in, ta), step)
		}
		step++
		return false
	}

	endsWithB := func(_ []string, ta string) bool {
		return ta == "b"
	}

	take := TakeStat(assert, stat, "take(tail eq b)", endsWithB)

	Of(2, take, "1", "b")

	stat.Print(log(t))
}

func logSink[T any](t *testing.T) func(_ []T, ta T) bool {
	return func(in []T, ta T) bool {
		t.Log(append(in, ta))
		return false
	}
}

func log(t *testing.T) func(name string, count int) {
	return func(name string, count int) { t.Log(name, count) }
}

func TestTakeDropStat(t *testing.T) { // delete me
	sink := CollectF(func(in []string) {
		t.Log(in)
	})

	stat := NewStat()

	drop := DropStat(sink, stat, "drop(starts with b)", BeginsWithFn("b"))
	take := TakeStat(drop, stat, "take(length of two)", Or(HasLen[string](2), HasLen[string](1)))

	Of(2, take, "1", "b")

	stat.Print(log(t))
}

func TestTakeDropOr(t *testing.T) {
	want := [][]string{
		{"1", "1"},
		{"1", "b"},
	}

	var c int
	sink := CollectF(func(in []string) {
		if c > 1 {
			t.Fatal("expected to stop at step 2, got", c+1)
		}
		if !slices.Equal(in, want[c]) {
			t.Fatalf("invalid input, want %s, got %s", want[c], in)
		}
		c++
	})

	var count int
	counter := Count(sink, &count)

	take := Take(counter, HasLen[string](2))

	filter := Drop(take, Or(BeginsWithFn("b"), BeginsWithFn("c")))

	Of(2, filter, "1", "b")

	if count != 2 {
		t.Fatal("want count 2, got", count)
	}
}
