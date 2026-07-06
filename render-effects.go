package facepaint

import (
	"image"

	"github.com/cufee/facepaint/style"
	"github.com/nao1215/imaging"
)

func renderWithFilter(
	layers *layerContext,
	computed style.Style,
	dims contentDimensions,
	pos Position,
	fn func(ctx *layer, pos Position) error,
) error {
	if len(computed.Filter) == 0 {
		ctx, err := layers.layer(computed.ZIndex)
		if err != nil {
			return err
		}
		return fn(ctx, pos)
	}

	temp := newLayer(dims.Width, dims.Height)
	if err := fn(temp, Position{0, 0}); err != nil {
		return err
	}

	out := image.Image(temp.Image())
	for _, e := range computed.Filter {
		ectx := newEffectContext(out, dims, computed)
		res, err := e.Apply(ectx)
		if err != nil {
			return err
		}
		out = res
	}

	parent, err := layers.layer(computed.ZIndex)
	if err != nil {
		return err
	}
	parent.DrawImage(out, ceil(pos.X), ceil(pos.Y))
	return nil
}

func registerBackdrop(layers *layerContext, computed style.Style, dims contentDimensions, pos Position) {
	if len(computed.Backdrop) == 0 {
		return
	}
	blockPos := pos
	blockDims := dims
	blockStyle := computed
	layers.registerHook(computed.ZIndex, layerHookBeforeRender(func(frame, layer *layer) {
		crop := image.Rect(ceil(blockPos.X), ceil(blockPos.Y), ceil(blockPos.X)+blockDims.Width, ceil(blockPos.Y)+blockDims.Height)
		var backdrop image.Image = imaging.Crop(frame.Image(), crop)
		for _, e := range blockStyle.Backdrop {
			ectx := newEffectContext(backdrop, blockDims, blockStyle)
			res, err := e.Apply(ectx)
			if err != nil {
				return
			}
			backdrop = res
		}
		if blockStyle.BackgroundColor != nil {
			backdrop = blendBackgroundColor(backdrop, blockDims, blockStyle)
		}
		drawBackgroundPath(frame, blockStyle, blockDims, blockPos)
		frame.Clip()
		frame.DrawImage(backdrop, ceil(blockPos.X), ceil(blockPos.Y))
		frame.ResetClip()
	}))
}

func blendBackgroundColor(backdrop image.Image, dims contentDimensions, st style.Style) image.Image {
	rgba := toRGBA(backdrop)
	c := st.BackgroundColor
	cr, cg, cb, ca := c.RGBA()
	if ca == 0 {
		return rgba
	}
	af := float64(ca) / 65535
	r := float64(cr) / 65535 * 255
	g := float64(cg) / 65535 * 255
	b := float64(cb) / 65535 * 255

	parallelRows(dims.Height, func(y0, y1 int) {
		for y := y0; y < y1; y++ {
			for x := 0; x < dims.Width; x++ {
				i := (y*dims.Width + x) * 4
				for c := range 3 {
					var tc float64
					switch c {
					case 0:
						tc = r
					case 1:
						tc = g
					case 2:
						tc = b
					}
					rgba.Pix[i+c] = uint8(clamp8(float64(rgba.Pix[i+c])*(1-af) + tc*af))
				}
			}
		}
	})
	return rgba
}
