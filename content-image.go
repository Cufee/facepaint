package facepaint

import (
	"image"
	"math"

	"github.com/cufee/facepaint/style"
	"github.com/fogleman/gg"
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
func MustNewImageContent(style style.StyleOptions, image image.Image) *Block {
	b, _ := NewImageContent(style, image)
	return b
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
		Width:           ceil(max(computed.Width, computed.MinWidth)),
		Height:          ceil(max(computed.Height, computed.MinHeight)),
		paddingAndGapsX: computed.PaddingLeft + computed.PaddingRight,
		paddingX:        computed.PaddingLeft + computed.PaddingRight,
		paddingAndGapsY: computed.PaddingTop + computed.PaddingBottom,
		paddingY:        computed.PaddingTop + computed.PaddingBottom,
	}

	if dimensions.Width == 0 && dimensions.Height == 0 {
		dimensions.Width = content.image.Bounds().Dx() + ceil(dimensions.paddingX)
		dimensions.Height = content.image.Bounds().Dy() + ceil(dimensions.paddingY)
	}

	// if new width or height is 0 then preserve aspect ratio, minimum 1px.
	if dimensions.Width == 0 {
		tmpW := float64(dimensions.Height) * float64(content.image.Bounds().Dx()) / float64(content.image.Bounds().Dy())
		dimensions.Width = int(max(1.0, math.Floor(tmpW+0.5)))
	}
	if dimensions.Height == 0 {
		tmpH := float64(dimensions.Width) * float64(content.image.Bounds().Dy()) / float64(content.image.Bounds().Dx())
		dimensions.Height = int(math.Max(1.0, math.Floor(tmpH+0.5)))
	}

	dimensions.Width = max(dimensions.Width, ceil(computed.MinWidth))
	dimensions.Height = max(dimensions.Height, ceil(computed.MinHeight))
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

func (content *contentImage) Render(layers *layerContext, pos Position) error {
	computed := content.style.Computed()
	dimensions := content.dimensions()
	ctx, err := layers.layer(computed.ZIndex)
	if err != nil {
		return err
	}

	childCtx := newLayer(dimensions.Width, dimensions.Height)

	// apply border radius
	drawBackgroundPath(childCtx, computed, dimensions, Position{0, 0})
	childCtx.Clip()

	if computed.BackgroundColor != nil {
		childCtx.SetColor(computed.BackgroundColor)
		childCtx.Clear()
	}
	if computed.Debug {
		childCtx.SetColor(getDebugColor())
		childCtx.DrawRectangle(0, 0, float64(dimensions.Width), float64(dimensions.Height))
		childCtx.Stroke()
	}

	image := imaging.Fill(content.image, dimensions.Width, dimensions.Height, computed.BackgroundPosition.Imaging(), imaging.Lanczos)
	if computed.Blur > 0 {
		image = imaging.Blur(image, computed.Blur)
	}

	if computed.Color != nil {
		mask := gg.NewContextForImage(image)
		childCtx.SetMask(mask.AsMask())
		childCtx.SetColor(computed.Color)
		drawBackgroundPath(childCtx, computed, dimensions, Position{0, 0})
		childCtx.Fill()

		ctx.DrawImage(childCtx.Image(), ceil(pos.X), ceil(pos.Y))
		return nil
	}

	childCtx.DrawImage(image, 0, 0)
	ctx.DrawImage(childCtx.Image(), ceil(pos.X), ceil(pos.Y))

	return nil
}
