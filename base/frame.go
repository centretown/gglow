package glow

import (
	"fmt"
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

func (frame *Frame) Setup(length, rows uint16, interval uint32) error {
	frame.Length = length
	frame.Rows = rows
	frame.Interval = interval
	return frame.Validate()
}

func (frame *Frame) Spin(light Light) {
	for _, layer := range frame.Layers {
		layer.Spin(light)
	}
}

func (frame *Frame) AddLayers(layers ...Layer) {
	frame.Layers = append(frame.Layers, layers...)
	frame.updateLayers()
}
