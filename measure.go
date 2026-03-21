package facepaint

import (
	"strings"

	"github.com/cufee/facepaint/style"
	"github.com/fogleman/gg"
)

type StringSize struct {
	TotalWidth  float64
	TotalHeight float64
	LineOffset  float64
	LineHeight  float64
}

func MeasureString(text string, font style.Font) StringSize {
	if !font.Valid() {
		return StringSize{}
	}
	if text == "" {
		return StringSize{}
	}

	face, close := font.Face()
	defer close()

	measureCtx := gg.NewContext(1, 1)
	measureCtx.SetFontFace(face)

	var result StringSize
	// Account for font descender height
	result.LineOffset = float64(face.Metrics().Descent>>6) * 2

	for _, line := range strings.Split(text, "\n") {
		w, h := measureCtx.MeasureString(line)
		h += result.LineOffset
		w += 1

		if w > result.TotalWidth {
			result.TotalWidth = w
		}
		if h > result.LineHeight {
			result.LineHeight = h
		}

		result.TotalHeight += h
	}

	return result
}

// MeasureStringWidth returns the TotalWidth of a measured string.
func MeasureStringWidth(text string, font style.Font) float64 {
	return MeasureString(text, font).TotalWidth
}

// MeasureStringHeight returns the TotalHeight of a measured string.
func MeasureStringHeight(text string, font style.Font) float64 {
	return MeasureString(text, font).TotalHeight
}

// MaxStringWidth returns the maximum TotalWidth across all provided texts measured with the given font.
func MaxStringWidth(font style.Font, texts ...string) float64 {
	var maxWidth float64
	for _, text := range texts {
		maxWidth = max(maxWidth, MeasureStringWidth(text, font))
	}
	return maxWidth
}

// MaxStringHeight returns the maximum TotalHeight across all provided texts measured with the given font.
func MaxStringHeight(font style.Font, texts ...string) float64 {
	var maxHeight float64
	for _, text := range texts {
		maxHeight = max(maxHeight, MeasureStringHeight(text, font))
	}
	return maxHeight
}

// MeasureBlockWidth returns the width a text block would occupy with the given style (text width + horizontal padding).
func MeasureBlockWidth(text string, s style.Style) float64 {
	return MeasureStringWidth(text, s.Font) + s.PaddingLeft + s.PaddingRight
}

// MeasureBlockHeight returns the height a text block would occupy with the given style (text height + vertical padding).
func MeasureBlockHeight(text string, s style.Style) float64 {
	return MeasureStringHeight(text, s.Font) + s.PaddingTop + s.PaddingBottom
}
