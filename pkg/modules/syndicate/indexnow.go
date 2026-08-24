package syndicate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// IndexNow tells Bing and Yandex about a page the moment it is published,
// instead of waiting to be crawled. One endpoint serves both, and Yandex is the
// reason this is worth having: it visited the site once in twelve days while
// Googlebot came 468 times, so for a Russian-language site aimed at Kazakhstan
// the largest search engine in the region effectively does not know we exist.
//
// The protocol is a shared secret proving control of the domain: the key is
// published at a URL on this host, and every submission points back at it.
// keyLocation lets that URL be anything on the domain, so the key lives at a
// fixed path rather than needing a route named after the key itself.
const (
	indexNowEndpoint = "https://api.indexnow.org/indexnow"
	indexNowKeyPath  = "/indexnow.txt"
	// indexNowBatch is the protocol's limit on addresses in one submission.
	indexNowBatch = 10000
)

// handleIndexNowKey serves the ownership key as plain text. Without it every
// submission is rejected.
func (m *Module) handleIndexNowKey(w http.ResponseWriter, r *http.Request) {
	if m.indexNowKey == "" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(m.indexNowKey))
}

// submitIndexNow pushes one article's three language URLs to IndexNow.
func (m *Module) submitIndexNow(slug string) {
	if slug == "" {
		return
	}
	urls := make([]string, 0, len(rssLangOrder))
	for _, lang := range rssLangOrder {
		urls = append(urls, m.articleURL(slug, lang))
	}
	m.SubmitURLs(urls, slug)
}

// SubmitURLs tells search engines about changed addresses. It suits any page of the
// site, not only an article: sitemaps wait to be crawled, and this is the only way
// to say so at once.
//
// The send is deliberately detached from the caller: a search engine must be able
// neither to slow a publication down nor to bring it down. An error is logged and
// forgotten — the reliable route is still the sitemap, this one is only the fast one.
func (m *Module) SubmitURLs(urls []string, label string) {
	if m.indexNowKey == "" || len(urls) == 0 {
		return
	}
	host := m.baseURL
	if u, err := url.Parse(m.baseURL); err == nil && u.Host != "" {
		host = u.Host
	}
	// The protocol accepts no more than ten thousand addresses at a time; we will never
	// have that many, but the list still has to be cut where it can be seen.
	if len(urls) > indexNowBatch {
		m.log.Warn("indexnow batch trimmed", zap.String("label", label),
			zap.Int("urls", len(urls)), zap.Int("kept", indexNowBatch))
		urls = urls[:indexNowBatch]
	}
	body, err := json.Marshal(map[string]any{
		"host":        host,
		"key":         m.indexNowKey,
		"keyLocation": m.baseURL + indexNowKeyPath,
		"urlList":     urls,
	})
	if err != nil {
		m.log.Warn("indexnow payload", zap.Error(err))
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, indexNowEndpoint, bytes.NewReader(body))
		if err != nil {
			m.log.Warn("indexnow request", zap.Error(err))
			return
		}
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		resp, err := m.http.Do(req)
		if err != nil {
			m.log.Warn("indexnow submit", zap.String("label", label), zap.Error(err))
			return
		}
		defer resp.Body.Close()
		// 200 accepted, 202 accepted but key still being validated. Anything else
		// is worth seeing in the log: 403 means the key file is unreachable.
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			m.log.Warn("indexnow rejected", zap.String("label", label), zap.Int("status", resp.StatusCode))
			return
		}
		m.log.Info("indexnow submitted", zap.String("label", label), zap.Int("urls", len(urls)))
	}()
}

// articleSlug looks up the slug for a published article id.
func (m *Module) articleSlug(ctx context.Context, id uuid.UUID) (string, error) {
	var slug string
	err := m.db.QueryRow(ctx, `SELECT slug FROM articles WHERE id = $1`, id).Scan(&slug)
	if err != nil {
		return "", fmt.Errorf("article slug: %w", err)
	}
	return slug, nil
}

// normalizeIndexNowKey trims and validates a configured key. IndexNow requires
// 8–128 characters drawn from hex digits and dashes; a key that fails this is
// dropped with a warning rather than being sent and silently rejected.
func normalizeIndexNowKey(raw string) (string, bool) {
	k := strings.TrimSpace(raw)
	if len(k) < 8 || len(k) > 128 {
		return "", k == ""
	}
	for _, r := range k {
		switch {
		case r >= '0' && r <= '9', r >= 'a' && r <= 'f', r >= 'A' && r <= 'F', r == '-':
		default:
			return "", false
		}
	}
	return k, true
}
