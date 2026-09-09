# CSRF and XSS: three guards everything gets

_Лид (summary):_ **The thirty-ninth lesson of the Go course. A stranger's script on your page, and another site submitting a form in your reader's name. Measured here: how html/template escapes by context — in text, in a link and inside a script — and how a token tells your form from somebody else's.**

## Why this matters

The blog can tell people apart and keeps what belongs to others untouchable. What remains is closing the two holes that go around all of it without guessing a single password.

**XSS** is a stranger's script that got onto your page. It runs in a reader's browser in the name of your site, which means it can do everything the reader can: read the page, submit a form, take the session away.

**CSRF** is another site submitting a form to yours. The browser obediently attaches your reader's cookie, and the server sees an ordinary request that is entirely "its own".

Neither hole is closed by vigilance. They are closed by two techniques and three headers.

## The whole thing at once

A new folder, `go mod init sabaq36`:

```go
package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
)

// tokens holds one token per session. In a real blog that is a row in the
// database beside the session; a map here keeps the lesson about the guard
// rather than about storage.
var (
	mu     sync.Mutex
	tokens = map[string]string{}
)

func main() {
	srv := httptest.NewServer(secure(routes()))
	defer srv.Close()

	fmt.Println("== the headers the wrapper sets")
	for _, h := range []string{"Content-Security-Policy", "X-Content-Type-Options", "Referrer-Policy"} {
		fmt.Printf("%-24s %s\n", h, header(srv.URL+"/form", h))
	}

	fmt.Println()
	fmt.Println("== escaping the output")
	fmt.Println(get(srv.URL + "/read?title=" + url.QueryEscape(`<script>alert(1)</script>`)))

	fmt.Println()
	fmt.Println("== the form and the token")
	token := tokenFrom(get(srv.URL + "/form"))
	fmt.Println("token in the form: ", token[:8]+"…")
	fmt.Println("no token:          ", post(srv.URL+"/delete", nil))
	fmt.Println("somebody else's:   ", post(srv.URL+"/delete", url.Values{"csrf": {"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"}}))
	fmt.Println("its own:           ", post(srv.URL+"/delete", url.Values{"csrf": {token}}))
}

func routes() http.Handler {
	mux := http.NewServeMux()

	// The article page: the title came from a form, which means from somebody else.
	mux.HandleFunc("GET /read", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.New("read").Parse(`<h1>{{ .Title }}</h1>`))
		t.Execute(w, map[string]any{"Title": r.URL.Query().Get("title")})
	})

	mux.HandleFunc("GET /form", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.New("form").Parse(
			`<form method="post" action="/delete">` +
				`<input type="hidden" name="csrf" value="{{ .Token }}">` +
				`<button>delete</button></form>`))
		t.Execute(w, map[string]any{"Token": issue(sessionOf(r))})
	})

	mux.HandleFunc("POST /delete", func(w http.ResponseWriter, r *http.Request) {
		if !valid(sessionOf(r), r.FormValue("csrf")) {
			http.Error(w, "the token does not fit", http.StatusForbidden)
			return
		}
		fmt.Fprint(w, "deleted")
	})

	return mux
}

// secure is the wrapper that puts three headers on every answer.
func secure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// What the page may load and run. default-src 'self' means from our own
		// site only, and a stranger's script will not run even if it reached
		// the markup.
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		// The browser must not guess the type of a file: if it says text/plain,
		// then it is not a script.
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// Another site is not told the full address somebody came from.
		w.Header().Set("Referrer-Policy", "same-origin")
		next.ServeHTTP(w, r)
	})
}

// issue hands out the token of a session and remembers it. One per session:
// it is no secret from its owner, it is a secret from another site.
func issue(session string) string {
	mu.Lock()
	defer mu.Unlock()
	if t, ok := tokens[session]; ok {
		return t
	}
	raw := make([]byte, 32)
	rand.Read(raw)
	t := base64.RawURLEncoding.EncodeToString(raw)
	tokens[session] = t
	return t
}

// valid compares the tokens in constant time: an ordinary == stops at the
// first byte that differs, and the time it takes gives the token away.
func valid(session, got string) bool {
	mu.Lock()
	want := tokens[session]
	mu.Unlock()
	if want == "" || got == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(want), []byte(got)) == 1
}

func sessionOf(r *http.Request) string {
	if c, err := r.Cookie("session"); err == nil {
		return c.Value
	}
	return "a guest"
}

func tokenFrom(page string) string {
	const mark = `name="csrf" value="`
	i := strings.Index(page, mark)
	rest := page[i+len(mark):]
	return rest[:strings.IndexByte(rest, '"')]
}
```

The helpers are in a second file, `helpers.go`. The output:

```
== the headers the wrapper sets
Content-Security-Policy  default-src 'self'
X-Content-Type-Options   nosniff
Referrer-Policy          same-origin

