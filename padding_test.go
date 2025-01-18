package facepaint

import (
	"testing"

	"github.com/cufee/facepaint/style"
	"github.com/matryer/is"
)

func TestApplyPadding(t *testing.T) {
	is := is.New(t)

	content := NewEmptyContent(style.NewStyle(style.Parent(style.Style{Width: contentSize, Height: contentSize, BackgroundColor: contentColor})))

	t.Run("uniform", func(t *testing.T) {
		is := is.New(t)

		wrapper := NewBlocksContent(style.NewStyle(
			style.SetPadding(10),
		), content)

		d := wrapper.Dimensions()
		is.True(d.Width == ceil(contentSize)+20)
		is.True(d.Height == ceil(contentSize)+20)

		img, err := wrapper.Render()
		is.NoErr(err)

		{
			_, _, _, a := img.At(9, 9).RGBA()
			is.True(a == 0)
		}
		{
			_, _, _, a := img.At(10, 10).RGBA()
			is.True(a == contentColorAlphaValue)
		}
	})

	t.Run("X", func(t *testing.T) {
		wrapper := NewBlocksContent(style.NewStyle(
			style.SetPaddingX(10),
		), content)

		d := wrapper.Dimensions()
		is.True(d.Width == ceil(contentSize)+20)
		is.True(d.Height == ceil(contentSize))
	})

	t.Run("Y", func(t *testing.T) {
		wrapper := NewBlocksContent(style.NewStyle(
			style.SetPaddingY(10),
		), content)

		d := wrapper.Dimensions()
		is.True(d.Width == ceil(contentSize))
		is.True(d.Height == ceil(contentSize)+20)
	})

	t.Run("overwrite", func(t *testing.T) {
		wrapper := NewBlocksContent(style.NewStyle(
			style.SetPadding(10),
			style.SetPadding(0),
		), content)

		d := wrapper.Dimensions()
		is.True(d.Width == ceil(contentSize))
		is.True(d.Height == ceil(contentSize))
	})

	t.Run("left", func(t *testing.T) {
		wrapper := NewBlocksContent(style.NewStyle(
			style.Parent(style.Style{
				PaddingLeft: 10,
			}),
		), content)

		d := wrapper.Dimensions()
		is.True(d.Width == ceil(contentSize)+10)
		is.True(d.Height == ceil(contentSize))

		img, err := wrapper.Render()
		is.NoErr(err)

		{
			_, _, _, a := img.At(9, 0).RGBA()
			is.True(a == 0)
		}
		{
			_, _, _, a := img.At(10, 0).RGBA()
			is.True(a == contentColorAlphaValue)
		}
	})

	t.Run("top", func(t *testing.T) {
		wrapper := NewBlocksContent(style.NewStyle(
			style.Parent(style.Style{
				PaddingTop: 10,
			}),
		), content)

		d := wrapper.Dimensions()
		is.True(d.Width == ceil(contentSize))
		is.True(d.Height == ceil(contentSize)+10)

		img, err := wrapper.Render()
		is.NoErr(err)

		{
			_, _, _, a := img.At(0, 9).RGBA()
			is.True(a == 0)
		}
		{
			_, _, _, a := img.At(0, 10).RGBA()
			is.True(a == contentColorAlphaValue)
		}
	})
}
