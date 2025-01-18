package facepaint

import (
	"image"
	"math"

	"github.com/cufee/facepaint/style"
	"github.com/nao1215/imaging"
	"github.com/pkg/errors"
)

var _ BlockContent = &contentEmpty{}

func NewImageContent(style style.StyleOptions, image image.Image) (*Block, error) {
	if image == nil {
		return nil, errors.New("image cannot be nil")
	}
	return NewBlock(&contentImage{
		style: style,
		image: image,
	}), nil
}

type contentImage struct {
	style style.StyleOptions
	image image.Image
}

func (content *contentImage) setStyle(style style.StyleOptions) {
	content.style = style
}

func (content *contentImage) dimensions() contentDimensions {
	computed := content.Style().Computed()
	dimensions := contentDimensions{
		width:           ceil(computed.Width),
		height:          ceil(computed.Height),
		paddingAndGapsX: computed.PaddingLeft + computed.PaddingRight,
		paddingX:        computed.PaddingLeft + computed.PaddingRight,
		paddingAndGapsY: computed.PaddingTop + computed.PaddingBottom,
		paddingY:        computed.PaddingTop + computed.PaddingBottom,
	}

	if dimensions.width == 0 && dimensions.height == 0 {
		dimensions.width = content.image.Bounds().Dx() + ceil(dimensions.paddingX)
		dimensions.height = content.image.Bounds().Dy() + ceil(dimensions.paddingY)
	}

	// if new width or height is 0 then preserve aspect ratio, minimum 1px.
	if dimensions.width == 0 {
		tmpW := float64(dimensions.height) * float64(content.image.Bounds().Dx()) / float64(content.image.Bounds().Dy())
		dimensions.width = int(max(1.0, math.Floor(tmpW+0.5)))
	}
	if dimensions.height == 0 {
		tmpH := float64(dimensions.width) * float64(content.image.Bounds().Dy()) / float64(content.image.Bounds().Dx())
		dimensions.height = int(math.Max(1.0, math.Floor(tmpH+0.5)))
	}

	return dimensions
}

func (content *contentImage) Type() blockContentType {
	return BlockContentTypeImage
}

func (content *contentImage) Layers() map[int]struct{} {
	return map[int]struct{}{content.style.Computed().ZIndex: {}}
}

func (content *contentImage) Style() style.StyleOptions {
	return content.style
}

func (content *contentImage) Render(layers layerContext, pos Position) error {
	computed := content.style.Computed()
	dimensions := content.dimensions()
	ctx, err := layers.layer(computed.ZIndex)
	if err != nil {
		return err
	}

	if computed.BackgroundColor != nil {
		ctx.SetColor(computed.BackgroundColor)
		drawBackgroundPath(ctx, computed, dimensions, pos)
		ctx.Fill()
	}
	if computed.Debug {
		ctx.SetColor(getDebugColor())
		ctx.DrawRectangle(pos.X, pos.Y, float64(dimensions.width), float64(dimensions.height))
		ctx.Stroke()
	}

	image := imaging.Resize(content.image, dimensions.width, dimensions.height, imaging.Lanczos)
	ctx.DrawImage(image, ceil(pos.X), ceil(pos.Y))

	return nil
}
