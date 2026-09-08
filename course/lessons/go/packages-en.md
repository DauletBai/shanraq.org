# Packages and modules in Go: go mod and the capital letter

_Лид (summary):_ **The fifteenth lesson of the Go course. A package is a folder, and a capital letter in a name is the door out of it. How to split a program across folders, what go.mod and go.sum hold, how to bring in someone else's package with go get, what the internal folder is for, and why go mod tidy is run before the code goes anywhere.**

## Why this matters

`main.go` has grown: the type, the methods, the store, the errors — all in one file. Time to spread it out.

And the other half of the lesson: sooner or later you will want somebody else's code. Writing your own identifier generator or JSON parser is a waste of a life when it has already been written and tried by thousands of people.

## The whole thing at once

A new folder `sabaq-14`, and another one inside it called `blog`. Three files:

```
sabaq-14/
  go.mod
  main.go
  blog/
    article.go
    store.go
```

`blog/article.go`:

```go
package blog

import "fmt"

type Article struct {
	Title string
	Words int
}

func (a Article) ReadingTime() int {
	return (a.Words + 199) / 200
}

func (a Article) String() string {
	return fmt.Sprintf("%s (%d min)", a.Title, a.ReadingTime())
}
```

`blog/store.go`:

```go
package blog

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("article not found")

type Store struct {
	items map[string]Article
}

func NewStore() *Store {
	return &Store{items: map[string]Article{}}
}

func (s *Store) Add(slug string, a Article) {
	s.items[slug] = a
}

func (s *Store) Get(slug string) (Article, error) {
	a, ok := s.items[slug]
	if !ok {
		return Article{}, fmt.Errorf("get %q: %w", slug, ErrNotFound)
	}
	return a, nil
}

func (s *Store) Len() int { return len(s.items) }
```

`main.go`:

```go
package main

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"sabaq14/blog"
)

func main() {
	s := blog.NewStore()
	s.Add("dala", blog.Article{Title: "About the steppe", Words: 400})

	a, err := s.Get("dala")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(a)
	fmt.Println("articles:", s.Len())

	if _, err := s.Get("kokek"); errors.Is(err, blog.ErrNotFound) {
		fmt.Println("404:", err)
	}

	fmt.Println("new id:", uuid.NewSHA1(uuid.NameSpaceURL, []byte("shanraq.org/dala")))
}
```

Three commands, in order:

```
go mod init sabaq14
go get github.com/google/uuid
go run .
```

Before you run it, and without looking below, say out loud what the program will print. Then run it and compare.

```
About the steppe (2 min)
articles: 1
404: get "kokek": article not found
new id: 7046abf7-a993-54e5-8960-558b4e3d88d9
```

You will get exactly the same identifier: it is computed from the address rather than from chance.

## Taking it apart

### A package is a folder

Every file in a folder belongs to one package, and the first line in each of them is the same: `package blog`. How you split a package across files is up to you — the compiler does not care, and a person finds it easier by meaning: the type in one file, the store in another.

The package name usually matches the folder name. Inside a package the files see each other whole, with no imports at all: `store.go` uses the `Article` type from the file next door and does nothing to get at it.

> **Picture it.** A room in a house. Inside the room everything is to hand. What goes out into the corridor is only what you carried to the door yourself.

### A capital letter is the door out

Now it is clear why the fields in the structs lesson were called `Title` and `Words`.

There is one rule and no exceptions: **a name starting with a capital letter is visible from other packages, one starting with a small letter is not.** It applies to everything: types, functions, methods, fields, variables.

In our store the field `items` is written in lower case, and that is not an accident. Try this from `main.go`:

```go
len(s.items)
```

The program will not build:

```
s.items undefined (cannot refer to unexported field items)
```

That is exactly why `Store` has a `Len()` method: what leaves the package is the ability, not the arrangement. Replace the map with a database tomorrow and nothing outside breaks.

### `go mod init`, and what is in go.mod

```
go mod init sabaq14
```

This command declares a **module** — the unit that can be built and published. A module is a tree of folders with a `go.mod` at its root; the packages are the folders inside it.

After `go get` the file looks like this:

```
module sabaq14

go 1.27.1

require github.com/google/uuid v1.6.0
```

The first line is the module's name, which is also the start of every internal import path. The second is the language version you are building with; yours will be your own. The third is the list of what you took from other people, with exact versions.

### Importing your own package

```go
import "sabaq14/blog"
```

The path starts with the module's name, then the folder's path inside it. Not a relative `./blog`, not an absolute path on disk: only this way.

Notice the order in the `import` block: the standard library first, then other people's packages, then your own, with a blank line between the groups. That is not a rule of the language but a shared habit, and the editor will arrange it for you on save.

### Someone else's package: `go get`, `go.sum`, `go mod tidy`

`go get github.com/google/uuid` does three things: it downloads the package, adds a line to `go.mod`, and writes a **fingerprint** of what it downloaded into `go.sum`.

`go.sum` is not a list of dependencies but a list of checksums. On the next build Go checks that what was downloaded matches what was recorded. Both files go into the repository.

