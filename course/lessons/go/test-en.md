# Your first test in Go: go test and the _test.go file

_Лид (summary):_ **The tenth lesson of the Go course. A check that stays written down: a _test.go file beside the code, a TestXxx function and the go test command. How to read a failure — the test name, the file, the line, "got" and "wanted" — how t.Errorf differs from t.Fatalf, and why you test the edges of a list rather than its middle.**

## Why this matters

You have already written `count` and `longer`. How do you know they are right? You ran them, looked at the output and nodded. A check like that lives for one minute: fix the function next week and you will have to check it again by hand, and you will not.

A test is that same check, only written down. A machine performs it, and performs it every time.

> **Picture it.** A pilot's checklist before take-off. It is not there because pilots forget, but because memory does not grow along with the number of things to remember. A list does.

## The whole thing at once

A new folder, `go mod init sabaq09`. This time there are two files. The first is `main.go`, the function you know from the previous lesson:

```go
package main

import "fmt"

func count(tags []string) map[string]int {
	m := make(map[string]int)
	for _, t := range tags {
		m[t]++
	}
	return m
}

func main() {
	fmt.Println(count([]string{"go", "web", "go"}))
}
```

The second is `main_test.go`, beside it, in the same folder:

```go
package main

import "testing"

func TestCount(t *testing.T) {
	got := count([]string{"go", "web", "go"})

	if got["go"] != 2 {
		t.Errorf("go: got %d, wanted 2", got["go"])
	}
	if len(got) != 2 {
		t.Errorf("tags: got %d, wanted 2", len(got))
	}
}
```

Before you run it, and without looking below, say out loud what `go test` will print. Then run it and compare.

```
PASS
ok  	sabaq09
```

The run time is removed from the excerpts: everybody's differs.

The last number is how long it took; yours will differ.

## Taking it apart

### A test is a file beside the code

The name must end in `_test.go`. That is not a tidiness convention but a rule of the language: Go builds such files only under `go test` and **leaves them out of the finished program**. You can write as many tests, and as thorough ones, as you like — they do not affect the size of what you hand to people.

> **Picture it.** Scaffolding. While the house is going up people walk on it every day; it is not part of the finished house, and nobody is surprised to find it gone.

The package is the same one — `package main` — so the test sees `count` directly, with nothing to import.

### What a test looks like

```go
func TestCount(t *testing.T) {
```

Three requirements, all of them compulsory. The name starts with `Test`. A capital letter follows — `TestCount`, not `Testcount`. And the only parameter is `t *testing.T`.

`t` is what the test reports a failure through. For now read the star as part of the spelling; the lesson on pointers takes apart what it means.

The function returns nothing. A test that said nothing counts as passed: silence here means everything matched.

### `go test` and `go test -v`

`go test` prints the verdict: `PASS` or `FAIL`. When there are many tests and you want to see each by name, add `-v`:

```
=== RUN   TestCount
--- PASS: TestCount 
PASS
ok  	sabaq09
```

### A failure reads as a report

Break it on purpose: change the `2` in the first comparison to `3` and fix the message to match. Run it:

```
--- FAIL: TestCount 
    main_test.go:9: go: got 2, wanted 3
    main_test.go:12: tags: got 2, wanted 3
FAIL
```

Everything you need is there: which test failed, in which file, on which line, what came back and what was expected. That is why the message carries both numbers. A message that says "did not work" says nothing and sends you into the code by hand.

A red `FAIL` is not a punishment but a report. Over a working life you will see it a hundred times more often than `PASS`.

### `t.Errorf` and `t.Fatalf`

`t.Errorf` records the failure and **carries on** — which is why you saw two lines at once. `t.Fatalf` records it and **stops** the test: the second check never runs.

The rule is simple. `Fatalf` when going on is pointless: the value you needed is not there and the next line will fail anyway. `Errorf` everywhere else — the more mismatches one run shows, the fewer runs you need.

### What to test first

Not the middle but the edges. The middle usually works for everyone.

For `count` the edges are these: an empty list, a single element, every element the same. For `longer` from the slices lesson — a list where no title qualifies, and a list where all of them do.

> **Picture it.** A road. Cars leave it at the verge, not in the middle of the lane. The verge is what you check.

There are usually three or four edges, and writing out a whole `if` for each will tire you quickly. That is not how it is done either: the cases go into one list, a loop walks it, and each case gets a name of its own in the report. This is called a table-driven test, and we will write one in the second lesson on testing — once you have the structs such a list is built from. For now, keep in mind that one test per function is a beginning, not a model.

### Why this is not doing the work twice

It looks as though you are writing the same thing twice. In fact you are writing it two different ways. In `count` you describe **how** to count. In the test you describe **what** should come out. A mistake in the method does not repeat itself in the description of the result, which is exactly how one catches the other.

And a test is the only check that never gets tired. You will fix the function a month from now, having forgotten half of it, and run `go test` — and it will remember on your behalf.

## The lesson map

![The lesson map: the test stands beside the code, calls it and gives one of two answers](/static/course/go/map-test-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why must a test file end in `_test.go`? What depends on it?
2. A test printed nothing and returned nothing. Did it pass?
3. One test has four checks and the second one fails. How many messages do you see with `t.Errorf`, and how many with `t.Fatalf`?

## Exercise

**Required.** Take the `longer` function from the slices lesson and write a test for it. Check two cases: the list `"Дала"`, `"Шаңырақ"`, `"Go"`, `"Көш"` with `n = 3` should give two titles, and an empty list with any `n` should give none. Run `go test` and get to `PASS`.

**Optional.**

- Break `longer` on purpose — put `>=` where `>` is, say — and make sure the test goes red. A check that cannot fail is checking nothing.
- Run `go test -v` and look at the difference.
- Replace `t.Errorf` with `t.Fatalf` in your test and break both checks at once. How many lines do you see now?

## Where this goes in your blog

Remember `readingTime` from the functions lesson, the one that works out how long an article takes to read? A mistake in it would be a quiet one: at exactly 400 words it is easy to get three minutes instead of two, and nobody notices until they sit down and count. That is precisely the kind of mistake a test catches: it checks 199, 200 and 201 words in a second, and does it every time.

From here tests travel with you through the rest of the course. The full lesson on testing — with a table of cases and coverage — comes later, once structs exist. Today's is enough to check your own exercises yourself.

> **Build your own blog.** In [step-4](https://github.com/DauletBai/shanraq.org/tree/main/course/go-blog/step-4) the same files gain a `main_test.go`: the checks sit on the edges — no words at all, exactly two hundred, two hundred and one. That is where formulas like reading time break.

## Answers

1. Because it is a rule of the language, not a habit. Files ending that way are built only by `go test` and never reach the finished program, so tests do not add to what you hand to people.
2. It passed. A test reports only failures, through `t.Errorf` or `t.Fatalf`. Silence means everything matched.
3. With `t.Errorf` you see a message for every check that failed, because the test carries on after each one. With `t.Fatalf` you see one, about the second check: the test stops there, and the third and fourth never run.

## Sources

- [Package testing — documentation](https://pkg.go.dev/testing)
- [go test — the command reference](https://pkg.go.dev/cmd/go#hdr-Test_packages)
- [Add a test — the official Go tutorial](https://go.dev/doc/tutorial/add-a-test)
