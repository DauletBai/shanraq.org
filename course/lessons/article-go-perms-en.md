# Permissions: who may change what

_Лид (summary):_ **The thirty-eighth lesson of the Go course. Hiding a button is not a guard: it is measured here how a stranger edits your article with a direct request while never seeing the button. The check lives in the handler and in one place, and the refusals differ: a guest is sent to sign in, a stranger gets 403, and what does not exist gets 404.**

## Why this matters

In the last lesson the blog learned to recognise whoever signed in. But recognising is half the job: what remains is deciding what that person may do.

The difference between two words is worth remembering. **Authentication** answers "who are you" — that is signing in with an address and a password. **Authorisation** answers something else: "what may you do" — and that is this lesson.

The rule in a blog is simple: anyone may read; only the person who wrote an article may change or delete it.

## The whole thing at once

A new folder, `go mod init sabaq35`. The program starts a server twice: first the way it is usually done at the start, then the way it should be.

```go
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	_ "modernc.org/sqlite"
)

// ErrNotAllowed means the person is signed in but the article is not theirs.
var ErrNotAllowed = errors.New("this article is not yours")

func main() {
	os.Remove("blog.db")
	db := open()
	defer db.Close()

	fmt.Println("== the button is hidden, nothing is checked")
	show(db, false)

	fmt.Println()
	fmt.Println("== the handler checks")
	show(db, true)
}

func show(db *sql.DB, guarded bool) {
	reset(db)
	srv := httptest.NewServer(routes(db, guarded))
	defer srv.Close()

	fmt.Println("owner sees the button:", hasButton(srv.URL, "aigul"))
	fmt.Println("stranger sees it:     ", hasButton(srv.URL, "bota"))

	fmt.Println("owner edits:          ", edit(srv.URL, "aigul", "Steppe: the owner edited"))
	fmt.Println("stranger edits:       ", edit(srv.URL, "bota", "Steppe: a stranger edited"))
	fmt.Println("guest edits:          ", edit(srv.URL, "", "Steppe: a guest edited"))
	fmt.Println("in the database after:", title(db, "step"))
}

func routes(db *sql.DB, guarded bool) http.Handler {
	mux := http.NewServeMux()

	// The article page: only the owner sees the edit link.
	mux.HandleFunc("GET /read/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		owner, err := ownerOf(db, slug)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "<h1>%s</h1>", title(db, slug))
		if who(r) == owner {
			fmt.Fprint(w, `<a href="/edit/`+slug+`">edit</a>`)
		}
	})

	mux.HandleFunc("POST /edit/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")

		if guarded {
			if err := mayEdit(db, who(r), slug); err != nil {
				answer(w, err)
				return
			}
		}

		if _, err := db.Exec(`update articles set title = ? where slug = ?`,
			r.FormValue("title"), slug); err != nil {
			http.Error(w, "did not work", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusSeeOther)
	})

	return mux
}

// mayEdit is the only place where it is decided who may do what. Two such
// places would drift apart one day.
func mayEdit(db *sql.DB, user, slug string) error {
	if user == "" {
		return ErrNoUser
	}
	owner, err := ownerOf(db, slug)
	if err != nil {
		return err
	}
	if owner != user {
		return ErrNotAllowed
	}
	return nil
}

// answer turns a permission error into an HTTP answer. A guest is sent to
// sign in, a stranger is refused, and what does not exist stays that way.
func answer(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNoUser):
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusSeeOther)
	case errors.Is(err, ErrNotAllowed):
		http.Error(w, "this article is not yours", http.StatusForbidden)
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}

func ownerOf(db *sql.DB, slug string) (string, error) {
	var owner string
	err := db.QueryRow(`select u.name from articles a
		join users u on u.id = a.author_id
		where a.slug = ?`, slug).Scan(&owner)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNoArticle
	}
	if err != nil {
		log.Fatal(err)
	}
	return owner, nil
}
```

The helpers live in a second file, `helpers.go`: it creates two people and an article with an owner, and sends requests as different people.

The output:

```
== the button is hidden, nothing is checked
owner sees the button: yes
stranger sees it:      no
owner edits:           303 
stranger edits:        303 
guest edits:           303 
in the database after: Steppe: a guest edited

== the handler checks
owner sees the button: yes
stranger sees it:      no
owner edits:           303 
stranger edits:        403 this article is not yours
guest edits:           303 → /login
in the database after: Steppe: the owner edited
```

## Taking it apart

### Hiding a button is not a guard

Look at the first half of the output. The stranger does not see the edit button — the template honestly did not show it. And the article was still edited, by a **guest**, who is not even signed in.

Because a button is a hint to the browser, not a lock. The request can be sent without it: from the address bar, from the browser console, from a one-line `curl`. Nothing has to be broken into.

> **Picture it.** A "no entry" sign on a door with no lock. It warns honestly, and it stops nobody who decides to walk in.

This mistake has a name — an insecure direct object reference: the address is known, no check is made, and everything belonging to somebody else is open to whoever substitutes their number. It has been near the top of the list of common holes for years.

The rule: **whatever the template decides, the handler must decide too**. The template hides the button for convenience; the handler checks the right for safety.

