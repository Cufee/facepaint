package facepaint

import (
	"testing"

	"github.com/cufee/facepaint/internal/tests"
	"github.com/cufee/facepaint/style"
	"github.com/matryer/is"
)

func TestMeasureStringWidth(t *testing.T) {
	is := is.New(t)
	font := tests.Font()

	full := MeasureString("hello", font)
	width := MeasureStringWidth("hello", font)
	is.True(width > 0)
	is.Equal(width, full.TotalWidth)
}

func TestMeasureStringHeight(t *testing.T) {
	is := is.New(t)
	font := tests.Font()

	full := MeasureString("hello", font)
	height := MeasureStringHeight("hello", font)
	is.True(height > 0)
	is.Equal(height, full.TotalHeight)
}

func TestMeasureStringWidthEmpty(t *testing.T) {
	is := is.New(t)
	font := tests.Font()

	is.Equal(MeasureStringWidth("", font), 0.0)
	is.Equal(MeasureStringHeight("", font), 0.0)
}

func TestMaxStringWidth(t *testing.T) {
	is := is.New(t)
	font := tests.Font()

	short := MeasureStringWidth("hi", font)
	long := MeasureStringWidth("hello world", font)
	maxW := MaxStringWidth(font, "hi", "hello world")
	is.Equal(maxW, long)
	is.True(maxW > short)
}

func TestMaxStringWidthEmpty(t *testing.T) {
	is := is.New(t)
	font := tests.Font()

	is.Equal(MaxStringWidth(font), 0.0)
}

func TestMaxStringHeight(t *testing.T) {
	is := is.New(t)
	font := tests.Font()

	single := MeasureStringHeight("hello", font)
	multi := MeasureStringHeight("hello\nworld", font)
	maxH := MaxStringHeight(font, "hello", "hello\nworld")
	is.Equal(maxH, multi)
	is.True(maxH > single)
}

func TestMeasureBlockWidth(t *testing.T) {
	is := is.New(t)
	font := tests.Font()

	s := style.Style{
		Font:         font,
		PaddingLeft:  10,
		PaddingRight: 15,
	}
	textWidth := MeasureStringWidth("test", font)
	blockWidth := MeasureBlockWidth("test", s)
	is.Equal(blockWidth, textWidth+25)
}

func TestMeasureBlockHeight(t *testing.T) {
	is := is.New(t)
	font := tests.Font()

	s := style.Style{
		Font:          font,
		PaddingTop:    5,
		PaddingBottom: 10,
	}
	textHeight := MeasureStringHeight("test", font)
	blockHeight := MeasureBlockHeight("test", s)
	is.Equal(blockHeight, textHeight+15)
}
