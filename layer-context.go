package facepaint

import (
	"image"
	"slices"

	"github.com/fogleman/gg"
	"github.com/pkg/errors"
)

type layerContext map[int]*gg.Context

func (ctx layerContext) layer(idx int) (*gg.Context, error) {
	layer := ctx[idx]
	if layer == nil {
		return nil, errors.New("layer context is nil")
	}
	return layer, nil
}

func (ctx layerContext) Image() image.Image {
	var layers []int
	for idx := range ctx {
		layers = append(layers, idx)
	}
	slices.Sort(layers)

	var frame *gg.Context
	for _, idx := range layers {
		layer := ctx[idx]
		if layer == nil {
			continue
		}
		if frame == nil {
			frame = layer
			continue
		}

		frame.DrawImage(layer.Image(), 0, 0)
	}

	return frame.Image()
}
