package content

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type Page struct {
	Title string
	HTML  string 
}

type Store struct {
	root string 
	md   goldmark.Markdown
}

func NewStore(root string) *Store {
	return &Store{
		root: root,
		md: goldmark.New(
			goldmark.WithExtensions(extension.GFM),
			goldmark.WithParserOptions(parser.WithAutoHeadingID()),
			goldmark.WithRendererOptions(html.WithHardWraps()),
		),
	}
}

func (s *Store) Get(slug string) (Page, error) {
	path := filepath.Join(s.root, slug+".md")
	b, err := os.ReadFile(path)
	if err != nil {
		return Page{}, fmt.Errorf("read %s: %w", path, err)
	}
	var out []byte
	var buf = &out
	if err := s.md.Convert(b, (*writer)(buf)); err != nil {
		return Page{}, fmt.Errorf("render md: %w", err)
	}
	return Page{Title: slug, HTML: string(out)}, nil
}

// tiny helper to satisfy goldmark
type writer []byte
func (w *writer) Write(p []byte) (int, error) { *w = append(*w, p...); return len(p), nil }