# Where to go next: generics and the road after the course

_Лид (summary):_ **The fiftieth and last lesson of the Go course. Generics are one function instead of three, and the condition on the type is checked by the compiler: measured, a `string` does not get through before the program runs. Measured too is what they gain over `any`: 534 nanoseconds against 1801. Plus where to go from here.**

## Why this is needed

The course is ending. The blog works, sits on a domain of its own under HTTPS, stops tidily and survives a reboot. Two conversations are left: about what we did not cover, and about where to go next.

What we did not cover is **generics** — a way of writing one piece of code for many types. That they come last is not an accident: you can live without them for a long time, and they are easier to understand once you have written three nearly identical functions by hand.

## The whole thing at once

A new folder, `go mod init sabaq47`:

```go
package main

import (
	"fmt"
	"sort"
)

// Number is not a type but a condition: anything that adds with a plus and
// behaves like a number. The tilde means "and your own types made of these".
type Number interface {
	~int | ~int64 | ~float64
}

// Sum adds a slice of any numbers. Before generics this was three functions —
// for int, for int64 and for float64 — with one and the same body.
func Sum[T Number](xs []T) T {
	var total T
	for _, x := range xs {
		total += x
	}
	return total
}

// Keys returns the keys of a map. K comparable, because only what can be
// compared may be a key; V any, because the value is of no interest here.
func Keys[K comparable, V any](m map[K]V) []K {
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// Views is the counter from the goroutines lesson, except that it no longer
// knows what it counts: articles, tags or readers.
type Views[K comparable] struct {
	counts map[K]int
}

func NewViews[K comparable]() *Views[K] {
	return &Views[K]{counts: make(map[K]int)}
}

func (v *Views[K]) Add(key K)       { v.counts[key]++ }
func (v *Views[K]) Count(key K) int { return v.counts[key] }

// A score is a type of its own with an int inside. The tilde in Number lets it in.
type Score int

func main() {
	fmt.Println("== one function instead of three")
	fmt.Println("whole:      ", Sum([]int{2, 3, 4}))
	fmt.Println("fractional: ", Sum([]float64{0.5, 0.25}))
	fmt.Println("own type:   ", Sum([]Score{10, 20}))

	fmt.Println()
	fmt.Println("== the keys of any map")
	ages := map[string]int{"Austen": 41, "Orwell": 46}
	names := Keys(ages)
	sort.Strings(names)
	fmt.Println("strings:    ", names)
	years := Keys(map[int]string{1775: "Austen", 1903: "Orwell"})
	sort.Ints(years)
	fmt.Println("numbers:    ", years)

	fmt.Println()
	fmt.Println("== a counter that does not care what it counts")
	byArticle := NewViews[string]()
	byArticle.Add("first-entry")
	byArticle.Add("first-entry")
	fmt.Println("by article: ", byArticle.Count("first-entry"))

	byYear := NewViews[int]()
	byYear.Add(2026)
	fmt.Println("by year:    ", byYear.Count(2026))
}
```

The output:

```
== one function instead of three
whole:       9
fractional:  0.75
own type:    30

== the keys of any map
strings:     [Austen Orwell]
numbers:     [1775 1903]

== a counter that does not care what it counts
by article:  2
by year:     1
```

## The walk-through

### A type parameter is one more argument, except it is a type

An ordinary function takes values. A generic takes a **type** as well — in square brackets before the ordinary ones:

```go
func Sum[T Number](xs []T) T
```

`T` is the name of a type that becomes known at the call site. Called with `[]int`, `T` is `int`; called with `[]float64`, it is `float64`. The compiler fills the type in itself, which is why nowhere in the example is there a `Sum[int](...)`: it is visible from the argument.

Without generics this is three functions with one body: `SumInt`, `SumInt64`, `SumFloat`. The first mistake in that body is three corrections, and one of them will be forgotten some day.

### The condition is the heart of a generic

`Number` is not a type, and not an interface in the old sense. It is a **condition**: the list of what `T` is allowed to be.

```go
type Number interface {
	~int | ~int64 | ~float64
}
```

The tilde means "and any type of your own with this inside". That is why the example also adds up `Score`, declared as `type Score int`.

The condition is not decoration. Without it the body would not compile: the compiler does not know that `+` may be applied to `T`. And it checks the call as well:

```
$ go build ./...
./main.go:20:17: string does not satisfy Number (string missing in ~int | ~int64 | ~float64)
```

Strings add with a plus, but they are not numbers. The error comes **before the program runs** rather than on the first unlucky evening — and that is the whole point of the exercise.

### What this gains over `any`

Before generics, mixed things were added through `any` and a type switch on the spot. It works, and it is paid for. Measured on a slice of a thousand numbers:

```
BenchmarkSumGeneric-8    2238902     533.7 ns/op     0 B/op   0 allocs/op
BenchmarkSumAny-8         663471    1801   ns/op     0 B/op   0 allocs/op
```

**534 nanoseconds against 1801** — three and a half times. The difference is not memory (neither version allocates) but work: the `any` version examines the type of every element, while the generic knows the type in advance and adds the numbers directly.

This is worth measuring yourself: `go test -bench . -benchmem` — the same command as in the lesson on tests, asking about speed instead.

