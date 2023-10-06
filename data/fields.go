package data

import (
	"glow-gui/glow"

	"fyne.io/fyne/v2/data/binding"
)

type Fields struct {
	HueShift    binding.Float
	Scan        binding.Float
	Begin       binding.Float
	End         binding.Float
	Origin      binding.Int
	Orientation binding.Int
}

func NewFields() *Fields {
	inp := &Fields{
		HueShift:    binding.NewFloat(),
		Scan:        binding.NewFloat(),
		Begin:       binding.NewFloat(),
		End:         binding.NewFloat(),
		Origin:      binding.NewInt(),
		Orientation: binding.NewInt(),
	}
	return inp
}

func (inp *Fields) FromLayer(layer *glow.Layer) {
	inp.HueShift.Set(float64(layer.HueShift))
	inp.Scan.Set(float64(layer.Scan))
	inp.Begin.Set(float64(layer.Begin))
	inp.End.Set(float64(layer.End))
	inp.Origin.Set(int(layer.Grid.Origin))
	inp.Orientation.Set(int(layer.Grid.Orientation))
}

func (inp *Fields) ToLayer(layer *glow.Layer) {
	var (
		floatTemp float64
		intTemp   int
	)

	floatTemp, _ = inp.HueShift.Get()
	layer.HueShift = int16(floatTemp)

	floatTemp, _ = inp.Scan.Get()
	layer.Scan = uint16(floatTemp)

	floatTemp, _ = inp.Begin.Get()
	layer.Begin = uint16(floatTemp)

	floatTemp, _ = inp.End.Get()
	layer.End = uint16(floatTemp)

	intTemp, _ = inp.Origin.Get()
	layer.Grid.Origin = glow.Origin(intTemp)

	intTemp, _ = inp.Orientation.Get()
	layer.Grid.Orientation = glow.Orientation(intTemp)
}
