# Registration: what is kept instead of a password

_Лид (summary):_ **The thirty-sixth lesson of the Go course. A password is not kept in the database; a bcrypt hash is, with the salt inside it and the cost outside. Measured: sha256 does nine million hashes a second, bcrypt does ten. Plus the 72-byte limit a Kazakh password runs into at its thirty-seventh letter.**

## Why this matters

Start from the assumption that the database will leak one day. Not necessarily through a break-in: a forgotten copy on a laptop, somebody else's access to the server, a mistake in a setting. The question is not whether it happens but what exactly ends up in a stranger's hands.

If the column holds `qwerty123`, they have your blog and the mail of half your readers at once — people reuse passwords. If it holds `$2a$10$...`, they have a list of meaningless strings, and working passwords out of them costs months of work and a bill for electricity.

The difference between those two cases is one library and one lesson.

## The whole thing at once

A new folder, `go mod init sabaq33`. Besides the database driver, one more library is needed:

```
go get golang.org/x/crypto/bcrypt
```

```go
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
	sqlite "modernc.org/sqlite"
)

const uniqueViolation = 2067

// What the rest of the program sees. It knows nothing about bcrypt or about
// the numbers the database answers with.
var (
	ErrEmailTaken = errors.New("the address is taken")
	ErrBadLogin   = errors.New("address or password does not fit")
	ErrShort      = errors.New("password shorter than eight characters")
	ErrLong       = errors.New("password longer than 72 bytes")
)

// dummyHash is a real hash of a string nobody knows. It makes an address that
// does not exist cost the same time as one that does: otherwise a stopwatch
// would let an outsider collect the addresses that are registered.
const dummyHash = "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"

func main() {
	os.Remove("blog.db")
	db, err := sql.Open("sqlite", "blog.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(`create table users (
		id            integer primary key,
		email         text not null unique,
		password_hash text not null,
		created_at    text not null default (datetime('now'))
	) strict`); err != nil {
		log.Fatal(err)
	}

	fmt.Println("== registration")
	fmt.Println("first time:     ", register(db, "Aigul@Example.KZ", "secretpass123"))
	fmt.Println("same address:   ", register(db, "aigul@example.kz", "anotherpass"))
	fmt.Println("too short:      ", register(db, "bota@example.kz", "1234"))
	fmt.Println("too long:       ", register(db, "bota@example.kz", strings.Repeat("қ", 37)))

	fmt.Println("== what is in the database")
	var hash string
	db.QueryRow(`select password_hash from users where email = ?`, "aigul@example.kz").Scan(&hash)
	fmt.Println(hash)

	fmt.Println("== signing in")
	fmt.Println("right password: ", timed(func() error { return login(db, "aigul@example.kz", "secretpass123") }))
	fmt.Println("wrong password: ", timed(func() error { return login(db, "aigul@example.kz", "guessing") }))
	fmt.Println("no such address:", timed(func() error { return login(db, "joq@example.kz", "guessing") }))
}

// register files a person. The password never reaches the database; its hash does.
func register(db *sql.DB, email, password string) error {
	email = strings.ToLower(strings.TrimSpace(email))

	// The length is checked before hashing: characters for a person, bytes for
	// bcrypt, whose limit is exactly 72 of them.
	if utf8.RuneCountInString(password) < 8 {
		return ErrShort
	}
	if len(password) > 72 {
		return ErrLong
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec(`insert into users (email, password_hash) values (?, ?)`,
		email, string(hash))
	var serr *sqlite.Error
	if errors.As(err, &serr) && serr.Code() == uniqueViolation {
		return ErrEmailTaken
	}
	return err
}

// login answers a wrong password and an unknown address in the same way.
func login(db *sql.DB, email, password string) error {
	email = strings.ToLower(strings.TrimSpace(email))

	var hash string
	err := db.QueryRow(`select password_hash from users where email = ?`, email).Scan(&hash)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		bcrypt.CompareHashAndPassword([]byte(dummyHash), []byte(password))
		return ErrBadLogin
	case err != nil:
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return ErrBadLogin
	}
	return nil
}

// timed prints the answer together with how long it took.
func timed(f func() error) string {
	start := time.Now()
	err := f()
	return fmt.Sprintf("%v (%s)", err, time.Since(start).Round(10*time.Millisecond))
}
```

The output:

