package facepaint

import (
	"github.com/cufee/facepaint/style"
)

var _ BlockContent = &contentEmpty{}

func NewEmptyContent(style style.StyleOptions) *Block {
	return NewBlock(&contentEmpty{
		style: style,
	})
}

type contentEmpty struct {
	style style.StyleOptions
}

func (content *contentEmpty) setStyle(style style.StyleOptions) {
	content.style = style
}

func (content *contentEmpty) dimensions() contentDimensions {
	computed := content.Style().Computed()
	return contentDimensions{
		width:           int(computed.Width),
		height:          int(computed.Height),
		paddingAndGapsY: computed.PaddingTop + computed.PaddingBottom,
		paddingAndGapsX: computed.PaddingLeft + computed.PaddingRight,
	}
}

func (content *contentEmpty) Type() blockContentType {
	return BlockContentTypeEmpty
}

func (content *contentEmpty) Layers() map[int]struct{} {
	return map[int]struct{}{content.style.Computed().ZIndex: {}}
}

func (content *contentEmpty) Style() style.StyleOptions {
	return content.style
}

func (content *contentEmpty) Render(layers layerContext, pos Position) error {
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
	return nil
}
