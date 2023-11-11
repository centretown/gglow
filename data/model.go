package data

import (
	"glow-gui/glow"
	"glow-gui/store"

	"fyne.io/fyne/v2/data/binding"
)

type Model struct {
	EffectName       string
	Frame            binding.Untyped
	LayerSummaryList binding.StringList
	Layer            binding.Untyped
	LayerIndex       int
	Store            *store.Store
}

func NewModel(store *store.Store) *Model {
	m := &Model{
		Store:            store,
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

func (m *Model) LoadFrame(key string) error {
	frame := &glow.Frame{}
	err := m.Store.LoadFrame(key, frame)
	if err != nil {
		return err
	}
	m.LayerIndex = 0
	m.EffectName = key
	m.Frame.Set(frame)
	return nil
}
