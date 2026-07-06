package tests

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/png"

	xdraw "golang.org/x/image/draw"
)

//go:embed assets/promo-invite.jpg
var promoData []byte

func PromoBackground() image.Image {
	img, _ := jpeg.Decode(bytes.NewReader(promoData))
	return img
}

// SplitBackground returns an image where the left half is the promo
// photo (scaled to fill) and the right half is a checkerboard pattern.
func SplitBackground(w, h, cell int) image.Image {
	halfW := w / 2
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	// Left half: scale promo to fill halfW x h
	promo := PromoBackground()
	xdraw.CatmullRom.Scale(img, image.Rect(0, 0, halfW, h), promo, promo.Bounds(), xdraw.Src, nil)

	// Right half: checker pattern
	c1 := color.RGBA{200, 200, 210, 255}
	c2 := color.RGBA{50, 50, 70, 255}
	for y := range h {
		for x := halfW; x < w; x++ {
			if (x/cell+y/cell)%2 == 0 {
				img.SetRGBA(x, y, c1)
			} else {
				img.SetRGBA(x, y, c2)
			}
		}
	}

	return img
}

func checkerPattern(w, h, cell int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	c1 := color.RGBA{200, 200, 210, 255}
	c2 := color.RGBA{50, 50, 70, 255}
	for y := range h {
		for x := range w {
			if (x/cell+y/cell)%2 == 0 {
				img.SetRGBA(x, y, c1)
			} else {
				img.SetRGBA(x, y, c2)
			}
		}
	}
	return img
}
