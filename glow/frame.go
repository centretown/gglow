package glow

import (
	"fmt"

	"github.com/barkimedes/go-deepcopy"
)

const (
	DefaultInterval = 48
	MinimumInterval = 16
	MaximumInterval = 10000
)

type Frame struct {
	Length   uint16   `yaml:"length" json:"length"`
	Rows     uint16   `yaml:"rows" json:"rows"`
	Interval uint32   `yaml:"interval" json:"interval"`
	Layers   []*Layer `yaml:"layers" json:"layers"`
}

func NewFrame() (frame *Frame) {
	frame = &Frame{}
	frame.Interval = 48
	frame.Layers = append(frame.Layers, NewLayer())
	return
}

func (frame *Frame) AppendLayer(layer *Layer) {
	frame.Layers = append(frame.Layers, layer)
}

func (frame *Frame) InsertLayer(position int, layer *Layer) {
	if position < 0 {
		position = 0
	}
	if layer == nil {
		layer = &Layer{}
	}

	layerLength := len(frame.Layers)

	if position >= layerLength {
		frame.Layers = append(frame.Layers, layer)
		return
	}

	ll := make([]*Layer, 0, layerLength+1)
	i := 0
	for i < layerLength {
		if position == i {
			position = -1
			ll = append(ll, layer)
			continue
		}
		ll = append(ll, frame.Layers[i])
		i++
	}
	frame.Layers = ll
}

func (frame *Frame) updateLayers() {
	for i := range frame.Layers {
		frame.Layers[i].SetupLength(frame.Length, frame.Rows)
	}
}

func (frame *Frame) Validate() (err error) {
	if frame.Length == 0 {
		return fmt.Errorf("Frame.Setup zero length")
	}
	if frame.Rows == 0 {
		return fmt.Errorf("Frame.Setup zero rows")
	}
	frame.updateLayers()
	return err
}

func (frame *Frame) Setup(length, rows uint16) error {
	frame.Length = length
	frame.Rows = rows
	return frame.Validate()
}

func (frame *Frame) SetInterval(interval uint32) {
	frame.Interval = interval
}

func (frame *Frame) Spin(light Light) {
	for i := range frame.Layers {
		frame.Layers[i].Spin(light)
	}
}

func (frame *Frame) AddLayers(layers ...*Layer) {
	frame.Layers = append(frame.Layers, layers...)
	frame.updateLayers()
}

func (frame *Frame) LoadImages() {
	for _, layer := range frame.Layers {
		if len(layer.ImageName) > 0 {
			fmt.Println("load image", layer.ImageName, int(layer.Grid.Rows), int(layer.Grid.columns))
			layer.LoadImage(int(layer.Grid.Rows), int(layer.Grid.columns))
		}
	}
}

func FrameDeepCopy(source *Frame) (frame *Frame, err error) {

	var (
		deepCopy interface{}
		ok       bool
	)

	deepCopy, err = deepcopy.Anything(source)
	if err != nil {
		return
	}

	frame, ok = deepCopy.(*Frame)
	if !ok {
		err = fmt.Errorf("frame deepcopy not ok")
	}
	return
}

func (frame *Frame) MakeCode() string {
	layers := func() string {
		var s string
		for _, layer := range frame.Layers {
			s += layer.MakeCode()
		}
		return s
	}

	s := fmt.Sprintf("{%d,%d,%d,{%s}},\n",
		frame.Length, frame.Rows, frame.Interval, layers())
	return s
}
