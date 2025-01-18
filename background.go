package facepaint

import (
	"github.com/cufee/facepaint/style"
	"github.com/fogleman/gg"
)

func drawBackgroundPath(ctx *gg.Context, style style.Style, dimensions contentDimensions, pos Position) {
	width, height := float64(dimensions.width), float64(dimensions.height)

	// Offset the position
	offsetX, offsetY := pos.X, pos.Y

	// Top-left corner
	topLeftRadius := style.BorderRadiusTopLeft
	ctx.MoveTo(offsetX+topLeftRadius, offsetY)
	ctx.LineTo(offsetX, offsetY+topLeftRadius)
	ctx.DrawEllipticalArc(offsetX+topLeftRadius, offsetY+topLeftRadius, topLeftRadius, topLeftRadius, gg.Radians(270), gg.Radians(180))

	// Bottom-left corner
	bottomLeftRadius := style.BorderRadiusBottomLeft
	ctx.LineTo(offsetX, offsetY+height-bottomLeftRadius)
	ctx.DrawEllipticalArc(offsetX+bottomLeftRadius, offsetY+height-bottomLeftRadius, bottomLeftRadius, bottomLeftRadius, gg.Radians(180), gg.Radians(90))

	// Bottom-right corner
	bottomRightRadius := style.BorderRadiusBottomRight
	ctx.LineTo(offsetX+width-bottomRightRadius, offsetY+height)
	ctx.DrawEllipticalArc(offsetX+width-bottomRightRadius, offsetY+height-bottomRightRadius, bottomRightRadius, bottomRightRadius, gg.Radians(90), gg.Radians(0))

	// Top-right corner
	topRightRadius := style.BorderRadiusTopRight
	ctx.LineTo(offsetX+width, offsetY+topRightRadius)
	ctx.DrawEllipticalArc(offsetX+width-topRightRadius, offsetY+topRightRadius, topRightRadius, topRightRadius, gg.Radians(0), gg.Radians(-90))

	// Close the path
	ctx.LineTo(offsetX+topLeftRadius, offsetY)
	ctx.ClosePath()

	// // top left
	// ctx.DrawEllipticalArc(borderRadius, borderRadius, borderRadius, borderRadius, gg.Radians(270), gg.Radians(180))

	// // bottom left
	// ctx.LineTo(0, height-borderRadius)
	// ctx.DrawEllipticalArc(borderRadius, height-borderRadius, borderRadius, borderRadius, gg.Radians(180), gg.Radians(90))

	// // bottom right
	// ctx.LineTo(width-borderRadius, height)
	// ctx.DrawEllipticalArc(width-borderRadius, height-borderRadius, borderRadius, borderRadius, gg.Radians(90), gg.Radians(0))

	// // top right
	// ctx.LineTo(width, borderRadius)
	// ctx.DrawEllipticalArc(width-borderRadius, borderRadius, borderRadius, borderRadius, gg.Radians(0), gg.Radians(-90))

}
