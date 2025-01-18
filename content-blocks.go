package facepaint

import (
	"errors"

	"github.com/cufee/facepaint/style"
	"github.com/fogleman/gg"
	"github.com/nao1215/imaging"
)

var _ BlockContent = &contentBlocks{}

func NewBlocksContent(style style.StyleOptions, value ...*Block) *Block {
	return NewBlock(&contentBlocks{
		value: value,
		style: style,
	})
}

type contentBlocks struct {
	style style.StyleOptions
	value []*Block
}

func (content *contentBlocks) setStyle(style style.StyleOptions) {
	content.style = style
}

func (content *contentBlocks) dimensions() contentDimensions {
	if len(content.value) == 0 {
		return contentDimensions{}
	}

	computed := content.style.Computed()
	dimensions := contentDimensions{
		width:           ceil(computed.Width),
		height:          ceil(computed.Height),
		paddingAndGapsX: computed.PaddingLeft + computed.PaddingRight,
		paddingX:        computed.PaddingLeft + computed.PaddingRight,
		paddingAndGapsY: computed.PaddingTop + computed.PaddingBottom,
		paddingY:        computed.PaddingTop + computed.PaddingBottom,
	}

	var relativeBlocks = 0.0
	for _, block := range content.value {
		if block.Style().Computed().Position == style.PositionAbsolute {
			continue
		}
		relativeBlocks++
	}

	switch computed.Direction {
	case style.DirectionHorizontal:
		gaps := max(0, computed.Gap*(relativeBlocks-1))
		dimensions.paddingAndGapsX += gaps
		dimensions.gapsX += gaps

	case style.DirectionVertical:
		gaps := max(0, computed.Gap*(relativeBlocks-1))
		dimensions.paddingAndGapsY += gaps
		dimensions.gapsY += gaps
	}

	if dimensions.width > 0 && dimensions.height > 0 {
		return dimensions
	}

	// add content dimensions of each block to the total
	var blockWidthTotal, blockWidthMax, blockHeightTotal, blockHeightMax int
	for _, block := range content.value {
		blockDimensions := block.Dimensions()

		if block.Style().Computed().Position == style.PositionAbsolute {
			continue
		}

		blockWidthTotal += blockDimensions.width
		blockWidthMax = max(blockWidthMax, blockDimensions.width)

		blockHeightTotal += blockDimensions.height
		blockHeightMax = max(blockHeightMax, blockDimensions.height)
	}

	// calculate final block width if it was not set already
	if dimensions.width == 0 {
		dimensions.width = ceil(dimensions.paddingAndGapsX)

		switch computed.Direction {
		case style.DirectionHorizontal:
			dimensions.width += blockWidthTotal

		case style.DirectionVertical:
			dimensions.width += blockWidthMax
		}
	}
	// calculate final block height if it was not set already
	if dimensions.height == 0 {
		dimensions.height = ceil(dimensions.paddingAndGapsY)

		switch computed.Direction {
		case style.DirectionHorizontal:
			dimensions.height += blockHeightMax
		case style.DirectionVertical:
			dimensions.height += blockHeightTotal
		}
	}

	return dimensions
}

func (content *contentBlocks) Type() blockContentType {
	return BlockContentTypeBlocks
}

func (content *contentBlocks) Layers() map[int]struct{} {
	var layers = make(map[int]struct{}, len(content.value))
	for _, block := range content.value {
		for i, v := range block.Layers() {
			layers[i] = v
		}
	}
	return layers
}

func (content *contentBlocks) Style() style.StyleOptions {
	return content.style
}

