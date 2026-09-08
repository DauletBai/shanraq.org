# Tests: knowing what an edit broke

_Лид (summary):_ **The forty-first lesson of the Go course. The table promised back in lesson 10: the cases go into a list, t.Run gives each of them a name, and a failure says which case broke. Plus checking a handler without starting a server, and an honest word about what -cover really shows.**

## Why this matters

The test in lesson 10 checked one function with one case. That was enough right up to the day there were five cases.

Copying the `if` five times is possible, and people do it — and then forget to change the expected value in the fourth copy, and the test spends years checking the wrong thing. Worse: when such a test fails, the report shows the name of the function but not **which** case did not match.

Today comes the technique almost every test in Go is written with. It is called table-driven, and in lesson 10 I promised to show it once you had structs. You have them.

## The whole thing at once

A new folder, `go mod init sabaq38`. The first file is `main.go`, three familiar things:

```go
package main

import (
	"fmt"
	"net/http"
	"strings"
)

// parseTags turns what was typed in a form into a list with no repeats.
func parseTags(s string) []string {
	var out []string
	seen := map[string]bool{}
	for _, part := range strings.Split(s, ",") {
		t := strings.ToLower(strings.TrimSpace(part))
		if t == "" || seen[t] {
			continue
		}
		seen[t] = true
		out = append(out, t)
	}
	return out
}

// readingTime rounds up: a minute left unfinished still counts.
func readingTime(words int) int {
	return (words + 199) / 200
}

func articles() http.Handler {
	titles := map[string]string{"steppe": "The steppe", "hello": "Hello"}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.URL.Path, "/read/")
		title, ok := titles[slug]
		if !ok {
			http.Error(w, "no such article", http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, "<h1>%s</h1>", title)
	})
}

func main() {
	fmt.Println(parseTags("  Steppe , steppe,, Blog "))
	fmt.Println(readingTime(400), readingTime(401))
}
```

The second is `main_test.go`, next to it:

```go
package main

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

func TestParseTags(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string
	}{
		{"empty", "", nil},
		{"one", "steppe", []string{"steppe"}},
		{"spaces around", "  steppe  ", []string{"steppe"}},
		{"mixed case", "Steppe, steppe", []string{"steppe"}},
		{"empty between commas", "steppe,,blog", []string{"steppe", "blog"}},
		{"order is kept", "blog, steppe", []string{"blog", "steppe"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := parseTags(c.in)
			if !slices.Equal(got, c.want) {
				t.Errorf("parseTags(%q) = %v, wanted %v", c.in, got, c.want)
			}
		})
	}
}

func TestReadingTime(t *testing.T) {
	cases := []struct {
		name  string
		words int
		want  int
	}{
		{"zero", 0, 0},
		{"one word", 1, 1},
		{"exactly two hundred", 200, 1},
		{"two hundred and one", 201, 2},
		{"exactly four hundred", 400, 2},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := readingTime(c.words); got != c.want {
				t.Errorf("readingTime(%d) = %d, wanted %d", c.words, got, c.want)
			}
		})
	}
}

func TestArticles(t *testing.T) {
	cases := []struct {
		name string
		path string
		code int
		body string
	}{
		{"the article is there", "/read/steppe", http.StatusOK, "The steppe"},
		{"no article", "/read/joq", http.StatusNotFound, "no such article"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// No server is needed: a handler is an ordinary function, and the
			// recorder stands in for a browser.
			req := httptest.NewRequest(http.MethodGet, c.path, nil)
			rec := httptest.NewRecorder()

			articles().ServeHTTP(rec, req)

			if rec.Code != c.code {
				t.Errorf("code %d, wanted %d", rec.Code, c.code)
			}
			if !strings.Contains(rec.Body.String(), c.body) {
				t.Errorf("the answer has no %q, answer: %q", c.body, rec.Body.String())
			}
		})
	}
}
```

Running it:

```
$ go test -v -run TestParseTags
=== RUN   TestParseTags
=== RUN   TestParseTags/empty
=== RUN   TestParseTags/one
=== RUN   TestParseTags/spaces_around
=== RUN   TestParseTags/mixed_case
=== RUN   TestParseTags/empty_between_commas
=== RUN   TestParseTags/order_is_kept
--- PASS: TestParseTags 
    --- PASS: TestParseTags/empty 
    --- PASS: TestParseTags/one 
    --- PASS: TestParseTags/spaces_around 
    --- PASS: TestParseTags/mixed_case 
    --- PASS: TestParseTags/empty_between_commas 
    --- PASS: TestParseTags/order_is_kept 
PASS
ok  	sabaq38
```

The run time is removed from the excerpts: everybody's differs.

## Taking it apart

### A table instead of copies

The cases sit in a slice of structs, one line per case:

```go
cases := []struct {
	name string
	in   string
	want []string
}{
	{"empty", "", nil},
	{"one", "steppe", []string{"steppe"}},
	{"spaces around", "  steppe  ", []string{"steppe"}},
}
```

The struct here has no name: the type is needed in exactly one place, and inventing a name for it would be pointless. Three fields answer three questions: what the case is called, what goes in, what is expected.

Adding a seventh case means writing one line. Not copying ten lines of checking — writing one. That is why the table is liked: it makes adding a case cheap, and so there come to be more cases.

> **Picture it.** A shopping list instead of ten notes in ten pockets. The list shows everything at once, adding a line is easy, and nothing goes missing on the way.

### `t.Run`: every case has a name