```
== registration
first time:      <nil>
same address:    the address is taken
too short:       password shorter than eight characters
too long:        password longer than 72 bytes
== what is in the database
$2a$10$jkez1EKdb2rCq8kj4VJ3UefVFdRi2U50jN8BEe3QeU/tm8k6KVZl2
== signing in
right password:  <nil> (100ms)
wrong password:  address or password does not fit (100ms)
no such address: address or password does not fit (100ms)
```

Your hash will be different — that is not a fault but the main property of bcrypt, which comes below.

## Taking it apart

### A hash is not a cipher

A cipher can be deciphered: whoever holds the key has a way back. A hash has no way back — it is a function that grinds an input into a string of a fixed length and can restore nothing.

Hence the way it is checked. The password at sign-in is not "deciphered"; it goes through the same function and the result is compared with what the database holds. A match means the password is right. The system never needs to know the password itself again, **ever**.

> **Picture it.** A mincer. Mince is easy to make; a piece of meat cannot be got back out of it. But if you put a second piece through and the mince comes out the same, the meat was the same.

Hence the way to spot a bad site: if it can send you your old password, it stores it in a way anyone reaching the database can read.

### Why `sha256` will not do

The first thing that comes to mind is `sha256` out of the standard library. It is indeed a hash. And it will not do, for two reasons, both measurable.

```
sha256 twice:
  68e8e213d72156bb64a8cfb39c2ec4d7750af7577c8306f3aefcdf19bd84945e
  68e8e213d72156bb64a8cfb39c2ec4d7750af7577c8306f3aefcdf19bd84945e
bcrypt twice:
  $2a$10$BUxXmib0/QAPQvlPrVrErOM9Ntx7p1vbw2N/eP6de2Khafr.k0QQS
  $2a$10$TalADwpj40zTG.rGlMmG.uJA4E8BSipe7avq1T3KKnSq6LyBz69gu
per second: sha256 — 9441422 hashes
per second: bcrypt — 10 hashes
```

The first: **the same password gives the same hash**. So in a leaked database it is immediately visible whose passwords match, and the "hash → word" tables for billions of popular passwords have long been computed for us.

The second: **speed**. Over nine million hashes a second on an ordinary laptop — and on a graphics card the count runs into billions. A dictionary of a million popular passwords is worked through faster than you read this paragraph.

For passwords a fast hash is a defect. And this is the one place in programming where we pay money for slowness.

### bcrypt: the salt inside, the cost outside

A word about the choice first. OWASP today recommends **Argon2id** for new systems and calls bcrypt the fallback — for places where Argon2 or scrypt is not available. We take bcrypt for two reasons: it is in `golang.org/x/crypto` with no extra dependency, and it has none of the four tuning parameters a beginner gets wrong. For a teaching blog that is enough; for a system with real money in it, read about Argon2id.

`bcrypt` solves both troubles at once.

The same password gives **different** hashes: look at the output above again. That is the salt at work — random bytes mixed into the password and stored right there, inside the string. A precomputed table is useless: it would have to be recomputed for every salt.

The string is built like this:

```
$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy
 │   │  └── the salt and the hash itself
 │   └───── the cost: 10
 └───────── the version of the algorithm
```

The cost is a power of two: `10` means 2¹⁰ passes. Every extra unit doubles the time. Mine came out at ten hashes a second against nine million for `sha256` — nearly a million times slower, and that is exactly what we are paying for.

Checking a hash is what `bcrypt.CompareHashAndPassword` does: it takes the version, the cost and the salt out of the string itself, so nothing extra has to be stored.

The cost is chosen for the actual server rather than from a number in an article: enough that checking a password takes roughly 0.1 to 0.5 seconds. In three years the hardware gets faster, the cost goes up by one — and the old hashes keep working, because the cost is written into every one of them.

### 72 bytes: the limit Kazakh runs into

bcrypt has a hard limit on the length of a password, and it is counted **in bytes**:

```
72 letters,  72 bytes -> <nil>
73 letters,  73 bytes -> bcrypt: password length exceeds 72 bytes
36 letters,  72 bytes -> <nil>
37 letters,  74 bytes -> bcrypt: password length exceeds 72 bytes
```

In Latin script that is 72 characters, which nobody runs into. But a Kazakh or Russian letter takes two bytes — the same two bytes as in the lesson on runes — and the limit arrives at the **37th letter**.