func (content *contentBlocks) Render(layers layerContext, pos Position) error {
	computed := content.style.Computed()
	dimensions := content.dimensions()
	ctx, err := layers.layer(computed.ZIndex)
	if err != nil {
		return err
	}

	if computed.Position == style.PositionAbsolute {
		if computed.Left != 0 {
			pos.X += computed.Left
		} else if computed.Right != 0 {
			pos.X += float64(dimensions.width) - computed.Right
		}
		if computed.Top != 0 {
			pos.Y += computed.Top
		} else if computed.Bottom != 0 {
			pos.Y += float64(dimensions.height) - computed.Bottom
		}
	}

	if computed.Blur > 0 {
		// replace the context
		parentPosition := pos
		pos = Position{X: 0, Y: 0}
		ctx = gg.NewContext(dimensions.width, dimensions.height)
		defer func() {
			// blur the result and paste onto the parent layer
			parent, _ := layers.layer(computed.ZIndex)
			img := imaging.Blur(ctx.Image(), computed.Blur)
			parent.DrawImage(img, ceil(parentPosition.X), ceil(parentPosition.Y))
		}()
	}

	if computed.BackgroundColor != nil {
		ctx.SetColor(computed.BackgroundColor)
		drawBackgroundPath(ctx, computed, dimensions, pos)
		ctx.Fill()
	}
	if computed.BackgroundImage != nil {
		background := imaging.Fill(computed.BackgroundImage, dimensions.width, dimensions.height, imaging.Center, imaging.Lanczos)
		ctx.DrawImage(background, ceil(pos.X), ceil(pos.Y))
	}

	if computed.Debug {
		ctx.SetColor(getDebugColor())
		ctx.DrawRectangle(pos.X, pos.Y, float64(dimensions.width), float64(dimensions.height))
		ctx.Stroke()
	}

	applyBlocksGrowth(computed, dimensions, content.value...)

	var originX, originY = pos.X + computed.PaddingLeft, pos.Y + computed.PaddingTop
	return renderBlocksContent(layers, computed, dimensions, Position{X: originX, Y: originY}, content.value...)
}

func renderBlocksContent(ctx layerContext, containerStyle style.Style, container contentDimensions, pos Position, blocks ...*Block) error {
	if len(blocks) < 1 {
		return errors.New("no blocks to render")
	}

	// calculate true content dimensions
	var relativeBlocks int
	var contentWidthTotal, contentHeightTotal int
	for _, block := range blocks {
		if block.Style().Computed().Position == style.PositionAbsolute {
			continue
		}
		relativeBlocks++

		blockDimensions := block.Dimensions()
		contentWidthTotal += blockDimensions.width
		contentHeightTotal += blockDimensions.height
	}
	// add gaps as content width
	switch containerStyle.Direction {
	case style.DirectionHorizontal:
		contentWidthTotal += ceil(container.paddingAndGapsX)
	case style.DirectionVertical:
		contentHeightTotal += ceil(container.paddingAndGapsY)
	}

	var lastX, lastY float64 = pos.X, pos.Y
	for _, block := range blocks {
		blockStyle := block.Style().Computed()
		blockSize := block.Dimensions()
		posX, posY := lastX, lastY

		// apply absolute position margins
		if blockStyle.Position == style.PositionAbsolute {
			if blockStyle.Left != 0 {
				posX += blockStyle.Left
			} else if blockStyle.Right != 0 {
				posX += float64(container.width-blockSize.width) - blockStyle.Right
			}
			if blockStyle.Top != 0 {
				posY += blockStyle.Top
			} else if blockStyle.Bottom != 0 {
				posY += float64(container.height-blockSize.height) - blockStyle.Bottom
			}
		}

		switch containerStyle.Direction {
		case style.DirectionVertical:
			// align content vertically
			switch containerStyle.JustifyContent {
			case style.JustifyContentCenter:
				posY += float64(container.height-contentHeightTotal) / 2
			case style.JustifyContentEnd:
				posY += float64(container.height - contentHeightTotal)
			case style.JustifyContentSpaceAround:
				if relativeBlocks > 0 {
					posY += float64((container.height - contentHeightTotal) / (relativeBlocks + 1))
				}
			case style.JustifyContentSpaceBetween:
				if relativeBlocks > 0 {
					posY += float64((container.height - contentHeightTotal) / (relativeBlocks - 1))
				}
			}

			// align content horizontally
			switch containerStyle.AlignItems {
			case style.AlignItemsCenter:
				posX += float64(container.width-ceil(container.paddingX)-blockSize.width) / 2
			case style.AlignItemsEnd:
				posX += float64(container.width - ceil(container.paddingX) - blockSize.width)
			}
		default: // DirectionHorizontal
			// align content horizontally
			switch containerStyle.JustifyContent {
			case style.JustifyContentCenter:
				posX += float64(container.width-contentWidthTotal) / 2
			case style.JustifyContentEnd:
				posX += float64(container.width - contentWidthTotal)
			case style.JustifyContentSpaceAround:
				if relativeBlocks > 0 {
					posX += float64((container.width - contentWidthTotal) / (relativeBlocks + 1))
				}
			case style.JustifyContentSpaceBetween:
				if relativeBlocks > 0 {
					posX += float64((container.width - contentWidthTotal) / (relativeBlocks - 1))
				}
			}

			// align content vertically
			switch containerStyle.AlignItems {
			case style.AlignItemsCenter:
				posY += float64(container.height-ceil(container.paddingY)-blockSize.height) / 2
			case style.AlignItemsEnd:
				posY += float64(container.height - ceil(container.paddingY) - blockSize.height)
			}

		}

		err := block.content.Render(ctx, Position{posX, posY})
		if err != nil {
			return err
		}

		if block.Style().Computed().Position == style.PositionAbsolute {
			continue
		}

		// save the position we rendered at
		switch containerStyle.Direction {
		case style.DirectionVertical:
			lastY = posY + float64(blockSize.height) + containerStyle.Gap
		default:
			lastX = posX + float64(blockSize.width) + containerStyle.Gap
		}
	}

	return nil
}

