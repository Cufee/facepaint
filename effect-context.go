package facepaint

import (
	"image"
	"math"
	"sync"

	"github.com/cufee/facepaint/style"
)

type effectContext struct {
	img   image.Image
	dims  contentDimensions
	style style.Style

	localBounds image.Rectangle

	sdfOnce sync.Once
	sdf     func(x, y float64) float64
}

func newEffectContext(img image.Image, dims contentDimensions, st style.Style) *effectContext {
	return &effectContext{
		img:         img,
		dims:        dims,
		style:       st,
		localBounds: image.Rect(0, 0, dims.Width, dims.Height),
	}
}

func (c *effectContext) Image() image.Image           { return c.img }
func (c *effectContext) LocalBounds() image.Rectangle { return c.localBounds }

func (c *effectContext) SDF() func(x, y float64) float64 {
	c.sdfOnce.Do(func() { c.sdf = buildRoundedBoxSDF(c.dims, c.style) })
	return c.sdf
}

func buildRoundedBoxSDF(dims contentDimensions, st style.Style) func(x, y float64) float64 {
	w, h := float64(dims.Width), float64(dims.Height)
	if w <= 0 || h <= 0 {
		return nil
	}
	hx, hy := w/2, h/2
	rtl := clampRadius(st.BorderRadiusTopLeft, w, h)
	rtr := clampRadius(st.BorderRadiusTopRight, w, h)
	rbl := clampRadius(st.BorderRadiusBottomLeft, w, h)
	rbr := clampRadius(st.BorderRadiusBottomRight, w, h)
	return func(x, y float64) float64 {
		return sdRoundedBox(x-hx, y-hy, hx, hy, rtl, rtr, rbl, rbr)
	}
}

func clampRadius(r, w, h float64) float64 {
	max := math.Min(w, h) / 2
	return clamp(r, 0, max)
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// sdRoundedBox is the 2D SDF for a rounded box with per-corner radii.
// Negative inside, positive outside. From https://iquilezles.org/articles/distfunctions2d/
func sdRoundedBox(px, py, hx, hy, rtl, rtr, rbl, rbr float64) float64 {
	var r float64
	if px > 0 {
		if py > 0 {
			r = rtr
		} else {
			r = rbr
		}
	} else {
		if py > 0 {
			r = rtl
		} else {
			r = rbl
		}
	}

	ax := math.Abs(px) - hx + r
	ay := math.Abs(py) - hy + r

	outside := math.Hypot(math.Max(ax, 0), math.Max(ay, 0))
	inside := math.Min(math.Max(ax, ay), 0)
	return outside + inside - r
}
