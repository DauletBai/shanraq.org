package httpserver

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
)

// etagLimit caps how much of a response is buffered to compute an ETag. Pages
// on this site run well under 200 KB compressed source; anything larger is
// almost certainly a file download or a stream, and those are passed through
// untouched rather than held in memory.
const etagLimit = 1 << 20 // 1 MiB

// Since scripts moved to a per-request nonce, an HTML body differs on every
// response and its ETag can no longer match on a revisit. The 304 path is
// kept because it still serves everything else this wraps, and because a
// static-hash policy would bring it back for HTML -- but a repeat visit to
// the same page now re-sends the body. At this site's traffic that is a few
// megabytes a day, which is the price of the policy and worth naming here.
//
// conditionalGet gives HTML responses an ETag and answers repeat requests with
// 304 Not Modified.
//
// Without it every crawler fetch downloads the whole page again: the site was
// serving no ETag, no Last-Modified, and answering If-Modified-Since with a
// full 200. With bots making several thousand requests a fortnight across ~324
// URLs, that is a large share of the crawl budget spent re-reading bytes the
// crawler already has.
//
// The tag is a hash of the exact bytes sent, so it stays correct for pages that
// differ between viewers — a signed-in reader and an anonymous one simply get
// different tags. Cache-Control: no-cache is deliberate: it permits storing the
// response but requires revalidation, which is precisely the behaviour that
// makes a 304 possible.
func conditionalGet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			next.ServeHTTP(w, r)
			return
		}
		rec := &etagRecorder{ResponseWriter: w, buf: &bytes.Buffer{}}
		next.ServeHTTP(rec, r)

		// Streamed, oversized, non-HTML or non-200 responses were written
		// straight through and there is nothing left to do.
		if rec.passthrough {
			return
		}
		body := rec.buf.Bytes()
		sum := sha256.Sum256(body)
		etag := `"` + hex.EncodeToString(sum[:16]) + `"`

		h := rec.ResponseWriter.Header()
		h.Set("ETag", etag)
		if h.Get("Cache-Control") == "" {
			h.Set("Cache-Control", "no-cache")
		}

		if matchesETag(r.Header.Get("If-None-Match"), etag) {
			h.Del("Content-Type")
			h.Del("Content-Length")
			rec.ResponseWriter.WriteHeader(http.StatusNotModified)
			return
		}
		rec.ResponseWriter.WriteHeader(rec.status)
		if r.Method != http.MethodHead {
			_, _ = rec.ResponseWriter.Write(body)
		}
	})
}

// matchesETag reports whether an If-None-Match header covers tag. It accepts
// the "*" wildcard, a comma-separated list, and weak tags, per RFC 9110.
func matchesETag(header, tag string) bool {
	header = strings.TrimSpace(header)
	if header == "" {
		return false
	}
	if header == "*" {
		return true
	}
	for _, candidate := range strings.Split(header, ",") {
		candidate = strings.TrimSpace(candidate)
		candidate = strings.TrimPrefix(candidate, "W/")
		if candidate == tag {
			return true
		}
	}
	return false
}

// etagRecorder buffers a response so its bytes can be hashed. It gives up and
// switches to passthrough as soon as the response turns out to be something an
// ETag should not be computed for.
type etagRecorder struct {
	http.ResponseWriter
	buf         *bytes.Buffer
	status      int
	passthrough bool
	wroteHeader bool
}

func (e *etagRecorder) WriteHeader(status int) {
	if e.wroteHeader || e.passthrough {
		return
	}
	e.wroteHeader = true
	e.status = status
	// Only cache-validate plain successful HTML. Redirects, errors, downloads
	// and anything already carrying its own validator are left alone.
	ct := e.Header().Get("Content-Type")
	if status != http.StatusOK || !strings.HasPrefix(ct, "text/html") || e.Header().Get("ETag") != "" {
		e.passthrough = true
		e.ResponseWriter.WriteHeader(status)
	}
}

func (e *etagRecorder) Write(p []byte) (int, error) {
	if !e.wroteHeader {
		e.WriteHeader(http.StatusOK)
	}
	if e.passthrough {
		return e.ResponseWriter.Write(p)
	}
	if e.buf.Len()+len(p) > etagLimit {
		// Too big to hold: flush what was buffered and stream the rest.
		e.passthrough = true
		e.ResponseWriter.WriteHeader(e.status)
		if _, err := e.ResponseWriter.Write(e.buf.Bytes()); err != nil {
			return 0, err
		}
		e.buf.Reset()
		return e.ResponseWriter.Write(p)
	}
	return e.buf.Write(p)
}

// Flush marks the response as streamed: a handler that flushes wants its bytes
// out now, and buffering them to hash would defeat that.
func (e *etagRecorder) Flush() {
	if !e.passthrough {
		if !e.wroteHeader {
			e.WriteHeader(http.StatusOK)
		}
		if !e.passthrough {
			e.passthrough = true
			e.ResponseWriter.WriteHeader(e.status)
			_, _ = e.ResponseWriter.Write(e.buf.Bytes())
			e.buf.Reset()
		}
	}
	if f, ok := e.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}
