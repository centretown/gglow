package data

import (
	"glow-gui/glow"
	"glow-gui/resources"
	"glow-gui/store"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

type Model struct {
	EffectName       string
	Frame            binding.Untyped
	LayerSummaryList binding.StringList
	Layer            binding.Untyped
	LayerIndex       int
}

func NewModel() *Model {
	m := &Model{
		Frame:            binding.NewUntyped(),
		LayerSummaryList: binding.NewStringList(),
		Layer:            binding.NewUntyped(),
	}

	m.Frame.Set(&glow.Frame{})
	m.Layer.Set(&glow.Layer{})
	m.Frame.AddListener(binding.NewDataListener(m.onChangeFrame))
	return m
}

func (m *Model) onChangeFrame() {
	frame := m.GetFrame()

	summaries := make([]string, 0, len(frame.Layers))
	for i, layer := range frame.Layers {
		summaries = append(summaries, Summarize(&layer, i+1))
	}
	m.LayerSummaryList.Set(summaries)
	m.SetCurrentLayer(m.LayerIndex)
}

func (m *Model) GetFrame() *glow.Frame {
	frame, _ := m.Frame.Get()
	return frame.(*glow.Frame)
}

func (m *Model) GetCurrentLayer() *glow.Layer {
	layer, _ := m.Layer.Get()
	return layer.(*glow.Layer)
}

func (m *Model) SetCurrentLayer(i int) {
	frame := m.GetFrame()
	var layer *glow.Layer
	if i < len(frame.Layers) {
		m.LayerIndex = i
		layer = &frame.Layers[i]
	} else {
		m.LayerIndex = 0
		layer = &glow.Layer{}
	}
	m.Layer.Set(layer)
}

func (m *Model) UpdateFrame() {
	current := m.GetFrame()
	frame := *current
	m.Frame.Set(&frame)
}

func (m *Model) LoadFrame(frameName string) error {
	var uri fyne.URI
	uri, err := store.LookupURI(frameName)
	if err != nil {
		resources.MsgGetEffectLookup.Log(frameName, err)
		return err
	}

	m.EffectName = frameName
	frame := &glow.Frame{}
	err = store.LoadFrameURI(uri, frame)
	if err != nil {
		resources.MsgGetEffectLoad.Log(uri.Name(), err)
		return err
	}

	m.LayerIndex = 0
	m.Frame.Set(frame)
	return nil
}
