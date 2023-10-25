package glow

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func compareFrames(t *testing.T, frame, frame2 *Frame) {
	if frame2.Length != frame.Length {
		t.Fatalf("Frame.Length want: %d got %d",
			frame.Length, frame2.Length)
	}
	if frame2.Rows != frame.Rows {
		t.Fatalf("Frame.Rows want: %d got %d",
			frame.Rows, frame2.Rows)
	}
	if frame2.Interval != frame.Interval {
		t.Fatalf("Frame.Interval want: %d got %d",
			frame.Interval, frame2.Interval)
	}
	if len(frame2.Layers) != len(frame.Layers) {
		t.Fatalf("Frame.Interval want: %d got %d",
			frame.Interval, frame2.Interval)
	}
	for i := range frame.Layers {
		compareLayers(t, &frame.Layers[i], &frame2.Layers[i])
	}
}

func TestFrameBasic(t *testing.T) {
	var chroma Chroma
	chroma.AddColors(HSV{HueRed, 1, 1}, HSV{HueBlue, 1, 1})

	var grid Grid
	grid.Orientation = Diagonal
	grid.Origin = TopLeft

	var layer1 Layer
	layer1.Grid = grid
	layer1.Chroma = chroma

	var layer2 Layer
	layer2.Grid = grid
	layer2.Chroma = chroma
	layer2.Scan = 4
	layer2.HueShift = -1

	var frame Frame
	frame.AddLayers(layer1, layer2)
	err := frame.Setup(36, 4)
	if err != nil {
		t.Fatalf(err.Error())
	}
	frame.SetInterval(16)

	layer1.Scan = 1
	layer2.Scan = 3
	frame.AddLayers(layer1, layer2)

	var buffer []byte
	buffer, err = yaml.Marshal(&frame)
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf(string(buffer))

	var frame2 Frame
	err = yaml.Unmarshal(buffer, &frame2)
	if err != nil {
		t.Fatalf(err.Error())
	}
	compareFrames(t, &frame, &frame2)

	buffer, err = json.Marshal(&frame)
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf(string(buffer))

	err = json.Unmarshal(buffer, &frame2)
	if err != nil {
		t.Fatalf(err.Error())
	}
	compareFrames(t, &frame, &frame2)
}
