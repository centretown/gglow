package codeio

import (
	"gglow/glow"
	"testing"
)

func TestGenerate(t *testing.T) {
	list := []*EffectItem{
		{"black white scan", "", &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []glow.Layer{{Length: 16, Rows: 2,
				Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
				Chroma: glow.Chroma{Length: 16, HueShift: -1,
					Colors: []glow.HSV{{Hue: 0, Saturation: 255, Value: 127}}}}}}},
		{"complementary scan", "", &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []glow.Layer{{Length: 16, Rows: 2,
				Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
				Chroma: glow.Chroma{Length: 16, HueShift: -1,
					Colors: []glow.HSV{{Hue: 0, Saturation: 255, Value: 127}}}}}}},
		{"double scan", "", &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []glow.Layer{{Length: 16, Rows: 2,
				Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
				Chroma: glow.Chroma{Length: 16, HueShift: -1,
					Colors: []glow.HSV{{Hue: 0, Saturation: 255, Value: 127}}}}}}},
		{"gradient scan", "", &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []glow.Layer{{Length: 16, Rows: 2,
				Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
				Chroma: glow.Chroma{Length: 16, HueShift: -1,
					Colors: []glow.HSV{{Hue: 0, Saturation: 255, Value: 127}}}}}}},
		{"rainbow diagonal", "", &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []glow.Layer{{Length: 16, Rows: 2,
				Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
				Chroma: glow.Chroma{Length: 16, HueShift: -1,
					Colors: []glow.HSV{{Hue: 0, Saturation: 255, Value: 127}}}}}}},
		{"rainbow horizontal", "", &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []glow.Layer{{Length: 16, Rows: 2,
				Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
				Chroma: glow.Chroma{Length: 16, HueShift: -1,
					Colors: []glow.HSV{{Hue: 0, Saturation: 255, Value: 127}}}}}}},
		{"rainbow vertical", "", &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []glow.Layer{{Length: 16, Rows: 2,
				Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
				Chroma: glow.Chroma{Length: 16, HueShift: -1,
					Colors: []glow.HSV{{Hue: 0, Saturation: 255, Value: 127}}}}}}},
		{"split in three", "", &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []glow.Layer{{Length: 16, Rows: 2,
				Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
				Chroma: glow.Chroma{Length: 16, HueShift: -1,
					Colors: []glow.HSV{{Hue: 0, Saturation: 255, Value: 127}}}}}}},
		{"split in two", "", &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []glow.Layer{{Length: 16, Rows: 2,
				Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
				Chroma: glow.Chroma{Length: 16, HueShift: -1,
					Colors: []glow.HSV{{Hue: 0, Saturation: 255, Value: 127}}}}}}},
		{"spotlight", "", &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []glow.Layer{
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
	list2 := []*EffectItem{
		{"black white scan", "", &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []glow.Layer{{Length: 16, Rows: 2,
				Grid: glow.Grid{Length: 16, Rows: 2, Origin: glow.TopLeft, Orientation: glow.Horizontal},
				Chroma: glow.Chroma{Length: 16, HueShift: -1,
					Colors: []glow.HSV{{Hue: 0, Saturation: 255, Value: 127}}}}}}},
		{"spotlight", "", &glow.Frame{Length: 16, Rows: 2, Interval: 48,
			Layers: []glow.Layer{
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
	var folderList []*FolderList = make([]*FolderList, 0)

	folderList = append(folderList, NewFolderList("a", list), NewFolderList("b", list2))
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
