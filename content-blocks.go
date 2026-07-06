package facepaint

import (
	"errors"
	"maps"

	"github.com/cufee/facepaint/style"
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
		Width:           ceil(computed.Width),
		Height:          ceil(computed.Height),
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

	if dimensions.Width > 0 && dimensions.Height > 0 {
		return dimensions
	}

	// add content dimensions of each block to the total
	var blockWidthTotal, blockWidthMax, blockHeightTotal, blockHeightMax int
	for _, block := range content.value {
		blockDimensions := block.Dimensions()

		if block.Style().Computed().Position == style.PositionAbsolute {
			continue
		}

		blockWidthTotal += blockDimensions.Width
		blockWidthMax = max(blockWidthMax, blockDimensions.Width)

		blockHeightTotal += blockDimensions.Height
		blockHeightMax = max(blockHeightMax, blockDimensions.Height)
	}

	// calculate final block width if it was not set already
	if dimensions.Width == 0 {
		dimensions.Width = ceil(dimensions.paddingAndGapsX)

		switch computed.Direction {
		case style.DirectionHorizontal:
			dimensions.Width += blockWidthTotal

		case style.DirectionVertical:
			dimensions.Width += blockWidthMax
		}
	}
	// calculate final block height if it was not set already
	if dimensions.Height == 0 {
		dimensions.Height = ceil(dimensions.paddingAndGapsY)

		switch computed.Direction {
		case style.DirectionHorizontal:
			dimensions.Height += blockHeightMax
		case style.DirectionVertical:
			dimensions.Height += blockHeightTotal
		}
	}

	dimensions.Width = max(dimensions.Width, ceil(computed.MinWidth))
	dimensions.Height = max(dimensions.Height, ceil(computed.MinHeight))
	return dimensions
}

func (content *contentBlocks) Type() blockContentType {
	return BlockContentTypeBlocks
}

func (content *contentBlocks) Layers() map[int]struct{} {
	var layers = make(map[int]struct{}, len(content.value))
	for _, block := range content.value {
		maps.Copy(layers, block.Layers())
	}
	return layers
}

func (content *contentBlocks) Style() style.StyleOptions {
	return content.style
}

func (content *contentBlocks) Render(layers *layerContext, pos Position) error {
	computed := content.style.Computed()
	dimensions := content.dimensions()

	if computed.Position == style.PositionAbsolute {
		if computed.Left != 0 {
			pos.X += computed.Left
		} else if computed.Right != 0 {
			pos.X += float64(dimensions.Width) - computed.Right
		}
		if computed.Top != 0 {
			pos.Y += computed.Top
		} else if computed.Bottom != 0 {
			pos.Y += float64(dimensions.Height) - computed.Bottom
		}
	}

	registerBackdrop(layers, computed, dimensions, pos)

	return renderWithFilter(layers, computed, dimensions, pos, func(ctx *layer, pos Position) error {
		if computed.BackgroundColor != nil && len(computed.Backdrop) == 0 {
			ctx.SetColor(computed.BackgroundColor)
			drawBackgroundPath(ctx, computed, dimensions, pos)
			ctx.Fill()
		}
		if computed.BackgroundImage != nil {
			background := imaging.Fill(computed.BackgroundImage, dimensions.Width, dimensions.Height, imaging.Center, imaging.Lanczos)
			ctx.DrawImage(background, ceil(pos.X), ceil(pos.Y))
		}

		if computed.Debug {
			ctx.SetColor(getDebugColor())
			ctx.DrawRectangle(pos.X, pos.Y, float64(dimensions.Width), float64(dimensions.Height))
			ctx.Stroke()
		}

		applyBlocksGrowth(computed, dimensions, content.value...)

		var originX, originY = pos.X + computed.PaddingLeft, pos.Y + computed.PaddingTop
		return renderBlocksContent(layers, computed, dimensions, Position{X: originX, Y: originY}, content.value...)
	})
}

