# Sessions and cookies: how a server remembers who signed in

_Лид (summary):_ **The thirty-seventh lesson of the Go course. HTTP remembers nothing: every request arrives from a stranger. A cookie holding a user number is forged in a second. So the browser carries a random token, the database keeps its hash, and three words close the rest: HttpOnly, Secure, SameSite.**

## Why this matters

In the last lesson a person registered. Now they want to sign in — and here the unpleasant part shows up: **HTTP remembers nothing**.

Every request arrives on its own. The server answered one and forgot it; the next request from the same person is indistinguishable from one by a passer-by. There is no connection between them in the protocol at all.

The connection has to be made by us, and it is made with a cookie — a small piece the server asks the browser to keep and send back with every request.

## The whole thing at once

A new folder, `go mod init sabaq34`. The program starts a server and calls it itself, so both the header and what happens on each request are visible:

```go
package main

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

const cookieName = "session"

// sessionLife is how long a sign-in lasts. A week: long enough not to sign in
// daily, short enough that a forgotten session does not live for years.
const sessionLife = 7 * 24 * time.Hour

func main() {
	os.Remove("blog.db")
	db := open()
	defer db.Close()

	srv := httptest.NewServer(routes(db))
	defer srv.Close()

	fmt.Println("== signing in")
	fmt.Println("wrong password: ", login(srv.URL, "aigul@example.kz", "guessing"))

	header, token := loginOK(srv.URL, "aigul@example.kz", "secretpass123456")
	fmt.Println("right password: ", header)

	fmt.Println("== who came")
	fmt.Println("no cookie:           ", ask(srv.URL, ""))
	fmt.Println("with cookie:         ", ask(srv.URL, token))
	fmt.Println("forged:              ", ask(srv.URL, "user_id=1"))
	fmt.Println("someone else's token:", ask(srv.URL, strings.Repeat("A", 43)))

	fmt.Println("== signing out")
	fmt.Println("signing out:         ", logout(srv.URL, token))
	fmt.Println("the same cookie:     ", ask(srv.URL, token))
}

func routes(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /me", func(w http.ResponseWriter, r *http.Request) {
		email, ok := current(db, r)
		if !ok {
			fmt.Fprint(w, "a guest")
			return
		}
		fmt.Fprintf(w, "signed in as %s", email)
	})

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		id, err := authenticate(db, r.FormValue("email"), r.FormValue("password"))
		if err != nil {
			http.Error(w, "address or password does not fit", http.StatusUnauthorized)
			return
		}

		token, err := newSession(db, id)
		if err != nil {
			http.Error(w, "did not work", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(sessionLife),
			HttpOnly: true,
			Secure:   overHTTPS(r),
			SameSite: http.SameSiteLaxMode,
		})
		w.WriteHeader(http.StatusSeeOther)
	})

	mux.HandleFunc("POST /logout", func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie(cookieName); err == nil {
			db.Exec(`delete from sessions where token_hash = ?`, hashToken(c.Value))
		}
		// The browser is asked to forget the cookie: same name, empty value, a
		// date in the past. The row is already deleted, which matters more.
		http.SetCookie(w, &http.Cookie{
			Name: cookieName, Value: "", Path: "/",
			Expires: time.Unix(0, 0), MaxAge: -1,
			HttpOnly: true, Secure: overHTTPS(r), SameSite: http.SameSiteLaxMode,
		})
		fmt.Fprint(w, "signed out")
	})

	return mux
}

// newSession hands out a random token and keeps its hash. The cookie carries
// the token; the database keeps only the hash, so a leaked database lets
// nobody sign in as anybody.
func newSession(db *sql.DB, userID int64) (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	token := base64.RawURLEncoding.EncodeToString(raw)

	_, err := db.Exec(
		`insert into sessions (token_hash, user_id, expires_at) values (?, ?, ?)`,
		hashToken(token), userID, time.Now().Add(sessionLife).UTC().Format(time.RFC3339))
	return token, err
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// current says who sent this request, and never takes the cookie's word for
// it: the cookie holds a random string and nothing else.
func current(db *sql.DB, r *http.Request) (string, bool) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return "", false
	}
	var email string
	err = db.QueryRow(`select u.email from sessions s
		join users u on u.id = s.user_id
		where s.token_hash = ? and s.expires_at > ?`,
		hashToken(c.Value), time.Now().UTC().Format(time.RFC3339)).Scan(&email)
	if err != nil {
		return "", false
	}
	return email, true
}

func authenticate(db *sql.DB, email, password string) (int64, error) {
	var id int64
	var hash string
	err := db.QueryRow(`select id, password_hash from users where email = ?`,
		strings.ToLower(strings.TrimSpace(email))).Scan(&id, &hash)
	if err != nil {
		return 0, errors.New("does not fit")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return 0, errors.New("does not fit")
	}
	return id, nil
}

// overHTTPS reports whether the request really arrived over TLS. The Secure
// flag depends on it: a cookie marked Secure does not travel over plain http,
// so on your own machine it would be set once and never sent back.
func overHTTPS(r *http.Request) bool { return r.TLS != nil }
```

