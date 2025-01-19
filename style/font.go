package style

import (
	"sync"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type Font interface {
	Size() float64
	Valid() bool
	Face() (font.Face, func() error)
}

type fontType struct {
	size float64
	face font.Face
	mx   *sync.Mutex
}

func (f *fontType) Size() float64 {
	return f.size
}

func (f *fontType) Valid() bool {
	return f.face != nil
}

func (f *fontType) Face() (font.Face, func() error) {
	f.mx.Lock()
	return f.face, func() error { f.mx.Unlock(); return nil }
}

func NewFont(data []byte, size float64) (Font, error) {
	ttf, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(ttf, &truetype.Options{
		Size: size,
	})
	return &fontType{size, face, &sync.Mutex{}}, nil
}
