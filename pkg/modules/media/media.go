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
	"net/url"

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
		svg, err := fs.ReadFile(web.StaticFS(), "brand/shanraq-mark-light.svg")
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
		r.Use(sameOriginOnly)
		r.Post("/media/upload", m.handleUpload)
		r.Post("/media/upload-doc", m.handleUploadDoc)
	})
}

// sameOriginOnly rejects cross-origin uploads. These two endpoints authenticate
// by session cookie, which makes them the browser surface, and the browser
// surface is where CSRF lives: without this, a form on any other site could
// spend a logged-in author's storage quota with a file of the attacker's
// choosing, and land it under our domain to be served back to readers.
//
// The cookie is SameSite=Lax, so a cross-site POST does not carry it and the
// attack already fails — this is the same second, explicit layer the articles
// module has had (see its verifyOrigin). Media was simply the one cookie-authed
// surface left without it, and an inconsistency in a defence is how the defence
// gets lost during the next refactor.
func sameOriginOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		ok := false
		if o := r.Header.Get("Origin"); o != "" && o != "null" {
			if u, err := url.Parse(o); err == nil {
				ok = u.Host == host
			}
		} else if ref := r.Header.Get("Referer"); ref != "" {
			// No Origin (older clients): fall back to the Referer host.
			if u, err := url.Parse(ref); err == nil {
				ok = u.Host == host
			}
		}
		if !ok {
			writeJSONError(w, http.StatusForbidden, "cross-origin request blocked")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// allowUpload is the gate both upload endpoints share. A valid signature is not
// enough here: these two routes are the only ones a session can use to spend
// something we cannot get back, so they ask whether the account still exists
// and still holds the session it claims, and how recently it last asked.
//
// The cookie is a JWT with no session row to delete, so a deleted or suspended
// account kept uploading until its token expired — two hours of writes from
// someone the site no longer knows.
func (m *Module) allowUpload(w http.ResponseWriter, r *http.Request) bool {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "authentication required")
		return false
	}
	if !m.auth.SessionStillValid(r.Context(), claims) {
		writeJSONError(w, http.StatusUnauthorized, "authentication required")
		return false
	}
	if !m.auth.AllowUpload(r, claims.Subject) {
		w.Header().Set("Retry-After", "60")
		writeJSONError(w, http.StatusTooManyRequests, "too many uploads; try again shortly")
		return false
	}
	return true
}

type uploadResponse struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

func (m *Module) handleUpload(w http.ResponseWriter, r *http.Request) {
	if !m.allowUpload(w, r) {
		return
	}

	limit := m.cfg.MaxUploadBytes
	if limit <= 0 {
		limit = 25 << 20
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

// handleUploadDoc accepts a listing document — a PDF (floor plan, technical
// passport, contract) stored as-is, or an image plan/scheme that goes through
// the normal image pipeline. Same auth and size limits as image upload.
func (m *Module) handleUploadDoc(w http.ResponseWriter, r *http.Request) {
	if !m.allowUpload(w, r) {
		return
	}
	limit := m.cfg.MaxUploadBytes
	if limit <= 0 {
		limit = 25 << 20
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

	// A PDF is stored verbatim; anything else is treated as an image plan and
	// runs through the same processing as listing photos.
	var (
		data        = raw
		contentType = "application/pdf"
		ext         = ".pdf"
	)
	if len(raw) < 5 || string(raw[:5]) != "%PDF-" {
		img, perr := m.processImage(raw)
		if perr != nil {
			writeJSONError(w, http.StatusUnsupportedMediaType, "unsupported file (PDF, JPG or PNG)")
			return
		}
		data, contentType, ext = img, "image/jpeg", ".jpg"
	}

	sum := sha256.Sum256(data)
	h := hex.EncodeToString(sum[:])
	key := h[:2] + "/" + h + ext
	if err := m.store.Put(r.Context(), key, data, contentType); err != nil {
		m.logger.Error("media put doc", zap.Error(err))
		writeJSONError(w, http.StatusInternalServerError, "storage error")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(uploadResponse{URL: m.store.URL(key), Key: key})
}

// ProcessAndSaveAvatar decodes a raw uploaded image, center-crops it to a
// square (no watermark), stores it, and returns the public URL. Other modules
// (the studio cabinet) call this to let a user set their profile photo.
func (m *Module) ProcessAndSaveAvatar(ctx context.Context, raw []byte) (string, error) {
	data, err := m.processAvatar(raw)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	h := hex.EncodeToString(sum[:])
	key := "avatar/" + h[:2] + "/" + h + ".jpg"
	if err := m.store.Put(ctx, key, data, "image/jpeg"); err != nil {
		return "", fmt.Errorf("store avatar: %w", err)
	}
	return m.store.URL(key), nil
}

// MaxUploadBytes is the configured upload size cap (for callers that read the
// multipart body themselves, like the avatar endpoint), with a sane default.
func (m *Module) MaxUploadBytes() int64 {
	if m.cfg.MaxUploadBytes > 0 {
		return m.cfg.MaxUploadBytes
	}
	return 10 << 20
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
