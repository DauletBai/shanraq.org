package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
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
//
// webp is not here on purpose. Below we do not believe the first bytes on
// their own -- we decode the picture's header -- and the standard library has
// no webp decoder. Accepting a format we cannot check would mean trusting the
// sender again, which is what this file exists to avoid.
var allowed = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
	"image/gif":  ".gif",
}

// maxSide caps how large a picture may claim to be. A file of a few kilobytes
// can announce twenty thousand pixels a side, and decoding it asks for
// gigabytes of memory: a "decompression bomb". The header is cheap to read,
// the pixels are not, so the size is checked before anything is decoded.
const maxSide = 8000

var (
	errNotImage = errors.New("бұл сурет емес")
	errTooLarge = errors.New("сурет тым үлкен")
)

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

	// The first bytes are a hint, not a proof: they are trivial to forge by
	// writing a real picture's opening into any file. So the picture's own
	// header is decoded -- that either works or it does not -- and the size it
	// announces is checked before we agree to keep the file.
	cfg, format, err := image.DecodeConfig(file)
	if err != nil {
		return "", fmt.Errorf("%w: %v", errNotImage, err)
	}
	if cfg.Width > maxSide || cfg.Height > maxSide {
		return "", fmt.Errorf("%dx%d: %w", cfg.Width, cfg.Height, errTooLarge)
	}
	if allowed["image/"+format] == "" && !(format == "jpeg" && ext == ".jpg") {
		return "", fmt.Errorf("%s: %w", format, errNotImage)
	}

	// The header alone still proves little: a real picture's opening followed
	// by rubbish reads as a valid header. So the whole picture is decoded once
	// -- after the size above made that safe -- and only what decodes is kept.
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	if _, _, err := image.Decode(file); err != nil {
		return "", fmt.Errorf("%w: %v", errNotImage, err)
	}
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
	path := filepath.Join(app.uploads, name)
	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		// A half-written file is worse than none: it stays on the disk, looks
		// like a picture and shows a reader a broken image.
		dst.Close()
		os.Remove(path)
		return "", err
	}
	return name, nil
}
