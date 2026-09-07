# JSON and an API of your own: a blog is not read by browsers alone

_Лид (summary):_ **The forty-fourth lesson of the Go course. The tags on a struct decide what a field is called outside and what is hidden altogether; an extra field in somebody's JSON stays silent unless you ask otherwise; a number with no type becomes a float64. And a handler whose errors are JSON too.**

## Why this matters

The blog hands out pages, and that is right up to the day its content is needed by something other than a person.

A mobile app wants the same articles without the markup. A neighbouring site wants the list of titles. You want a script that checks everything is in place. Handing any of them an `<h1>` is pointless: they parse the answer with a program, and they need a format a program understands.

That format is **JSON**, and in Go it lives in the standard library.

## The whole thing at once

A new folder, `go mod init sabaq41`:

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

// Article is what goes out. The tags decide what a field is called in JSON;
// omitempty drops what is empty, and a dash hides a field altogether.
type Article struct {
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	Tags      []string  `json:"tags,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
	Draft     bool      `json:"-"`
	secret    string    // lower case — it will never leave the package
}

func main() {
	fmt.Println("== what goes out")
	a := Article{
		Slug: "dala", Title: "The steppe", Draft: true, secret: "not for strangers",
		UpdatedAt: time.Date(2026, time.September, 6, 12, 0, 0, 0, time.UTC),
	}
	out, _ := json.MarshalIndent(a, "", "  ")
	fmt.Println(string(out))

	fmt.Println()
	fmt.Println("== empty fields")
	empty, _ := json.Marshal(Article{Slug: "salem", Title: "Hello"})
	fmt.Println(string(empty))

	fmt.Println()
	fmt.Println("== what comes in")
	body := `{"slug":"koktem","title":"Spring","tags":["steppe"],"extra":1}`
	var got Article
	fmt.Println("plain decoding:  ", json.Unmarshal([]byte(body), &got), "->", got.Slug, got.Tags)

	strict := json.NewDecoder(strings.NewReader(body))
	strict.DisallowUnknownFields()
	fmt.Println("strict decoding: ", strict.Decode(&Article{}))

	fmt.Println()
	fmt.Println("== a number with no type")
	var any1 any
	json.Unmarshal([]byte(`{"n":7}`), &any1)
	n := any1.(map[string]any)["n"]
	fmt.Printf("%v is a %T\n", n, n)

	fmt.Println()
	fmt.Println("== the handler's answer")
	srv := httptest.NewServer(routes())
	defer srv.Close()
	show(srv.URL + "/api/articles/dala")
	show(srv.URL + "/api/articles/joq")
}

func routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/articles/{slug}", func(w http.ResponseWriter, r *http.Request) {
		if r.PathValue("slug") != "dala" {
			// An error is JSON too: whoever reads an API parses the answer
			// rather than looking at a page with a picture on it.
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "no such article"})
			return
		}
		writeJSON(w, http.StatusOK, Article{Slug: "dala", Title: "The steppe", Tags: []string{"steppe"}})
	})

	return mux
}

// writeJSON sets the header before the status: after WriteHeader the headers
// are already on their way and changing them is too late.
func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// The status has gone out; the log is all that is left.
		fmt.Println("could not write the answer:", err)
	}
}
```

The `show` from the second file simply prints the answer. The output:

```
== what goes out
{
  "slug": "dala",
  "title": "The steppe",
  "updated_at": "2026-09-06T12:00:00Z"
}

== empty fields
{"slug":"salem","title":"Hello"}

== what comes in
plain decoding:   <nil> -> koktem [steppe]
strict decoding:  json: unknown field "extra"

== a number with no type
7 is a float64

== the handler's answer
200 application/json; charset=utf-8 -> {"slug":"dala","title":"The steppe","tags":["steppe"]}
404 application/json; charset=utf-8 -> {"error":"no such article"}
```

## Taking it apart

### Only what starts with a capital letter goes out

The first block of output has three fields; the struct has six. Where did the rest go?

`secret` is written in lower case, which makes it invisible outside its package — to `encoding/json` as well. That is the rule from the packages lesson, and here it has a pleasant consequence: an internal field cannot be handed out by accident.

`Draft` disappeared for a different reason: it carries the tag `json:"-"`, which means "never output". That is how what nobody outside needs to know is hidden: a draft flag, internal notes, housekeeping fields.

> **Picture it.** A shop window. Not everything in the storeroom goes into it, and it is not the storeroom that decides.

### Tags: the name, `omitempty` and `omitzero`

A tag is the string in back quotes after a field, and it tells `encoding/json` what to do:

```go
Slug      string    `json:"slug"`
Tags      []string  `json:"tags,omitempty"`
UpdatedAt time.Time `json:"updated_at,omitzero"`
```

Without a tag a field goes out under its Go name — `Slug`, with a capital. Outside that is not how it is done: JSON says `slug` or `updated_at`, and the tag brings the name into line.

`omitempty` drops a field whose value is empty: an empty string, a zero, an empty slice. The second block of output shows it — an article with no tags has no `tags` field at all.

For a `time.Time` `omitempty` **does not work**: it is a struct, and a struct is never "empty" in that sense. Hence `omitzero`, which looks at the zero value of the type — and in the output the zero time went with it.

### Decoding somebody's JSON: the extra stays silent

The third block is about what comes in. The body has a field `extra` that the struct does not have, and plain decoding **said nothing at all**: `<nil>`, the article decoded.

That is by design, and it is convenient until it is dangerous: a typo in a field name turns into a silently skipped value. A reader sends `titel` instead of `title` — you accept an article with no title and never learn of it.

So wherever somebody's JSON is accepted, strictness goes on:

```go
dec := json.NewDecoder(r.Body)
dec.DisallowUnknownFields()
```

Then an extra field becomes an error: `json: unknown field "extra"` — and the sender learns of their typo straight away rather than a week later.

### A number with no type becomes a `float64`

The fourth block is short, and people pay for it in hours of debugging.

Decode into an `any` (or a `map[string]any`) and every number becomes a `float64`. Not the `int` it looks like in the text, but a `float64` — JSON has no separate integer type, and Go picks the most general one.

Hence the `interface conversion: interface {} is float64, not int` panic for those who decoded into an `any`. The cure is simple: decode into a struct. Then an `int` stays an `int`, and the extra fields fall away by type as well.

### The header goes before the status

`writeJSON` does three things in a strict order:

```go
w.Header().Set("Content-Type", "application/json; charset=utf-8")
w.WriteHeader(code)
json.NewEncoder(w).Encode(v)
```

The order is not accidental. `WriteHeader` sends the status line **and every header at once**; after it `w.Header().Set` changes nothing — Go quietly ignores it. That is how a `Content-Type` gets lost and a browser is left guessing what it was sent.

The same goes for the body: the first write to `w` without a `WriteHeader` means `200`. If another status is wanted, it goes in before a single byte is written.

### An error is JSON too

The last two lines of the output: both the `200` and the `404` came with `application/json`.

That is a rule rather than a detail. Whoever reads an API parses the answer with a program: if an HTML page with a picture arrives in place of an error, the parsing breaks at the very first character, and the log says "unexpected `<`" instead of "no such article".

So an API has one shape for every answer. The status code says what happened, the body gives the detail: `{"error": "no such article"}`.

### A stream instead of a string: `Encoder` and `Decoder`

`json.Marshal` assembles the whole answer in memory, while `json.NewEncoder(w).Encode(v)` writes **straight into the stream**. For one article that makes no difference; for a list of ten thousand it is tens of megabytes that need not be held at once.

Decoding is the same: `json.NewDecoder(r.Body).Decode(&v)` reads from the body of the request without gathering it into a string. And only the decoder has `DisallowUnknownFields`, which is why we reached for it.

Do not forget the cap: the body comes from outside, so the `http.MaxBytesReader` from the upload lesson belongs here too.

## The lesson map

![Lesson map: a struct, its tags and a stream](/static/course/go/map-json-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why does a field in lower case not reach the JSON?
2. What happens to a typo in a field name under plain decoding?
3. Why is `Content-Type` set before `WriteHeader`?

## Exercise

**Required.** Give the blog an API: `GET /api/articles` hands back the list, `GET /api/articles/{slug}` one article. Answers and errors are JSON with the right header, and drafts and internal fields never leave.

All of it is done in [step-29](https://github.com/DauletBai/shanraq.org/tree/main/course/go-blog/step-29) — compare against it once you have written your own.

**Optional.**

- Add a `POST /api/articles` that takes an article with a strict decoder and answers `201` with the address of what it made.
- Add a `reading_time` field to the answer — a method of a type does not reach JSON, so it has to be computed while the answer is assembled.
- Check your API from a terminal: `curl -s localhost:8090/api/articles | jq`.

## Where this goes in your blog

The blog gets a second way of being read — and the first that does not depend on markup. From now on a page can be rewritten without breaking whoever reads the data.

Debts. The API hands everything back at once, with no paging. Nothing limits it: a hundred requests a second from one script go unstopped. And there is no writing through the API — which would mean deciding how to guard it.

## Answers

1. Because `encoding/json` sees exported fields only — the same rule as for any other package. A field in lower case does not exist for it, and cannot be handed out by accident.
2. Nothing: plain decoding silently skips fields the struct does not have. The sender is sure they sent a title; you get an empty string. `DisallowUnknownFields` cures it.
3. Because `WriteHeader` sends the status and the headers as one piece. After it a header is too late to change — Go ignores it, and the reader gets an answer with no `Content-Type`.

## Sources

- [Go: the encoding/json package](https://pkg.go.dev/encoding/json)
- [The Go blog: JSON and Go](https://go.dev/blog/json)
- [MDN: HTTP status codes](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status)
