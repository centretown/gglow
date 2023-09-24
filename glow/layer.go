package glow

import "fmt"

type Layer struct {
	Length uint16 `yaml:"length" json:"length"`
	Rows   uint16 `yaml:"rows" json:"rows"`
	Grid   Grid   `yaml:"grid" json:"grid"`
	Chroma Chroma `yaml:"chroma" json:"chroma"`

	HueShift int16  `yaml:"hue_shift" json:"hue_shift"`
	Scan     uint16 `yaml:"scan" json:"scan"`
	Begin    uint16 `yaml:"begin" json:"begin"`
	End      uint16 `yaml:"end" json:"end"`

	position uint16
	first    uint16
	last     uint16
}

func (layer *Layer) Setup(length, rows uint16,
	grid *Grid, chroma *Chroma, hueShift int16,
	scan uint16, begin uint16, end uint16) error {

	layer.Length = length
	layer.Rows = rows
	layer.Grid = *grid
	layer.Chroma = *chroma
	layer.HueShift = hueShift
	layer.Scan = scan
	layer.Begin = begin
	layer.End = end

	return layer.Validate()
}

func (layer *Layer) SetupLength(length, rows uint16) error {
	layer.Length = length
	layer.Rows = rows
	layer.Grid.SetupLength(length, rows)
	layer.Chroma.SetupLength(length, layer.HueShift)
	return layer.Validate()
}

func (layer *Layer) Validate() error {
	if layer.Length == 0 {
		return fmt.Errorf("Layer.Setup zero length")
	}
	if layer.Rows == 0 {
		return fmt.Errorf("Layer.Setup zero rows")
	}
	if err := layer.Grid.SetupLength(layer.Length, layer.Rows); err != nil {
		return err
	}
	if err := layer.Chroma.SetupLength(layer.Length, layer.HueShift); err != nil {
		return err
	}
	if layer.Scan > layer.Length {
		layer.Scan = layer.Length
	}
	if layer.End == 0 {
		layer.End = 100
	}
	layer.setBounds()

	return nil
}

func (layer *Layer) setBounds() {
	ratio := func(offset, length uint16) float32 {
		if offset > 100 {
			offset %= 100
		}
		return float32(offset) / 100.0 * float32(length)
	}

	layer.first = layer.Grid.AdjustBounds(ratio(layer.Begin, layer.Length))
	layer.last = layer.Grid.AdjustBounds(ratio(layer.End, layer.Length))

	if layer.last < layer.first {
		layer.first, layer.last = layer.last, layer.first
	}
}

func (layer *Layer) Spin(light Light) {
	startAt := layer.first
	endAt := layer.last
	if layer.Scan > 0 {
		startAt, endAt = layer.updateScanPosition()
	}

	for i := startAt; i < endAt; i++ {
		light.Set(layer.Grid.Map(i), layer.Chroma.Map(i))
	}

	layer.Chroma.Update()
}

func (layer *Layer) updateScanPosition() (startAt, endAt uint16) {
	startAt = layer.position
	endAt = layer.position + layer.Scan
	layer.position++
	if layer.position >= layer.last {
		layer.position = layer.first
	}
	return startAt, endAt
}
