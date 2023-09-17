package glow

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func compareLayers(t *testing.T, want, got *Layer) {
	if got.Length != want.Length {
		t.Fatalf("Layer Length got %d want %d",
			got.Length, want.Length)
	}
	if got.Chroma.Length != want.Length {
		t.Fatalf("Chroma Length not Layer Length got %d want %d",
			got.Chroma.Length, want.Length)
	}
	if got.Chroma.Length != want.Chroma.Length {
		t.Fatalf("Chroma Length got %d want %d",
			got.Chroma.Length, want.Chroma.Length)
	}
	if got.Chroma.HueShift != want.HueShift {
		t.Fatalf("Chroma HueShift got %d want %d",
			got.Chroma.HueShift, want.HueShift)
	}
	if got.Chroma.HueShift != want.Chroma.HueShift {
		t.Fatalf("Chroma HueShift got %d want %d",
			got.Chroma.HueShift, want.Chroma.HueShift)
	}
	if got.Grid.Length != want.Length {
		t.Fatalf("Grid Length not Layer Length got %d want %d",
			got.Grid.Length, want.Length)
	}
	if got.Grid.Rows != want.Rows {
		t.Fatalf("Grid Rows not Layer Length got %d want %d",
			got.Grid.Rows, want.Rows)
	}
	if got.Grid.Length != want.Grid.Length {
		t.Fatalf("Grid Length got %d want %d",
			got.Grid.Length, want.Grid.Length)
	}
	if got.Grid.Rows != want.Grid.Rows {
		t.Fatalf("Grid Length got %d want %d",
			got.Grid.Rows, want.Grid.Rows)
	}
	if got.Grid.Orientation != want.Grid.Orientation {
		t.Fatalf("Grid Orientation got %d want %d",
			got.Grid.Orientation, want.Grid.Orientation)
	}
	if got.Grid.Origin != want.Grid.Origin {
		t.Fatalf("Grid Origin got %d want %d",
			got.Grid.Origin, want.Grid.Origin)
	}
}

func TestLayerBasic(t *testing.T) {
	var chroma Chroma
	chroma.AddColors(HSV{HueRed, 1, 1}, HSV{HueBlue, 1, 1})

	var grid Grid
	grid.Orientation = Diagonal
	grid.Origin = TopLeft

	var layer Layer
	if err := layer.Setup(20, 4, &grid, &chroma, -1, 0, 0, 100); err != nil {
		t.Fatalf(err.Error())
	}

}

func TestLayerYAML(t *testing.T) {
	var chroma Chroma
	chroma.AddColors(HSV{HueRed, 1, 1}, HSV{HueBlue, 1, 1})

	var grid Grid
	grid.Orientation = Diagonal
	grid.Origin = TopLeft

	var layer Layer
	if err := layer.Setup(20, 4, &grid, &chroma, -1, 0, 0, 100); err != nil {
		t.Fatalf(err.Error())
	}

	var (
		err error
		d   []byte
	)

	d, err = json.Marshal(&layer)
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("\n%s", string(d))

	d, err = yaml.Marshal(&layer)
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("\n%s", string(d))

	err = yaml.Unmarshal(d, &layer)
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("\n%v\n", layer)
}
