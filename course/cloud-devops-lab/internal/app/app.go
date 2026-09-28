package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type Config struct {
	Addr        string
	ServiceName string
	Environment string
	DataFile    string
	Version     string
}

func ConfigFromEnv(version string) (Config, error) {
	cfg := Config{
		Addr:        env("CLOUDLAB_ADDR", ":8080"),
		ServiceName: env("CLOUDLAB_NAME", "CloudLab"),
		Environment: env("CLOUDLAB_ENV", "local"),
		DataFile:    env("CLOUDLAB_DATA", "data/notes.json"),
		Version:     version,
	}
	if !strings.HasPrefix(cfg.Addr, ":") {
		return Config{}, fmt.Errorf("CLOUDLAB_ADDR must look like :8080")
	}
	if cfg.ServiceName == "" || cfg.Environment == "" {
		return Config{}, fmt.Errorf("name and environment must not be empty")
	}
	return cfg, nil
}

func env(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}

type noteStore struct {
	mu    sync.RWMutex
	path  string
	notes []string
}

func loadStore(path string) (*noteStore, error) {
	store := &noteStore{path: path}
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &store.notes); err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}
	return store, nil
}

func (s *noteStore) list() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]string(nil), s.notes...)
}

func (s *noteStore) add(note string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	next := append(append([]string(nil), s.notes...), note)
	b, err := json.MarshalIndent(next, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o750); err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, append(b, '\n'), 0o640); err != nil {
		return err
	}
	if err := os.Rename(tmp, s.path); err != nil {
		return err
	}
	s.notes = next
	return nil
}

type service struct {
	cfg      Config
	log      *slog.Logger
	store    *noteStore
	requests atomic.Uint64
	ready    atomic.Bool
}

var page = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="en"><meta charset="utf-8"><meta name="viewport" content="width=device-width">
<title>{{.Name}}</title><style>body{font:18px system-ui;max-width:760px;margin:3rem auto;padding:0 1rem}code{background:#eee;padding:.15rem .35rem}li{margin:.5rem 0}</style>
<h1>{{.Name}}</h1><p>Environment: <code>{{.Environment}}</code> · Version: <code>{{.Version}}</code></p>
<p>This is the service operated throughout the Cloud & DevOps course.</p>
<h2>Deployment notes</h2><ul>{{range .Notes}}<li>{{.}}</li>{{else}}<li>No notes yet.</li>{{end}}</ul>
<form method="post" action="/notes"><label>New note <input name="note" maxlength="120" required></label> <button>Add</button></form>
<p><a href="/healthz">health</a> · <a href="/readyz">readiness</a> · <a href="/metrics">metrics</a></p></html>`))

func Run(cfg Config, logger *slog.Logger) error {
	store, err := loadStore(cfg.DataFile)
	if err != nil {
		return err
	}
	s := &service{cfg: cfg, log: logger, store: store}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.home)
	mux.HandleFunc("POST /notes", s.addNote)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	mux.HandleFunc("GET /readyz", s.readiness)
	mux.HandleFunc("GET /metrics", s.metrics)

	server := &http.Server{Addr: cfg.Addr, Handler: s.accessLog(mux), ReadHeaderTimeout: 5 * time.Second, IdleTimeout: 30 * time.Second}
	s.ready.Store(true)
	logger.Info("service starting", "address", cfg.Addr, "environment", cfg.Environment, "version", cfg.Version)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		s.ready.Store(false)
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdown); err != nil {
			logger.Error("graceful shutdown", "error", err)
		}
	}()
	err = server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (s *service) home(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = page.Execute(w, struct {
		Name, Environment, Version string
		Notes                      []string
	}{s.cfg.ServiceName, s.cfg.Environment, s.cfg.Version, s.store.list()})
}

func (s *service) addNote(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	note := strings.TrimSpace(r.FormValue("note"))
	if note == "" || len([]rune(note)) > 120 {
		http.Error(w, "note must contain 1–120 characters", http.StatusBadRequest)
		return
	}
	if err := s.store.add(note); err != nil {
		s.log.Error("save note", "error", err)
		http.Error(w, "cannot save note", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *service) readiness(w http.ResponseWriter, _ *http.Request) {
	if !s.ready.Load() {
		http.Error(w, "not ready", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *service) metrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	_, _ = fmt.Fprintf(w, "# TYPE cloudlab_http_requests_total counter\ncloudlab_http_requests_total %s\n", strconv.FormatUint(s.requests.Load(), 10))
}

func (s *service) accessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		s.requests.Add(1)
		next.ServeHTTP(w, r)
		s.log.Info("http request", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(started).Milliseconds())
	})
}
