package facepaint

import (
	"image"
	"image/draw"
	"runtime"
	"sync"

	"github.com/nao1215/imaging"
)

type BevelProfile byte

const (
	BevelProfileCircular BevelProfile = iota
	BevelProfileSmoothstep
)

func toRGBA(img image.Image) *image.RGBA {
	if rgba, ok := img.(*image.RGBA); ok {
		return rgba
	}
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	return rgba
}

func clamp8(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return v
}

func imagingBlur(src *image.RGBA, sigma float64) *image.RGBA {
	return toRGBA(imaging.Blur(src, sigma))
}

func parallelRows(rows int, fn func(y0, y1 int)) {
	if rows <= 0 {
		return
	}
	procs := runtime.GOMAXPROCS(0)
	if rows < 64 || procs <= 1 {
		fn(0, rows)
		return
	}

	stripes := min(procs, rows)
	chunk := (rows + stripes - 1) / stripes

	var wg sync.WaitGroup
	for i := 0; i < stripes; i++ {
		y0 := i * chunk
		y1 := min(y0+chunk, rows)
		if y0 >= y1 {
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn(y0, y1)
		}()
	}
	wg.Wait()
}
