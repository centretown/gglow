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
	fld := &Fields{
		HueShift:    binding.NewFloat(),
		Scan:        binding.NewFloat(),
		Begin:       binding.NewFloat(),
		End:         binding.NewFloat(),
		Origin:      binding.NewInt(),
		Orientation: binding.NewInt(),
	}
	return fld
}

func (fld *Fields) FromLayer(layer *glow.Layer) {
	fld.HueShift.Set(float64(layer.HueShift))
	fld.Scan.Set(float64(layer.Scan))
	fld.Begin.Set(float64(layer.Begin))
	fld.End.Set(float64(layer.End))
	fld.Origin.Set(int(layer.Grid.Origin))
	fld.Orientation.Set(int(layer.Grid.Orientation))
}

func (fld *Fields) ToLayer(layer *glow.Layer) {
	var (
		f float64
		i int
	)

	f, _ = fld.HueShift.Get()
	layer.HueShift = int16(f)

	f, _ = fld.Scan.Get()
	layer.Scan = uint16(f)

	f, _ = fld.Begin.Get()
	layer.Begin = uint16(f)

	f, _ = fld.End.Get()
	layer.End = uint16(f)

	i, _ = fld.Origin.Get()
	layer.Grid.Origin = glow.Origin(i)

	i, _ = fld.Orientation.Get()
	layer.Grid.Orientation = glow.Orientation(i)
}
