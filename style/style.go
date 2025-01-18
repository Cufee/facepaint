package style

import (
	"image"
	"image/color"
)

type alignItemsValue byte
type justifyContentValue byte

const (
	AlignItemsStart alignItemsValue = iota
	AlignItemsCenter
	AlignItemsEnd

	JustifyContentStart justifyContentValue = iota
	JustifyContentCenter
	JustifyContentEnd
	JustifyContentSpaceBetween // Spacing between each element is the same
	JustifyContentSpaceAround  // Spacing around all element is the same
)

type directionValue byte

const (
	DirectionHorizontal directionValue = iota
	DirectionVertical
)

type positionValue byte

const (
	PositionRelative positionValue = iota
	PositionAbsolute
)

type overflowValue byte

const (
	OverflowVisible overflowValue = iota
	OverflowHidden
)

type Style struct {
	Debug bool

	Width  float64
	Height float64

	Blur float64

	Font Font

	Color           color.Color
	BackgroundColor color.Color
	BackgroundImage image.Image

	Overflow overflowValue

	JustifyContent justifyContentValue
	AlignItems     alignItemsValue // Depends on Direction
	Direction      directionValue
	Position       positionValue

	Gap float64

	PaddingLeft   float64
	PaddingRight  float64
	PaddingTop    float64
	PaddingBottom float64

	Left   float64
	Right  float64
	Top    float64
	Bottom float64

	GrowHorizontal bool
	GrowVertical   bool

	BorderRadiusTopLeft     float64
	BorderRadiusTopRight    float64
	BorderRadiusBottomLeft  float64
	BorderRadiusBottomRight float64

	ZIndex int
}

func (s *Style) Options() StyleOptions {
	return StyleOptions{Parent(*s)}
}
