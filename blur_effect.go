package facepaint

import (
	"image"

	"github.com/cufee/facepaint/style"
	"github.com/nao1215/imaging"
)

// Blur is an Effect that applies a Gaussian blur.
type Blur struct{ Sigma float64 }

func (b *Blur) Apply(ctx style.EffectContext) (image.Image, error) {
	if b.Sigma <= 0 {
		return ctx.Image(), nil
	}
	return imaging.Blur(ctx.Image(), b.Sigma), nil
}
