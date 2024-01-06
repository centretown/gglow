package codeio

import (
	"gglow/glow"
	"gglow/iohandler"
	"testing"
)

func TestGenerate(t *testing.T) {
	list := []*iohandler.EffectItem{
		{Title: "black white scan", Constant: "", Frame: &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []*glow.Layer{{Length: 16, Rows: 2,
				Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
				Chroma: glow.Chroma{Length: 16, HueShift: -1,
					Colors: []glow.HSV{{Hue: 0, Saturation: 255, Value: 127}}}}}}},
		{Title: "complementary scan", Constant: "", Frame: &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []*glow.Layer{{Length: 16, Rows: 2,
				Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
				Chroma: glow.Chroma{Length: 16, HueShift: -1,
					Colors: []glow.HSV{{Hue: 0, Saturation: 255, Value: 127}}}}}}},
		{Title: "double scan", Constant: "", Frame: &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []*glow.Layer{{Length: 16, Rows: 2,
				Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
				Chroma: glow.Chroma{Length: 16, HueShift: -1,
					Colors: []glow.HSV{{Hue: 0, Saturation: 255, Value: 127}}}}}}},
		{Title: "gradient scan", Constant: "", Frame: &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []*glow.Layer{{Length: 16, Rows: 2,
				Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
				Chroma: glow.Chroma{Length: 16, HueShift: -1,
					Colors: []glow.HSV{{Hue: 0, Saturation: 255, Value: 127}}}}}}},
		{Title: "rainbow diagonal", Constant: "", Frame: &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []*glow.Layer{{Length: 16, Rows: 2,
				Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
				Chroma: glow.Chroma{Length: 16, HueShift: -1,
					Colors: []glow.HSV{{Hue: 0, Saturation: 255, Value: 127}}}}}}},
		{Title: "rainbow horizontal", Constant: "", Frame: &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []*glow.Layer{{Length: 16, Rows: 2,
				Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
				Chroma: glow.Chroma{Length: 16, HueShift: -1,
					Colors: []glow.HSV{{Hue: 0, Saturation: 255, Value: 127}}}}}}},
		{Title: "rainbow vertical", Constant: "", Frame: &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []*glow.Layer{{Length: 16, Rows: 2,
				Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
				Chroma: glow.Chroma{Length: 16, HueShift: -1,
					Colors: []glow.HSV{{Hue: 0, Saturation: 255, Value: 127}}}}}}},
		{Title: "split in three", Constant: "", Frame: &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []*glow.Layer{{Length: 16, Rows: 2,
				Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
				Chroma: glow.Chroma{Length: 16, HueShift: -1,
					Colors: []glow.HSV{{Hue: 0, Saturation: 255, Value: 127}}}}}}},
		{Title: "split in two", Constant: "", Frame: &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []*glow.Layer{{Length: 16, Rows: 2,
				Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
				Chroma: glow.Chroma{Length: 16, HueShift: -1,
					Colors: []glow.HSV{{Hue: 0, Saturation: 255, Value: 127}}}}}}},
		{Title: "spotlight", Constant: "", Frame: &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []*glow.Layer{
				{Length: 16, Rows: 2, HueShift: -1, Begin: 0, End: 100, Rate: 0,
					Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
					Chroma: glow.Chroma{Length: 16, HueShift: -1,
						Colors: []glow.HSV{{Hue: 0, Saturation: 255, Value: 127}, {Hue: 0, Saturation: 255, Value: 127}}}},
				{Length: 16, Rows: 2, HueShift: -1, Scan: 10, Begin: 0, End: 100, Rate: 0,
					Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
					Chroma: glow.Chroma{Length: 16, HueShift: -1,
						Colors: []glow.HSV{{Hue: 0, Saturation: 127, Value: 196}, {Hue: 0, Saturation: 255, Value: 127}}}},
			},
		}},
	}
	list2 := []*iohandler.EffectItem{
		{Title: "black white scan", Constant: "", Frame: &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []*glow.Layer{{Length: 16, Rows: 2,
				Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
				Chroma: glow.Chroma{Length: 16, HueShift: -1,
					Colors: []glow.HSV{{Hue: 0, Saturation: 255, Value: 127}}}}}}},
		{Title: "spotlight", Constant: "", Frame: &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []*glow.Layer{
				{Length: 16, Rows: 2, HueShift: -1, Begin: 0, End: 100, Rate: 0,
					Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
					Chroma: glow.Chroma{Length: 16, HueShift: -1,
						Colors: []glow.HSV{{Hue: 0, Saturation: 255, Value: 127}, {Hue: 0, Saturation: 255, Value: 127}}}},
				{Length: 16, Rows: 2, HueShift: -1, Scan: 10, Begin: 0, End: 100, Rate: 0,
					Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
					Chroma: glow.Chroma{Length: 16, HueShift: -1,
						Colors: []glow.HSV{{Hue: 0, Saturation: 127, Value: 196}, {Hue: 0, Saturation: 255, Value: 127}}}},
			},
		}},
	}

	var hg HeaderGenerator
	var folderList []*iohandler.EffectItems = make([]*iohandler.EffectItems, 0)

	folderList = append(folderList, iohandler.NewFolderList("a", list), iohandler.NewFolderList("b", list2))
	err := hg.Write(folderList)
	if err != nil {
		t.Fatal(err)
	}

	var sg SourceGenerator
	err = sg.Write(folderList)
	if err != nil {
		t.Fatal(err)
	}

	var eg EffectGenerator
	err = eg.Write(folderList)
	if err != nil {
		t.Fatal(err)
	}
}