func renderBlocksContent(ctx *layerContext, containerStyle style.Style, container contentDimensions, pos Position, blocks ...*Block) error {
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
		contentWidthTotal += blockDimensions.Width
		contentHeightTotal += blockDimensions.Height
	}
	// add gaps as content width
	switch containerStyle.Direction {
	case style.DirectionHorizontal:
		contentWidthTotal += ceil(container.paddingAndGapsX)
	case style.DirectionVertical:
		contentHeightTotal += ceil(container.paddingAndGapsY)
	}

	var lastX, lastY float64 = pos.X, pos.Y
	for i, block := range blocks {
		blockStyle := block.Style().Computed()
		blockSize := block.Dimensions()
		posX, posY := lastX, lastY

		// apply absolute position margins
		if blockStyle.Position == style.PositionAbsolute {
			if blockStyle.Left != 0 {
				posX += blockStyle.Left
			} else if blockStyle.Right != 0 {
				posX += float64(container.Width-blockSize.Width) - blockStyle.Right
			}
			if blockStyle.Top != 0 {
				posY += blockStyle.Top
			} else if blockStyle.Bottom != 0 {
				posY += float64(container.Height-blockSize.Height) - blockStyle.Bottom
			}
		} else {
			switch containerStyle.Direction {
			case style.DirectionVertical:
				switch containerStyle.JustifyContent {
				case style.JustifyContentCenter:
					posY += float64(container.Height-contentHeightTotal) / 2
				case style.JustifyContentEnd:
					posY += float64(container.Height - contentHeightTotal)
				case style.JustifyContentSpaceAround:
					if relativeBlocks == 1 {
						posY += float64(container.Height-contentHeightTotal) / 2
					}
					if relativeBlocks > 1 {
						posY += float64(container.Height-contentHeightTotal) / float64(relativeBlocks+1)
					}
				case style.JustifyContentSpaceBetween:
					if i > 0 && relativeBlocks > 1 {
						posY += float64(container.Height-contentHeightTotal) / float64(relativeBlocks-1)
					}
				}

				switch containerStyle.AlignItems {
				case style.AlignItemsCenter:
					posX += float64(container.Width-ceil(container.paddingX)-blockSize.Width) / 2
				case style.AlignItemsEnd:
					posX += float64(container.Width - ceil(container.paddingX) - blockSize.Width)
				}
			default: // DirectionHorizontal
				switch containerStyle.JustifyContent {
				case style.JustifyContentCenter:
					posX += float64(container.Width-contentWidthTotal) / 2
				case style.JustifyContentEnd:
					posX += float64(container.Width - contentWidthTotal)
				case style.JustifyContentSpaceAround:
					if relativeBlocks == 1 {
						posX += float64(container.Width-contentWidthTotal) / 2
					}
					if relativeBlocks > 1 {
						posX += float64(container.Width-contentWidthTotal) / float64(relativeBlocks+1)
					}
				case style.JustifyContentSpaceBetween:
					if i > 0 && relativeBlocks > 1 {
						posX += float64(container.Width-contentWidthTotal) / float64(relativeBlocks-1)
					}
				}

				switch containerStyle.AlignItems {
				case style.AlignItemsCenter:
					posY += float64(container.Height-ceil(container.paddingY)-blockSize.Height) / 2
				case style.AlignItemsEnd:
					posY += float64(container.Height - ceil(container.paddingY) - blockSize.Height)
				}
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
			lastY = posY + float64(blockSize.Height) + containerStyle.Gap
		default:
			lastX = posX + float64(blockSize.Width) + containerStyle.Gap
		}
	}

	return nil
}

