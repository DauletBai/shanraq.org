// Step 14 — after the lesson "Configuration".
//
// The address stops being written into the code. It comes from a flag, or from
// the environment, or from the default here — in that order, because the
// variable is what the flag falls back to.
//
// The database path is not here yet: the blog has no database until the next
// module. When it arrives, -db and BLOG_DB_PATH go in this struct beside the
// address, and nowhere else.
package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"
	"unicode/utf8"

	"myblog/blog"
)

//go:embed templates/*.html static
var files embed.FS

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

// config holds everything the blog can be told from outside. Assembled once at
// start-up; after that it is only read.
type config struct {
	addr string
	dev  bool
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

// Parsed once at start: a broken template stops the server now rather than
// surprising a reader later.
var tpl = template.Must(template.ParseFS(files, "templates/*.html"))

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

// render builds the page in full before a single byte goes out.
func render(w http.ResponseWriter, code int, name string, data any) {
	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, name, data); err != nil {
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

// front is what list.html asks for. Draft carries back what the reader typed,
// because making someone retype eighty letters is rudeness, not validation.
type front struct {
	Articles []blog.Article
	Draft    string
	Err      string
}

func seed() *blog.Store {
	store := blog.NewStore()
	store.Add(blog.Article{Slug: "salem", Title: "Сәлем, әлем", Words: 150, Lang: "kz",
		Body: "Бірінші жазба. Блог осыдан басталады."})
	store.Add(blog.Article{Slug: "dala", Title: "Дала туралы", Words: 400, Lang: "kz",
		Body: "Дала жазда сары, көктемде жасыл болады."})
	store.Add(blog.Article{Slug: "shanyraq", Title: "Шаңырақ деген не", Words: 1000, Lang: "kz",
		Body: "Шаңырақ — киіз үйдің төбесіндегі шеңбер."})
	return store
}

func routes(store *blog.Store) http.Handler {
	mux := http.NewServeMux()

	showList := func(w http.ResponseWriter, code int, p front) {
		p.Articles = store.All()
		render(w, code, "list.html", p)
	}

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		showList(w, http.StatusOK, front{})
	})

	mux.HandleFunc("POST /add", func(w http.ResponseWriter, r *http.Request) {
		title := strings.TrimSpace(r.FormValue("title"))

		switch {
		case title == "":
			showList(w, http.StatusBadRequest, front{Err: "тақырып бос"})
			return
		case utf8.RuneCountInString(title) > 80:
			showList(w, http.StatusBadRequest, front{Draft: title, Err: "тақырып тым ұзын"})
			return
		}

		// Making an address out of a Kazakh title is a job for a package we
		// have not met yet, so the number does for now.
		slug := fmt.Sprintf("maqala-%d", store.Len()+1)
		store.Add(blog.Article{Slug: slug, Title: title, Lang: "kz",
			Body: draftBody, Words: len(strings.Fields(draftBody))})

		// A handler that changed something does not draw a page: it sends the
		// browser to one that can be refreshed as often as the reader likes.
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	mux.HandleFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
		render(w, http.StatusOK, "about.html", nil)
	})

	mux.HandleFunc("GET /read/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		a, err := store.Get(slug)
		if err != nil {
			render(w, http.StatusNotFound, "notfound.html", map[string]any{"Slug": slug})
			return
		}
		render(w, http.StatusOK, "article.html", map[string]any{"Article": a})
	})

	// No StripPrefix: the embedded tree already begins with static/, so the
	// request path and the path inside files are the same.
	mux.Handle("GET /static/", http.FileServerFS(files))

	return mux
}

func main() {
	// logging sits outside recoverPanic on purpose: a panic unwinds past the
	// logging wrapper before it can write its line, so recovering first means
	// every request gets exactly one request line, panics included.
	handler := withID(logging(recoverPanic(routes(seed()))))

	c := load()

	logger.Info("сервер іске қосылды", "мекенжай", c.addr, "әзірлеу", c.dev)
	if err := http.ListenAndServe(c.addr, handler); err != nil {
		logger.Error("сервер тоқтады", "қате", err)
		os.Exit(1)
	}
}
