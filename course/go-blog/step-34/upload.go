package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// maxUpload is how many bytes we agree to read at all. Without this number one
// request can take the server's memory.
const maxUpload = 2 << 20 // 2 МиБ

// allowed lists what we accept. A white list rather than a black one: the
// imagination of a stranger is richer than ours.
var allowed = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

var errNotImage = errors.New("бұл сурет емес")

// saveCover writes what came in the form into the uploads folder and returns
// the name it was given. The name is ours: what a reader sent is good for
// showing back to them and for nothing else.
func (app *app) saveCover(file multipart.File) (string, error) {
	// The first 512 bytes decide what this really is. Neither the extension
	// nor the Content-Type the browser sent deserves any trust.
	head := make([]byte, 512)
	n, err := io.ReadFull(file, head)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
		return "", err
	}
	kind := strings.SplitN(http.DetectContentType(head[:n]), ";", 2)[0]
	ext, ok := allowed[kind]
	if !ok {
		return "", fmt.Errorf("%s: %w", kind, errNotImage)
	}

	// Those 512 bytes are already read; without rewinding the file would be
	// saved with its head missing.
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	raw := make([]byte, 12)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	name := base64.RawURLEncoding.EncodeToString(raw) + ext

	if err := os.MkdirAll(app.uploads, 0o755); err != nil {
		return "", err
	}
	dst, err := os.Create(filepath.Join(app.uploads, name))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}
	return name, nil
}
