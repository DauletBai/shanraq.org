package web

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"io/fs"
	"net/http"
	"strings"
	"sync"
)

//go:embed static/*
var staticFiles embed.FS

var staticRoot fs.FS

func init() {
	sub, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic("web: failed to load static assets: " + err.Error())
	}
	staticRoot = sub
}

// StaticFS returns the compiled-in static file tree.
func StaticFS() fs.FS {
	return staticRoot
}

// StaticHandler provides an http.Handler that serves static assets.
func StaticHandler() http.Handler {
	return http.FileServer(http.FS(staticRoot))
}

// assetVersions caches one content hash per static path. The tree is embedded
// in the binary, so a hash is fixed for the life of a build and is computed at
// most once per path.
var assetVersions sync.Map

// AssetURL appends a content hash to a static path: /static/css/shanraq.css →
// /static/css/shanraq.css?v=1f4c9a2b70.
//
// Without it a deploy is invisible to anyone who has been here before. The
// assets go out with Cache-Control: max-age=86400 and no ETag, so at a stable
// URL a returning reader keeps yesterday's stylesheet for a full day — new
// markup styled by old CSS, which looks like a broken site rather than a stale
// cache. The hash changes only when the file does, so caching stays aggressive.
func AssetURL(p string) string {
	if v, ok := assetVersions.Load(p); ok {
		return p + "?v=" + v.(string)
	}
	b, err := fs.ReadFile(staticRoot, strings.TrimPrefix(p, "/static/"))
	if err != nil {
		return p // unknown path: better unversioned than broken
	}
	sum := sha256.Sum256(b)
	v := hex.EncodeToString(sum[:])[:10]
	assetVersions.Store(p, v)
	return p + "?v=" + v
}