func applyBlocksGrowth(containerStyle style.Style, container contentDimensions, blocks ...*Block) {
	var blockWidthTotal, blockWidthMax, blockHeightTotal, blockHeightMax int
	var growBlocksX, growBlocksY int
	var evenBlocksX, evenBlocksY int

	for _, block := range blocks {
		blockComputed := block.Style().Computed()
		if blockComputed.Position == style.PositionAbsolute {
			continue
		}

		blockDimensions := block.Dimensions()
		blockWidthTotal += blockDimensions.Width
		blockWidthMax = max(blockWidthMax, blockDimensions.Width)
		blockHeightTotal += blockDimensions.Height
		blockHeightMax = max(blockHeightMax, blockDimensions.Height)

		if blockComputed.GrowHorizontal {
			growBlocksX++
			if blockComputed.Basis == style.BasisEven {
				evenBlocksX++
			}
		}
		if blockComputed.GrowVertical {
			growBlocksY++
			if blockComputed.Basis == style.BasisEven {
				evenBlocksY++
			}
		}
	}

	if growBlocksX == 0 && growBlocksY == 0 {
		return
	}

	// apply BasisEven growth first -- even blocks split available space equally
	applyEvenGrowth(containerStyle, container, blocks, evenBlocksX, evenBlocksY)

	// recalculate totals after even growth since block sizes changed
	if evenBlocksX > 0 || evenBlocksY > 0 {
		blockWidthTotal, blockWidthMax, blockHeightTotal, blockHeightMax = 0, 0, 0, 0
		for _, block := range blocks {
			if block.Style().Computed().Position == style.PositionAbsolute {
				continue
			}
			d := block.Dimensions()
			blockWidthTotal += d.Width
			blockWidthMax = max(blockWidthMax, d.Width)
			blockHeightTotal += d.Height
			blockHeightMax = max(blockHeightMax, d.Height)
		}
	}

	// additive growth for BasisNone blocks
	noneGrowX := growBlocksX - evenBlocksX
	noneGrowY := growBlocksY - evenBlocksY
	blockGrowX := max(0, container.Width-ceil(container.paddingAndGapsX)-blockWidthTotal) / max(1, noneGrowX)
	blockGrowY := max(0, container.Height-ceil(container.paddingAndGapsY)-blockHeightTotal) / max(1, noneGrowY)

	for _, block := range blocks {
		blockStyle := block.Style()
		blockComputed := blockStyle.Computed()

		if blockComputed.Position == style.PositionAbsolute {
			applyAbsoluteGrowth(containerStyle, container, block)
			continue
		}
		if !blockComputed.GrowHorizontal && !blockComputed.GrowVertical {
			continue
		}

		blockSize := block.Dimensions()

		switch containerStyle.Direction {
		case style.DirectionHorizontal:
			if blockComputed.GrowHorizontal && blockComputed.Basis == style.BasisNone {
				blockStyle.Merge(style.SetWidth(float64(blockSize.Width) + float64(blockGrowX)))
				block.content.setStyle(blockStyle)
			}
			if blockComputed.GrowVertical {
				blockStyle.Merge(style.SetHeight(float64(blockHeightMax)))
				block.content.setStyle(blockStyle)
			}
		case style.DirectionVertical:
			if blockComputed.GrowHorizontal {
				blockStyle.Merge(style.SetWidth(float64(blockWidthMax)))
				block.content.setStyle(blockStyle)
			}
			if blockComputed.GrowVertical && blockComputed.Basis == style.BasisNone {
				blockStyle.Merge(style.SetHeight(float64(blockSize.Height) + float64(blockGrowY)))
				block.content.setStyle(blockStyle)
			}
		}
	}
}

func applyAbsoluteGrowth(containerStyle style.Style, container contentDimensions, block *Block) {
	blockStyle := block.Style()
	blockComputed := blockStyle.Computed()

	switch containerStyle.Direction {
	case style.DirectionHorizontal:
		if blockComputed.GrowHorizontal {
			blockStyle.Merge(style.SetWidth(float64(container.Width) - containerStyle.PaddingLeft - containerStyle.PaddingRight))
			block.content.setStyle(blockStyle)
		}
	case style.DirectionVertical:
		if blockComputed.GrowVertical {
			blockStyle.Merge(style.SetHeight(float64(container.Height) - containerStyle.PaddingTop - containerStyle.PaddingBottom))
			block.content.setStyle(blockStyle)
		}
	}
}

func applyEvenGrowth(containerStyle style.Style, container contentDimensions, blocks []*Block, evenBlocksX, evenBlocksY int) {
	if evenBlocksX == 0 && evenBlocksY == 0 {
		return
	}

	switch containerStyle.Direction {
	case style.DirectionHorizontal:
		if evenBlocksX > 0 {
			fixedWidth := 0
			for _, block := range blocks {
				c := block.Style().Computed()
				if c.Position == style.PositionAbsolute {
					continue
				}
				if c.GrowHorizontal && c.Basis == style.BasisEven {
					continue
				}
				fixedWidth += block.Dimensions().Width
			}
			evenWidth := max(0, container.Width-ceil(container.paddingAndGapsX)-fixedWidth) / evenBlocksX
			for _, block := range blocks {
				bs := block.Style()
				c := bs.Computed()
				if c.GrowHorizontal && c.Basis == style.BasisEven && c.Position != style.PositionAbsolute {
					bs.Merge(style.SetWidth(float64(evenWidth)))
					block.content.setStyle(bs)
				}
			}
		}

	case style.DirectionVertical:
		if evenBlocksY > 0 {
			fixedHeight := 0
			for _, block := range blocks {
				c := block.Style().Computed()
				if c.Position == style.PositionAbsolute {
					continue
				}
				if c.GrowVertical && c.Basis == style.BasisEven {
					continue
				}
				fixedHeight += block.Dimensions().Height
			}
			evenHeight := max(0, container.Height-ceil(container.paddingAndGapsY)-fixedHeight) / evenBlocksY
			for _, block := range blocks {
				bs := block.Style()
				c := bs.Computed()
				if c.GrowVertical && c.Basis == style.BasisEven && c.Position != style.PositionAbsolute {
					bs.Merge(style.SetHeight(float64(evenHeight)))
					block.content.setStyle(bs)
				}
			}
		}
	}
}