== escaping the output
<h1>&lt;script&gt;alert(1)&lt;/script&gt;</h1>

== the form and the token
token in the form:  S3tZqJ-E…
no token:           403 Forbidden — the token does not fit
somebody else's:    403 Forbidden — the token does not fit
its own:            200 OK — deleted
```

## Taking it apart

### `html/template` escapes by context

Start with the good news: Go guards against XSS in templates by itself, and does it more cleverly than one might think. It looks at **where** a value lands and escapes accordingly:

```
<h1>{{ .Title }}</h1>                      -> <h1>&lt;script&gt;alert(1)&lt;/script&gt;</h1>
<a href="{{ .URL }}">link</a>              -> <a href="#ZgotmplZ">link</a>
<h1>{{ .Safe }}</h1>                       -> <h1><script>alert(1)</script></h1>
<script>var t = "{{ .Title }}";</script>   -> <script>var t = "\u003cscript\u003ealert(1)\u003c\/script\u003e";</script>
```

The first line is ordinary text: the angle brackets became `&lt;` and `&gt;`, and the script turned into a visible string.

The second is the interesting one. The address `javascript:alert(1)` landed in an `href`, and Go **refused to print it**, putting `#ZgotmplZ` there instead. That is neither an error nor rubbish: it is how `html/template` marks a value that is unsafe in this place. Seeing it on a page means something reached a link that had no business being there.

The fourth is a value inside `<script>`: escaped by the rules of JavaScript now, as `<`. The same template engine did three different things in three places, because it parses markup rather than gluing strings.

> **Picture it.** An interpreter who knows where you have come. In one place they say "hello", in another "your honour", in a third they stay silent — because they understand the setting instead of merely swapping words.

Hence the rule: **do not build HTML out of strings**. `fmt.Sprintf("<h1>%s</h1>", title)` knows about no places at all and will insert whatever it is given.

### `template.HTML` is the hole you open yourself

The third line of the output: the same value wrapped in `template.HTML` comes out as it is — a script on the page.

That is by design: the type means "I take responsibility, this is real markup". It is needed — we used it in the search lesson to show the highlighting — and it is also the only way to open an XSS hole for yourself in Go.

The rule is simple: `template.HTML` may only go on what **you assembled yourself**, and never on what came from a person. In the search lesson that is why the fragment was escaped by hand and only our own markers were put back.

### CSRF: another site with your cookie

Now the second hole, built the other way round: no script is needed at all.

A reader signs in to your blog and opens somebody else's page in the next tab. On that page:

```html
<form method="post" action="https://your-blog/delete">
  <input type="hidden" name="slug" value="dala">
</form>
<script>document.forms[0].submit()</script>
```

The browser submits the form and **attaches your blog's cookie** — because a cookie is attached by where the request is going, not by where it came from. The server sees a signed-in person and obediently deletes the article.

You put up the first barrier in the sessions lesson: `SameSite=Lax` tells the browser not to attach the cookie to requests started on another site. It is good, and it cannot be the only thing: it knows nothing about your own subdomains, and old browsers do not have it.

And the consequence people forget: `Lax` lets ordinary link navigation through. Which means **`GET` must change nothing** — otherwise luring a reader to `<img src="your-blog/delete?slug=dala">` is enough.

### The token: what another site does not have

