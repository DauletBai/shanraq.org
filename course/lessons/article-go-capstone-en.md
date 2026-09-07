# Putting it together: a console notebook out of what you know

_Лид (summary):_ **The eighteenth lesson of the Go course and the end of the module on the language. One finished program made of familiar parts: a struct, a slice, a map, methods with a pointer receiver, errors, a package of your own and tests. Nothing new but reading commands from the keyboard: putting what you learned into a working thing is a skill of its own.**

## Why this matters

For fifteen lessons you learned the parts one at a time. Today they come together into a program that works as a whole: it takes commands, keeps notes, searches, counts and deletes.

Putting things together is a skill of its own, and it is practised separately. Someone who knows every construction in a language may still be unable to assemble a working thing out of them: where each piece was clear on its own, together they raise the question of where each one goes.

> **Picture it.** A bicycle built from the parts already lying on the table. Not one new part was delivered. The skill is not in finding another one but in assembling these.

## The whole thing at once

A folder `sabaq-16`, and a `notes` folder inside it. Four files.

```
sabaq-16/
  go.mod
  main.go
  notes/
    note.go
    store.go
    store_test.go
```

`notes/note.go` — the note itself:

```go
package notes

import (
	"fmt"
	"unicode/utf8"
)

type Note struct {
	ID    int
	Text  string
	Words int
}

func (n Note) Letters() int {
	return utf8.RuneCountInString(n.Text)
}

func (n Note) String() string {
	return fmt.Sprintf("%d. %s (%d words · %d letters)", n.ID, n.Text, n.Words, n.Letters())
}
```

`notes/store.go` — the store, with everything it can do:

```go
package notes

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var (
	ErrEmpty    = errors.New("the note is empty")
	ErrNotFound = errors.New("note not found")
)

type Store struct {
	next  int
	items []Note
}

func NewStore() *Store {
	return &Store{next: 1}
}

func (s *Store) Add(text string) (Note, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return Note{}, ErrEmpty
	}
	n := Note{ID: s.next, Text: text, Words: len(strings.Fields(text))}
	s.items = append(s.items, n)
	s.next++
	return n, nil
}

func (s *Store) All() []Note { return s.items }

func (s *Store) Find(part string) []Note {
	var out []Note
	for _, n := range s.items {
		if strings.Contains(strings.ToLower(n.Text), strings.ToLower(part)) {
			out = append(out, n)
		}
	}
	return out
}

func (s *Store) Delete(id int) error {
	for i, n := range s.items {
		if n.ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("delete %d: %w", id, ErrNotFound)
}

// Top counts how often each word occurs and returns them in alphabetical
// order. The map counts, the slice puts the answer in an order a map does not
// promise -- both lessons in four lines.
func (s *Store) Top() []string {
	count := map[string]int{}
	for _, n := range s.items {
		for _, w := range strings.Fields(strings.ToLower(n.Text)) {
			count[w]++
		}
	}
	words := make([]string, 0, len(count))
	for w := range count {
		words = append(words, w)
	}
	slices.Sort(words)

	out := make([]string, 0, len(words))
	for _, w := range words {
		out = append(out, fmt.Sprintf("%s — %d", w, count[w]))
	}
	return out
}

func (s *Store) Len() int { return len(s.items) }
```

`main.go` — the conversation with a person:

```go
package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"sabaq16/notes"
)

func main() {
	store := notes.NewStore()
	in := bufio.NewScanner(os.Stdin)

	fmt.Println("notebook. commands: add, list, find, top, del, help, exit")

	for {
		fmt.Print("> ")
		if !in.Scan() {
			return
		}
		cmd, arg, _ := strings.Cut(strings.TrimSpace(in.Text()), " ")

		switch cmd {
		case "add":
			n, err := store.Add(arg)
			if err != nil {
				fmt.Println("error:", err)
				continue
			}
			fmt.Println("added:", n)

		case "list":
			if store.Len() == 0 {
				fmt.Println("empty")
				continue
			}
			for _, n := range store.All() {
				fmt.Println(n)
			}

		case "find":
			found := store.Find(arg)
			if len(found) == 0 {
				fmt.Println("nothing found")
				continue
			}
			for _, n := range found {
				fmt.Println(n)
			}

		case "del":
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("error: a number is needed")
				continue
			}
			if err := store.Delete(id); err != nil {
				if errors.Is(err, notes.ErrNotFound) {
					fmt.Println("error: no such number")
					continue
				}
				fmt.Println("error:", err)
				continue
			}
			fmt.Println("deleted:", id)

		case "top":
			for _, line := range store.Top() {
				fmt.Println(line)
			}

		case "help":
			fmt.Println("add text · list · find part · top · del number · exit")

		case "exit", "":
			return

		default:
			fmt.Printf("unknown command: %q, type help\n", cmd)
		}
	}
}
```

Build it and run it:

```
go mod init sabaq16
go run .
```

Before typing any commands, say out loud what the program will answer to `add a steppe note`. Then type it and compare:

```
notebook. commands: add, list, find, top, del, help, exit
> add a steppe note
added: 1. a steppe note (3 words · 13 letters)
> add what a shanyraq is
added: 2. what a shanyraq is (4 words · 18 letters)
> add
error: the note is empty
> top
a — 2
is — 1
note — 1
shanyraq — 1
steppe — 1
what — 1
> del 9
error: no such number
> exit
```

## Taking it apart

### Where each piece came from

There is not one construction in this program you are seeing for the first time. Check yourself against the list:

