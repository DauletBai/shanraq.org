// Step 22 — after the lesson "Sessions and cookies".
//
// The blog can tell people apart between requests. What the browser carries is
// a random token and nothing else; the database holds its hash, the person it
// belongs to and when it runs out.
//
// Signing out deletes the row rather than asking the browser nicely, so the
// same cookie stops working for everyone holding a copy of it.
package main

import (
	"bytes"
	"embed"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"strings"
	"time"
	"unicode/utf8"

	"myblog/blog"
)

//go:embed templates static
var files embed.FS

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

// config holds everything the blog can be told from outside. Assembled once at
// start-up; after that it is only read.
type config struct {
	addr   string
	dbPath string
	dev    bool
}

// env is the value of the variable, or the fallback when there is no such
// variable. LookupEnv rather than Getenv, so an empty value stays empty.
func env(name, fallback string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	return fallback
}

func load() config {
	var c config
	flag.StringVar(&c.addr, "addr", env("BLOG_ADDR", ":8080"), "мекенжай мен порт")
	flag.StringVar(&c.dbPath, "db", env("BLOG_DB_PATH", "blog.db"), "SQLite дерекқор файлы")
	flag.BoolVar(&c.dev, "dev", env("BLOG_DEV", "") == "1", "әзірлеуге арналған толық жауаптар")
	flag.Parse()
	return c
}

// reqID reads the number back out of the header set by withID. The request
// context is where this belongs, and that is a later lesson.
func reqID(w http.ResponseWriter) string { return w.Header().Get("X-Request-Id") }

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

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &recorder{ResponseWriter: w, code: http.StatusOK}

		next.ServeHTTP(rec, r)

		logger.Info("сұраныс", "id", reqID(w), "әдіс", r.Method, "жол", r.URL.Path,
			"код", rec.code, "уақыт", time.Since(start).Round(time.Millisecond).String())
	})
}

// currentUser answers who sent this request. Every handler asks for itself
// for now; the right place for this is a wrapper that does it once and puts
// the answer in the request context, and that is a later lesson.
func currentUser(store *blog.Store, r *http.Request) (blog.User, bool) {
	c, err := r.Cookie(blog.CookieName)
	if err != nil {
		return blog.User{}, false
	}
	u, err := store.UserByToken(c.Value)
	if err != nil {
		return blog.User{}, false
	}
	return u, true
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
	// Шапка одна на все страницы, поэтому вошедшего кладут здесь, а не в
	// каждом обработчике.
	if u, ok := currentUser(app.store, r); ok {
		data["User"] = u
	}
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
func seed(store *blog.Store) error {
	n, err := store.Len()
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
		if err := store.AddWithTags(s.article, s.tags); err != nil {
			return err
		}
	}
	return nil
}

// app — то, что нужно обработчикам: хранилище и всё, что появится дальше.
// Раньше это передавали параметром, но теперь их стало двое, и метод удобнее.
type app struct {
	store *blog.Store
}

