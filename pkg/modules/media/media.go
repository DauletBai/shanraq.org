package media

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io"
	"io/fs"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
	ledger *ledger
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
	if rt.DB != nil {
		m.ledger = &ledger{db: rt.DB}
	}

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
func (m *Module) allowUpload(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "authentication required")
		return uuid.Nil, false
	}
	if !m.auth.SessionStillValid(r.Context(), claims) {
		writeJSONError(w, http.StatusUnauthorized, "authentication required")
		return uuid.Nil, false
	}
	owner, err := uuid.Parse(claims.Subject)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "authentication required")
		return uuid.Nil, false
	}
	if !m.auth.AllowUpload(r, claims.Subject) {
		w.Header().Set("Retry-After", "60")
		writeJSONError(w, http.StatusTooManyRequests, "too many uploads; try again shortly")
		return uuid.Nil, false
	}
	return owner, true
}

// reserve decides whether these bytes may be stored. It is checked before the
// write and recorded after it, which leaves a window where two simultaneous
// uploads can both be admitted: a quota is a bound on how much an account can
// accumulate, not a ledger that has to balance to the byte, and the upload rate
// limiter caps how far past it a burst can get.
func (m *Module) reserve(ctx context.Context, owner uuid.UUID, key string, size int64) error {
	if m.ledger == nil {
		return nil // no database wired (tests, and the fs-only smoke image)
	}
	// Storing what you already have costs nothing and must not be charged.
	if held, err := m.ledger.held(ctx, owner, key); err != nil {
		return err
	} else if held {
		return nil
	}
	if quota := m.cfg.QuotaBytes; quota > 0 {
		used, err := m.ledger.usage(ctx, owner)
		if err != nil {
			return err
		}
		if used+size > quota {
			return ErrQuotaExceeded
		}
	}
	if cap := m.cfg.MaxTotalBytes; cap > 0 {
		// An object already on disk under someone else's name adds no bytes.
		known, err := m.ledger.known(ctx, key)
		if err != nil {
			return err
		}
		if !known {
			total, err := m.ledger.total(ctx)
			if err != nil {
				return err
			}
			if total+size > cap {
				return ErrStoreFull
			}
		}
	}
	return nil
}

// keep files the stored object against its owner, and takes the bytes back off
// disk if it cannot.
//
// The failure used to be logged and swallowed, on the reasoning that the bytes
// were already written and a reader's page should not break over an accounting
// row. But the sweep only ever removes keys it has a row for, so a file with no
// row is a file nothing will ever remove -- it is charged to nobody, counted in
// no quota, and stays until the disk fills. Losing an upload is a retry;
// leaking one is permanent.
//
// So the object is deleted and the caller told, which turns a silent leak into
// a visible failure the author can act on.
func (m *Module) keep(ctx context.Context, owner uuid.UUID, key string, size int64, contentType string) error {
	if m.ledger == nil {
		return nil
	}
	if err := m.ledger.record(ctx, key, size, contentType, owner); err != nil {
		m.logger.Warn("media ledger record", zap.Error(err), zap.String("key", key))
		// Best effort in its turn: if the delete fails too there is nothing
		// further to try, and the log now carries both halves of the story.
		if derr := m.store.Delete(ctx, key); derr != nil {
			m.logger.Error("orphaned media object left on disk",
				zap.Error(derr), zap.String("key", key))
		}
		return err
	}
	return nil
}

// refuse maps a storage refusal onto a status. Being over quota is the caller's
// problem to fix by deleting something; the store being full is ours, and
// saying so plainly beats a generic error that reads like a bug.
func (m *Module) refuse(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrQuotaExceeded):
		writeJSONError(w, http.StatusRequestEntityTooLarge, "storage quota reached; delete some files first")
	case errors.Is(err, ErrStoreFull):
		m.logger.Error("media store is full")
		writeJSONError(w, http.StatusInsufficientStorage, "storage is full; the site operator has been alerted")
	default:
		m.logger.Error("media reserve", zap.Error(err))
		writeJSONError(w, http.StatusInternalServerError, "storage error")
	}
}

type uploadResponse struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