> **Picture it.** A receipt with a fingerprint on it. The list of books you borrowed is `go.mod`. `go.sum` is the mark that shows the book handed back is the same one, not a similar-looking copy with different text inside.

`go mod tidy` puts things in order: it adds what you imported and forgot to declare, and removes what is no longer used. Run it before sending your code anywhere: right after `go get` our package was marked `// indirect`, and only `tidy` removed that note.

### The package name is read at the call site

From outside people write `blog.Article` and `blog.NewStore`. So there is no need to repeat the package name inside it: `blog.BlogArticle` is a stutter, while `blog.Article` reads as a phrase.

For the same reason package names are short and single-word, lower case, with no underscores: `http`, `fmt`, `errors`.

### `internal/` — the folder closed from outside

Name a folder `internal` and only its neighbours within the module can import what is in it. Another project that pulls in your module cannot reach it — the compiler enforces this itself.

It takes a minute to check. Two modules: in the first the package sits in `internal`, and the second tries to import it.

```
$ cd lib && go build ./...          # свой модуль читает свой internal
свой модуль: собралось

$ cd ../app && go build ./...       # чужой модуль лезет в тот же пакет
package example.com/app
	main.go:6:2: use of internal package example.com/lib/internal/store not allowed
```

The compiler does not advise and does not warn — it **refuses to build**. `internal` is the one folder name Go understands as a rule of the language rather than a habit.

Handy when you want packages but do not want to promise strangers that they will never change.

> **Picture it.** A door marked "staff only". It is in the same building and not locked to the staff, but a visitor does not go through it.

### The tree of a project: `cmd`, `internal` and `pkg`

Our blog has a flat tree: `main.go` beside the `blog` folder, and no `cmd` or `internal` at all. That is not a simplification for the sake of a lesson — it is what Go itself advises for a program of this size: while the module is one program, extra folders only lengthen the import paths and add nothing.

Three names you will meet in other people's repositories:

- **`cmd/`** — for when a module holds more than one program. `cmd/app`, `cmd/migrate`, `cmd/export`: each has its own `main`, while the shared code sits beside them and is imported by all.
- **`internal/`** — what cannot be imported from outside. The rule of the language we have just measured.
- **`pkg/`** — a folder for code meant for other people's projects. This is **not** a rule of Go but a habit of part of the community: the compiler knows nothing about `pkg`.

The site you are reading this lesson on is laid out like that: four programs in `cmd/` — `app`, `migrate`, `export` and `adminctl` — five packages in `internal/` and seventeen in `pkg/`. That tree grew out of a flat one, and it grew when the second program appeared rather than in advance.

Does the course's blog suffer for not having any of it? No: it has one program and one shared package, and any of those names would add an empty folder. But when your blog gains a second program — a separate `migrate`, say — you already know where to put it.

And one last thing. Do not start with folders. While there is one file with two hundred lines in it, leave it alone. Splitting is worth doing when you can already see where the boundary runs, not before.

## The lesson map

![The lesson map: a package is a folder, a capital letter is the door out](/static/course/go/map-packages-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why can `main.go` not reach `s.items`, when both files are part of one program?
2. How does `go.mod` differ from `go.sum`, and why are both needed?
3. Why `blog.Article` rather than `blog.BlogArticle`?
4. How does `internal/` differ from `pkg/`?

## Exercise

**Required.** Take your program from the errors lesson and split it into two packages: `blog`, holding the type, the store and the errors, and `main`, which only calls them. Leave the data field in lower case and add a `Len()` method. Check that `go run .` works and that reaching for the internal field from `main` does not build.

**Optional.**

- Bring in `github.com/google/uuid` and give every article an identifier.
- Delete the `require` line from `go.mod` and run `go mod tidy`. See what happens.
- Rename the `blog` folder to `internal/blog` and fix the import. The program should work as before.

## Where this goes in your blog

This site is built the same way: the server in one package, the articles in another, the store in a third, each talking to the rest through a handful of names that start with a capital letter. The state of the blog after this lesson is [step-6](https://github.com/DauletBai/shanraq.org/tree/main/course/go-blog/step-6): the same blog, spread across folders.

## Answers

1. Because `items` is written in lower case, so it is visible only inside the `blog` package. There is one program but two packages, and the boundary runs through the first letter of the name. What leaves the package is the `Len()` method — the ability, not the arrangement.
2. `go.mod` is the list of what you took and in which versions; a person reads it. `go.sum` holds the checksums of what was downloaded; Go reads it to notice a substitution. Both files are kept in the repository.
3. Because the package name is read at the call site: from outside it comes out as `blog.Article`. Repeating the package name inside it is a stutter, and it shows up on every line of your code that other people read.
4. `internal` is a rule of the language: another module that reaches into it simply does not build, and it is the compiler that says so. `pkg` is an ordinary folder name and a habit of part of the community; Go knows nothing about it and forbids nothing through it.

## Sources

- [How to Write Go Code — go.dev](https://go.dev/doc/code)
- [Managing dependencies — go.dev](https://go.dev/doc/modules/managing-dependencies)
- [Go Modules Reference](https://go.dev/ref/mod)
- [Effective Go: package names](https://go.dev/doc/effective_go#package-names)
