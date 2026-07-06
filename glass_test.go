package facepaint

import (
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/cufee/facepaint/internal/tests"
	"github.com/cufee/facepaint/style"
	"github.com/matryer/is"
)

func TestGlassVisual(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("local visual test")
	}
	is := is.New(t)

	bg := tests.SplitBackground(700, 400, 20)

	card := NewBlocksContent(style.NewStyle(
		style.Parent(style.Style{
			Width:  500,
			Height: 300,
			ZIndex: 1,
		}),
		style.SetBorderRadius(25),
		style.SetBackdrop(NewGlass()),
		func(s *style.Style) { s.BackgroundColor = color.NRGBA{0, 0, 0, 80} },
	), MustNewTextContent(style.NewStyle(
		style.Parent(style.Style{ZIndex: 1}),
		style.SetFont(tests.Font(), color.RGBA{220, 220, 230, 255}),
	), "Glass"))

	bgBlock, _ := NewImageContent(style.NewStyle(
		style.Parent(style.Style{
			Width:    700,
			Height:   400,
			ZIndex:   0,
			Position: style.PositionAbsolute,
		}),
	), bg)

	root := NewBlocksContent(style.NewStyle(
		style.Parent(style.Style{
			Width:          700,
			Height:         400,
			ZIndex:         0,
			Direction:      style.DirectionHorizontal,
			AlignItems:     style.AlignItemsCenter,
			JustifyContent: style.JustifyContentCenter,
		}),
	), bgBlock, card)

	img, err := root.Render()
	is.NoErr(err)

	dir := filepath.Join(tests.Root(), "tmp")
	is.NoErr(os.MkdirAll(dir, 0o755))
	f, err := os.Create(filepath.Join(dir, "glass.png"))
	is.NoErr(err)
	defer f.Close()
	is.NoErr(png.Encode(f, img))
}

func TestGlassRefraction(t *testing.T) {
	dims := contentDimensions{Width: 200, Height: 100}
	st := style.Style{
		BorderRadiusTopLeft: 20, BorderRadiusTopRight: 20,
		BorderRadiusBottomLeft: 20, BorderRadiusBottomRight: 20,
	}
	sdf := buildRoundedBoxSDF(dims, st)
	if sdf == nil {
		t.Fatal("sdf is nil")
	}

	glass := NewGlass()
	R := glass.Width

	edgeDX, edgeDY := computeDisplacement(sdf, 100, 3, R, glass)
	centerDX, centerDY := computeDisplacement(sdf, 100, 50, R, glass)
	edgeMag := edgeDX*edgeDX + edgeDY*edgeDY
	centerMag := centerDX*centerDX + centerDY*centerDY

	if edgeMag <= centerMag {
		t.Errorf("edge should exceed center: edge=%.4f center=%.4f", edgeMag, centerMag)
	}
	if edgeMag < 1.0 {
		t.Errorf("edge displacement too small: %.4f", edgeMag)
	}
}

func computeDisplacement(sdf func(x, y float64) float64, x, y, R float64, g *Glass) (float64, float64) {
	sInside := -sdf(x, y)
	if sInside < 0 {
		return 0, 0
	}
	hpx := g.profileHeight(sInside, R)
	fPrime := g.profileDeriv(sInside, R, hpx)
	gx := (sdf(x+1, y) - sdf(x-1, y)) / 2
	gy := (sdf(x, y+1) - sdf(x, y-1)) / 2
	nx, ny, nz := fPrime*gx, fPrime*gy, 1.0
	nlen := math.Sqrt(nx*nx + ny*ny + nz*nz)
	nx, ny = nx/nlen, ny/nlen
	refrPow := 1.0 - 1.0/g.IOR
	displace := refrPow * g.RefractionStrength * hpx
	return -nx * displace, -ny * displace
}
