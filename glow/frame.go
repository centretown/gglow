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
	Length   uint16  `yaml:"length" json:"length"`
	Rows     uint16  `yaml:"rows" json:"rows"`
	Interval uint32  `yaml:"interval" json:"interval"`
	Layers   []Layer `yaml:"layers" json:"layers"`
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
	frame.Interval = interval
}

func (frame *Frame) Spin(light Light) {
	for i := range frame.Layers {
		frame.Layers[i].Spin(light)
	}
}

func (frame *Frame) AddLayers(layers ...Layer) {
	frame.Layers = append(frame.Layers, layers...)
	frame.updateLayers()
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
