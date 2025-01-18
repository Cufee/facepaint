package render

import (
	"fmt"
	"image"

	"github.com/cufee/facepaint/style"
	"github.com/fogleman/gg"
)

func NewBlock(content BlockContent) *Block {
	return &Block{
		content: content,
	}
}

type blockContentType int

func (t blockContentType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", t)), nil
}
func (t blockContentType) String() string {
	return fmt.Sprintf("%d", t)
}

const (
	BlockContentTypeEmpty blockContentType = iota
	BlockContentTypeBlocks
	BlockContentTypeImage
	BlockContentTypeText
)

type Position struct {
	X float64
	Y float64
}

type BlockContent interface {
	Type() blockContentType

	// Renders the block onto an image
	Render(layerContext, Position) error

	Style() style.StyleOptions
	setStyle(style.StyleOptions)

	Layers() map[int]struct{}

	// returns final block image dimensions without rendering
	dimensions() contentDimensions
}

type Block struct {
	content BlockContent
}

func (b *Block) Layers() map[int]struct{} {
	return b.content.Layers()
}

func (b *Block) Style() style.StyleOptions {
	return b.content.Style()
}

func (b *Block) Type() blockContentType {
	return b.content.Type()
}

func (b *Block) Render() (image.Image, error) {
	dimensions := b.Dimensions()

	layers := b.Layers()
	ctx := make(layerContext, len(layers))
	for idx := range layers {
		ctx[idx] = gg.NewContext(dimensions.width, dimensions.height)
	}

	err := b.content.Render(ctx, Position{0, 0})
	if err != nil {
		return nil, err
	}
	return ctx.Image(), nil
}

func (b *Block) Dimensions() contentDimensions {
	return b.content.dimensions()
}

type contentDimensions struct {
	width           int
	height          int
	paddingAndGapsX float64
	paddingAndGapsY float64
}
