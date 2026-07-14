package media

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"shanraq.org/internal/config"
	"shanraq.org/pkg/modules/auth"
	"shanraq.org/pkg/shanraq"
	"shanraq.org/web"
)

// Module owns user-uploaded media: it processes images (resize, watermark,
// EXIF strip), stores them behind a pluggable Store, and serves them. Video is
// handled via third-party embeds elsewhere, not stored here.
type Module struct {
	auth   *auth.Module
	cfg    config.MediaConfig
	store  Store
	mark   *image.RGBA
	maxDim int
	logger *zap.Logger
}

// New builds the media module. It depends on auth to gate the upload endpoint.
func New(authModule *auth.Module) *Module { return &Module{auth: authModule} }

func (m *Module) Name() string { return "media" }

// Init selects the storage backend and rasterizes the brand watermark once.
func (m *Module) Init(_ context.Context, rt *shanraq.Runtime) error {
	m.cfg = rt.Config.Media
	m.logger = rt.Logger
	m.maxDim = m.cfg.MaxDimension

	switch m.cfg.Backend {
	case "", "fs":
		store, err := NewFSStore(m.cfg.Dir, m.cfg.PublicPrefix)
		if err != nil {
			return fmt.Errorf("media: init fs store: %w", err)
		}
		m.store = store
	default:
		return fmt.Errorf("media: unsupported backend %q", m.cfg.Backend)
	}

	if m.cfg.Watermark {
		svg, err := fs.ReadFile(web.StaticFS(), "brand/shanraq-mark.svg")
		if err != nil {
			return fmt.Errorf("media: read brand watermark: %w", err)
		}
		mark, err := rasterizeSVG(svg, watermarkPx, watermarkPx)
		if err != nil {
			return fmt.Errorf("media: rasterize watermark: %w", err)
		}
		m.mark = mark
	}
	return nil
}

// Routes serves stored objects and registers the auth-gated upload endpoint.
func (m *Module) Routes(r chi.Router) {
	if fsStore, ok := m.store.(*FSStore); ok {
		r.Handle(fsStore.Prefix()+"/*", fsStore.FileServer())
	}
	r.Group(func(r chi.Router) {
		r.Use(m.auth.LoadSession)
		r.Post("/media/upload", m.handleUpload)
	})
}

type uploadResponse struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

func (m *Module) handleUpload(w http.ResponseWriter, r *http.Request) {
	if _, ok := auth.ClaimsFromContext(r.Context()); !ok {
		writeJSONError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	limit := m.cfg.MaxUploadBytes
	if limit <= 0 {
		limit = 10 << 20
	}
	r.Body = http.MaxBytesReader(w, r.Body, limit)
	if err := r.ParseMultipartForm(limit); err != nil {
		writeJSONError(w, http.StatusRequestEntityTooLarge, "file too large")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "missing file")
		return
	}
	defer file.Close()

	raw, err := io.ReadAll(file)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "cannot read file")
		return
	}

	data, err := m.processImage(raw)
	if err != nil {
		// Decode failure means it was not a supported image.
		writeJSONError(w, http.StatusUnsupportedMediaType, "unsupported image")
		return
	}

	sum := sha256.Sum256(data)
	h := hex.EncodeToString(sum[:])
	key := h[:2] + "/" + h + ".jpg"

	if err := m.store.Put(r.Context(), key, data, "image/jpeg"); err != nil {
		m.logger.Error("media put", zap.Error(err))
		writeJSONError(w, http.StatusInternalServerError, "storage error")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(uploadResponse{URL: m.store.URL(key), Key: key})
}

func writeJSONError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

var _ interface {
	shanraq.Module
	shanraq.RouterModule
	shanraq.InitializerModule
} = (*Module)(nil)
