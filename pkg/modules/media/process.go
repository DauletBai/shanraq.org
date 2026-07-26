package media

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"

	// Decoders for the formats we accept on upload. Re-encoding to JPEG below
	// intentionally drops all metadata (EXIF/GPS) and any animation.
	_ "image/gif"
	_ "image/png"

	_ "golang.org/x/image/webp"

	xdraw "golang.org/x/image/draw"
)

// jpegQuality balances size and fidelity for photographic cover images.
const jpegQuality = 82

// maxDecodePixels caps the declared pixel count before full decode. A small
// compressed file can declare enormous dimensions (a decompression bomb) that
// would allocate gigabytes; reject those from the cheap header first. 50 MP is
// well above any legitimate phone photo.
const maxDecodePixels = 50_000_000

// processImage decodes an uploaded image, corrects its size, stamps the brand
// watermark, and returns re-encoded JPEG bytes. Decoding then re-encoding from
// raw pixels strips EXIF (including GPS) and any animation frames.
func (m *Module) processImage(raw []byte) ([]byte, error) {
	// Header-only check first so a decompression bomb can't force a huge alloc.
	cfg, _, err := image.DecodeConfig(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("decode image config: %w", err)
	}
	if int64(cfg.Width)*int64(cfg.Height) > maxDecodePixels {
		return nil, fmt.Errorf("image too large: %dx%d exceeds %d pixels", cfg.Width, cfg.Height, maxDecodePixels)
	}

	src, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	img := flattenAndResize(src, m.maxDim)
	if m.mark != nil {
		stampTopRight(img, m.mark)
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: jpegQuality}); err != nil {
		return nil, fmt.Errorf("encode jpeg: %w", err)
	}
	return buf.Bytes(), nil
}

// avatarDim is the square side length (px) for a stored profile avatar.
const avatarDim = 512

// processAvatar decodes an uploaded image, center-crops it to a square, scales
// it to avatarDim, and re-encodes as JPEG. Unlike processImage it does NOT stamp
// the brand watermark — an avatar is a personal photo, not shared content.
// Re-encoding from raw pixels strips EXIF/GPS the same way.
func (m *Module) processAvatar(raw []byte) ([]byte, error) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("decode image config: %w", err)
	}
	if int64(cfg.Width)*int64(cfg.Height) > maxDecodePixels {
		return nil, fmt.Errorf("image too large: %dx%d exceeds %d pixels", cfg.Width, cfg.Height, maxDecodePixels)
	}
	src, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}
	img := squareCrop(src, avatarDim)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: jpegQuality}); err != nil {
		return nil, fmt.Errorf("encode jpeg: %w", err)
	}
	return buf.Bytes(), nil
}

// squareCrop center-crops src to its largest centred square and scales that to
// a dim×dim RGBA on an opaque white base (so transparent PNGs stay clean).
func squareCrop(src image.Image, dim int) *image.RGBA {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	side := w
	if h < side {
		side = h
	}
	if side < 1 {
		side = 1
	}
	ox := b.Min.X + (w-side)/2
	oy := b.Min.Y + (h-side)/2
	crop := image.Rect(ox, oy, ox+side, oy+side)

	dst := image.NewRGBA(image.Rect(0, 0, dim, dim))
	xdraw.Draw(dst, dst.Bounds(), image.NewUniform(color.White), image.Point{}, xdraw.Src)
	xdraw.CatmullRom.Scale(dst, dst.Bounds(), src, crop, xdraw.Over, nil)
	return dst
}

// flattenAndResize composes src onto an opaque white canvas (so transparent
// PNGs don't turn black in JPEG) and scales it down so neither side exceeds max.
// The result is always an *image.RGBA anchored at (0,0).
func flattenAndResize(src image.Image, max int) *image.RGBA {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	nw, nh := w, h
	if max > 0 && (w > max || h > max) {
		if w >= h {
			nw, nh = max, h*max/w
		} else {
			nh, nw = max, w*max/h
		}
	}
	if nw < 1 {
		nw = 1
	}
	if nh < 1 {
		nh = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, nw, nh))
	xdraw.Draw(dst, dst.Bounds(), image.NewUniform(color.White), image.Point{}, xdraw.Src)
	if nw == w && nh == h {
		xdraw.Draw(dst, dst.Bounds(), src, b.Min, xdraw.Over)
	} else {
		xdraw.CatmullRom.Scale(dst, dst.Bounds(), src, b, xdraw.Over, nil)
	}
	return dst
}