func applyBlocksGrowth(containerStyle style.Style, container contentDimensions, blocks ...*Block) {
	// calculate content dimensions before growth
	var blockWidthTotal, blockWidthMax, blockHeightTotal, blockHeightMax int
	var growBlocksX, growBlocksY = 0, 0
	for _, block := range blocks {
		blockDimensions := block.Dimensions()

		blockWidthTotal += blockDimensions.width
		blockWidthMax = max(blockWidthMax, blockDimensions.width)

		blockHeightTotal += blockDimensions.height
		blockHeightMax = max(blockHeightMax, blockDimensions.height)

		blockStyle := block.Style().Computed()
		switch {
		case blockStyle.Position == style.PositionAbsolute:
			// absolute blocks do not "consume" grow space
		case blockStyle.GrowHorizontal:
			growBlocksX++
		case blockStyle.GrowVertical:
			growBlocksY++
		}
	}

	blockGrowX := max(0, container.width-ceil(container.paddingAndGapsX)-blockWidthTotal) / max(1, growBlocksX)
	blockGrowY := max(0, container.height-ceil(container.paddingAndGapsY)-blockHeightTotal) / max(1, growBlocksY)

	// apply growth to blocks
	for _, block := range blocks {
		blockStyle := block.Style()
		blockComputed := blockStyle.Computed()
		blockSize := block.Dimensions()

		if !blockComputed.GrowHorizontal && !blockComputed.GrowVertical {
			continue
		}

		switch containerStyle.Direction {
		case style.DirectionHorizontal:
			// update the block width
			if blockComputed.GrowHorizontal && blockComputed.Position == style.PositionAbsolute {
				blockStyle.Add(style.SetWidth(float64(container.width) - containerStyle.PaddingLeft - containerStyle.PaddingRight))
				block.content.setStyle(blockStyle)
			} else if blockComputed.GrowHorizontal {
				blockStyle.Add(style.SetWidth(float64(blockSize.width) + float64(blockGrowX)))
				block.content.setStyle(blockStyle)
			}
			// update the block height
			if blockComputed.GrowVertical {
				blockStyle.Add(style.SetHeight(float64(blockHeightMax)))
				block.content.setStyle(blockStyle)
			}
		case style.DirectionVertical:
			// update the block width
			if blockComputed.GrowHorizontal {
				blockStyle.Add(style.SetWidth(float64(blockWidthMax)))
				block.content.setStyle(blockStyle)
			}
			// update the block height
			if blockComputed.GrowVertical && blockComputed.Position == style.PositionAbsolute {
				blockStyle.Add(style.SetWidth(float64(container.height) - containerStyle.PaddingTop - containerStyle.PaddingBottom))
				block.content.setStyle(blockStyle)
			} else if blockComputed.GrowVertical {
				blockStyle.Add(style.SetHeight(float64(blockSize.height) + float64(blockGrowY)))
				block.content.setStyle(blockStyle)
			}
		}
	}
}
