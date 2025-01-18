package tests

import (
	_ "embed"

	"github.com/cufee/facepaint/style"
)

//go:embed font.ttf
var fontData []byte

func Font() style.Font {
	return style.NewFont(fontData, 24)
}
