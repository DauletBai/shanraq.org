// Step 32 — after the lesson "A domain of your own and HTTPS".
//
// A proxy now stands in front of the blog and holds the certificate, which
// changes one thing inside: every connection arrives from 127.0.0.1, so
// r.RemoteAddr is the proxy for everybody. The real address comes in a header
// the proxy sets — and a header anyone can send, which is why clientIP trusts
// it only when the connection itself came from the proxy.
package main

import (
	"bytes"
	"context"
	cryptorand "crypto/rand"
	"crypto/subtle"
	"embed"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"math/rand/v2"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime/debug"
	"strings"
	"time"
	"unicode/utf8"

	"myblog/blog"
)

// The tree splits in two here. What is ours is embedded, so a binary is still
// one file. uploads is not: embed reads at build time, and what a reader will
// send does not exist yet.
//
//go:embed templates static/brand static/css static/js static/images
var files embed.FS

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

// config holds everything the blog can be told from outside. Assembled once at
// start-up; after that it is only read.
type config struct {
	addr    string
	dbPath  string
	uploads string
	name    string
	dev     bool
}

// env is the value of the variable, or the fallback when there is no such
// variable. LookupEnv rather than Getenv, so an empty value stays empty.
func env(name, fallback string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	return fallback
}

// defaultAddr is where a host has its say. A free plan hands the port over in
// PORT and expects the service on every interface: ":" with no host in front
// means exactly that, while "127.0.0.1:" would answer nobody but ourselves.
func defaultAddr() string {
	if p := os.Getenv("PORT"); p != "" {
		return ":" + p
	}
	return ":8080"
}

func load() config {
	var c config
	flag.StringVar(&c.addr, "addr", env("BLOG_ADDR", defaultAddr()), "мекенжай мен порт")
	flag.StringVar(&c.dbPath, "db", env("BLOG_DB_PATH", "blog.db"), "SQLite дерекқор файлы")
	flag.StringVar(&c.uploads, "uploads", env("BLOG_UPLOADS", "static/uploads"), "жүктелген суреттер қалтасы")
	flag.StringVar(&c.name, "name", env("BLOG_NAME", "Менің блогым"), "блогтың аты")
	flag.BoolVar(&c.dev, "dev", env("BLOG_DEV", "") == "1", "әзірлеуге арналған толық жауаптар")
	flag.Parse()
	return c
}

// reqID reads the number back out of the header set by withID. The request
// context is where this belongs, and that is a later lesson.
func reqID(w http.ResponseWriter) string { return w.Header().Get("X-Request-Id") }

// secure puts the three headers on every answer. One wrapper, the whole blog.
func secure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 'self' keeps a stranger's script from running even if it somehow
		// reached the markup. unsafe-inline is here only because our styles
		// are inline; a blog without them should drop it.
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' 'unsafe-inline'")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "same-origin")
		next.ServeHTTP(w, r)
	})
}

func withID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", fmt.Sprintf("%08x", rand.Uint32()))
		next.ServeHTTP(w, r)
	})
}

func recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				w.Header().Set("Connection", "close")
				logger.Error("паника", "id", reqID(w), "жол", r.URL.Path,
					"себебі", fmt.Sprint(v), "стек", string(debug.Stack()))
				serverError(w)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// serverError says one thing to every reader and keeps the detail in the log.
func serverError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Бізде бірдеңе бұзылды. Өтініш нөмірі: %s\n", reqID(w))
}

// One set per page, built once at start: a broken template stops the server
// now rather than surprising a reader later.
var tpl = pages()

// funcs are what a template may call by itself. The year in the footer is a
// good example: passing it in the data of every page would mean every handler
// remembering to, and forgetting once is enough for a footer saying 2026 in
// 2027.
var funcs = template.FuncMap{
	"year": func() int { return time.Now().Year() },
}

