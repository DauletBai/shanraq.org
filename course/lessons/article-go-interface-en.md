# Interfaces in Go: a contract about what a thing can do

_Лид (summary):_ **The thirteenth lesson of the Go course. An interface is a list of abilities: one function works with any store as long as the store can do what the contract says. Why implementing is never declared but simply happens, how to have it checked at build time, and why a stand-in store is what makes tests possible without a database.**

## Why this matters

You are writing a function that hands back a list of articles. Today they live in memory, in a month they will live in a database, and in tests they should not exist at all. Write the function three times?

No. The function has no business knowing **what** the store is. It is enough for it to know **what the store can do**.

You have used this already. In your very first server the handler received `w http.ResponseWriter` — and `w` was not a struct, it was an interface: "I can accept bytes written to me". Who was on the other end — a browser, a file, a test — your code did not care.

> **Picture it.** A job advert. It does not say who you are or where you come from. It says what you must be able to do: take an article in and hand a list back. A librarian or a robot, it makes no difference, as long as they can.

## The whole thing at once

A new folder, `go mod init sabaq12`, `main.go`:

```go
package main

import "fmt"

type Article struct {
	Title string
	Words int
}

func (a Article) String() string {
	return fmt.Sprintf("%s (%d words)", a.Title, a.Words)
}

type Store interface {
	Add(a Article)
	All() []Article
}

type MemoryStore struct {
	items []Article
}

func (s *MemoryStore) Add(a Article)  { s.items = append(s.items, a) }
func (s *MemoryStore) All() []Article { return s.items }

type LastTwoStore struct {
	items []Article
}

func (s *LastTwoStore) Add(a Article) {
	s.items = append(s.items, a)
	if len(s.items) > 2 {
		s.items = s.items[len(s.items)-2:]
	}
}
func (s *LastTwoStore) All() []Article { return s.items }

func report(name string, s Store) {
	s.Add(Article{Title: "Hello, world", Words: 150})
	s.Add(Article{Title: "About the steppe", Words: 400})
	s.Add(Article{Title: "What a shanyraq is", Words: 1000})

	fmt.Println(name, "—", len(s.All()), "articles")
	for _, a := range s.All() {
		fmt.Println("  ", a)
	}
}

func main() {
	report("memory", &MemoryStore{})
	report("the last two", &LastTwoStore{})
}
```

Before you run it, and without looking below, say out loud what it will print. Then run it and compare.

```
memory — 3 articles
   Hello, world (150 words)
   About the steppe (400 words)
   What a shanyraq is (1000 words)
the last two — 2 articles
   About the steppe (400 words)
   What a shanyraq is (1000 words)
```

`report` was written once and did not change. The stores are different.

## Taking it apart

### An interface is a list of abilities, not of fields

```go
type Store interface {
	Add(a Article)
	All() []Article
}
```

Inside a struct you list fields — what it has. Inside an interface you list methods — what it can do. An interface has no fields at all.

`Store` is a type like `Article`. You can declare a variable of it, pass it to a function, put it in a slice. What sits in the variable is not a struct but "something that can do `Add` and `All`".

### Implementing is not declared — it simply happens

Find the word in the program where `MemoryStore` declares that it implements `Store`. There is no such word.

In Go it is enough to **have the methods**. Write `Add` and `All` with the same names and signatures and the type already qualifies, and nobody needs to be told.

> **Picture it.** Nobody brings a certificate. In other languages you have to write the paper: "I implement this contract." In Go you are taken on because you can do the work, not because you produced a document.

There is a useful consequence: an interface can be declared **after** the types, and even in someone else's package. You describe what you need from other people's code, and it fits without knowing anything about you.

### One function, any store

```go
func report(name string, s Store) {
```

`report` knows nothing about `MemoryStore` or about `LastTwoStore`. It knows the contract. When a third store appears, the function will not have to change.