### The check lives in one place

In the second half of the output `mayEdit` appears, and the result is different: the stranger is refused, and what stays in the database is the owner's edit.

What matters is not only that the check exists but that there is one of it:

```go
func mayEdit(db *sql.DB, user, slug string) error {
	if user == "" {
		return ErrNoUser
	}
	owner, err := ownerOf(db, slug)
	if err != nil {
		return err
	}
	if owner != user {
		return ErrNotAllowed
	}
	return nil
}
```

When there are two such places — in editing and in deleting — they drift apart one day: one learns about editors and the other does not. A single function called from everywhere prevents that.

Notice what it returns: not a `bool` but an error. `true` and `false` do not tell "sign in" from "not yours", and the handler needs to tell them apart.

### Three different refusals

`answer` turns the error into a response, and all three cases differ:

- **a guest — 303 to `/login`**. They did nothing wrong; they simply need to sign in. The output shows exactly where they were sent;
- **a stranger — 403**. They are signed in, we know who they are, and the answer means "you are recognised and you may not". Asking them to sign in again is pointless;
- **something missing — 404**. There is no article, and permissions have nothing to do with it.

The difference between 401 and 403 is what confuses people most. `401 Unauthorized` means "identify yourself", though the name suggests the opposite; `403 Forbidden` means "identified, and still no". Browser applications usually send a person to the sign-in page instead of a 401, as here.

There is a fourth approach worth knowing: answering `404` instead of `403` for somebody else's thing. Then an outsider does not even learn that such an article exists. For a blog that is too much — the reader sees all the articles anyway — but for other people's drafts or private messages it is the right choice.

### The check comes after the thing is found

The order in `mayEdit` is not accidental: first "are they signed in", then "does the article exist", and only then "is it theirs".

Swap the last two and something odd appears: somebody else's non-existent article answers 403 instead of 404, and that can be used to work out which addresses are taken. A small thing, and small things like it are what a picture is assembled from.

### The database knows the owner, not the form

`ownerOf` asks the database. That is the only right source: anything arriving with the request can be forged, as the cookie lesson showed.

There is one link in the schema: an article has an `author_id` pointing at `users`. Checking a permission is an ordinary `join`, and permissions in a blog need neither a table of their own nor a clever system.

When there are more roles — an editor, a moderator — a roles table appears, and a check of the form "the owner **or** an editor". But there is no need to start there: complexity is worth adding when it can no longer be avoided.

### A wrapper for pages a guest may not see

Checking for a sign-in in every handler is the same trouble as checking a permission in two places. That is what the wrappers from the middleware lesson are for:

```go
func requireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if who(r) == "" {
			http.Redirect(w, r, "/login?next="+url.QueryEscape(r.URL.Path), http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
```

The wrapper answers "are they signed in at all", and `mayEdit` answers "is this theirs". The first is the same for every closed page; the second depends on the article, which is why they live apart.

The `next` in the link is a small courtesy: after signing in a person is returned to where they were going rather than to the front page. Just check afterwards that `next` leads to your own site: the address came from outside, so it cannot be trusted.

## The lesson map

![Lesson map: signed in, found, and theirs](/static/course/go/map-perms-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why does a hidden button not protect an article?
2. How does 403 differ from 401, and when is 404 answered instead of 403?
3. Why does `mayEdit` return an error rather than `true`/`false`?

## Exercise

**Required.** Add permissions to the blog. A migration adds `author_id` to the articles and fills it in for the ones already there. The form files an article in the name of whoever is signed in. Editing and deleting check the owner: a guest is sent to sign in, a stranger gets a 403. Only the owner still sees the edit button — but that is convenience now, not a guard.

All of it is done in [step-23](https://github.com/DauletBai/shanraq.org/tree/main/course/go-blog/step-23) — compare against it once you have written your own.

**Optional.**

- Take the check out of the handler, leaving it only in the template, and edit somebody else's article with `curl` again. Then put the check back.
- Add `next` to the sign-in page and check that a person lands where they were going.
- Make somebody else's draft answer `404` rather than `403`, and explain why that is better for a draft.

## Where this goes in your blog

From this lesson on, the blog stops being one person's: every article has an owner, and what belongs to somebody else is untouchable.

Debts. There are no roles — only "owner or not". There is no limit on sign-in attempts either. And the forms can still be submitted from another site: that is the CSRF hole, and the next lesson is about exactly that.

## Answers

1. Because the button lives in the browser and the request can be sent without it — from the address bar, from the console, or from `curl`. The template hides the button for convenience; the handler is what decides.
2. `401` means "identify yourself" — the person is not recognised; `403` means "recognised, and still no". Browser applications usually send a person to the sign-in page instead of a `401`. A `404` in place of a `403` is for when the very existence of the thing is a secret: other people's drafts, private messages.
3. Because a `bool` carries no reason, and the handler has to answer differently: send a guest to sign in, give a stranger a 403, answer 404 for what is missing. An error carries the reason; `true`/`false` does not.

## Sources

- [MDN: the 403 status](https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/403)
- [OWASP: broken access control](https://owasp.org/Top10/A01_2021-Broken_Access_Control/)
