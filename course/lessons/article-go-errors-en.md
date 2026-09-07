# Errors in Go: not an exception, just an ordinary value

_Лид (summary):_ **The fourteenth lesson of the Go course. An error here does not interrupt the program; it comes back as the second value, and it is the first thing you check. error as an interface with one method, errors.New and fmt.Errorf, wrapping with %w and errors.Is, your own error as a value — and why panic is not for any of this.**

## Why this matters

A reader opened `/read/kokek` and there is no such article. That is not the program breaking — it is an ordinary turn of events that has to be answered with a 404.

Go has no `try` and no `catch` for this. A function simply returns **two** values: the result and an error. You look at the second and decide what to do.

> **Picture it.** Not an alarm but a receipt. An alarm puts the whole building on its feet. A receipt is a slip of paper handed to you: read it and decide for yourself whether to go on or turn back.

## The whole thing at once

A new folder, `go mod init sabaq13`, `main.go`:

```go
package main

import (
	"errors"
	"fmt"
)

type Article struct {
	Title string
	Words int
}

var ErrNotFound = errors.New("article not found")

type Store struct {
	items map[string]Article
}

func (s *Store) Get(slug string) (Article, error) {
	a, ok := s.items[slug]
	if !ok {
		return Article{}, fmt.Errorf("get %q: %w", slug, ErrNotFound)
	}
	return a, nil
}

func show(s *Store, slug string) {
	a, err := s.Get(slug)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			fmt.Println("404:", err)
			return
		}
		fmt.Println("500:", err)
		return
	}
	fmt.Println(a.Title, "·", a.Words, "words")
}

func main() {
	s := &Store{items: map[string]Article{
		"dala": {Title: "About the steppe", Words: 400},
	}}

	show(s, "dala")
	show(s, "kokek")

	_, err := s.Get("kokek")
	fmt.Println(errors.Is(err, ErrNotFound))
	fmt.Println(errors.Unwrap(err) == ErrNotFound)
}
```

Before you run it, and without looking below, say out loud what it will print. Then run it and compare.

```
About the steppe · 400 words
404: get "kokek": article not found
true
true
```

## Taking it apart

### `error` is an interface with one method

Yesterday's lesson made today's almost obvious. `error` is not a special construction of the language but an ordinary interface with a single method:

```go
type error interface {
	Error() string
}
```

Anything that can do `Error() string` will serve as an error. No hierarchy of classes, no "exception types" — a type with one method.

### A function returns two values

```go
func (s *Store) Get(slug string) (Article, error)
```

The result on the left, the error on the right — Go keeps that order everywhere. When there is no error the second value is `nil` and the first one means something. When there is an error the first value is usually the zero one and must not be relied on.

Hence a rule for reading other people's code: when you see a call with two results, look for the check on the second.

### `if err != nil`, and why it is not ugly

```go
a, err := s.Get(slug)
if err != nil {
	return
}
```

These three lines occur in Go more often than any others, and it is a habit to complain about them. Here is the case for the defence: they make the error's path **visible**. In a language with exceptions you cannot tell from the code which line might cut the function short; here you can tell at a glance.

The price is verbosity. The return is that you are never surprised by where something came from.

### `errors.New` and `fmt.Errorf`

`errors.New("text")` makes an error with a message. `fmt.Errorf` does the same but can substitute values:

```go
fmt.Errorf("get %q: %w", slug, ErrNotFound)
```

There are two verbs here. `%q` puts the article's address in quotes. `%w` is the special one: it **wraps** another error inside the new one and keeps it whole.

That is why it printed `get "kokek": article not found` — first where, then what.

> **Picture it.** An envelope inside an envelope. Everyone whose hands the letter passed through puts it in an envelope of their own and writes "received from so-and-so". Open them all and you see both the route and the original letter.

The convention is simple: an error message starts in lower case and ends without a full stop, because it will almost always end up in the middle of somebody else's line.

### Your own error as a value, and `errors.Is`

```go
var ErrNotFound = errors.New("article not found")
```

An error is declared as a package-level variable named starting with `Err`. Then it can be **compared** instead of having its text picked apart:

```go
if errors.Is(err, ErrNotFound) {
```

`errors.Is` opens every envelope and checks whether that particular error is inside. That is why our `%w`-wrapped error was recognised even though its text is now different.

Never compare errors by their text: `err.Error() == "article not found"` breaks on the first rewording, and breaks silently.

### What to do with an error

Three respectable options, and no more.

1. **Handle it.** If you know what to do, do it: a 404 instead of the article, a default value, a retry.
2. **Wrap it and pass it up.** If you do not know, add your envelope with `%w` and hand it to someone who does.
3. **Log it.** Once, at the very top, not in every layer.

And one disreputable option: **swallow it**. `a, _ := s.Get(slug)` is legal Go, the compiler stays quiet, and the program carries on with an empty article. An underscore where the error goes is acceptable only where you can say out loud why it does not matter.

### `panic` is not for errors

`panic` stops the program. It is for situations with no sensible way out: the configuration would not load, a required template failed to parse at start-up.

A missing article is not one of those. There will be thousands of them a day.

> **Picture it.** A fire alarm. Every building has one, and rightly so. But nobody pulls it because a visitor could not find the right office.

## The lesson map

![The lesson map: a function has two ways out, a value and an error](/static/course/go/map-errors-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. What is `error` as far as the language is concerned, and why is that no longer surprising after the interfaces lesson?
2. How does `%w` differ from `%v` in `fmt.Errorf`, and what depends on the difference?
3. Why is `err.Error() == "article not found"` a poor way to identify an error while `errors.Is` is a good one?

## Exercise

**Required.** Give the type `Store` a method `Add(slug string, a Article) error`. It must return the error `ErrEmptyTitle` when the title is empty and `nil` when the article is accepted. Write a test: add a normal article and confirm there is no error; add one with an empty title and check `errors.Is(err, ErrEmptyTitle)`; confirm that the second never reached the store.

**Optional.**

- Wrap `ErrEmptyTitle` with `%w`, adding the article's address, and check that `errors.Is` still recognises it.
- Print `errors.Unwrap(err)` and see what is inside the envelope.
- Replace the `errors.Is` check in `show` with a text comparison, then reword `ErrNotFound` and watch what breaks.

## Where this goes in your blog

This is how the site you are reading works too. The store returns `ErrNotFound`, the handler calls `errors.Is` and answers 404, and everything else becomes a 500 and one line in the log. No exceptions, no catching ten levels up — a second value and one check.

## Answers

1. `error` is an ordinary interface with the single method `Error() string`. After the previous lesson that is no surprise: any type with that method will serve as an error, and nothing has to be declared.
2. `%v` inserts the **text** of the error and the link to the original is lost. `%w` wraps **the error itself** inside the new one, so `errors.Is` and `errors.Unwrap` can reach it later. Whether the caller can tell "not found" from "everything is broken" depends on it.
3. The text is a message for a human being, and it gets changed without a second thought: translated, given one more detail — and the comparison stops working, silently at that. `errors.Is` compares values rather than letters and survives any rewording.

## Sources

- [Errors — A Tour of Go](https://go.dev/tour/methods/19)
- [Working with Errors in Go 1.13 — the go.dev blog](https://go.dev/blog/go1.13-errors)
- [Package errors — documentation](https://pkg.go.dev/errors)
- [Effective Go: errors](https://go.dev/doc/effective_go#errors)
