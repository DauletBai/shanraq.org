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

// processImage decodes an uploaded image, corrects its size, stamps the brand
// watermark, and returns re-encoded JPEG bytes. Decoding then re-encoding from
// raw pixels strips EXIF (including GPS) and any animation frames.
func (m *Module) processImage(raw []byte) ([]byte, error) {
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
