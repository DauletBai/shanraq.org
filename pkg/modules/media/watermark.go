package media

import (
	"bytes"
	"image"
	"image/color"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	xdraw "golang.org/x/image/draw"
)

const (
	// watermarkPx is the rendered size of the brand mark, per spec: a 32px icon.
	watermarkPx = 32
	// watermarkMargin is the inset from the top and right edges, in pixels.
	watermarkMargin = 12
	// watermarkAlpha makes the mark semi-transparent (0–255): visible as a
	// brand stamp but not obscuring the photo underneath. The mark is white, so
	// it reads on most photos; a touch more opacity keeps it present.
	watermarkAlpha = 175
)

// supersample renders the mark at this multiple of the target size, then
// downscales — the logo's thick crossing strokes merge into a blob if rendered
// straight at 32px, but area-averaging a 4× render keeps the star readable.
const supersample = 4

// rasterizeSVG renders an SVG document to a w×h RGBA image (transparent
// background), supersampled for crispness. Called once at startup.
func rasterizeSVG(svg []byte, w, h int) (*image.RGBA, error) {
	icon, err := oksvg.ReadIconStream(bytes.NewReader(svg))
	if err != nil {
		return nil, err
	}
	hiW, hiH := w*supersample, h*supersample
	icon.SetTarget(0, 0, float64(hiW), float64(hiH))
	hi := image.NewRGBA(image.Rect(0, 0, hiW, hiH))
	scanner := rasterx.NewScannerGV(hiW, hiH, hi, hi.Bounds())
	raster := rasterx.NewDasher(hiW, hiH, scanner)
	icon.Draw(raster, 1.0)
	if supersample == 1 {
		return hi, nil
	}
	out := image.NewRGBA(image.Rect(0, 0, w, h))
	xdraw.CatmullRom.Scale(out, out.Bounds(), hi, hi.Bounds(), xdraw.Over, nil)
	return out, nil
}

// stampTopRight composites a static watermark into the top-right corner of dst.
// The mark keeps its own shape (transparent background) but is dimmed uniformly
// by watermarkAlpha. No animation is involved — it is a single still overlay.
func stampTopRight(dst *image.RGBA, mark image.Image) {
	mb := mark.Bounds()
	b := dst.Bounds()
	// Skip if the image is too small to carry the mark cleanly.
	if b.Dx() < mb.Dx()+2*watermarkMargin || b.Dy() < mb.Dy()+2*watermarkMargin {
		return
	}
	x0 := b.Max.X - mb.Dx() - watermarkMargin
	y0 := b.Min.Y + watermarkMargin
	r := image.Rect(x0, y0, x0+mb.Dx(), y0+mb.Dy())
	mask := image.NewUniform(color.Alpha{A: watermarkAlpha})
	xdraw.DrawMask(dst, r, mark, mb.Min, mask, image.Point{}, xdraw.Over)
}
