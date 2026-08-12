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
//
// Fire-and-forget on purpose: search engines must never be able to slow down or
// fail a publish. It runs detached with its own timeout, and a failure is
// logged and dropped — the sitemap remains the reliable path, this is only the
// fast one.
func (m *Module) submitIndexNow(slug string) {
	if m.indexNowKey == "" || slug == "" {
		return
	}
	host := m.baseURL
	if u, err := url.Parse(m.baseURL); err == nil && u.Host != "" {
		host = u.Host
	}
	urls := make([]string, 0, len(rssLangOrder))
	for _, lang := range rssLangOrder {
		urls = append(urls, m.articleURL(slug, lang))
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
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, indexNowEndpoint, bytes.NewReader(body))
		if err != nil {
			m.log.Warn("indexnow request", zap.Error(err))
			return
		}
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		resp, err := m.http.Do(req)
		if err != nil {
			m.log.Warn("indexnow submit", zap.String("slug", slug), zap.Error(err))
			return
		}
		defer resp.Body.Close()
		// 200 accepted, 202 accepted but key still being validated. Anything else
		// is worth seeing in the log: 403 means the key file is unreachable.
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			m.log.Warn("indexnow rejected", zap.String("slug", slug), zap.Int("status", resp.StatusCode))
			return
		}
		m.log.Info("indexnow submitted", zap.String("slug", slug), zap.Int("urls", len(urls)))
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
