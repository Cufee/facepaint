package style

import (
	"image"
	"image/color"

	"github.com/nao1215/imaging"
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

type positionAnchor byte

const (
	BackgroundCenter positionAnchor = iota

	BackgroundLeft
	BackgroundRight

	BackgroundTop
	BackgroundBottom

	BackgroundTopLeft
	BackgroundTopRight

	BackgroundBottomLeft
	BackgroundBottomRight
)

func (p positionAnchor) Imaging() imaging.Anchor {
	switch p {
	default:
		return imaging.Center

	case BackgroundLeft:
		return imaging.Left
	case BackgroundRight:
		return imaging.Right
	case BackgroundTop:
		return imaging.Top
	case BackgroundBottom:
		return imaging.Bottom
	case BackgroundTopLeft:
		return imaging.TopLeft
	case BackgroundTopRight:
		return imaging.TopRight
	case BackgroundBottomLeft:
		return imaging.BottomLeft
	case BackgroundBottomRight:
		return imaging.BottomRight
	}
}

type BasisValue byte

const (
	BasisNone BasisValue = iota // default: additive growth
	BasisEven                   // equal share of available space
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

	MinWidth  float64
	MinHeight float64

	Font Font

	Color              color.Color
	BackgroundColor    color.Color
	BackgroundImage    image.Image
	BackgroundPosition positionAnchor

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
	Basis          BasisValue

	BorderRadiusTopLeft     float64
	BorderRadiusTopRight    float64
	BorderRadiusBottomLeft  float64
	BorderRadiusBottomRight float64

	// Backdrop effects run on the frame behind the block before compositing.
	Backdrop []Effect
	// Filter effects run on the block's content before compositing.
	Filter []Effect

	ZIndex int
}

func (s *Style) Options() StyleOptions {
	return StyleOptions{Parent(*s)}
}
