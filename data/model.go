package data

import (
	"glow-gui/glow"
	"glow-gui/resources"
	"glow-gui/store"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

type Model struct {
	Frame            binding.Untyped
	Title            binding.String
	LayerList        binding.UntypedList
	LayerSummaryList binding.StringList
	LayerIndex       binding.Int
	Layer            binding.Untyped
	Fields           *Fields
}

func NewModel() *Model {
	m := &Model{
		Frame:            binding.NewUntyped(),
		Title:            binding.NewString(),
		LayerList:        binding.NewUntypedList(),
		LayerSummaryList: binding.NewStringList(),
		LayerIndex:       binding.NewInt(),
		Layer:            binding.NewUntyped(),
		Fields:           NewFields(),
	}

	m.Frame.Set(&glow.Frame{})
	m.Title.Set("")
	m.LayerList.Set(make([]interface{}, 0))
	layer := &glow.Layer{}
	m.Layer.Set(layer)
	m.Fields.FromLayer(layer)
	m.Frame.AddListener(binding.NewDataListener(m.onChangeFrame))
	return m
}

func (m *Model) onChangeFrame() {
	frame := m.getFrame()
	listLength := len(frame.Layers)
	list := make([]interface{}, 0, listLength)
	for i := range frame.Layers {
		list = append(list, &frame.Layers[i])
	}

	summaries := make([]string, 0, m.LayerList.Length())
	for i, layer := range frame.Layers {
		summaries = append(summaries, Summarize(&layer, i+1))
	}

	m.LayerList.Set(list)
	m.LayerSummaryList.Set(summaries)
	m.SetCurrentLayer(0)
}

func (m *Model) getFrame() *glow.Frame {
	frame, _ := m.Frame.Get()
	return frame.(*glow.Frame)
}

func (m *Model) SetCurrentLayer(i int) {
	frame := m.getFrame()
	var layer *glow.Layer
	if i < len(frame.Layers) {
		layer = &frame.Layers[i]
		m.LayerIndex.Set(i)
	} else {
		layer = &glow.Layer{}
		// m.LayerIndex.Set(0)
		m.LayerIndex.Set(-1)
	}
	m.Layer.Set(layer)
	m.Fields.FromLayer(layer)
}

func (m *Model) LoadFrame(frameName string) error {
	var uri fyne.URI
	uri, err := store.LookupURI(frameName)
	if err != nil {
		resources.MsgGetEffectLookup.Log(frameName, err)
		return err
	}

	frame := &glow.Frame{}
	err = store.LoadFrameURI(uri, frame)
	if err != nil {
		resources.MsgGetEffectLoad.Log(uri.Name(), err)
		return err
	}

	m.Frame.Set(frame)
	return nil
}

func (m *Model) GetTitle() string {
	str, _ := m.Title.Get()
	return str
}
