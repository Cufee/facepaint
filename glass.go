package facepaint

import (
	"image"
	"image/color"
	"math"

	"github.com/cufee/facepaint/style"
)

// Glass is a backdrop effect that simulates a glass lens.
type Glass struct {
	Width              float64
	IOR                float64
	RefractionStrength float64
	ChromaticStrength  float64
	LensStrength       float64
	SpecularStrength   float64
	SpecularExponent   float64
	FresnelStrength    float64
	BlurSigma          float64
	Azimuth            float64
	Elevation          float64
	LightingColor      color.Color
	Profile            BevelProfile
}

func NewGlass() *Glass {
	return &Glass{
		Width:              18,
		IOR:                1.5,
		RefractionStrength: 1.5,
		ChromaticStrength:  0.6,
		LensStrength:       0.06,
		SpecularStrength:   0.4,
		SpecularExponent:   8,
		FresnelStrength:    0.3,
		BlurSigma:          3,
		Azimuth:            225,
		Elevation:          60,
		LightingColor:      color.RGBA{255, 255, 255, 255},
		Profile:            BevelProfileCircular,
	}
}

func (g *Glass) Apply(ctx style.EffectContext) (image.Image, error) {
	sdf := ctx.SDF()
	if sdf == nil {
		return ctx.Image(), nil
	}

	src := toRGBA(ctx.Image())
	bounds := ctx.LocalBounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w <= 0 || h <= 0 {
		return src, nil
	}

	g.applyDefaults()
	R := g.Width

	azRad := g.Azimuth * math.Pi / 180
	elRad := g.Elevation * math.Pi / 180
	lx2d := math.Cos(azRad)
	ly2d := math.Sin(azRad)
	Lx := lx2d * math.Cos(elRad)
	Ly := ly2d * math.Cos(elRad)
	Lz := math.Sin(elRad)
	Hx, Hy, Hz := Lx, Ly, Lz+1
	Hlen := math.Sqrt(Hx*Hx + Hy*Hy + Hz*Hz)
	if Hlen > 0 {
		Hx, Hy, Hz = Hx/Hlen, Hy/Hlen, Hz/Hlen
	}

	lr, lg, lb, _ := g.LightingColor.RGBA()
	lrf, lgf, lbf := float64(lr)/65535, float64(lg)/65535, float64(lb)/65535

	refrPow := 1.0 - 1.0/g.IOR
	specScale := g.SpecularStrength

	dist := make([]float64, w*h)
	for y := range h {
		for x := range w {
			dist[y*w+x] = sdf(float64(x), float64(y))
		}
	}

	out := image.NewRGBA(bounds)

	parallelRows(h, func(y0, y1 int) {
		for y := y0; y < y1; y++ {
			for x := range w {
				idx := y*w + x
				sInside := -dist[idx]

				si := idx * 4

				coverage, ok := edgeCoverage(sInside, R)
				if !ok {
					copy(out.Pix[si:], src.Pix[si:si+4])
					continue
				}

				hpx := g.profileHeight(sInside, R)
				fPrime := g.profileDeriv(sInside, R, hpx)

				var gx, gy float64
				if x > 1 && x < w-2 {
					gx = (dist[idx+2] - dist[idx-2]) / 4
				} else if x > 0 && x < w-1 {
					gx = (dist[idx+1] - dist[idx-1]) / 2
				} else if x == 0 && w > 1 {
					gx = dist[idx+1] - dist[idx]
				} else if x == w-1 && w > 1 {
					gx = dist[idx] - dist[idx-1]
				}
				if y > 1 && y < h-2 {
					gy = (dist[idx+2*w] - dist[idx-2*w]) / 4
				} else if y > 0 && y < h-1 {
					gy = (dist[idx+w] - dist[idx-w]) / 2
				} else if y == 0 && h > 1 {
					gy = dist[idx+w] - dist[idx]
				} else if y == h-1 && h > 1 {
					gy = dist[idx] - dist[idx-w]
				}

				nx, ny, nz := fPrime*gx, fPrime*gy, 1.0
				nlen := math.Sqrt(nx*nx + ny*ny + nz*nz)
				if nlen > 0 {
					nx, ny, nz = nx/nlen, ny/nlen, nz/nlen
				}

				displace := refrPow * g.RefractionStrength * hpx
				dx := -nx * displace
				dy := -ny * displace

				caDx := dx * g.ChromaticStrength
				caDy := dy * g.ChromaticStrength

				bx, by := float64(x), float64(y)
				if g.LensStrength > 0 {
					cx, cy := float64(w)/2, float64(h)/2
					lensRange := math.Min(float64(w), float64(h)) / 2
					dt := sInside / lensRange
					if dt > 1 {
						dt = 1
					}
					if dt < 0 {
						dt = 0
					}
					dt = dt * dt * (3 - 2*dt)
					scale := 1 - g.LensStrength*dt
					bx = cx + (float64(x)-cx)*scale
					by = cy + (float64(y)-cy)*scale
				}

				sx := bx + dx
				sy := by + dy

				rCh := sampleBilinearChannel(src, w, h, sx+caDx, sy+caDy, 0)
				gCh := sampleBilinearChannel(src, w, h, sx, sy, 1)
				bCh := sampleBilinearChannel(src, w, h, sx-caDx, sy-caDy, 2)

				t := sInside / R
				if t > 1 {
					t = 1
				}
				edgeFade := (1 - t) * (1 - t)

				nxNorm := float64(x)/float64(w) - 0.5
				nyNorm := float64(y)/float64(h) - 0.5
				diag := nxNorm*lx2d + nyNorm*ly2d
				diagFactor := 0.2 + 0.8*math.Min(1, 2*diag*diag)

				ndh := math.Max(0, nx*Hx+ny*Hy+nz*Hz)
				spec := specScale * math.Pow(ndh, g.SpecularExponent) * edgeFade * diagFactor

				rim := g.FresnelStrength * edgeFade

				lightMix := math.Min(spec+rim, 0.6)

				outR := rCh*(1-lightMix) + lrf*255*lightMix
				outG := gCh*(1-lightMix) + lgf*255*lightMix
				outB := bCh*(1-lightMix) + lbf*255*lightMix

				finalR := float64(src.Pix[si])*(1-coverage) + outR*coverage
				finalG := float64(src.Pix[si+1])*(1-coverage) + outG*coverage
				finalB := float64(src.Pix[si+2])*(1-coverage) + outB*coverage

				out.Pix[si] = uint8(clamp8(finalR))
				out.Pix[si+1] = uint8(clamp8(finalG))
				out.Pix[si+2] = uint8(clamp8(finalB))
				out.Pix[si+3] = src.Pix[si+3]
			}
		}
	})

	g.applyInternalBlur(out, src, dist, w, h, R, g.BlurSigma)

	return out, nil
}

