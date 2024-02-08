package glow

import (
	"encoding/json"
	"image/color"
	"testing"

	"gopkg.in/yaml.v3"
)

var test_colors = []color.NRGBA{
	{255, 0, 0, 255},
	{127, 0, 0, 255},
	{63, 0, 0, 255},
	{0, 255, 0, 255},
	{0, 127, 0, 255},
	{0, 63, 0, 255},
	{0, 0, 255, 255},
	{0, 0, 127, 255},
	{0, 0, 63, 255},
	{255, 255, 0, 255},
	{255, 127, 0, 255},
	{255, 63, 0, 255},
	{127, 63, 0, 255},
	{63, 31, 0, 255},
	{255, 0, 255, 255},
	{255, 0, 127, 255},
	{255, 0, 63, 255},
	{64, 0, 32, 255},
	{255, 255, 0, 255},
	{0, 255, 255, 255},
	{55, 10, 11, 255},
}

func TestHSVfConversions(t *testing.T) {
	var hsvf HSV
	for i, item := range test_colors {
		hsvf.FromRGB(item)
		var color color.NRGBA = hsvf.ToRGB()
		if color != item {
			t.Errorf("%d want(%d,%d,%d) hsv(%f,%f,%f) got(%d,%d,%d)",
				i, item.R, item.G, item.B,
				hsvf.Hue, hsvf.Saturation, hsvf.Value,
				color.R, color.G, color.B)
		}
	}
}

func TestHSVToGradient(t *testing.T) {
	hsv := HSV{0, 1, 1}
	expected := HSV{180, .5, .5}
	length := uint16(5)
	var actual HSV
	for i := uint16(0); i <= length; i++ {
		actual = hsv.ToGradient(expected, i, length)
	}

	if actual != expected {
		t.Logf("hsv(%f,%f,%f)",
			hsv.Hue, hsv.Saturation, hsv.Value)
		t.Logf("target(%f,%f,%f)",
			expected.Hue, expected.Saturation, expected.Value)
		t.Errorf("actual(%f,%f,%f)",
			actual.Hue, actual.Saturation, actual.Value)
	}
}

func TestYamlHSV(t *testing.T) {
	var (
		err     error
		buffer  []byte
		want    HSV
		got     HSV
		hsv_set = make([]HSV, len(test_colors))
	)

	for i, rgba := range test_colors {
		want.FromRGB(rgba)
		hsv_set[i] = want

		buffer, err = yaml.Marshal(&rgba)
		if err != nil {
			t.Fatalf(err.Error())
		}
		t.Logf("\n%s", string(buffer))

		got_rgba := want.ToRGB()
		if got_rgba != rgba {
			t.Logf("hsv(%f,%f,%f)",
				want.Hue, want.Saturation, want.Value)
			t.Logf("target(%d,%d,%d)",
				rgba.R, rgba.G, rgba.B)
			t.Errorf("actual(%d,%d,%d)",
				got_rgba.R, got_rgba.G, got_rgba.B)
		}

		buffer, err = yaml.Marshal(&want)
		if err != nil {
			t.Fatalf(err.Error())
		}

		err = yaml.Unmarshal(buffer, &got)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if got != want {
			t.Logf("want(%f,%f,%f)",
				want.Hue, want.Saturation, want.Value)
			t.Fatalf("got(%f,%f,%f)",
				got.Hue, got.Saturation, got.Value)
		}

		buffer, err = json.Marshal(&want)
		if err != nil {
			t.Fatalf(err.Error())
		}

		err = json.Unmarshal(buffer, &got)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if got != want {
			t.Logf("want(%f,%f,%f)",
				want.Hue, want.Saturation, want.Value)
			t.Fatalf("got(%f,%f,%f)",
				got.Hue, got.Saturation, got.Value)
		}
	}

	buffer, err = yaml.Marshal(hsv_set)
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("\n%s", string(buffer))

	buffer, err = json.MarshalIndent(hsv_set, "", "  ")
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("\n%s", string(buffer))
}
