package media

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Store persists processed media objects and knows their public URL. The
// filesystem backend below is the media-1 implementation; an S3/MinIO backend
// is a drop-in replacement behind this same interface (media-2).
type Store interface {
	// Put stores data under key (a forward-slash path like "ab/abcd….jpg").
	Put(ctx context.Context, key string, data []byte, contentType string) error
	// Delete removes an object; missing objects are not an error.
	Delete(ctx context.Context, key string) error
	// URL returns the public URL a browser can fetch key from.
	URL(key string) string
}

// FSStore keeps objects on local disk under Dir and serves them from Prefix.
type FSStore struct {
	dir    string
	prefix string
}

// NewFSStore creates the storage directory if needed.
func NewFSStore(dir, prefix string) (*FSStore, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &FSStore{dir: dir, prefix: "/" + strings.Trim(prefix, "/")}, nil
}

func (s *FSStore) Put(_ context.Context, key string, data []byte, _ string) error {
	p := filepath.Join(s.dir, filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}

func (s *FSStore) Delete(_ context.Context, key string) error {
	err := os.Remove(filepath.Join(s.dir, filepath.FromSlash(key)))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (s *FSStore) URL(key string) string { return s.prefix + "/" + key }

// FileServer serves stored objects. http.FileServer already rejects path
// traversal, and keys are content hashes, so no user string reaches the path.
func (s *FSStore) FileServer() http.Handler {
	return http.StripPrefix(s.prefix+"/", http.FileServer(http.Dir(s.dir)))
}

// Prefix is the URL prefix objects are served under.
func (s *FSStore) Prefix() string { return s.prefix }
