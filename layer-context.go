package facepaint

import (
	"image"
	"slices"

	"github.com/fogleman/gg"
	"github.com/pkg/errors"
)

type layerHookTarget byte

const (
	layerHookTargetBeforeLayerRender = iota
	layerHookTargetAfterLayerRender  = iota
)

type layerHook struct {
	target  layerHookTarget
	execute func(frame, layer *layer)
}

func layerHookBeforeRender(fn func(frame, layer *layer)) layerHook {
	return layerHook{
		target:  layerHookTargetBeforeLayerRender,
		execute: fn,
	}
}
func layerHookAfterRender(fn func(frame, layer *layer)) layerHook {
	return layerHook{
		target:  layerHookTargetAfterLayerRender,
		execute: fn,
	}
}

type layerContext struct {
	layers     map[int]*layer
	layerHooks map[int][]layerHook
}

func newLayerContext(size int) *layerContext {
	return &layerContext{
		layers:     make(map[int]*layer, size),
		layerHooks: make(map[int][]layerHook, size),
	}
}

type layer struct {
	*gg.Context
}

func newLayer(width, height int) *layer {
	return &layer{Context: gg.NewContext(width, height)}
}

func (ctx *layerContext) registerHook(idx int, hook layerHook) {
	ctx.layerHooks[idx] = append(ctx.layerHooks[idx], hook)
}

func (ctx *layerContext) layer(idx int) (*layer, error) {
	layer := ctx.layers[idx]
	if layer == nil {
		return nil, errors.New("layer context is nil")
	}
	return layer, nil
}

func (ctx *layerContext) Image() image.Image {
	var layers []int
	for idx := range ctx.layers {
		layers = append(layers, idx)
	}
	slices.Sort(layers)

	var frame *layer
	for _, idx := range layers {
		layer := ctx.layers[idx]
		if layer == nil {
			continue
		}
		if frame == nil {
			frame = layer
			continue
		}

		// run before render hooks
		for _, hook := range ctx.layerHooks[idx] {
			if hook.target == layerHookTargetBeforeLayerRender {
				hook.execute(frame, layer)
			}
		}

		frame.DrawImage(layer.Image(), 0, 0)

		// run after render hooks
		for _, hook := range ctx.layerHooks[idx] {
			if hook.target == layerHookTargetAfterLayerRender {
				hook.execute(frame, layer)
			}
		}
	}

	return frame.Image()
}
