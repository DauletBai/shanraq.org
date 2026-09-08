# Routing in Go: who handles which address

_Лид (summary):_ **The twentieth lesson of the Go course. ServeMux sends requests to handlers: the method and the path are written on one line, {slug} pulls a piece out of the address, and {$} stops the front page from catching everything. The router answers 404 and 405 on its own, and two patterns that conflict are caught at start-up rather than in production.**

## Why this matters

In the previous lesson the handler checked the path and the method itself:

```go
if r.Method != http.MethodGet { ... }
if r.URL.Path != "/" { ... }
```

That works exactly as far as the second page. By the fifth, such a handler turns into a ladder of `if`s where it is easy to forget a branch and serve the front page instead of a 404.

Let the addresses be sorted out by the thing whose job that is — the router.

## The whole thing at once

A new folder, `go mod init sabaq18`, `main.go`:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "the home page")
	})

	mux.HandleFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "about the author")
	})

	mux.HandleFunc("GET /read/{slug}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "article:", r.PathValue("slug"))
	})

	mux.HandleFunc("POST /read/{slug}/like", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "liked:", r.PathValue("slug"))
	})

	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

Run it and check every case. Before looking below, say out loud what the answer to `GET /read/dala/like` will be:

```
GET  /                     → 200  the home page
GET  /about                → 200  about the author
GET  /read/dala            → 200  article: dala
POST /read/dala/like       → 200  liked: dala
GET  /read/dala/like       → 405  Allow: POST
GET  /kokek                → 404
POST /about                → 405  Allow: GET, HEAD
```

There is not a single path or method check left in the code — and the answers came out more precise than before.

## Taking it apart

### Method and path on one line

```go
mux.HandleFunc("GET /about", ...)
```

Before Go 1.22 a pattern held only the path and the method was checked inside the handler. Now the method is part of the pattern, and that `if` is gone.

Hence the rule: **one pattern, one action.** `GET /read/{slug}` shows an article, `POST /read/{slug}/like` marks it. Two lines and two functions, not one with a fork inside.

> **Picture it.** Labelled doors. In an office you do not ask one window for everything: each has a sign saying what is done there. You read the sign and go to the right one.

### `{slug}` — a piece of the address as a value

```go
mux.HandleFunc("GET /read/{slug}", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "article:", r.PathValue("slug"))
})
```

Braces in a pattern mean: any value goes here, call it `slug`. You fetch it with `r.PathValue("slug")`.

The name inside the braces is yours — call it `{id}` or `{name}`, as long as `PathValue` asks for the same thing. And an article's address on this site is built exactly that way: `/read/go-marshruttar`.

### `{$}`, or the front page catches everything

In Go the pattern `GET /` does not mean "the front page" but "everything nobody else claimed". The very case from the previous lesson.

See for yourself: replace `GET /{$}` with `GET /` and `/kokek` will get the front page with a 200 instead of a 404.

A `{$}` at the end of a pattern means "the path ends here". `GET /{$}` is the front page itself and nothing else.

> **Picture it.** A tray marked "everything else". Such a tray is useful while everything else is what goes in it. Label it "incoming" and both the post and the rubbish land there, and you will not notice for a while.

### Which one wins when two patterns fit

The order of registration makes no difference. The **more specific** pattern wins:

```go
mux.HandleFunc("GET /read/{slug}", ...)   // any article
mux.HandleFunc("GET /read/dala", ...)     // this one exactly
```

`/read/dala` goes to the second, `/read/kokek` to the first, and it does not matter which order you wrote them in.

And if two patterns catch **the same requests**, the program will not start:

```
panic: pattern "GET /read/{name}" conflicts with pattern "GET /read/{slug}":
	GET /read/{name} matches the same requests as GET /read/{slug}
```

A panic at start-up looks harsh, but it is a gift: the routing conflict was found by you on your own machine rather than by a reader in production.

> **Picture it.** Two identical signposts at a fork. Until one comes down you cannot drive on — and it is right that the car stops there rather than guessing.

### The router answers 404 and 405 itself

Not one handler contains a `WriteHeader(404)` or an `Allow`. Yet the answers are right, because `ServeMux` knows every pattern: no pattern matched the path — 404; the path matched but the method did not — 405 with the list of what is allowed.

Look at that list: `Allow: GET, HEAD`. `HEAD` appeared on its own — it is the same `GET` without a body, and Go answers it for you.

### When there are many routes

After a dozen lines of registration `main` stops being readable. The usual move is to lift the routes into a function of their own that returns a ready router:

```go
func routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /about", about)
	mux.HandleFunc("GET /read/{slug}", article)
	return mux
}
```

Then `main` keeps one line, and every address of the site is visible in one place as a list. That same list can be handed to a test without bringing up a real server.

### The `nil` that was in the second argument

```go
http.ListenAndServe(":8080", mux)
```

Since the third lesson there was a `nil` there. Now it can be said what it meant: "use the default router" — the one `http.HandleFunc` without a `mux` writes into.

A `mux` of your own is better for three reasons: it is visible in the code, it can be handed to a test, and two different servers in one program do not get in each other's way.

## The lesson map

![The lesson map: one entrance, several labelled doors, the router chooses](/static/course/go/map-routes-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. How does `GET /` differ from `GET /{$}`?
2. You registered `GET /read/{slug}` and `GET /read/dala`. Where does `/read/dala` go, and does the order of the lines matter?
3. Where does `Allow: GET, HEAD` come from if you never wrote it?

## Exercise

**Required.** Split the blog across routes: `GET /{$}` for the list of articles, `GET /read/{slug}` for one article by its address, `GET /about` for yourself. Take the articles from a slice of structs, and answer 404 with your own text when there is no such address. Check every case with `curl -i`, `POST /about` included.

**Optional.**

- Replace `GET /{$}` with `GET /` and see what becomes of `/kokek`. Then put it back.
- Register two patterns that catch the same requests and read the panic in full.
- Add `GET /read/{slug}/{part}` and see what `PathValue` gives for both names.

## Where this goes in your blog

The blog's addresses are nearly real now: `/`, `/about`, `/read/{slug}`. Only the content is missing — for the moment strings travel in the response rather than pages.

The next lesson brings a wrapper that does the same thing for every route at once: a log, a response time, a check. After that come templates, and HTML rides in the body instead of `fmt.Fprintln`.

## Answers

1. `GET /` catches every path no other pattern matched, non-existent ones included. `GET /{$}` catches only the root itself: `{$}` means the path ends there. Without it any typo in an address would return the front page with a 200.
2. To `GET /read/dala`, because it is more specific. The order of the lines makes no difference: the router picks the most precise pattern rather than the first one that fits.
3. `ServeMux` adds it. It knows every pattern for that path and answers 405 with the list of allowed methods. `HEAD` appears in the list because it is the same `GET` without a body, and Go handles it itself.

## Sources

- [http.ServeMux — routing patterns](https://pkg.go.dev/net/http#ServeMux)
- [Routing Enhancements for Go 1.22 — the go.dev blog](https://go.dev/blog/routing-enhancements)
- [r.PathValue — documentation](https://pkg.go.dev/net/http#Request.PathValue)
