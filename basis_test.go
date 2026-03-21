package facepaint

import (
	"testing"

	"github.com/cufee/facepaint/style"
	"github.com/matryer/is"
)

func TestBasisEvenHorizontal(t *testing.T) {
	t.Run("equal width blocks", func(t *testing.T) {
		is := is.New(t)

		block1 := NewEmptyContent(style.NewStyle(style.Parent(style.Style{
			Width: 10, Height: contentSize,
			GrowHorizontal: true,
			Basis:          style.BasisEven,
		})))
		block2 := NewEmptyContent(style.NewStyle(style.Parent(style.Style{
			Width: 30, Height: contentSize,
			GrowHorizontal: true,
			Basis:          style.BasisEven,
		})))

		wrapper := NewBlocksContent(style.NewStyle(
			style.SetWidth(100),
		), block1, block2)

		d := wrapper.Dimensions()
		is.Equal(d.Width, 100)

		_, err := wrapper.Render()
		is.NoErr(err)

		// after render, both blocks should have equal width: 100 / 2 = 50
		is.Equal(block1.Dimensions().Width, 50)
		is.Equal(block2.Dimensions().Width, 50)
	})

	t.Run("even blocks with gap", func(t *testing.T) {
		is := is.New(t)

		block1 := NewEmptyContent(style.NewStyle(style.Parent(style.Style{
			Width: 10, Height: contentSize,
			GrowHorizontal: true,
			Basis:          style.BasisEven,
		})))
		block2 := NewEmptyContent(style.NewStyle(style.Parent(style.Style{
			Width: 10, Height: contentSize,
			GrowHorizontal: true,
			Basis:          style.BasisEven,
		})))

		wrapper := NewBlocksContent(style.NewStyle(
			style.Parent(style.Style{Gap: 20}),
			style.SetWidth(100),
		), block1, block2)

		_, err := wrapper.Render()
		is.NoErr(err)

		// available = 100 - 20 (gap) = 80, each gets 40
		is.Equal(block1.Dimensions().Width, 40)
		is.Equal(block2.Dimensions().Width, 40)
	})

	t.Run("even blocks with padding", func(t *testing.T) {
		is := is.New(t)

		block1 := NewEmptyContent(style.NewStyle(style.Parent(style.Style{
			Width: 10, Height: contentSize,
			GrowHorizontal: true,
			Basis:          style.BasisEven,
		})))
		block2 := NewEmptyContent(style.NewStyle(style.Parent(style.Style{
			Width: 10, Height: contentSize,
			GrowHorizontal: true,
			Basis:          style.BasisEven,
		})))

		wrapper := NewBlocksContent(style.NewStyle(
			style.SetWidth(100),
			style.SetPaddingX(10),
		), block1, block2)

		_, err := wrapper.Render()
		is.NoErr(err)

		// available = 100 - 20 (padding) = 80, each gets 40
		is.Equal(block1.Dimensions().Width, 40)
		is.Equal(block2.Dimensions().Width, 40)
	})

	t.Run("mixed even and non-grow", func(t *testing.T) {
		is := is.New(t)

		fixedBlock := NewEmptyContent(style.NewStyle(style.Parent(style.Style{
			Width: 20, Height: contentSize,
		})))
		evenBlock1 := NewEmptyContent(style.NewStyle(style.Parent(style.Style{
			Width: 5, Height: contentSize,
			GrowHorizontal: true,
			Basis:          style.BasisEven,
		})))
		evenBlock2 := NewEmptyContent(style.NewStyle(style.Parent(style.Style{
			Width: 15, Height: contentSize,
			GrowHorizontal: true,
			Basis:          style.BasisEven,
		})))

		wrapper := NewBlocksContent(style.NewStyle(
			style.SetWidth(100),
		), fixedBlock, evenBlock1, evenBlock2)

		_, err := wrapper.Render()
		is.NoErr(err)

		// fixed takes 20, even blocks split remaining 80: 40 each
		is.Equal(fixedBlock.Dimensions().Width, 20)
		is.Equal(evenBlock1.Dimensions().Width, 40)
		is.Equal(evenBlock2.Dimensions().Width, 40)
	})

	t.Run("absolute blocks excluded from even growth", func(t *testing.T) {
		is := is.New(t)

		absBlock := NewEmptyContent(style.NewStyle(style.Parent(style.Style{
			Width: 50, Height: contentSize,
			Position: style.PositionAbsolute,
		})))
		evenBlock1 := NewEmptyContent(style.NewStyle(style.Parent(style.Style{
			Width: 10, Height: contentSize,
			GrowHorizontal: true,
			Basis:          style.BasisEven,
		})))
		evenBlock2 := NewEmptyContent(style.NewStyle(style.Parent(style.Style{
			Width: 10, Height: contentSize,
			GrowHorizontal: true,
			Basis:          style.BasisEven,
		})))

		wrapper := NewBlocksContent(style.NewStyle(
			style.SetWidth(100),
		), absBlock, evenBlock1, evenBlock2)

		_, err := wrapper.Render()
		is.NoErr(err)

		// absolute block doesn't count; even blocks split full 100: 50 each
		is.Equal(evenBlock1.Dimensions().Width, 50)
		is.Equal(evenBlock2.Dimensions().Width, 50)
	})
}

func TestBasisEvenVertical(t *testing.T) {
	t.Run("equal height blocks", func(t *testing.T) {
		is := is.New(t)

		block1 := NewEmptyContent(style.NewStyle(style.Parent(style.Style{
			Width: contentSize, Height: 10,
			GrowVertical: true,
			Basis:        style.BasisEven,
		})))
		block2 := NewEmptyContent(style.NewStyle(style.Parent(style.Style{
			Width: contentSize, Height: 30,
			GrowVertical: true,
			Basis:        style.BasisEven,
		})))

		wrapper := NewBlocksContent(style.NewStyle(
			style.Parent(style.Style{Direction: style.DirectionVertical}),
			style.SetHeight(100),
		), block1, block2)

		_, err := wrapper.Render()
		is.NoErr(err)

		// both should have equal height: 100 / 2 = 50
		is.Equal(block1.Dimensions().Height, 50)
		is.Equal(block2.Dimensions().Height, 50)
	})

	t.Run("mixed even and fixed", func(t *testing.T) {
		is := is.New(t)

		fixedBlock := NewEmptyContent(style.NewStyle(style.Parent(style.Style{
			Width: contentSize, Height: 40,
		})))
		evenBlock1 := NewEmptyContent(style.NewStyle(style.Parent(style.Style{
			Width: contentSize, Height: 5,
			GrowVertical: true,
			Basis:        style.BasisEven,
		})))
		evenBlock2 := NewEmptyContent(style.NewStyle(style.Parent(style.Style{
			Width: contentSize, Height: 5,
			GrowVertical: true,
			Basis:        style.BasisEven,
		})))

		wrapper := NewBlocksContent(style.NewStyle(
			style.Parent(style.Style{Direction: style.DirectionVertical}),
			style.SetHeight(100),
		), fixedBlock, evenBlock1, evenBlock2)

		_, err := wrapper.Render()
		is.NoErr(err)

		// fixed takes 40, even blocks split remaining 60: 30 each
		is.Equal(fixedBlock.Dimensions().Height, 40)
		is.Equal(evenBlock1.Dimensions().Height, 30)
		is.Equal(evenBlock2.Dimensions().Height, 30)
	})
}