| What is in the code | From which lesson |
|---|---|
| `type Note struct` | structs |
| `[]Note` and `append` | slices |
| `map[string]int` in `Top` | maps |
| `func (n Note) String()` | methods |
| `func (s *Store) Add(...)` | pointers |
| `ErrEmpty`, `errors.Is` | errors |
| `package notes`, capital letters | packages |
| `utf8.RuneCountInString` | runes |
| `store_test.go` | your first test |

If something in the table is not familiar, go back to that lesson rather than this one.

### The one new thing: reading commands

```go
in := bufio.NewScanner(os.Stdin)
for {
	fmt.Print("> ")
	if !in.Scan() {
		return
	}
	...
}
```

`os.Stdin` is what a person types in the terminal. `bufio.Scanner` breaks that stream into lines and hands them over one at a time: `Scan` reads the next one and returns `false` when the lines run out.

> **Picture it.** A turnstile. A crowd arrives as one mass, but goes through one at a time, and each person can be dealt with separately.

`strings.Cut` cuts the line at the first space: the command on the left, everything else on the right. For `add a steppe note` that gives `add` and `Дала туралы жазба` — text with spaces in it does not fall apart.

### `switch` instead of a ladder of `if`s

```go
switch cmd {
case "add":
	...
case "list":
	...
default:
	...
}
```

When one value is compared against several options, `switch` reads better than a chain of `if — else if`. In Go it does not fall through into the next case: no `break` is needed.

`default` is the answer to an unknown command. A program that says nothing in reply to a typo looks broken.

### The `notes` package and its door out

What leaves the package is `Note`, `Store`, `NewStore`, the actions and two errors. The `next` and `items` fields stay inside, and nothing outside can reach them.

Check for yourself: `main.go` never touches the list directly. It **asks** — `store.Add`, `store.All`, `store.Delete`. Which is why the list can become a database tomorrow without `main.go` noticing.

### Errors: two for the whole package

```go
var (
	ErrEmpty    = errors.New("the note is empty")
	ErrNotFound = errors.New("note not found")
)
```

`Delete` wraps `ErrNotFound` and adds the number, while `main.go` recognises it with `errors.Is` and prints the human "мұндай нөмір жоқ". A clear sentence for the user, the detail in the error text for the programmer.

### The test knows nothing about the keyboard

```go
package notes

import (
	"errors"
	"testing"
)

func TestAddAndDelete(t *testing.T) {
	s := NewStore()

	n, err := s.Add("the steppe")
	if err != nil {
		t.Fatalf("a valid note was rejected: %v", err)
	}
	if n.ID != 1 || n.Words != 2 {
		t.Errorf("got ID %d, %d words; wanted 1 and 2", n.ID, n.Words)
	}

	if _, err := s.Add("   "); !errors.Is(err, ErrEmpty) {
		t.Errorf("for an empty note got %v, wanted ErrEmpty", err)
	}

	if err := s.Delete(99); !errors.Is(err, ErrNotFound) {
		t.Errorf("for a missing number got %v, wanted ErrNotFound", err)
	}

	if err := s.Delete(1); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if s.Len() != 0 {
		t.Errorf("after the delete %d notes remain", s.Len())
	}
}

func TestLettersCountsLettersNotBytes(t *testing.T) {
	n := Note{Text: "Шаңырақ"}
	if got := n.Letters(); got != 7 {
		t.Errorf("got %d letters, wanted 7", got)
	}
}
```

The test checks `Store` without running `main`. Input, output and command parsing are none of its business — and that is not a trick but a consequence of the logic living in the package rather than in `main`.

> **Picture it.** An engine on a test bench. It is run apart from the car: faults are easier to find and nobody has to drive out onto the road each time.

Run `go test ./...` and both files of the package are checked at once.

### Why the notes live only until you quit

Close the program and the notes are gone. That is not an omission: files and a database come in later modules, and what matters today is that **all the logic is already in place** and will not change when the store becomes a real one.

That is exactly how you will add it later: `Store` keeps the same methods, and only what is inside changes.

## The lesson map

![The lesson map: six familiar parts come together into one program](/static/course/go/map-capstone-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why does the test live in the `notes` package rather than next to `main.go`?
2. What happens if `main.go` tries to write `store.items`?
3. Why does `Delete` wrap the error when it could have returned `ErrNotFound` as it is?

## Exercise

**Required.** Give notes a tag. A `Tag string` field in `Note`, a `tag number word` command that puts a tag on an existing note, and a `bytag word` command that shows the notes carrying it. Use the error you already have for a number that does not exist — `ErrNotFound`. Add a test for both commands and get `go test ./...` to pass.

**Optional.**

- A `stats` command: how many notes, words and letters there are in total.
- Make `find` search by tag as well.
- Put the finished notebook on GitHub as a repository of its own with a README — it is the first finished program you can show anyone.

## Where this goes in your blog

The module on the language ends here. Next comes the web: the same `Store`, but called by HTTP handlers instead of `add` and `list`, answering a browser instead of a terminal.

Everything you wrote today stays. Only who asks the questions changes.

## Answers

1. Because it checks `Store`, and `Store` lives in the `notes` package. A test in the same package sees the internal fields too, so nothing has to be opened up for the sake of testing. It also means the test does not depend on keyboard input.
2. The program will not build: `items` is written in lower case and is visible only inside the `notes` package. The compiler will say `cannot refer to unexported field items`.
3. So that the number that was not found stays in the message: `delete 9: жазба табылмады`. Wrapping adds the detail while keeping the error itself, which is why `errors.Is` still recognises it.

## Sources

- [Package bufio — documentation](https://pkg.go.dev/bufio)
- [strings.Cut — documentation](https://pkg.go.dev/strings#Cut)
- [Go spec: Switch statements](https://go.dev/ref/spec#Switch_statements)
- [Effective Go: switch](https://go.dev/doc/effective_go#switch)
