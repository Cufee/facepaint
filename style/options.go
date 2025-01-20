package style

import (
	"image"
	"image/color"
)

func NewStyle(opts ...styleOption) StyleOptions {
	return opts
}

type styleOption func(s *Style)

type StyleOptions []styleOption

func (o StyleOptions) Computed() Style {
	var s Style
	for _, apply := range o {
		apply(&s)
	}
	return s
}

func (arr *StyleOptions) Add(opt styleOption) {
	*arr = append(*arr, opt)
}

func Parent(parent Style) styleOption {
	return func(s *Style) { *s = parent }
}

func SetFont(value Font, color color.Color) styleOption {
	return func(s *Style) { s.Font = value; s.Color = color }
}

func SetDebug(value bool) styleOption {
	return func(s *Style) { s.Debug = value }
}

func SetPosition(value positionValue) styleOption {
	return func(s *Style) { s.Position = value }
}

func SetWidth(value float64) styleOption {
	return func(s *Style) { s.Width = value }
}
func SetHeight(value float64) styleOption {
	return func(s *Style) { s.Height = value }
}

func SetMinWidth(value float64) styleOption {
	return func(s *Style) { s.MinWidth = value }
}
func SetMinHeight(value float64) styleOption {
	return func(s *Style) { s.MinHeight = value }
}

func SetZIndex(value int) styleOption {
	return func(s *Style) { s.ZIndex = value }
}

func SetPadding(value float64) styleOption {
	return func(s *Style) {
		s.PaddingLeft = value
		s.PaddingRight = value
		s.PaddingTop = value
		s.PaddingBottom = value
	}
}
func SetPaddingX(value float64) styleOption {
	return func(s *Style) {
		s.PaddingLeft = value
		s.PaddingRight = value
	}
}
func SetPaddingY(value float64) styleOption {
	return func(s *Style) {
		s.PaddingTop = value
		s.PaddingBottom = value
	}
}

func SetGrow(value bool) styleOption {
	return func(s *Style) {
		s.GrowHorizontal = value
		s.GrowVertical = value
	}
}
func SetGrowX(value bool) styleOption {
	return func(s *Style) {
		s.GrowHorizontal = value
	}
}
func SetGrowY(value bool) styleOption {
	return func(s *Style) {
		s.GrowVertical = value
	}
}

func SetBlur(value float64) styleOption {
	return func(s *Style) {
		s.Blur = value
	}
}

func SetBackground(value image.Image) styleOption {
	return func(s *Style) {
		if value != nil {
			s.BackgroundImage = value
		}
	}
}

func SetBorderRadius(value float64) styleOption {
	return func(s *Style) {
		s.BorderRadiusTopLeft = value
		s.BorderRadiusTopRight = value
		s.BorderRadiusBottomLeft = value
		s.BorderRadiusBottomRight = value
	}
}
func SetBorderRadiusLeft(value float64) styleOption {
	return func(s *Style) {
		s.BorderRadiusTopLeft = value
		s.BorderRadiusBottomLeft = value
	}
}
func SetBorderRadiusRight(value float64) styleOption {
	return func(s *Style) {
		s.BorderRadiusTopRight = value
		s.BorderRadiusBottomRight = value
	}
}
func SetBorderRadiusTop(value float64) styleOption {
	return func(s *Style) {
		s.BorderRadiusTopLeft = value
		s.BorderRadiusTopRight = value
	}
}
func SetBorderRadiusBottom(value float64) styleOption {
	return func(s *Style) {
		s.BorderRadiusBottomLeft = value
		s.BorderRadiusBottomRight = value
	}
}
