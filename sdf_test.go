package facepaint

import (
	"math"
	"testing"

	"github.com/cufee/facepaint/style"
	"github.com/matryer/is"
)

func TestSDFRoundedBox(t *testing.T) {
	is := is.New(t)

	dims := contentDimensions{Width: 100, Height: 60}
	st := style.Style{
		BorderRadiusTopLeft: 20, BorderRadiusTopRight: 20,
		BorderRadiusBottomLeft: 20, BorderRadiusBottomRight: 20,
	}
	sdf := buildRoundedBoxSDF(dims, st)
	is.True(sdf != nil)

	if d := sdf(50, 30); d >= 0 {
		t.Fatalf("center should be negative, got %f", d)
	}
	if d := sdf(200, 30); d < 90 || d > 110 {
		t.Fatalf("outside point distance out of range: %f", d)
	}
	if d := sdf(50, 0); math.Abs(d) > 0.5 {
		t.Fatalf("top-edge boundary should be ~0, got %f", d)
	}
	if d := sdf(50, 1); d > 0 {
		t.Fatalf("just inside top should be negative, got %f", d)
	}
	a, b := sdf(49, 30), sdf(51, 30)
	if math.Abs(a-b) > 1.0 {
		t.Fatalf("SDF should be ~symmetric across centerline, got %f vs %f", a, b)
	}
}

func TestSDFNilForDegenerate(t *testing.T) {
	is := is.New(t)
	if sdf := buildRoundedBoxSDF(contentDimensions{}, style.Style{}); sdf != nil {
		is.Fail()
	}
}