The helpers (`open`, `login`, `ask` and the rest) live in a second file, `helpers.go`, and do what a browser normally does: [step-22](https://github.com/DauletBai/shanraq.org/tree/main/course/go-blog/step-22) shows the same thing in the blog, where the browser is real.

The output:

```
== signing in
wrong password:  401, no cookie
right password:  303, session=0HhwSnNtWOgQrUxJZv9YXCFz9-Ojj0XdOO6zNClfRe4; Path=/; Expires=Sun, 13 Sep 2026 13:08:58 GMT; HttpOnly; Secure; SameSite=Lax
== who came
no cookie:            a guest
with cookie:          signed in as aigul@example.kz
forged:               a guest
someone else's token: a guest
== signing out
signing out:          signed out
the same cookie:      a guest
```

## Taking it apart

### A cookie is a header, not magic

`http.SetCookie` writes a `Set-Cookie` header into the response, and the whole of it is there in the second line of the output. A browser that receives it saves the string and from then on adds a `Cookie` header to every request to this site.

That is all. A cookie needs no storage on the server and means nothing by itself: it is a string the browser carries back and forth.

> **Picture it.** A cloakroom ticket. There is no coat on it and no owner's name — only a number. The coat is with the attendant, and the number is how they find yours.

### Why `user_id` must not go into a cookie

The temptation is clear: put `user_id=1` in and read it on the server. The "forged" line of the output shows how that ends — or rather how it would end if the server believed it.

A cookie lives in the browser. Its owner can open the developer tools and write whatever they like: `user_id=1`, `admin=true`, `balance=1000000`. Nothing is signed, nothing is checked, and it takes two seconds.

The rule: **nothing the server intends to believe goes into a cookie**. What goes in is a ticket — a random string that means nothing on its own.

### Randomness comes in kinds

The token is made by `crypto/rand`, not `math/rand`, and that is not pedantry.

`math/rand` was built for something else: it is fast and predictable by design. Knowing a few of the numbers it gave out, the next one can be worked out. For shuffling a list that does not matter; for the ticket to somebody's account it matters a great deal.

`crypto/rand` takes its randomness from the operating system and exists for exactly this. Thirty-two bytes of it are so many possibilities that guessing is pointless on any hardware.

`base64.RawURLEncoding` turns the bytes into a string that can go in a header: 32 bytes make 43 characters, with no `+`, `/` or `=` to spoil a link.

### The database keeps the hash, not the token

Look at the `sessions` table: it holds `token_hash`, not `token`.

The reason is the one from the password lesson. If the database leaks with live tokens in it, whoever took it signs in as anybody without guessing a thing. A hash does not allow that: signing in needs the token itself, and it cannot be got out of a hash.

There is one difference from a password: `sha256` is enough here, not `bcrypt`. A password is invented by a person and may well be `qwerty123` — only slowness saves it from being guessed. The token was invented by us out of 32 random bytes; guessing it is pointless, so there is nothing to pay slowness for.

### Three words in the header

Three words stand in the `Set-Cookie` of the output, and each closes a hole of its own:

- **HttpOnly** — the cookie is invisible to JavaScript. If a stranger's script gets onto the page, it cannot read the token and send it away. This is the main defence against a session stolen through XSS;
- **Secure** — the browser sends the cookie over HTTPS only. Without it the token travels in plain text across any café wifi. It cannot be nailed to on: while the blog runs on `http://127.0.0.1` a cookie marked that way is accepted and never sent back — you sign in and the next page has forgotten you. So the flag follows the request: over TLS it is set, over plain http on your own machine it is not;
- **SameSite=Lax** — the cookie is not attached to cross-site `POST`s, to requests for images and styles, or to background calls. On an ordinary link followed from another site — a top-level `GET` — it is attached, and that is deliberate, so that you are not signed out every time you open the blog from a search. So `Lax` closes the forged form and does not replace a token check: there will be a separate lesson on CSRF.

Plus `Path=/` — for the whole site, and `Expires` — how long the browser keeps it. Without `Expires` a cookie lives until the browser closes.

### Signing out is deleting a row

Notice the last two lines of the output: after signing out **the very same cookie** no longer works.

That is how it should be. Asking the browser to forget a cookie is politeness, not a defence: a copy of the string could be anywhere. Signing out for real is `delete from sessions`, after which the token is found by nobody, whoever brings it.

That is also what the sessions table is good for: a person can be shown the list of their sign-ins and given a "sign out everywhere" button, which is simply deleting all their rows.

### What this lesson still leaves out

Every sign-in issues a new token — the code does exactly that — so a ticket slipped to you in advance does not become authorised afterwards. Session fixation is closed here.

What remains is the older sessions: sign in from somebody else's laptop, change your password, and that row in the table still works. A "sign out everywhere" button has to delete them, and that is not done yet.

Expired rows are cleaned up by nobody: a request will not find them, but they pile up. A real blog runs `delete from sessions where expires_at < ...` once a day.

And the signed-in person is looked up afresh in every handler. The place for that is a wrapper that does it once and puts what it found into the request context; context gets a lesson of its own.

## The lesson map

![Lesson map: a ticket in the browser, everything else in the database](/static/course/go/map-sessions-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why can `user_id` not go into a cookie and be believed?
2. Why is the token in the database hashed, but not with `bcrypt` the way a password is?
3. Why is signing out a deleted row rather than a request to the browser to forget the cookie?

## Exercise

**Required.** Add signing in and out to the blog. A migration creates the `sessions` table. `POST /login` checks the pair from the last lesson and hands out a cookie; `POST /logout` deletes the row. The header shows who is signed in and a link to sign out, and a guest gets a link to sign in.

All of it is done in [step-22](https://github.com/DauletBai/shanraq.org/tree/main/course/go-blog/step-22) — compare against it once you have written your own.

**Optional.**

- Hand out a fresh token on every sign-in and delete the old one. Check that the old cookie stops working.
- Show a list of sessions on the profile page, with times and a "sign out everywhere" button.
- Take `HttpOnly` out and look in the developer tools at what `document.cookie` now shows.

## Where this goes in your blog

The blog can tell people apart for the first time. Everything after this is built on it: permissions, other people's drafts, "my articles".

Still open. Sessions are neither renewed nor cleaned up. The signed-in person is looked up by every handler. And the sign-in form is not protected from guessing yet — no delay and no limit on attempts.

## Answers

1. Because the cookie sits in a person's browser and they can write whatever they like in it in two seconds. Only what the server keeps for itself can be believed; what goes into the cookie is a meaningless ticket.
2. Because `bcrypt` is for a secret a person invented, which may be weak: slowness saves it from guessing. The token was invented by us out of 32 random bytes; guessing it is pointless, so a fast `sha256` is enough.
3. Because a copy of the cookie could be anywhere, and a request to the browser binds nothing. While the row is alive the token works; delete the row and it is dead for everyone at once.

## Sources

- [Go: the http package and cookies](https://pkg.go.dev/net/http#Cookie)
- [Go: the crypto/rand package](https://pkg.go.dev/crypto/rand)
- [MDN: Set-Cookie and its attributes](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Set-Cookie)