func (m *Module) handleUpload(w http.ResponseWriter, r *http.Request) {
	owner, ok := m.allowUpload(w, r)
	if !ok {
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

	if err := m.reserve(r.Context(), owner, key, int64(len(data))); err != nil {
		m.refuse(w, err)
		return
	}
	if err := m.store.Put(r.Context(), key, data, "image/jpeg"); err != nil {
		m.logger.Error("media put", zap.Error(err))
		writeJSONError(w, http.StatusInternalServerError, "storage error")
		return
	}
	if err := m.keep(r.Context(), owner, key, int64(len(data)), "image/jpeg"); err != nil {
		writeJSONError(w, http.StatusServiceUnavailable, "could not record the upload; try again")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(uploadResponse{URL: m.store.URL(key), Key: key})
}

// handleUploadDoc accepts a listing document — a PDF (floor plan, technical
// passport, contract) stored as-is, or an image plan/scheme that goes through
// the normal image pipeline. Same auth and size limits as image upload.
func (m *Module) handleUploadDoc(w http.ResponseWriter, r *http.Request) {
	owner, ok := m.allowUpload(w, r)
	if !ok {
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
	} else if bad := activePDF(raw); bad != "" {
		m.logger.Warn("active pdf refused", zap.String("construct", bad))
		writeJSONError(w, http.StatusUnsupportedMediaType,
			"this PDF contains active content and cannot be published; export it as a plain document")
		return
	}

	sum := sha256.Sum256(data)
	h := hex.EncodeToString(sum[:])
	key := h[:2] + "/" + h + ext
	if err := m.reserve(r.Context(), owner, key, int64(len(data))); err != nil {
		m.refuse(w, err)
		return
	}
	if err := m.store.Put(r.Context(), key, data, contentType); err != nil {
		m.logger.Error("media put doc", zap.Error(err))
		writeJSONError(w, http.StatusInternalServerError, "storage error")
		return
	}
	if err := m.keep(r.Context(), owner, key, int64(len(data)), contentType); err != nil {
		writeJSONError(w, http.StatusServiceUnavailable, "could not record the upload; try again")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(uploadResponse{URL: m.store.URL(key), Key: key})
}

// ProcessAndSaveAvatar decodes a raw uploaded image, center-crops it to a
// square (no watermark), stores it, and returns the public URL. Other modules
// (the studio cabinet) call this to let a user set their profile photo.
func (m *Module) ProcessAndSaveAvatar(ctx context.Context, owner uuid.UUID, raw []byte) (string, error) {
	data, err := m.processAvatar(raw)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	h := hex.EncodeToString(sum[:])
	key := "avatar/" + h[:2] + "/" + h + ".jpg"
	if err := m.reserve(ctx, owner, key, int64(len(data))); err != nil {
		return "", err
	}
	if err := m.store.Put(ctx, key, data, "image/jpeg"); err != nil {
		return "", fmt.Errorf("store avatar: %w", err)
	}
	if err := m.keep(ctx, owner, key, int64(len(data)), "image/jpeg"); err != nil {
		return "", err
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
	shanraq.StarterModule
} = (*Module)(nil)

// activePDF names the first construct that makes a PDF more than a document, or
// "" when it is only pages.
//
// This is not virus scanning and does not pretend to be: there is no signature
// database here and a novel payload would pass. It refuses the categories that
// make a PDF dangerous at all -- a document that runs script on open, launches
// a program, submits a form somewhere, or carries another file inside it. A
// reader's browser will happily honour every one of those, and an article
// attachment has no use for any of them.
//
// The check reads the raw bytes rather than parsing the document. A parser is
// the thing being defended against, and a name split across an object stream
// costs a false refusal, which an author can fix by exporting the file again --
// where the other kind of mistake costs a reader.
func activePDF(raw []byte) string {
	// Names as they appear in a PDF's own syntax, including the hex escaping a
	// producer may use for any character of a name.
	for _, c := range []struct {
		label string
		forms []string
	}{
		{"JavaScript", []string{"/JavaScript", "/JS", "/J#61vaScript"}},
		{"an action on open", []string{"/OpenAction", "/AA"}},
		{"a program launch", []string{"/Launch"}},
		{"an embedded file", []string{"/EmbeddedFile", "/Filespec"}},
		{"a form submission", []string{"/SubmitForm"}},
		{"a remote reference", []string{"/RichMedia", "/GoToR"}},
	} {
		for _, f := range c.forms {
			if bytes.Contains(raw, []byte(f)) {
				return c.label
			}
		}
	}
	return ""
}
