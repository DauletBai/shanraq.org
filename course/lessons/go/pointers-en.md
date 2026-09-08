# Pointers and methods in Go: why the change never stuck

_Лид (summary):_ **The twelfth lesson of the Go course. The answer to last lesson's trap: a pointer is the address of a value, not a copy of it. How & and * actually work, how a method differs from a plain function, when the receiver takes a pointer and when a value, and why for i := range fixes the very thing that for _, a := range quietly breaks.**

## Why this matters

In the previous lesson a change made inside `range` vanished without trace, and the compiler said nothing. We said then that the real answer is pointers. Here it is.

A pointer is wanted exactly where you need to **change someone else's value rather than a copy of it**: publish an article, add words to it, fix its title.

## The whole thing at once

A new folder, `go mod init sabaq11`, `main.go`:

```go
package main

import "fmt"

type Article struct {
	Title     string
	Words     int
	Published bool
}

func (a Article) ReadingTime() int {
	return (a.Words + 199) / 200
}

func (a *Article) Publish() {
	a.Published = true
}

func rename(a *Article, title string) {
	a.Title = title
}

func main() {
	first := Article{Title: "About the steppe", Words: 400}

	fmt.Println(first.Title, "·", first.ReadingTime(), "min ·", first.Published)

	first.Publish()
	rename(&first, "About the steppe: the full version")
	fmt.Println(first.Title, "·", first.Published)

	blog := []Article{
		{Title: "Hello, world", Words: 150},
		{Title: "What a shanyraq is", Words: 1000},
	}

	for i := range blog {
		blog[i].Publish()
	}
	for _, a := range blog {
		fmt.Println(a.Title, "·", a.Published)
	}

	p := &first
	fmt.Printf("%T · %s\n", p, p.Title)
}
```

Before you run it, and without looking below, say out loud what it will print. Then run it and compare.

```
About the steppe · 2 min · false
About the steppe: the full version · true
Hello, world · true
What a shanyraq is · true
*main.Article · About the steppe: the full version
```

## Taking it apart

### A pointer is an address, not the thing itself

```go
p := &first        // take the address
fmt.Println(*p)    // follow the address back to the value
```

`&` in front of a value gives its address. `*` in front of a pointer follows the address back to the value. And the same `*` in a type declaration means "the address of": `*Article` is "the address of an article", not an article.

`%T` confirmed it: the type of `p` is `*main.Article`, star and all.

> **Picture it.** A key and a photograph of a flat. Given the photograph, draw on it all you like — nothing in the flat moves. Given the key, you can walk in and rearrange the furniture for real. A pointer is the key.

A copy is cheaper and safer, which is why Go copies by default. The key is handed over when something has to change.

### A function with a pointer changes the original

```go
func rename(a *Article, title string) {
	a.Title = title
}

rename(&first, "About the steppe: the full version")
```

The parameter is declared as `*Article`, and at the call we pass `&first`, the address. The title changed on the outside, which is what the second line of output shows.

Compare that with the previous lesson: there the function took an `Article`, got a copy, and its edits never arrived anywhere.

### A method is a function written onto a type

```go
func (a Article) ReadingTime() int {
```

The brackets before the name are the **receiver**. They say: this function belongs to the type `Article`, and it is called through a dot — `first.ReadingTime()`, not `ReadingTime(first)`.

> **Picture it.** Instructions printed on the thing itself. A washing machine has its programmes on the panel, not in a separate book you have to find first. A method is the same: the ability sits on the type.

### The receiver: a value or a pointer

The whole rule of the lesson fits in two lines.

```go
func (a Article) ReadingTime() int   // only reads — take a value
func (a *Article) Publish()          // changes it — take a pointer
```

`ReadingTime` counts and touches nothing: a copy is enough. `Publish` writes to a field, and without the pointer the write would go into a copy while `Published` stayed `false` — exactly the trap this lesson opened with.

One note for later: if one method of a type takes a pointer, the others usually take one too. Mixing is allowed, but anyone reading your code then has to remember which kind each method is.

### Go puts the stars in for you

`first` is an ordinary value, not a pointer. `first.Publish()` works all the same. And the other way round: `p` is a pointer, and `p.Title` is written with no star at all.

Go inserts `&` and `*` itself where it is unambiguous. That is why stars barely appear in everyday code — only in type declarations.

### How last lesson's trap gets fixed

```go
for i := range blog {
	blog[i].Publish()
}
```

Here `range` hands over only the index. `blog[i]` is the element of the slice itself rather than a copy of it, and Go takes its address automatically, because an element of a slice is addressable.

`for _, a := range blog { a.Publish() }` still does nothing at all: `a` is a copy, a copy has an address of its own, and it is the copy that gets published. The compiler says nothing here either.

It is easy to remember: **if you are changing it, go by index.**

### The nil pointer

A pointer may point nowhere. Its zero value is `nil`:

```go
var p *Article
fmt.Println(p.Title)   // panic: invalid memory address or nil pointer dereference
```

> **Picture it.** An address on a piece of paper with nothing built there. You arrive and find an empty lot.

This is neither rare nor a disaster: the message says plainly what happened. In the lesson on errors we learn to hand such situations back like human beings.

## The lesson map

![The lesson map: a copy is a dead end, an address leads back to the value itself](/static/course/go/map-pointers-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. How does `func (a Article) ReadingTime()` differ from `func (a *Article) Publish()`, and how do you choose between them?
2. `first` is not a pointer, yet `first.Publish()` works. Why?
3. `for _, a := range blog { a.Publish() }` publishes nothing while `for i := range blog { blog[i].Publish() }` publishes. What is the difference?

## Exercise

**Required.** Give the type `Article` a method `Shorten(limit int)` that trims the title to `limit` letters and puts an ellipsis at the end. Letters, not bytes, as in the runes lesson. Apply it to every article in the slice and print the result. Think about what the receiver has to be before you write it.

**Optional.**

- Call `Shorten` inside a `for _, a := range blog` loop and confirm that nothing changed.
- Print `%T` for `first` and for `&first`.
- Write a function `publishAll(blog []Article)` and check whether it works. If it does not, work out why, and how one character in the declaration fixes it.

## Where this goes in your blog

From here on almost everything comes with a pointer: a handler receives a `*http.Request` — you met that star in your very first server — the database hands back `*sql.Rows`, and a template takes a `*template.Template`. Now you know what it means there: you are being given the key, not the photograph.

## Answers

1. The first only reads the fields and gets a copy, which is enough. The second writes to a field, and without a pointer the write would land in a copy that disappears immediately. The rule: changes it — pointer; only reads it — value.
2. Because Go inserts the `&` itself when the value is addressable. `first` is a variable, it has an address, so the call becomes `(&first).Publish()`. Writing stars and ampersands by hand is rarely necessary.
3. In the first case `a` is a copy of the element; the method publishes the copy, and the copy disappears on the next turn of the loop. In the second, `blog[i]` is the element of the slice itself, it has an address, and a pointer method changes what is actually in the slice.

## Sources

- [Pointers — A Tour of Go](https://go.dev/tour/moretypes/1)
- [Methods and pointers — A Tour of Go](https://go.dev/tour/methods/4)
- [Go spec: Method declarations](https://go.dev/ref/spec#Method_declarations)
