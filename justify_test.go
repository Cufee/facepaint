package facepaint

import (
	"testing"

	"github.com/cufee/facepaint/style"
	"github.com/matryer/is"
)

func TestRenderJustify(t *testing.T) {
	is := is.New(t)

	content := NewEmptyContent(style.NewStyle(style.Parent(style.Style{Width: contentSize, Height: contentSize, BackgroundColor: contentColor})))

	t.Run("horizontal", func(t *testing.T) {
		t.Run("start", func(t *testing.T) {
			is := is.New(t)

			wrapper := NewBlocksContent(style.NewStyle(
				style.SetWidth(contentSize*2),
			), content)

			img, err := wrapper.Render()
			is.NoErr(err)

			{
				_, _, _, imgA := img.At(0, 0).RGBA()
				is.True(imgA == contentColorAlphaValue)
			}
			{
				_, _, _, imgA := img.At(int(contentSize*2-1), 0).RGBA()
				is.True(imgA == 0)
			}
		})

		t.Run("center", func(t *testing.T) {
			is := is.New(t)

			wrapper := NewBlocksContent(style.NewStyle(
				style.Parent(style.Style{
					JustifyContent: style.JustifyContentCenter,
				}),
				style.SetWidth(contentSize*2),
			), content)

			img, err := wrapper.Render()
			is.NoErr(err)

			{
				_, _, _, imgA := img.At(int(contentSize/3), 0).RGBA()
				is.True(imgA == 0)
			}
			{
				_, _, _, imgA := img.At(int(contentSize), 0).RGBA()
				is.True(imgA == contentColorAlphaValue)
			}
			{
				_, _, _, imgA := img.At(int(contentSize*2-contentSize/3), 0).RGBA()
				is.True(imgA == 0)
			}
		})

		t.Run("end", func(t *testing.T) {
			is := is.New(t)

			wrapper := NewBlocksContent(style.NewStyle(
				style.Parent(style.Style{
					JustifyContent: style.JustifyContentEnd,
				}),
				style.SetWidth(contentSize*2),
			), content)

			img, err := wrapper.Render()
			is.NoErr(err)

			{
				_, _, _, imgA := img.At(int(contentSize-1), 0).RGBA()
				is.True(imgA == 0)
			}
			{
				_, _, _, imgA := img.At(int(contentSize*2-1), 0).RGBA()
				is.True(imgA == contentColorAlphaValue)
			}
		})
	})

	t.Run("vertical", func(t *testing.T) {
		t.Run("start", func(t *testing.T) {
			is := is.New(t)

			wrapper := NewBlocksContent(style.NewStyle(
				style.Parent(style.Style{
					Direction: style.DirectionVertical,
				}),
				style.SetHeight(contentSize*2),
			), content)

			img, err := wrapper.Render()
			is.NoErr(err)

			{
				_, _, _, imgA := img.At(0, 0).RGBA()
				is.True(imgA == contentColorAlphaValue)
			}
			{
				_, _, _, imgA := img.At(0, int(contentSize*2-1)).RGBA()
				is.True(imgA == 0)
			}
		})

		t.Run("center", func(t *testing.T) {
			is := is.New(t)

			wrapper := NewBlocksContent(style.NewStyle(
				style.Parent(style.Style{
					JustifyContent: style.JustifyContentCenter,
					Direction:      style.DirectionVertical,
				}),
				style.SetHeight(contentSize*2),
			), content)

			img, err := wrapper.Render()
			is.NoErr(err)

			{
				_, _, _, imgA := img.At(0, int(contentSize/4)).RGBA()
				is.True(imgA == 0)
			}
			{
				_, _, _, imgA := img.At(0, int(contentSize)).RGBA()
				is.True(imgA == contentColorAlphaValue)
			}
			{
				_, _, _, imgA := img.At(0, int(contentSize*2-contentSize/4)).RGBA()
				is.True(imgA == 0)
			}
		})

		t.Run("end", func(t *testing.T) {
			is := is.New(t)

			wrapper := NewBlocksContent(style.NewStyle(
				style.Parent(style.Style{
					JustifyContent: style.JustifyContentEnd,
					Direction:      style.DirectionVertical,
				}),
				style.SetHeight(contentSize*2),
			), content)

			img, err := wrapper.Render()
			is.NoErr(err)

			{
				_, _, _, imgA := img.At(0, int(contentSize-1)).RGBA()
				is.True(imgA == 0)
			}
			{
				_, _, _, imgA := img.At(0, int(contentSize*2-1)).RGBA()
				is.True(imgA == contentColorAlphaValue)
			}
		})
	})
}
