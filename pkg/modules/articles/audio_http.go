package articles

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// maxNarrationBytes caps one upload. Eighteen minutes of speech at 48 kbit/s is
// about six megabytes, so this leaves room for a long article at a generous
// bitrate and still refuses anything that is obviously not narration.
const maxNarrationBytes = 64 << 20

// narrationTypes maps the audio containers we accept to the extension the file
// is stored under. Deliberately short: these are the formats every browser
// plays, and a format no browser plays is not narration, it is a bug report
// from the future.
var narrationTypes = map[string]string{
	"audio/mpeg":  ".mp3",
	"audio/mp4":   ".m4a",
	"audio/aac":   ".m4a",
	"audio/ogg":   ".ogg",
	"audio/opus":  ".opus",
	"audio/wav":   ".wav",
	"audio/x-wav": ".wav",
}

// UseAPIKeyAuth supplies the middleware that authenticates machine callers.
//
// Narration is produced off this server -- there is no speech synthesiser here
// and no reason to put a 122 MB voice model on a small VPS -- so the generator
// runs on someone's laptop and posts the finished audio. That caller has no
// browser session to present, which is what the API key is for. Without this
// wired the upload route is not registered at all: an unauthenticated writer of
// article audio is worse than no narration.
func (m *Module) UseAPIKeyAuth(mw func(http.Handler) http.Handler) { m.apiAuth = mw }

// routeNarration registers the machine-facing audio endpoints.
func (m *Module) routeNarration(r chi.Router) {
	if m.apiAuth == nil || m.auth == nil {
		return
	}
	r.Group(func(r chi.Router) {
		r.Use(m.apiAuth)
		r.Use(m.auth.RequireRoles("operator", "admin"))
		r.Put("/api/articles/{id}/audio/{lang}", m.handleNarrationPut)
		r.Delete("/api/articles/{id}/audio/{lang}", m.handleNarrationDelete)
	})
}

// handleNarrationPut stores one uploaded reading.
//
// The body is the audio itself rather than a multipart form: the only sender is
// a script, and a raw body is the shape a script already has. Everything about
// the recording that is not the bytes travels in the query string.
func (m *Module) handleNarrationPut(w http.ResponseWriter, r *http.Request) {
	id, lang, err := narrationTarget(r)
	if err != nil {
		httpFail(w, http.StatusBadRequest, err)
		return
	}
	ext, ok := narrationTypes[strings.ToLower(strings.TrimSpace(strings.Split(r.Header.Get("Content-Type"), ";")[0]))]
	if !ok {
		http.Error(w, "content-type must be audio/mpeg, audio/mp4, audio/ogg, audio/opus or audio/wav",
			http.StatusUnsupportedMediaType)
		return
	}

	data, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxNarrationBytes))
	if err != nil {
		http.Error(w, "audio too large", http.StatusRequestEntityTooLarge)
		return
	}
	if len(data) == 0 {
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	// The key carries the article, the language and the current time. Time is in
	// it so a re-render lands on a new URL: browsers and any cache in front of
	// them hold audio for a long time, and a listener who came back for the
	// corrected reading would otherwise be served the old one from disk.
	key := path.Join("audio", id.String(), fmt.Sprintf("%s-%d%s", lang, time.Now().Unix(), ext))
	url, err := m.media.SaveBlob(r.Context(), key, data, contentTypeFor(ext))
	if err != nil {
		httpFail(w, http.StatusInternalServerError, fmt.Errorf("store audio: %w", err))
		return
	}

	n := Narration{
		Lang:        lang,
		URL:         url,
		StorageKey:  key,
		Bytes:       int64(len(data)),
		DurationSec: atoiDefault(r.URL.Query().Get("duration"), 0),
		Voice:       strings.TrimSpace(r.URL.Query().Get("voice")),
		TextSHA256:  strings.TrimSpace(r.URL.Query().Get("digest")),
		// Cues ride in a header rather than the query string: a long article has
		// a cue per block, and that is kilobytes of JSON. Query strings are
		// truncated by proxies and written into access logs; a header is neither.
		Cues: cueJSON(r.Header.Get("X-Audio-Cues")),
	}
	replaced, err := m.audio.Upsert(r.Context(), id, n)
	if err != nil {
		// The row is what makes the file reachable. If it did not land, the file
		// is unreferenced from the moment it was written, so remove it now
		// rather than leave litter nothing will ever collect.
		_ = m.media.DeleteBlob(r.Context(), key)
		httpFail(w, http.StatusInternalServerError, err)
		return
	}
	// Only now, with the row pointing at the new file, is the old one safe to
	// drop. A failure here costs disk, not playback, so it is not fatal.
	if replaced != "" {
		_ = m.media.DeleteBlob(r.Context(), replaced)
	}

	writeJSONObj(w, map[string]any{
		"url": url, "bytes": n.Bytes, "duration": n.DurationSec, "lang": lang,
	})
}

// handleNarrationDelete removes a reading and the file behind it.
func (m *Module) handleNarrationDelete(w http.ResponseWriter, r *http.Request) {
	id, lang, err := narrationTarget(r)
	if err != nil {
		httpFail(w, http.StatusBadRequest, err)
		return
	}
	key, err := m.audio.Delete(r.Context(), id, lang)
	if err != nil {
		httpFail(w, http.StatusInternalServerError, err)
		return
	}
	if key != "" {
		_ = m.media.DeleteBlob(r.Context(), key)
	}
	writeJSONObj(w, map[string]any{"deleted": key != ""})
}

// narrationTarget reads and validates the article and language from the path.
func narrationTarget(r *http.Request) (uuid.UUID, string, error) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return uuid.Nil, "", errors.New("bad article id")
	}
	lang := strings.ToLower(strings.TrimSpace(chi.URLParam(r, "lang")))
	switch lang {
	case LangKZ, LangRU, LangEN:
		return id, lang, nil
	}
	return uuid.Nil, "", errors.New("lang must be kz, ru or en")
}

func contentTypeFor(ext string) string {
	switch ext {
	case ".mp3":
		return "audio/mpeg"
	case ".m4a":
		return "audio/mp4"
	case ".ogg":
		return "audio/ogg"
	case ".opus":
		return "audio/opus"
	default:
		return "audio/wav"
	}
}

func atoiDefault(s string, def int) int {
	if n, err := strconv.Atoi(strings.TrimSpace(s)); err == nil && n >= 0 {
		return n
	}
	return def
}

// httpFail answers with the status and nothing else. The caller is a script: it
// needs to know that the upload failed, not the shape of our database error.
func httpFail(w http.ResponseWriter, code int, _ error) {
	http.Error(w, http.StatusText(code), code)
}

// cueJSON accepts the timing map only if it is valid JSON. A malformed map must
// not reach the column: the page would then ship broken JSON to every reader,
// and one bad upload would take the player down for an article that already had
// working audio.
func cueJSON(s string) []byte {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if !json.Valid([]byte(s)) {
		return nil
	}
	return []byte(s)
}
