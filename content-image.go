package facepaint

import (
	"image"

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
		height:          ceil(computed.Width),
		paddingAndGapsY: computed.PaddingTop + computed.PaddingBottom,
		paddingX:        computed.PaddingTop + computed.PaddingBottom,
		paddingAndGapsX: computed.PaddingLeft + computed.PaddingRight,
		paddingY:        computed.PaddingLeft + computed.PaddingRight,
	}
	if dimensions.width == 0 {
		dimensions.width = content.image.Bounds().Dx() + ceil(dimensions.paddingX)
	}
	if dimensions.height == 0 {
		dimensions.height = content.image.Bounds().Dy() + ceil(dimensions.paddingY)
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

	var originX, originY float64 = pos.X + computed.PaddingLeft, pos.Y + computed.PaddingTop
	if computed.BackgroundColor != nil {
		ctx.SetColor(computed.BackgroundColor)
		ctx.DrawRectangle(originX, originY, float64(dimensions.width), float64(dimensions.height))
		ctx.Fill()
	}

	if computed.Debug {
		ctx.SetColor(getDebugColor())
		ctx.DrawRectangle(pos.X, pos.Y, float64(dimensions.width), float64(dimensions.height))
		ctx.Stroke()
	}

	image := imaging.Fill(content.image, dimensions.width, dimensions.height, computed.BackgroundPosition.Imaging(), imaging.Lanczos)
	ctx.DrawImage(image, ceil(pos.X), ceil(pos.Y))

	return nil
}