`t.Run(c.name, func(t *testing.T) { ... })` starts a **subtest** — a separate test inside a test, with a name and a result of its own.

The output shows the names became part of the report: `TestParseTags/mixed_case`. Go replaces spaces with underscores so a name can be passed on the command line.

Names are written for people: "spaces around", not "case 3". In six months you are the one reading the report, and "case 3" will tell you nothing.

A subtest has a `t` of its own. That matters more than it looks: `t.Errorf` in a subtest fails **that** case while the others keep running, so one run tells you about every breakage rather than the first.

### What is visible when something breaks

Let us break the rounding in `readingTime` by removing the `+ 199`:

```
$ go test
--- FAIL: TestReadingTime 
    --- FAIL: TestReadingTime/one_word 
        main_test.go:51: readingTime(1) = 0, wanted 1
    --- FAIL: TestReadingTime/two_hundred_and_one 
        main_test.go:51: readingTime(201) = 1, wanted 2
FAIL
exit status 1
FAIL	sabaq38
```

This is what it was all for. The report says: exactly two cases out of five broke, here are their names, here is what came out and what was expected. Not "the test failed" but "one word and two hundred and one words are counted wrong" — and from those two names it is plain that rounding up went missing.

Notice the shape of the message: `readingTime(1) = 0, wanted 1`. It holds the input, the result and the expectation. A message like "did not match" holds none of that, and it will have to be debugged instead of the code.

### One case by name

Since a case has a name, it can be run on its own:

```
$ go test -run 'TestReadingTime/two_hundred_and_one' -v
=== RUN   TestReadingTime
=== RUN   TestReadingTime/two_hundred_and_one
--- PASS: TestReadingTime 
    --- PASS: TestReadingTime/two_hundred_and_one 
PASS
ok  	sabaq38
```

The slash separates the name of the test from the name of the subtest. When there are two hundred tests and one thing to fix, this saves minutes on every edit.

### A handler is checked without a server

`TestArticles` checks a web handler, and no server is started for it:

```go
req := httptest.NewRequest(http.MethodGet, c.path, nil)
rec := httptest.NewRecorder()

articles().ServeHTTP(rec, req)
```

A handler in Go is an ordinary function of two arguments. `httptest.NewRequest` makes a request and `httptest.NewRecorder` is the sheet of paper a handler writes its answer onto instead of the network. After the call `rec.Code` holds the status and `rec.Body` the body.

This is fast (no ports, no waiting) and reliable (nothing can be busy). A real server is started in tests with `httptest.NewServer` when the whole chain is being checked — wrappers, cookies and redirects included.

### `-cover`: how many lines ran under the tests

```
$ go test -cover
coverage: 90.0% of statements
ok  	sabaq38
```

Ninety per cent is the share of **lines** that ran at least once during the tests. A useful number: it shows where the tests never looked.

And a poor goal. A line can be run without anything being checked: call a function and ignore the result — coverage grows, quality does not. A hundred per cent is easy to reach and does not mean the program works.

Use it as a map of blank spots, not as a grade.

### What to test and what not to

Test what is easy to get wrong and expensive to miss: boundaries (zero, one, exactly two hundred), parsing input from outside, permissions, money, anything to do with dates.

Do not test what the compiler checks for you, and do not test other people's libraries: `database/sql` was tested without you. And do not test markup line by line — it changes every week, and the test will fail on every comma.

Separately: a test must not depend on the order tests run in, on the clock, or on the internet. Such a test fails once a week for no visible reason, and in the end it gets switched off — and its usefulness with it.

## The lesson map

![Lesson map: a table of cases, names and a report](/static/course/go/map-tests-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. What does `t.Run` give you compared with five `if`s in a row?
2. Why can a handler be checked without starting a server?
3. Why is a hundred per cent coverage a poor goal?

## Exercise

**Required.** Write table-driven tests for the blog: for `ParseTags`, for `ReadingTime` and for the `/read/{slug}` handler — with the cases "the article is there" and "no article". Check that `go test ./...` passes, then break one function on purpose and look at the report.

All of it is done in [step-26](https://github.com/DauletBai/shanraq.org/tree/main/course/go-blog/step-26) — compare against it once you have written your own.

**Optional.**

- Add a store test: `t.TempDir()` gives a folder that removes itself — open a database in it, file an article and read it back.
- Find a function the tests never touch, with `go test -coverprofile=c.out` and `go tool cover -html=c.out`.
- Add a case that does not pass yet — a tag made only of spaces, say — and fix the code rather than the test.

## Where this goes in your blog

From this lesson on, an edit stops being frightening: `go test ./...` says in a second whether you broke what was working yesterday.

Debts. The blog has few tests, and they touch neither sessions nor permissions. There is no check for races — that comes in the lesson on goroutines, together with `-race`. And we have no test database that fills itself the same way before every test.

## Answers

1. A name and a result of its own for every case. A failure shows which case did not match while the rest keep running, and a case can be run on its own by name with `-run`.
2. Because a handler is an ordinary function of an `http.ResponseWriter` and an `*http.Request`. `httptest` hands it a request and a sheet of paper instead of the network; neither a port nor a server is needed.
3. Because coverage counts lines that ran, not assertions that were checked. A line can run with nothing checked — the number grows while breakages sail past. It is a map of blank spots, not a grade.

## Sources

- [Go: the testing package](https://pkg.go.dev/testing)
- [Go: the net/http/httptest package](https://pkg.go.dev/net/http/httptest)
- [Go: how to write tests](https://go.dev/doc/tutorial/add-a-test)
