package facepaint

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/cufee/facepaint/internal/tests"
	"github.com/cufee/facepaint/style"
	"github.com/matryer/is"
)

var _ = saveImage

var contentSize = 12.0
var contentColorAlphaValue uint32
var contentColor = color.RGBA{255, 255, 255, 255}

func init() {
	_, _, _, a := contentColor.RGBA()
	contentColorAlphaValue = a
}

func TestRenderV2(t *testing.T) {
	if os.Getenv("CI") == "true" {
		return // this is a local test for visual debugging
	}

	is := is.New(t)

	text1, err := NewTextContent(style.NewStyle(
		style.Parent(style.Style{
			Left: -5,
			Top:  -5,
			// Blur: 1,
			ZIndex: 1,
		}),
		style.SetDebug(true),
		style.SetPosition(style.PositionAbsolute),
		style.SetFont(tests.Font(), color.Black),
		// style.SetWidth(100),
		// style.SetGrowX(true),
		// style.SetGrowY(true),
	), "TEST - 1")
	is.NoErr(err)

	text2, err := NewTextContent(style.NewStyle(
		// style.SetDebug(true),
		// style.SetGrowX(true),
		style.SetGrowY(true),
		style.SetPadding(10),
		style.SetFont(tests.Font(), color.Black),
	), "TEST - 2")
	is.NoErr(err)

	block1 := NewBlocksContent(style.NewStyle(
		style.Parent(style.Style{
			// JustifyContent: style.JustifyContentCenter,
		}),
		style.SetDebug(true),
		style.SetPadding(20),
		// style.SetGrowX(true),
	), text2)

	block2 := NewBlocksContent(style.NewStyle(
		style.Parent(style.Style{
			Gap: 10,
		}),
		style.SetDebug(true),
		style.SetPadding(10),
		// style.SetWidth(300),
	), text1, block1)

	img, err := block2.Render()
	is.NoErr(err)

	saveImage(is, img)
}

func saveImage(is *is.I, img image.Image) {
	f, err := os.Create(filepath.Join(tests.Root(), "tmp", "test_render_blocks.png"))
	is.NoErr(err)

	err = png.Encode(f, img)
	is.NoErr(err)
}