// pages puts the frame, every partial and exactly one page into each set. A
// single set for everything would not work: names inside a set are shared, and
// the last main parsed would win for every page.
func pages() map[string]*template.Template {
	found, err := fs.Glob(files, "templates/pages/*.html")
	if err != nil || len(found) == 0 {
		panic("беттер табылмады: " + fmt.Sprint(err))
	}
	out := make(map[string]*template.Template, len(found))
	for _, page := range found {
		name := path.Base(page)
		out[name] = template.Must(template.New(name).Funcs(funcs).ParseFS(files,
			"templates/base.html", "templates/partials/*.html", page))
	}
	return out
}

type recorder struct {
	http.ResponseWriter
	code int
}

func (rec *recorder) WriteHeader(code int) {
	rec.code = code
	rec.ResponseWriter.WriteHeader(code)
}

// clientIP is who to write in the log. Behind a proxy every connection comes
// from the proxy itself, and the reader's address arrives in X-Forwarded-For --
// a header any reader can send too. So it is believed only when the connection
// came from the loopback, where our proxy is; a request straight from the
// network is trusted for what it is, not for what it says about itself.
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		return host
	}
	// The header is a list: the reader first, then every proxy that added
	// itself. Ours is the first entry.
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		if first, _, ok := strings.Cut(fwd, ","); ok {
			return strings.TrimSpace(first)
		}
		return strings.TrimSpace(fwd)
	}
	return host
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &recorder{ResponseWriter: w, code: http.StatusOK}

		next.ServeHTTP(rec, r)

		logger.Info("сұраныс", "id", reqID(w), "әдіс", r.Method, "жол", r.URL.Path,
			"оқырман", clientIP(r),
			"код", rec.code, "уақыт", time.Since(start).Round(time.Millisecond).String())
	})
}

// The token travels in a cookie of its own, so that a guest has one too: the
// sign-in and registration forms are submitted before there is any session.
const (
	csrfCookie = "csrf"
	csrfField  = "csrf"
)

// csrfToken hands back this visitor's token, setting the cookie the first time.
// It is not HttpOnly-secret from the reader — it sits in their own form. It is
// a secret from another site, which cannot read a page of ours.
func csrfToken(w http.ResponseWriter, r *http.Request) string {
	if c, err := r.Cookie(csrfCookie); err == nil && c.Value != "" {
		return c.Value
	}
	raw := make([]byte, 32)
	if _, err := cryptorand.Read(raw); err != nil {
		return ""
	}
	token := base64.RawURLEncoding.EncodeToString(raw)
	http.SetCookie(w, &http.Cookie{
		Name: csrfCookie, Value: token, Path: "/",
		HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode,
	})
	return token
}

