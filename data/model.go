package data

import (
	"glow-gui/glow"
	"glow-gui/res"
	"glow-gui/store"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

type Model struct {
	Frame     binding.Untyped
	Title     binding.String
	LayerList binding.UntypedList
	Layer     binding.Untyped
	Fields    binding.Struct
}

func NewModel() *Model {
	m := &Model{
		Frame:     binding.NewUntyped(),
		Title:     binding.NewString(),
		LayerList: binding.NewUntypedList(),
		Layer:     binding.NewUntyped(),
	}

	m.Frame.Set(&glow.Frame{})
	m.Title.Set("")
	m.LayerList.Set(make([]interface{}, 0))
	layer := &glow.Layer{}
	m.Layer.Set(layer)
	m.Fields = binding.BindStruct(layer)

	return m
}

func (m *Model) GetTitle() string {
	str, _ := m.Title.Get()
	return str
}

func (m *Model) getFrame() (frame *glow.Frame) {
	face, err := getUntyped("getFrame", m.Frame, res.MsgGetFrame)
	if err != nil {
		panic(res.MsgGetFrame) // panic
	}
	frame = face.(*glow.Frame) // panic
	return
}

func (m *Model) GetLayer() (layer *glow.Layer) {
	face, err := getUntyped("getLayer", m.Layer, res.MsgGetLayer)
	if err != nil {
		s := res.MsgGetFrame.String() + err.Error()
		panic(s) // panic
	}
	layer = face.(*glow.Layer) // panic
	return
}

func (m *Model) SetLayer(i int) {
	frame := m.getFrame()
	var layer *glow.Layer
	if i < len(frame.Layers) {
		layer = &frame.Layers[i]
	} else {
		layer = &glow.Layer{}
	}
	setUntyped(m.Layer, layer, res.MsgSetLayer)
	m.Fields = binding.BindStruct(layer)
}

func (m *Model) setFrame(frame *glow.Frame) (err error) {
	err = setUntyped(m.Frame, frame, res.MsgSetFrame)
	if err != nil {
		return
	}

	list := make([]interface{}, 0, len(frame.Layers))
	for i := range frame.Layers {
		list = append(list, &frame.Layers[i])
	}
	err = setUntypedList(m.LayerList, list, res.MsgSetLayerList)
	if err != nil {
		return
	}

	var layer *glow.Layer
	if len(frame.Layers) > 0 {
		layer = &frame.Layers[0]
	} else {
		layer = &glow.Layer{}
	}
	err = setUntyped(m.Layer, layer, res.MsgSetLayer)
	if err != nil {
		return
	}

	m.Fields = binding.BindStruct(layer)
	return
}

func (m *Model) LoadFrame(frameName string) (err error) {
	var uri fyne.URI
	uri, err = store.LookupURI(frameName)
	if err != nil {
		res.MsgGetEffectLookup.Log(frameName, err)
		return
	}

	frame := &glow.Frame{}
	err = store.LoadFrameURI(uri, frame)
	if err != nil {
		res.MsgGetEffectLoad.Log(uri.Name(), err)
		return
	}
	err = m.setFrame(frame)
	return
}