func (g *Glass) applyInternalBlur(out *image.RGBA, src *image.RGBA, dist []float64, w, h int, R, sigma float64) {
	if sigma <= 0 {
		return
	}
	blurred := imagingBlur(src, sigma)

	parallelRows(h, func(y0, y1 int) {
		for y := y0; y < y1; y++ {
			for x := range w {
				idx := y*w + x
				sInside := -dist[idx]
				coverage, ok := edgeCoverage(sInside, R)
				if !ok {
					continue
				}
				if sInside > R {
					continue
				}
				t := sInside / R
				fade := 4 * t * (1 - t)

				si := idx * 4
				mix := fade * 0.25 * coverage
				for c := range 3 {
					sharp := float64(out.Pix[si+c])
					soft := float64(blurred.Pix[si+c])
					out.Pix[si+c] = uint8(clamp8(sharp*(1-mix) + soft*mix))
				}
			}
		}
	})
}

func edgeCoverage(sInside, R float64) (float64, bool) {
	if sInside < 0 {
		c := math.Max(0, 1+sInside)
		if c <= 0 {
			return 0, false
		}
		return c, true
	}
	return 1, true
}

func (g *Glass) applyDefaults() {
	if g.Width <= 0 {
		g.Width = 12
	}
	if g.IOR <= 0 {
		g.IOR = 1.5
	}
	if g.RefractionStrength <= 0 {
		g.RefractionStrength = 1.5
	}
	if g.ChromaticStrength <= 0 {
		g.ChromaticStrength = 0.6
	}
	if g.LensStrength < 0 {
		g.LensStrength = 0
	}
	if g.SpecularStrength <= 0 {
		g.SpecularStrength = 0.4
	}
	if g.SpecularExponent <= 0 {
		g.SpecularExponent = 8
	}
	if g.FresnelStrength <= 0 {
		g.FresnelStrength = 0.3
	}
	if g.BlurSigma < 0 {
		g.BlurSigma = 0
	}
	if g.LightingColor == nil {
		g.LightingColor = color.RGBA{255, 255, 255, 255}
	}
	if g.Azimuth == 0 && g.Elevation == 0 {
		g.Azimuth = 225
		g.Elevation = 60
	}
}

func (g *Glass) profileHeight(s, R float64) float64 {
	if s >= R {
		return R
	}
	if s < 0 {
		return 0
	}
	switch g.Profile {
	case BevelProfileSmoothstep:
		t := s / R
		return (R - s) * (3*t*t - 2*t*t*t)
	default:
		v := R*R - (s-R)*(s-R)
		if v <= 0 {
			return 0
		}
		return math.Sqrt(v)
	}
}

func (g *Glass) profileDeriv(s, R, h float64) float64 {
	if s >= R {
		return 0
	}
	switch g.Profile {
	case BevelProfileSmoothstep:
		t := s / R
		return -(3*t*t - 2*t*t*t) + (1-t)*(6*t-6*t*t)
	default:
		return (R - s) / math.Max(h, 0.5)
	}
}

func sampleBilinearChannel(src *image.RGBA, w, h int, x, y float64, ch int) float64 {
	if w <= 0 || h <= 0 {
		return 0
	}

	xf := math.Floor(x)
	yf := math.Floor(y)
	x0 := clampInt(int(xf), 0, w-1)
	y0 := clampInt(int(yf), 0, h-1)
	x1 := clampInt(x0+1, 0, w-1)
	y1 := clampInt(y0+1, 0, h-1)

	fx := x - xf
	if fx < 0 {
		fx = 0
	}
	fy := y - yf
	if fy < 0 {
		fy = 0
	}

	i00 := (y0*w + x0) * 4
	i10 := (y0*w + x1) * 4
	i01 := (y1*w + x0) * 4
	i11 := (y1*w + x1) * 4

	v00 := float64(src.Pix[i00+ch])
	v10 := float64(src.Pix[i10+ch])
	v01 := float64(src.Pix[i01+ch])
	v11 := float64(src.Pix[i11+ch])

	top := v00*(1-fx) + v10*fx
	bot := v01*(1-fx) + v11*fx
	return top*(1-fy) + bot*fy
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