Libraries used to cut the password off at 72 bytes in silence, and someone with a long password would sign in with its beginning for years. Go now returns an error, which is better: cutting a password silently means lying about how strong it is.

In the blog it means one thing: check `len(password) > 72` before hashing and say something a person can act on rather than "internal server error".

### The same answer to a wrong password and an unknown address

Look at the last three lines of the output. All three answers are the same in words **and** in time.

The words are the same on purpose. If an unknown address gets "no such user" while a known one with a wrong password gets "that password is wrong", anyone can collect a list of your readers in an evening: feed in addresses and watch the answer. That is called user enumeration, and the defence is one shared sentence.

Time is harder. With no address there is nothing to compare against, so the function would return at once — while on an existing one it would spend its hundred milliseconds. The difference is visible with a stopwatch, and the list gets collected anyway. So in the "no such address" branch we compare the password against a hash nobody owns: the result is not wanted, the time spent is.

### What to check in a password, and what not to

The "at least one capital, one digit and one symbol" rules are out of date, and since 2017 they are not recommended even by NIST, the organisation that introduced them. People answer such rules identically: `Password1!`. The complexity is invented and the memory suffers.

What does help:

- **length**. NIST's 2024 revision asks for **at least fifteen characters** from a password that stands alone; eight are allowed only beside a second factor. Set the upper limit no lower than 64 and count characters rather than bytes: a long phrase of three words beats a short jumble;
- **a check against the leaked lists**. There are open databases of known passwords, and keeping one of those out is worth more than any symbols;
- **not forcing a change without a reason**. A compulsory change every three months produces `Qazaq2026`, then `Qazaq2027`.

And one thing people forget: a password must never go into a log. Not in full, and not "the first three characters, for debugging".

### What goes into the `users` table

Exactly as many columns as are needed:

```sql
create table users (
	id            integer primary key,
	email         text not null unique,
	password_hash text not null,
	created_at    text not null default (datetime('now'))
) strict
```

The `email` is lower-cased before the insert: `Aigul@Example.KZ` and `aigul@example.kz` are one person, and `unique` has to understand that. The output shows the second registration on that address being turned away, though it was written differently.

There is no `password` column in this table and there never will be. There is only `password_hash`, and the name says what is inside.

## The lesson map

![Lesson map: a password, a salt and a cost](/static/course/go/map-passwords-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why does a site that can send you your old password store it wrongly?
2. Why does bcrypt need to be slow when every other function tries to be fast?
3. Why does the program compare a password against a hash even for an address that does not exist?

## Exercise

**Required.** Add registration to the blog. A migration adds the `users` table. `GET /register` shows the form and `POST /register` files a person: the address lower-cased, the password from eight characters and no longer than 72 bytes, and only the hash into the database. A taken address answers with a sentence a person can act on rather than a five hundred.

All of it is done in [step-21](https://github.com/DauletBai/shanraq.org/tree/main/course/go-blog/step-21) — compare against it once you have written your own.

**Optional.**

- Measure how long `GenerateFromPassword` takes at cost 10, 12 and 14 on your machine. Pick the one where the check fits into 0.1–0.5 seconds — and remember the number will be different on a server.
- Write a test: register, sign in with the right password, sign in with the wrong one. The third case must return `ErrBadLogin` and nothing else.
- Try registering with a password of 37 Kazakh letters and make sure a person sees a message they can understand.

## Where this goes in your blog

The blog has its first person. So far they can only register — signing in, which remembers something, is the next lesson, on sessions.

Still open. Nobody confirms the address, so somebody else's can still be used. The number of sign-in attempts is not limited. And we check no list of leaked passwords — all of it is genuinely needed, and each pulls a lesson of its own.

## Answers

1. Because sending a password back is only possible if it is stored readably — in plain text, or encrypted with a key lying next to it. A hash cannot be sent: no password comes out of it.
2. Because speed here works for whoever is guessing. A fast function lets millions of candidates be tried per second; a slow one turns a dictionary run into months. One check in a hundredth of a second costs us nothing, and costs the guesser everything.
3. So that the time of the answer is the same. Otherwise a fast answer would mean "no such address", and a stopwatch would be enough to collect the list of people registered.

## Sources

- [Go: the bcrypt package](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [OWASP: how to store passwords](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
- [NIST SP 800-63B: password guidance](https://pages.nist.gov/800-63-3/sp800-63b.html)