func routes(store *blog.Store) http.Handler {
	app := &app{store: store}
	mux := http.NewServeMux()

	// Every page is handed a map, and render adds who is signed in. Draft
	// carries back what the reader typed: making someone retype eighty letters
	// is rudeness, not validation.
	showList := func(w http.ResponseWriter, r *http.Request, code int, p map[string]any) {
		list, err := store.All(fmt.Sprint(p["Sort"]))
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
		n, err := store.Len()
		if err != nil {
			logger.Error("санау", "id", reqID(w), "қате", err)
			serverError(w)
			return
		}
		slug := fmt.Sprintf("maqala-%d", n+1)
		a := blog.Article{Slug: slug, Title: title, Lang: "kz",
			Body: draftBody, Words: len(strings.Fields(draftBody))}

		// One call, two tables, one transaction inside the store.
		switch err := store.AddWithTags(a, tags); {
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
		found, err := store.Search(q)
		if err != nil {
			logger.Error("іздеу", "id", reqID(w), "сұрау", q, "қате", err)
			serverError(w)
			return
		}
		app.render(w, r, http.StatusOK, "search.html", map[string]any{"Query": q, "Found": found})
	})

	mux.HandleFunc("GET /tag/{name}", func(w http.ResponseWriter, r *http.Request) {
		tag := r.PathValue("name")
		list, err := store.ByTag(tag)
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

		id, err := store.Authenticate(email, r.FormValue("password"))
		if err != nil {
			// One sentence for a wrong password and for an address nobody
			// registered: anything else hands over a list of who is here.
			app.render(w, r, http.StatusUnauthorized, "login.html",
				map[string]any{"Email": email, "Err": "пошта немесе құпиясөз келмейді"})
			return
		}

		token, err := store.StartSession(id)
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
			if err := store.EndSession(c.Value); err != nil {
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

		switch err := store.Register(email, password); {
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
		if id, err := store.Authenticate(email, password); err == nil {
			if token, err := store.StartSession(id); err == nil {
				setSessionCookie(w, token)
			}
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	mux.HandleFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
		app.render(w, r, http.StatusOK, "about.html", nil)
	})

	mux.HandleFunc("GET /read/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		a, err := store.Get(slug)
		switch {
		case errors.Is(err, blog.ErrNotFound):
			app.render(w, r, http.StatusNotFound, "notfound.html", map[string]any{"Slug": slug})
			return
		case err != nil:
			logger.Error("оқу", "id", reqID(w), "жол", slug, "қате", err)
			serverError(w)
			return
		}
		tags, err := store.TagsOf(slug)
		if err != nil {
			logger.Error("тегтер", "id", reqID(w), "жол", slug, "қате", err)
			serverError(w)
			return
		}
		app.render(w, r, http.StatusOK, "article.html", map[string]any{"Article": a, "Tags": tags})
	})

	mux.HandleFunc("GET /edit/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		a, err := store.Get(slug)
		switch {
		case errors.Is(err, blog.ErrNotFound):
			app.render(w, r, http.StatusNotFound, "notfound.html", map[string]any{"Slug": slug})
			return
		case err != nil:
			logger.Error("оқу", "id", reqID(w), "жол", slug, "қате", err)
			serverError(w)
			return
		}
		app.render(w, r, http.StatusOK, "edit.html", map[string]any{"Article": a})
	})

	mux.HandleFunc("POST /edit/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		title := strings.TrimSpace(r.FormValue("title"))

		// The same two checks as the add form, and for the same reason: what
		// comes from outside is checked where it arrives.
		if title == "" || utf8.RuneCountInString(title) > 80 {
			a, err := store.Get(slug)
			if err != nil {
				app.render(w, r, http.StatusNotFound, "notfound.html", map[string]any{"Slug": slug})
				return
			}
			a.Title = title
			app.render(w, r, http.StatusBadRequest, "edit.html",
				map[string]any{"Article": a, "Err": "тақырып бос немесе тым ұзын"})
			return
		}

		// Update itself reports a missing article: it asks the database how
		// many rows it changed, and zero means there was nothing to change.
		switch err := store.Update(slug, title); {
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

	// No StripPrefix: the embedded tree already begins with static/, so the
	// request path and the path inside files are the same.
	mux.Handle("GET /static/", http.FileServerFS(files))

	return mux
}

func main() {
	c := load()

	store, err := blog.Open(c.dbPath)
	if err != nil {
		logger.Error("дерекқор", "қате", err)
		os.Exit(1)
	}
	defer store.Close()

	if err := seed(store); err != nil {
		logger.Error("бастапқы дерек", "қате", err)
		os.Exit(1)
	}

	// logging sits outside recoverPanic on purpose: a panic unwinds past the
	// logging wrapper before it can write its line, so recovering first means
	// every request gets exactly one request line, panics included.
	handler := withID(logging(recoverPanic(routes(store))))

	logger.Info("сервер іске қосылды", "мекенжай", c.addr, "дерекқор", c.dbPath, "әзірлеу", c.dev)
	if err := http.ListenAndServe(c.addr, handler); err != nil {
		logger.Error("сервер тоқтады", "қате", err)
		os.Exit(1)
	}
}