### When generics are not needed

The beginner's temptation is to make everything generic. Three rules that hold it back:

- **One function for one type is not a reason.** `func words(s string) int` stays as it is; a type parameter adds brackets and nothing else.
- **An interface is often better.** When different types need different behaviour, that is an interface, like `blog.Store` in the lesson on interfaces. A generic is about the same behaviour over different kinds of data.
- **Readability is worth more.** `func Map[T, U any](xs []T, f func(T) U) []U` is elegant, but three of those nested read worse than an ordinary loop. In Go a loop is not a disgrace.

A good sign: you have **already written** the second nearly identical function and are copying the third. Now.

### The standard library already has them

Half of what you would want to write does not need writing:

- `slices` — `slices.Sort`, `slices.Contains`, `slices.Index`, `slices.Max`;
- `maps` — `maps.Keys`, `maps.Values`. One subtlety: they return a sequence rather than a slice, so a slice comes out as `keys := slices.Collect(maps.Keys(m))`. Our `Keys` above is a teaching copy that returns the slice directly;
- `cmp` — comparison and `cmp.Or`;
- `sync.OnceValue` — compute once and hand back what was computed.

## Where to go next

The course is over; the language is not. Seven directions, each with one line about what it is for:

1. **`slices`, `maps`, `cmp`** — what you can now read: they are generics all the way down.
2. **Postgres instead of SQLite** — the same `database/sql`, another driver, and a real connection pool. SQLite in the blog was a choice, not a law.
3. **`errgroup` and `sync.WaitGroup`** — for when there are forty goroutines rather than two, and all of them must finish.
4. **`net/http/pprof`** — to know where a program spends its time instead of guessing.
5. **`go test -fuzz`** — tests that invent their own input and find what you did not think of.
6. **Metrics and tracing** — `expvar`, Prometheus, OpenTelemetry: how many requests, which are slow, where the waiting happened.
7. **A build at the push of a button** — GitHub Actions: tests on every commit and a deploy in one command.

Worth reading in this order: [go.dev/doc/effective_go](https://go.dev/doc/effective_go), then [Go by Example](https://gobyexample.com/), then other people's code — starting with the standard library, which is readable.

## What to build in the blog

The blog from this course is not a finished thing but a foundation. What asks to be next, in order of difficulty:

- **a picture for the social networks** — `og:image`, so that a link looks like a link rather than a line of text;
- **drafts** — an article not visible to everyone at once; one field and one condition in a query;
- **comments** — a form, a table, moderation, and CSRF once more;
- **several authors** — the permissions are already there; the author's page is not;
- **email** — confirming an address and resetting a password.

## The map of the lesson

![The map of the lesson: a type parameter, its condition and the road on](/static/course/go/map-next-en.svg)

## What you can do now

Fifty lessons ago we wrote a program of four lines. Behind you now is a blog that keeps its articles in a database with migrations, searches them, lets people in by password and holds their sessions, tells permissions apart, accepts pictures, serves JSON, is covered by tests run with the race detector, works over HTTPS on a domain of its own, stops tidily and can take a backup.

The thirty-four steps in the repository are that same blog at every lesson, and any of them starts with one command. Go back to them when something slips: the code in them is the code in the text.

## Say it in your own words

Without looking, answer aloud or on paper. The answers are at the end of the lesson.

1. Why does a generic need a condition on its type when the type is filled in at the call anyway?
2. Why is `any` with a type switch slower than a type parameter?
3. When is a generic not needed even though it is possible?

## The exercise

**Required.** Find two functions in your blog with nearly the same body and bring them together into one with a type parameter. If there are none, write `Filter[T any](xs []T, keep func(T) bool) []T` and use it on the list of articles.

**If you want more.**

- Measure your version against an `any` one: `go test -bench . -benchmem`.
- Rewrite the view counter from the goroutines lesson so that it counts anything at all — with the mutex inside.
- Replace your own small helpers with `slices` and `maps` and see how many lines went away.

## Where this goes in the blog

Nowhere — and that is the right ending. The core that was promised is done: the blog works, it is published, and it can be shown. Generics are not needed in it, and adding them for a lesson's sake would be dishonest.

There is still a list of improvements — every working system has one. From what the lessons named: rate-limiting the sign-in, clearing expired sessions and a "sign out everywhere" button, paging the list of articles, deleting pictures along with their article, Kazakh morphology in search, and the database's readiness in the health check. Those are not debts of the course but the ordinary operational queue.

From here the blog is changed by you, not by the course.

## The answers

1. Because without a condition the compiler does not know what may be done with `T`: addition, comparison, calling a method. The condition is both the permission the body needs and the check on whoever calls it.
2. Because the `any` version examines the type of every element while the program runs. A generic knows the type at build time — measured: 534 nanoseconds against 1801.
3. When the function is needed for one type; when different types need different behaviour, where an interface belongs; and when readability loses more to the type parameter than reuse gains.

## Sources

- [Go: an introduction to generics](https://go.dev/blog/intro-generics)
- [Go: package slices](https://pkg.go.dev/slices)
- [Go: Effective Go](https://go.dev/doc/effective_go)
