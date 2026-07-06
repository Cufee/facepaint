package style

import "image"

// Effect is a render pass that mutates a block's content or backdrop.
type Effect interface {
	Apply(ctx EffectContext) (image.Image, error)
}

// EffectContext is a lazy view of a block's render state.
type EffectContext interface {
	Image() image.Image
	SDF() func(x, y float64) float64
	LocalBounds() image.Rectangle
}
