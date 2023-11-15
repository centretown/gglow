package data

import (
	"glow-gui/glow"

	"fyne.io/fyne/v2/data/binding"
)

type LayerFields struct {
	HueShift    binding.Int
	Scan        binding.Int
	Rate        binding.Int
	Origin      binding.Int
	Orientation binding.Int
	Begin       binding.Int
	End         binding.Int
	Colors      []glow.HSV
}

func NewLayerFields() *LayerFields {
	fld := &LayerFields{
		HueShift:    binding.NewInt(),
		Scan:        binding.NewInt(),
		Rate:        binding.NewInt(),
		Origin:      binding.NewInt(),
		Orientation: binding.NewInt(),
		Begin:       binding.NewInt(),
		End:         binding.NewInt(),
	}
	return fld
}

func (fld *LayerFields) FromLayer(layer *glow.Layer) {
	fld.HueShift.Set(int(layer.HueShift))
	fld.Scan.Set(int(layer.Scan))
	fld.Rate.Set(int(layer.Rate))
	fld.Origin.Set(int(layer.Grid.Origin))
	fld.Orientation.Set(int(layer.Grid.Orientation))
	fld.Begin.Set(int(layer.Begin))
	fld.End.Set(int(layer.End))
	fld.Colors = make([]glow.HSV, len(layer.Chroma.Colors))
	copy(fld.Colors, layer.Chroma.Colors)
}

func (fld *LayerFields) ToLayer(layer *glow.Layer) {
	var i int
	i, _ = fld.HueShift.Get()
	layer.HueShift = int16(i)

	i, _ = fld.Scan.Get()
	layer.Scan = uint16(i)

	i, _ = fld.Rate.Get()
	layer.Rate = uint32(i)

	i, _ = fld.Origin.Get()
	layer.Grid.Origin = glow.Origin(i)

	i, _ = fld.Orientation.Get()
	layer.Grid.Orientation = glow.Orientation(i)

	i, _ = fld.Begin.Get()
	layer.Begin = uint16(i)

	i, _ = fld.End.Get()
	layer.End = uint16(i)

	layer.Chroma.Colors = make([]glow.HSV, len(fld.Colors))
	copy(layer.Chroma.Colors, fld.Colors)
}

// func (fld *Fields) IsDirty(layer *glow.Layer) bool {
// 	var i int

// 	i, _ = fld.HueShift.Get()
// 	if layer.HueShift != int16(i) {
// 		return true
// 	}

// 	i, _ = fld.Scan.Get()
// 	if layer.Scan != uint16(i) {
// 		return true
// 	}

// 	i, _ = fld.Rate.Get()
// 	if layer.Rate != uint32(i) {
// 		return true
// 	}

// 	i, _ = fld.Origin.Get()
// 	if layer.Grid.Origin != glow.Origin(i) {
// 		return true
// 	}

// 	i, _ = fld.Orientation.Get()
// 	if layer.Grid.Orientation != glow.Orientation(i) {
// 		return true
// 	}

// 	i, _ = fld.Begin.Get()
// 	if layer.Begin != uint16(i) {
// 		return true
// 	}

// 	i, _ = fld.End.Get()
// 	return layer.End != uint16(i)
// }
