package web

import (
	"embed"
	"io/fs"
	"net/http"
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