The second barrier is on the server. A random string tied to the session goes into the form and is compared when the form comes back:

```html
<input type="hidden" name="csrf" value="{{ .Token }}">
```

All three cases are in the output: with no token, 403; with somebody else's, 403; with its own, "deleted".

It works because another site **cannot read your page**: the browser will not let it see the content of a response from a different address. It can submit a form blind; it cannot learn the token.

The token is no secret from the reader: it sits in their own form. It is a secret from another page, and that is enough.

Then come the details worth knowing. The token is made per session rather than per form: it is simpler, and open tabs keep working. It has to be checked on every method that changes something, not only on deletion. And after signing in it is better to hand out a fresh one, for the same reason the session is handed out afresh.

### Comparing in constant time

The `subtle.ConstantTimeCompare` line in `valid` is not decoration.

An ordinary `==` on strings stops at the first byte that differs. So the time of the answer depends on how many leading bytes were guessed, and that time gives the secret away one character at a time. The attack is old, it is called a timing attack, and it catches more than tokens: any comparison of a secret.

`subtle.ConstantTimeCompare` compares the contents in a time that does not depend on where the strings diverge. One caveat worth knowing: when the lengths differ it returns zero at once — the constant time is promised for slices of equal length. That does not hurt us: our tokens are always the same length, and input of another length is rejected earlier. Write it where you compare a secret; for ordinary strings it is not needed.

### The three headers everything gets

The `secure` wrapper adds three headers to every answer — there they are in the output:

- **`Content-Security-Policy: default-src 'self'`** — what the page may load and run. With `'self'` a stranger's script will not run even if it somehow reached the markup. This is the last line when escaping has failed;
- **`X-Content-Type-Options: nosniff`** — the browser must not guess the type of a file from its content. Without it an uploaded "text" file can be executed as a script;
- **`Referrer-Policy: same-origin`** — another site is not told which address a visitor came from.

Three lines in one wrapper, and they apply to the whole blog at once — the very case middleware was invented for.

## The lesson map

![Lesson map: escaping, a token and headers](/static/course/go/map-csrf-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why did `html/template` put `#ZgotmplZ` in a link and `&lt;` in text?
2. Why can another site submit your form but not forge the token?
3. Why must `GET` change nothing?

## Exercise

**Required.** Close both holes in the blog. A wrapper puts three headers on every answer. Every form gets a hidden field with a token and every `POST` checks it; the check lives in a wrapper rather than in the handlers. Make sure a form without a token answers 403.

All of it is done in [step-24](https://github.com/DauletBai/shanraq.org/tree/main/course/go-blog/step-24) — compare against it once you have written your own.

**Optional.**

- File an article titled `<script>alert(1)</script>` and look at the page source. Then wrap the title in `template.HTML` and look again. Put it back the way it was.
- Save a file with a form pointing at your `/delete` and open it in the browser while signed in to the blog. First without the token check on the server, then with it.
- Add `img-src` to the CSP and allow images from your own domain only.

## Where this goes in your blog

After this lesson the blog stops being defenceless against the next tab. Three headers, a token in every form, and the escaping that was already there — that is the minimum below which a site is not opened to people.

Still open. The token lives in memory rather than in the database: restarting the server breaks every open form. The CSP allows inline styles, because we have them. And uploading files, which has dangers of its own, gets a lesson to itself.

## Answers

1. Because `html/template` parses markup and knows the place a value lands in. In text the angle brackets are dangerous; in an address it is the `javascript:` scheme, and the second it does not escape but rejects outright, putting `#ZgotmplZ` there.
2. Because the browser lets another site send a request but not read a response from your address. The token sits in your page, and another site cannot see it.
3. Because `SameSite=Lax` lets ordinary link navigation through, and any `<img src="...">` or a link in an email turns into an action. Only `POST` and its relatives change anything.

## Sources

- [Go: html/template and its contexts](https://pkg.go.dev/html/template#hdr-Contexts)
- [MDN: SameSite cookies](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Set-Cookie#samesitesamesite-value)
- [OWASP: preventing CSRF](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html)
