package media

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"shanraq.org/web"
)

func brandMark(t *testing.T) *image.RGBA {
	t.Helper()
	svg, err := fs.ReadFile(web.StaticFS(), "brand/shanraq-mark.svg")
	if err != nil {
		t.Fatalf("read brand svg: %v", err)
	}
	mark, err := rasterizeSVG(svg, watermarkPx, watermarkPx)
	if err != nil {
		t.Fatalf("rasterize: %v", err)
	}
	return mark
}

func whitePNG(t *testing.T, w, h int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := range img.Pix {
		img.Pix[i] = 0xff
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
	return buf.Bytes()
}

func TestProcessImageResizesAndWatermarks(t *testing.T) {
	m := &Module{maxDim: 2000, mark: brandMark(t)}

	out, err := m.processImage(whitePNG(t, 3000, 1000))
	if err != nil {
		t.Fatalf("process: %v", err)
	}

	img, format, err := image.Decode(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("decode output: %v", err)
	}
	if format != "jpeg" {
		t.Fatalf("output format = %q, want jpeg (re-encode strips EXIF)", format)
	}

	b := img.Bounds()
	if b.Dx() != 2000 || b.Dy() != 666 {
		t.Fatalf("output size = %dx%d, want 2000x666", b.Dx(), b.Dy())
	}

	// The brand mark must appear in the top-right corner: scan that box for a
	// pixel noticeably darker than the white background.
	found := false
	for y := watermarkMargin; y < watermarkMargin+watermarkPx; y++ {
		for x := b.Dx() - watermarkMargin - watermarkPx; x < b.Dx()-watermarkMargin; x++ {
			r, g, bl, _ := img.At(x, y).RGBA()
			if r < 55000 || g < 55000 || bl < 55000 {
				found = true
			}
		}
	}
	if !found {
		t.Fatal("no watermark pixels found in the top-right corner")
	}

	// The opposite corner must stay clean (white).
	rr, gg, bb, _ := img.At(20, b.Dy()-20).RGBA()
	if rr < 60000 || gg < 60000 || bb < 60000 {
		t.Fatalf("bottom-left unexpectedly not white: %d,%d,%d", rr, gg, bb)
	}
}

func TestProcessImageRejectsNonImage(t *testing.T) {
	m := &Module{maxDim: 2000}
	if _, err := m.processImage([]byte("this is not an image")); err == nil {
		t.Fatal("expected error for non-image input")
	}
}

func TestProcessImageWithoutWatermark(t *testing.T) {
	m := &Module{maxDim: 500} // mark nil
	out, err := m.processImage(whitePNG(t, 400, 400))
	if err != nil {
		t.Fatalf("process: %v", err)
	}
	img, _, err := image.Decode(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if img.Bounds().Dx() != 400 { // below maxDim: no upscaling
		t.Fatalf("width = %d, want 400", img.Bounds().Dx())
	}
	_ = color.White
}

func TestFSStorePutAndURL(t *testing.T) {
	dir := t.TempDir()
	store, err := NewFSStore(dir, "/media")
	if err != nil {
		t.Fatalf("new fs store: %v", err)
	}
	key := "ab/abcdef.jpg"
	if err := store.Put(context.Background(), key, []byte("data"), "image/jpeg"); err != nil {
		t.Fatalf("put: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(dir, "ab", "abcdef.jpg"))
	if err != nil || string(got) != "data" {
		t.Fatalf("read back = %q, %v", got, err)
	}
	if store.URL(key) != "/media/ab/abcdef.jpg" {
		t.Fatalf("url = %q", store.URL(key))
	}
	if err := store.Delete(context.Background(), key); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if err := store.Delete(context.Background(), key); err != nil {
		t.Fatalf("delete missing should be nil: %v", err)
	}
}
