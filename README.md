# perm

**Composable permutation and combination generation for Go — filter, transform, and collect sequences without materializing them first.**

`perm` lets you walk every permutation or combination of a set of values, shaping the traversal on the fly with chainable predicates. Instead of generating all results up front and filtering after, you prune, route, and collect *during* the walk — which keeps memory flat and makes complex rules easy to express.

---

## Table of Contents

- [Install](#install)
- [Core Concepts](#core-concepts)
- [Generating Permutations](#generating-permutations)
  - [`Of` — depth-first permutation walk](#of--depth-first-permutation-walk)
  - [`OfCtx` — stateful walk with a context](#ofctx--stateful-walk-with-a-context)
  - [`CombOf` — combinations as an iterator](#combof--combinations-as-an-iterator)
- [Predicates](#predicates)
- [Sinks & Combinators](#sinks--combinators)
- [Usage Examples](#usage-examples)
  - [Password policy checker](#password-policy-checker)
  - [Seating arrangements](#seating-arrangements)
  - [Combo meal builder](#combo-meal-builder)
  - [Running totals with context](#running-totals-with-context)
- [API Reference](#api-reference)

---

## Install

```bash
go get github.com/anatollupacescu/perm
```

Requires Go 1.23+ (uses `iter.Seq`).

---

## Core Concepts

`perm` is built around one idea: **a sink is a function that decides what to do with the next candidate element**.

```go
type Fn[T any] func(acc []T, next T) bool
```

- `acc` is the sequence built so far.
- `next` is the element being considered.
- Return `true` to **accept and stop recursing** into this branch (prune deeper).
- Return `false` to **decline** this candidate at this depth (allow recursion to continue).

Sinks and predicates are ordinary functions. You compose them by wrapping — `Take`, `Drop`, `And`, `Or` are all just higher-order functions that return another `Fn[T]`. There is no query builder, no DSL, no reflection.

---

## Generating Permutations

### `Of` — depth-first permutation walk

```go
func Of[T any](maxSize int, sink func([]T, T) bool, in ...T)
```

Walks all permutations of `in` up to length `maxSize`, calling `sink` at every step. Recursion into a branch is skipped when `sink` returns `true`.

### `OfCtx` — stateful walk with a context

```go
func OfCtx[X, T any](maxSize int, sink func(*X, []T, T) bool, in ...T)
```

Same as `Of` but threads a context value `X` through the walk. The context is copied at each branch, so sibling branches don't interfere with each other.

### `CombOf` — combinations as an iterator

```go
func CombOf[T any](minSize, maxSize int, in ...T) iter.Seq[[]T]
```

Returns a lazy `iter.Seq[[]T]` that yields every combination (order-independent subset) of `in` whose length falls in `[minSize, maxSize]`. Supports up to 64 input elements.

---

## Predicates

Predicates are `Fn[T]` values — they can be passed to `Take`, `Drop`, combined with `And`/`Or`, or used directly as a sink.

| Predicate | Description |
|---|---|
| `HasLen(n)` | True when the current sequence (including `next`) has exactly `n` elements |
| `Duplicate()` | True when `next` equals the last element in `acc` |
| `BeginsWithFn(v...)` | True when the sequence so far starts with the given prefix |
| `And(f...)` | True when **all** predicates match |
| `Or(f...)` | True when **any** predicate matches |
| `DevNull()` | Always false — swallows everything |

---

## Sinks & Combinators

| Function | Description |
|---|---|
| `Take(sink, pred)` | Forward to `sink` only when `pred` is true |
| `Drop(sink, pred)` | Forward to `sink` only when `pred` is **false** |
| `CollectF(fn)` | Materialise each complete sequence and call `fn([]T)` |
| `Collect(fn)` | Same as `CollectF` but `fn` has no return value |
| `Count(sink, *n)` | Increment `*n` for every call, then forward |
| `Filter(delegate, skip)` | Skip (return true) if `skip` matches, else delegate |
| `Peek(delegate, name, fn)` | Log/inspect without changing flow |
| `MinLen(n, delegate)` | Only forward once the sequence has at least `n` elements |
| `BeginsWith(delegate, name, v...)` | Forward only when the sequence starts with named values |
| `MutateCtx(delegate)` | Apply each element's `Mutate(*X)` to the context before forwarding |

---

## Usage Examples

### Password policy checker

Find every 3-character password from a small alphabet that satisfies a policy: must contain at least one digit and must not start with a special character.

```go
chars := []string{"a", "b", "1", "2", "!"}

var valid []string

collect := CollectF(func(seq []string) {
    valid = append(valid, strings.Join(seq, ""))
})

// Must be exactly 3 chars long
mustBe3 := Take(collect, HasLen[string](3))

// Must contain a digit somewhere
hasDigit := func(acc []string, next string) bool {
    for _, c := range append(acc, next) {
        if c == "1" || c == "2" {
            return true
        }
    }
    return false
}

// Must not start with "!"
noSpecialStart := Drop(mustBe3, BeginsWithFn("!"))

filter := Take(noSpecialStart, Or(hasDigit))

Of(3, filter, chars...)

fmt.Println(valid) // ["a1b", "1ab", "b2a", ...]
```

---

### Seating arrangements

You have five guests and a round table with four seats. Print every valid seating where Alice and Bob are never adjacent.

```go
guests := []string{"Alice", "Bob", "Carol", "Dave", "Eve"}

adjacent := func(acc []string, next string) bool {
    if len(acc) == 0 {
        return false
    }
    last := acc[len(acc)-1]
    pair := map[string]bool{last + next: true, next + last: true}
    return pair["AliceBob"] || pair["BobAlice"]
}

collect := CollectF(func(seq []string) {
    fmt.Println(seq)
})

// Drop any sequence where the incoming element sits next to a forbidden neighbour
noAdjacent := Drop(collect, adjacent)

// Only emit once we have a full table of 4
seated := Take(noAdjacent, HasLen[string](4))

Of(4, seated, guests...)
```

---

### Combo meal builder

A menu has mains, sides, and drinks. Build every 3-item combo that contains exactly one item from each category, and collect combos whose total calorie count is under 800.

```go
type Item struct {
    Name     string
    Category string
    Calories int
}

menu := []Item{
    {"Burger", "main", 550}, {"Salad", "main", 320},
    {"Fries", "side", 370},  {"Coleslaw", "side", 180},
    {"Cola", "drink", 150},  {"Water", "drink", 0},
}

var combos [][]Item

collect := CollectF(func(seq []Item) {
    var total int
    for _, it := range seq {
        total += it.Calories
    }
    if total < 800 {
        combos = append(combos, seq)
    }
})

// One of each category
onePerCategory := func(acc []Item, next Item) bool {
    for _, existing := range acc {
        if existing.Category == next.Category {
            return false // category already taken — drop this branch
        }
    }
    return true
}

filter := Take(
    collect,
    And(HasLen[Item](3), onePerCategory),
)

Of(3, filter, menu...)
```

---

### Running totals with context

Use `OfCtx` to accumulate a running sum as the walk progresses and prune branches that exceed a budget.

```go
type Budget struct{ Spent int }

type Product struct {
    Name  string
    Price int
}

func (p Product) Mutate(b *Budget) {
    b.Spent += p.Price
}

products := []Product{
    {"Coffee", 3}, {"Sandwich", 6}, {"Cookie", 2}, {"Juice", 4},
}

var baskets [][]Product

collect := CollectCtx(func(_ *Budget, seq []Product) {
    baskets = append(baskets, seq)
})

// Prune the moment spending exceeds £10
underBudget := func(b *Budget, _ []Product, _ Product) bool {
    return b.Spent > 10
}

sink := FilterCtx(collect, underBudget)
mut  := MutateCtx(sink)

OfCtx(3, mut, products...)

for _, basket := range baskets {
    fmt.Println(basket)
}
```

---

## API Reference

### Permutation generators

```go
// Walk permutations of in[] up to maxSize depth.
// sink returns true to prune recursion on a branch.
func Of[T any](maxSize int, sink func([]T, T) bool, in ...T)

// Walk with a per-branch context copied at each fork.
func OfCtx[X, T any](maxSize int, sink func(*X, []T, T) bool, in ...T)

// Lazy combination iterator (order-independent subsets).
// minSize/maxSize control subset cardinality. Supports up to 64 elements.
func CombOf[T any](minSize, maxSize int, in ...T) iter.Seq[[]T]
```

### Predicates

```go
func DevNull[T any]() Fn[T]
func HasLen[T any](exactSize int) Fn[T]
func Duplicate[T comparable]() Fn[T]
func BeginsWithFn[T comparable](values ...T) Fn[T]
func And[T any](filter ...Fn[T]) Fn[T]
func Or[T any](filter ...Fn[T]) Fn[T]
```

### Sinks & wiring

```go
func Take[T any](sink Fn[T], take Fn[T]) Fn[T]
func Drop[T any](sink Fn[T], drop Fn[T]) Fn[T]
func Filter[T any](delegate, skip func([]T, T) bool) func([]T, T) bool
func FilterCtx[X, T any](delegate, skip func(*X, []T, T) bool) func(*X, []T, T) bool

func CollectF[T any](sink func([]T)) func([]T, T) bool
func Collect[T any](sink func([]T)) func([]T, T)
func CollectCtx[X, T any](sink func(*X, []T)) func(*X, []T, T) bool

func Count[T any](delegate func([]T, T) bool, counter *int) func([]T, T) bool
func Peek[T any](delegate func([]T, T) bool, name func(T) string, peek func([]string)) func([]T, T) bool

func MinLen[T any](minSize int, delegate func([]T, T) bool) func([]T, T) bool
func MinLenCtx[X, T any](minSize int, delegate func(*X, []T, T) bool) func(*X, []T, T) bool

func BeginsWith[T any](delegate func([]T, T) bool, name func(T) string, values ...string) func([]T, T) bool
func MutateCtx[X any, T interface{ Mutate(*X) }](delegate func(*X, []T, T) bool) func(*X, []T, T) bool
```

---

## License

MIT