Notice the call: `&MemoryStore{}`, with the address sign. The methods are declared with a pointer receiver — `func (s *MemoryStore) Add(...)` — because `Add` changes the contents. So it is `*MemoryStore` that satisfies the contract, not `MemoryStore`. This carries straight on from the previous lesson, and here it becomes compulsory for the first time.

### The check at build time

While the store is not passed anywhere yet, a typo in a method name is easy to miss. There is a line that has the contract checked at compile time:

```go
var _ Store = (*MemoryStore)(nil)
```

It reads: "make sure `*MemoryStore` will do as a `Store`". No value is wanted — hence the `_`. If a method is missing, the build fails and says so plainly:

```
*BrokenStore does not implement Store (missing method All)
```

It is cheap insurance, and it goes next to the type.

### `fmt.Stringer` — an interface you were already using

```go
func (a Article) String() string {
	return fmt.Sprintf("%s (%d words)", a.Title, a.Words)
}
```

One `String() string` method, and `fmt.Println` starts printing the article your way instead of in braces. That is `fmt.Stringer` from the standard library: the `fmt` package asks the value whether it can do `String`, and calls it if it can.

You imported nothing and registered nowhere. You wrote a method with the right name.

### What sits in a variable of an interface type

Two things: the value itself and its real type. That is why `fmt.Printf("%T", s)` shows `*main.MemoryStore` rather than `Store`.

The zero value of an interface is `nil`: no value and no type. Calling a method on that gives the crash you met in the previous lesson.

### Do not build the interface too early

A rule that saves months: **accept interfaces, return structs.** A function declares what it needs; a constructor hands back a concrete type.

And do not introduce an interface while there is only one implementation. A contract is worth having where there is a choice; where there is none it only adds a layer someone has to read.

> **Picture it.** A tailor's dummy. Not a person, but the shoulders are right — and the jacket can be checked on it. A stand-in store in a test is exactly that: not a database, but it keeps the contract.

## The lesson map

![The lesson map: one function, any store, the contract is a list of abilities](/static/course/go/map-interface-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Where in the program does `MemoryStore` declare that it implements `Store`?
2. Why does `main` pass `&MemoryStore{}` rather than `MemoryStore{}`?
3. What is the line `var _ Store = (*MemoryStore)(nil)` for, when the program works without it?

## Exercise

**Required.** Write a function `total(s Store) int` giving the sum of words across every article in the store. Then, in `main_test.go`, declare a stand-in store of your own called `fakeStore` with the same two methods, put two articles of 150 and 400 words into it, and check with a test that `total` returns 550. The real store may not be used in the test.

**Optional.**

- Add a third store, `ReverseStore`, that hands the articles back in reverse order. `report` may not be changed.
- Remove the `All` method from `MemoryStore` and read the compiler's message in full.
- Give `MemoryStore` a `String()` method and print the store with `fmt.Println`.

## Where this goes in your blog

The interface is where the blog stops being an exercise. Later in the course `Store` becomes real: memory first, then Postgres, while the tests keep running on a stand-in with no database and no network. The functions you have already written will not change by a line — they only ever knew the contract.

## Answers

1. Nowhere. In Go implementing is not declared: a type satisfies an interface if it has all the required methods with the same signatures. That is why you can describe an interface even for someone else's type that knows nothing about you.
2. Because `Add` and `All` are declared with a pointer receiver, so it is `*MemoryStore` that satisfies `Store`. `MemoryStore` formally does not have those methods, and the compiler refuses.
3. It has the contract checked at build time rather than at first use. While the store is not passed anywhere, a typo in a method name goes unnoticed; that line catches it at once and costs nothing at run time.

## Sources

- [Interfaces — A Tour of Go](https://go.dev/tour/methods/9)
- [Effective Go: interfaces](https://go.dev/doc/effective_go#interfaces)
- [Go spec: Interface types](https://go.dev/ref/spec#Interface_types)
- [fmt.Stringer — documentation](https://pkg.go.dev/fmt#Stringer)