// checkCSRF guards every method that changes something. It sits in a wrapper
// rather than in the handlers, because a handler added later would be added
// without the check.
func checkCSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			next.ServeHTTP(w, r)
			return
		}

		c, err := r.Cookie(csrfCookie)
		sent := r.FormValue(csrfField)
		// Constant time: an ordinary == stops at the first byte that differs,
		// and the time it takes gives the token away character by character.
		if err != nil || sent == "" ||
			subtle.ConstantTimeCompare([]byte(c.Value), []byte(sent)) != 1 {
			logger.Info("csrf", "id", reqID(w), "жол", r.URL.Path)
			http.Error(w, "форма ескірген, бетті жаңартыңыз", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// denied turns an error from MayEdit into an answer. A guest is sent to sign
// in and brought back afterwards; a stranger is told no; anything else is a
// missing page.
func (app *app) denied(w http.ResponseWriter, r *http.Request, slug string, err error) {
	switch {
	case errors.Is(err, blog.ErrNeedLogin):
		http.Redirect(w, r, "/login?next="+url.QueryEscape(r.URL.Path), http.StatusSeeOther)
	case errors.Is(err, blog.ErrNotAllowed):
		app.render(w, r, http.StatusForbidden, "forbidden.html", map[string]any{"Slug": slug})
	default:
		app.render(w, r, http.StatusNotFound, "notfound.html", map[string]any{"Slug": slug})
	}
}

// currentUser now reads what the withUser wrapper has already found. No
// handler goes to the database for this any more.
func currentUser(store *blog.Store, r *http.Request) (blog.User, bool) {
	return userFrom(r.Context())
}

// setSessionCookie carries the token to the browser. HttpOnly keeps it away
// from JavaScript, Secure keeps it off plain http, SameSite keeps it out of
// requests another site starts.
func setSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     blog.CookieName,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(blog.SessionLife),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

// render builds the page in full before a single byte goes out. What is always
// executed is base: the page itself has neither html nor body, it only fills
// the frame's blocks.
func (app *app) render(w http.ResponseWriter, r *http.Request, code int, name string, data map[string]any) {
	if data == nil {
		data = map[string]any{}
	}
	// The name of the blog stands in the header, in the footer and in the
	// title of every page. It is added once here, so a new page cannot be
	// left with somebody else's name on it.
	data["SiteName"] = app.name
	// The header is the same on every page, so who is signed in is put in
	// here rather than remembered by each handler.
	if u, ok := currentUser(app.store, r); ok {
		data["User"] = u
	}
	// Every form needs the token, so it is put in beside the user rather than
	// remembered by each handler.
	data["CSRF"] = csrfToken(w, r)
	set, ok := tpl[name]
	if !ok {
		logger.Error("үлгі жоқ", "id", reqID(w), "аты", name)
		serverError(w)
		return
	}
	var buf bytes.Buffer
	if err := set.ExecuteTemplate(&buf, "base", data); err != nil {
		logger.Error("үлгі", "id", reqID(w), "аты", name, "қате", err)
		serverError(w)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	if _, err := buf.WriteTo(w); err != nil {
		logger.Error("жазу", "id", reqID(w), "қате", err)
	}
}

// draftBody stands in until the form learns to take the text as well — that
// is this lesson's exercise. Counting its words keeps the reading time honest
// instead of printing a zero.
const draftBody = "Мәтін әзірге жоқ."

// seed fills an empty database once, so a fresh checkout has something to show.
func seed(ctx context.Context, store *blog.Store) error {
	n, err := store.Len(ctx)
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	for _, s := range []struct {
		article blog.Article
		tags    []string
	}{
		{blog.Article{Slug: "salem", Title: "Сәлем, әлем", Words: 150, Lang: "kz",
			Body: "Бірінші жазба. Блог осыдан басталады."}, []string{"блог"}},
		{blog.Article{Slug: "dala", Title: "Дала туралы", Words: 400, Lang: "kz",
			Body: "Дала жазда сары, көктемде жасыл болады."}, []string{"дала", "табиғат"}},
		{blog.Article{Slug: "shanyraq", Title: "Шаңырақ деген не", Words: 1000, Lang: "kz",
			Body: "Шаңырақ — киіз үйдің төбесіндегі шеңбер."}, []string{"дала", "дәстүр"}},
	} {
		if err := store.AddWithTags(ctx, s.article, s.tags); err != nil {
			return err
		}
	}
	return nil
}

// app is what the handlers need: the store, and whatever joins it later. It
// used to travel as a parameter; with more than one thing to carry, a method
// reads better.
type app struct {
	store   *blog.Store
	uploads string
	name    string
	views   *views
}

func routes(store *blog.Store, c config) http.Handler {
	app := &app{store: store, uploads: c.uploads, name: c.name, views: newViews()}
	mux := http.NewServeMux()

	// Every page is handed a map, and render adds who is signed in. Draft
	// carries back what the reader typed: making someone retype eighty letters
	// is rudeness, not validation.
	showList := func(w http.ResponseWriter, r *http.Request, code int, p map[string]any) {
		list, err := store.All(r.Context(), fmt.Sprint(p["Sort"]))
		if err != nil {
			logger.Error("тізім", "id", reqID(w), "қате", err)
			serverError(w)
			return
		}
		p["Articles"] = list
		app.render(w, r, code, "list.html", p)
	}

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		// The key travels no further than blog.All, which checks it. Nothing
		// here has to trust it.
		showList(w, r, http.StatusOK, map[string]any{"Sort": r.URL.Query().Get("sort")})
	})

	mux.HandleFunc("POST /add", func(w http.ResponseWriter, r *http.Request) {
		// 32 KiB stays in memory; Go moves the rest into a temporary file by
		// itself and clears up afterwards.
		if err := r.ParseMultipartForm(32 << 10); err != nil {
			showList(w, r, http.StatusRequestEntityTooLarge,
				map[string]any{"Err": "файл тым үлкен"})
			return
		}
		title := strings.TrimSpace(r.FormValue("title"))
		tags := blog.ParseTags(r.FormValue("tags"))

		switch {
		case title == "":
			showList(w, r, http.StatusBadRequest, map[string]any{"Err": "тақырып бос"})
			return
		case utf8.RuneCountInString(title) > 80:
			showList(w, r, http.StatusBadRequest, map[string]any{"Draft": title, "Err": "тақырып тым ұзын"})
			return
		}

		// Making an address out of a Kazakh title is a job for a package we
		// have not met yet, so the number does for now.
		n, err := store.Len(r.Context())
		if err != nil {
			logger.Error("санау", "id", reqID(w), "қате", err)
			serverError(w)
			return
		}
		slug := fmt.Sprintf("maqala-%d", n+1)
		a := blog.Article{Slug: slug, Title: title, Lang: "kz",
			Body: draftBody, Words: len(strings.Fields(draftBody))}
		if u, ok := currentUser(store, r); ok {
			a.AuthorID = u.ID
		}

		// The cover is optional: no file is not an error.
		if file, _, err := r.FormFile("cover"); err == nil {
			defer file.Close()
			switch name, err := app.saveCover(file); {
			case errors.Is(err, errNotImage):
				showList(w, r, http.StatusUnsupportedMediaType,
					map[string]any{"Draft": title, "Err": "тек сурет: png, jpeg, gif, webp"})
				return
			case err != nil:
				logger.Error("сурет", "id", reqID(w), "қате", err)
				serverError(w)
				return
			default:
				a.Cover = name
			}
		}

		// One call, two tables, one transaction inside the store.
		switch err := store.AddWithTags(r.Context(), a, tags); {
		case errors.Is(err, blog.ErrTaken):
			showList(w, r, http.StatusBadRequest, map[string]any{
				"Draft": title, "DraftTag": r.FormValue("tags"), "Err": "мұндай мекенжай бар"})
			return
		case err != nil:
			logger.Error("қосу", "id", reqID(w), "қате", err)
			serverError(w)
			return
		}

		// A handler that changed something does not draw a page: it sends the
		// browser to one that can be refreshed as often as the reader likes.
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	mux.HandleFunc("GET /search", func(w http.ResponseWriter, r *http.Request) {
		q := strings.TrimSpace(r.URL.Query().Get("q"))
		found, err := store.Search(r.Context(), q)
		if err != nil {
			logger.Error("іздеу", "id", reqID(w), "сұрау", q, "қате", err)
			serverError(w)
			return
		}
		app.render(w, r, http.StatusOK, "search.html", map[string]any{"Query": q, "Found": found})
	})

	mux.HandleFunc("GET /tag/{name}", func(w http.ResponseWriter, r *http.Request) {
		tag := r.PathValue("name")
		list, err := store.ByTag(r.Context(), tag)
		if err != nil {
			logger.Error("тег", "id", reqID(w), "тег", tag, "қате", err)
			serverError(w)
			return
		}
		app.render(w, r, http.StatusOK, "tag.html", map[string]any{"Articles": list, "Tag": tag})
	})

	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		app.render(w, r, http.StatusOK, "login.html", map[string]any{})
	})

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		email := blog.NormalizeEmail(r.FormValue("email"))

		id, err := store.Authenticate(r.Context(), email, r.FormValue("password"))
		if err != nil {
			// One sentence for a wrong password and for an address nobody
			// registered: anything else hands over a list of who is here.
			app.render(w, r, http.StatusUnauthorized, "login.html",
				map[string]any{"Email": email, "Err": "пошта немесе құпиясөз келмейді"})
			return
		}

		token, err := store.StartSession(r.Context(), id)
		if err != nil {
			logger.Error("сессия", "id", reqID(w), "қате", err)
			serverError(w)
			return
		}
		setSessionCookie(w, token)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	mux.HandleFunc("POST /logout", func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie(blog.CookieName); err == nil {
			if err := store.EndSession(r.Context(), c.Value); err != nil {
				logger.Error("шығу", "id", reqID(w), "қате", err)
			}
		}
		// The row is already gone; this only asks the browser to forget.
		http.SetCookie(w, &http.Cookie{
			Name: blog.CookieName, Value: "", Path: "/",
			Expires: time.Unix(0, 0), MaxAge: -1,
			HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode,
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	mux.HandleFunc("GET /register", func(w http.ResponseWriter, r *http.Request) {
		app.render(w, r, http.StatusOK, "register.html", map[string]any{})
	})

	mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		email := blog.NormalizeEmail(r.FormValue("email"))
		password := r.FormValue("password")

		// The address comes back to the form; the password never does. Making
		// someone retype an address is rude, keeping their password in the
		// page is careless.
		fail := func(msg string) {
			app.render(w, r, http.StatusBadRequest, "register.html",
				map[string]any{"Email": email, "Err": msg})
		}

		switch {
		case !strings.Contains(email, "@"):
			fail("пошта дұрыс емес")
			return
		case utf8.RuneCountInString(password) < blog.MinPassword:
			fail("құпиясөз сегіз таңбадан қысқа")
			return
		case len(password) > blog.MaxPasswordBytes:
			fail("құпиясөз тым ұзын: 72 байт, қазақша 36 әріп")
			return
		}

		switch err := store.Register(r.Context(), email, password); {
		case errors.Is(err, blog.ErrEmailTaken):
			fail("бұл пошта тіркелген")
			return
		case err != nil:
			logger.Error("тіркеу", "id", reqID(w), "қате", err)
			serverError(w)
			return
		}

		logger.Info("жаңа адам", "id", reqID(w), "пошта", email)

		// Just registered means just signed in: making someone type the same
		// pair again immediately is rudeness, not security.
		if id, err := store.Authenticate(r.Context(), email, password); err == nil {
			if token, err := store.StartSession(r.Context(), id); err == nil {
				setSessionCookie(w, token)
			}
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	mux.HandleFunc("GET /api/articles", func(w http.ResponseWriter, r *http.Request) {
		list, err := store.All(r.Context(), r.URL.Query().Get("sort"))
		if err != nil {
			logger.Error("api тізім", "id", reqID(w), "қате", err)
			apiError(w, r, http.StatusInternalServerError, "ішкі қате")
			return
		}

		out := make([]apiArticle, 0, len(list))
		for _, a := range list {
			tags, err := store.TagsOf(r.Context(), a.Slug)
			if err != nil {
				logger.Error("api тегтер", "id", reqID(w), "қате", err)
				apiError(w, r, http.StatusInternalServerError, "ішкі қате")
				return
			}
			out = append(out, newAPIArticle(a, tags))
		}
		writeJSON(w, r, http.StatusOK, map[string]any{"articles": out})
	})

	mux.HandleFunc("GET /api/articles/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")

		a, err := store.Get(r.Context(), slug)
		switch {
		case errors.Is(err, blog.ErrNotFound):
			apiError(w, r, http.StatusNotFound, "мұндай мақала жоқ")
			return
		case err != nil:
			logger.Error("api оқу", "id", reqID(w), "жол", slug, "қате", err)
			apiError(w, r, http.StatusInternalServerError, "ішкі қате")
			return
		}

		tags, err := store.TagsOf(r.Context(), slug)
		if err != nil {
			logger.Error("api тегтер", "id", reqID(w), "қате", err)
			apiError(w, r, http.StatusInternalServerError, "ішкі қате")
			return
		}
		writeJSON(w, r, http.StatusOK, newAPIArticle(a, tags))
	})

	mux.HandleFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
		app.render(w, r, http.StatusOK, "about.html",
			map[string]any{"Views": app.views.Total()})
	})

	mux.HandleFunc("GET /read/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		a, err := store.Get(r.Context(), slug)
		switch {
		case errors.Is(err, blog.ErrNotFound):
			app.render(w, r, http.StatusNotFound, "notfound.html", map[string]any{"Slug": slug})
			return
		case err != nil:
			logger.Error("оқу", "id", reqID(w), "жол", slug, "қате", err)
			serverError(w)
			return
		}
		app.views.Add(slug)

		tags, err := store.TagsOf(r.Context(), slug)
		if err != nil {
			logger.Error("тегтер", "id", reqID(w), "жол", slug, "қате", err)
			serverError(w)
			return
		}
		u, _ := currentUser(store, r)
		app.render(w, r, http.StatusOK, "article.html", map[string]any{
			"Article": a, "Tags": tags, "Views": app.views.Count(slug),
			"MayEdit": store.MayEdit(r.Context(), u, slug) == nil})
	})

	mux.HandleFunc("GET /edit/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		u, _ := currentUser(store, r)
		if err := store.MayEdit(r.Context(), u, slug); err != nil {
			app.denied(w, r, slug, err)
			return
		}

		a, err := store.Get(r.Context(), slug)
		if err != nil {
			app.render(w, r, http.StatusNotFound, "notfound.html", map[string]any{"Slug": slug})
			return
		}
		app.render(w, r, http.StatusOK, "edit.html", map[string]any{"Article": a})
	})

	mux.HandleFunc("POST /edit/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		u, _ := currentUser(store, r)

		// The same call as on the page before it. Hiding the link is a
		// convenience; this line is the guard.
		if err := store.MayEdit(r.Context(), u, slug); err != nil {
			app.denied(w, r, slug, err)
			return
		}

		title := strings.TrimSpace(r.FormValue("title"))
		if title == "" || utf8.RuneCountInString(title) > 80 {
			a, err := store.Get(r.Context(), slug)
			if err != nil {
				app.render(w, r, http.StatusNotFound, "notfound.html", map[string]any{"Slug": slug})
				return
			}
			a.Title = title
			app.render(w, r, http.StatusBadRequest, "edit.html",
				map[string]any{"Article": a, "Err": "тақырып бос немесе тым ұзын"})
			return
		}

		switch err := store.Update(r.Context(), slug, title); {
		case errors.Is(err, blog.ErrNotFound):
			app.render(w, r, http.StatusNotFound, "notfound.html", map[string]any{"Slug": slug})
			return
		case err != nil:
			logger.Error("өңдеу", "id", reqID(w), "жол", slug, "қате", err)
			serverError(w)
			return
		}

		http.Redirect(w, r, "/read/"+slug, http.StatusSeeOther)
	})

	// The host asks for this one, not a reader. It touches neither the
	// database nor a template on purpose: a page that can break is no use for
	// deciding whether the service is alive.
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("ok"))
	})

	// A browser asks for /favicon.ico by itself, whatever the markup says.
	// Without this line every reader leaves a 404 in the log, so the request
	// is sent on to the mark we already have.
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/brand/logo.svg", http.StatusMovedPermanently)
	})

	// No StripPrefix: the embedded tree already begins with static/, so the
	// request path and the path inside files are the same.
	mux.Handle("GET /static/", http.FileServerFS(files))

	// Uploads come off the disk rather than out of the binary. nosniff from
	// the previous lesson matters here: without it a browser may decide our
	// "png" is HTML and run it.
	mux.Handle("GET /static/uploads/", http.StripPrefix("/static/uploads/",
		http.FileServer(http.Dir(app.uploads))))

	return app.withUser(mux)
}

func main() {
	c := load()

	store, err := blog.Open(c.dbPath)
	if err != nil {
		logger.Error("дерекқор", "қате", err)
		os.Exit(1)
	}
	defer store.Close()

	if err := seed(context.Background(), store); err != nil {
		logger.Error("бастапқы дерек", "қате", err)
		os.Exit(1)
	}

	// logging sits outside recoverPanic on purpose: a panic unwinds past the
	// logging wrapper before it can write its line, so recovering first means
	// every request gets exactly one request line, panics included.
	handler := withID(secure(logging(recoverPanic(checkCSRF(
		http.MaxBytesHandler(routes(store, c), maxUpload))))))

	logger.Info("сервер іске қосылды", "мекенжай", c.addr, "дерекқор", c.dbPath, "әзірлеу", c.dev)
	if err := http.ListenAndServe(c.addr, handler); err != nil {
		logger.Error("сервер тоқтады", "қате", err)
		os.Exit(1)
	}
}
