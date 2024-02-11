package glow

import (
	"fmt"
	"image/color"
	"testing"
)

func TestHSVConversions(t *testing.T) {
	var hsvf HSV
	var hsv HSVColor
	for i, item := range test_colors {
		hsv.FromRGB(item)
		hsvf.FromRGB(item)
		var colorf color.NRGBA = hsvf.ToRGB()
		var color color.NRGBA = hsv.ToRGB()
		if color != item {
			t.Errorf("%d want(%d,%d,%d,%d) hsv(%d,%d,%d) got(%d,%d,%d,%d) hsvf(%.0f,%.0f,%.0f) got(%d,%d,%d,%d)",
				i, item.R, item.G, item.B, item.A,
				hsv.Hue, hsv.Saturation, hsv.Value,
				color.R, color.G, color.B, color.A,
				hsvf.Hue,
				hsvf.Saturation,
				hsvf.Value,
				colorf.R, colorf.G, colorf.B, colorf.A,
			)
		}
	}
}
func TestToConversions(t *testing.T) {
	var hsvf HSV
	var hsv HSVColor
	var hsv2 HSVColor
	for i, item := range test_colors {
		hsv.FromRGB(item)
		hsvf.FromRGB(item)
		hsv2.Hue = uint16(hsvf.Hue * float32(hue_limit) / 360)
		hsv2.Saturation = uint8(hsvf.Saturation * 255)
		hsv2.Value = uint8(hsvf.Value * 255)

		if hsv != hsv2 {
			t.Errorf("%d want(%d,%d,%d,%d)  hsvf(%.0f,%.0f,%.0f) hsv2(%d,%d,%d) hsv(%d,%d,%d)",
				i,
				item.R, item.G, item.B, item.A,
				hsvf.Hue, hsvf.Saturation, hsvf.Value,
				hsv2.Hue, hsv2.Saturation, hsv2.Value,
				hsv.Hue, hsv.Saturation, hsv.Value,
			)
		}
		// }
	}
}

// func TestHSVToGradient(t *testing.T) {
// 	hsv := HSV{0, 1, 1}
// 	expected := HSV{180, .5, .5}
// 	length := uint16(5)
// 	var actual HSV
// 	for i := uint16(0); i <= length; i++ {
// 		actual = hsv.ToGradient(expected, i, length)
// 	}

// 	if actual != expected {
// 		t.Logf("hsv(%f,%f,%f)",
// 			hsv.Hue, hsv.Saturation, hsv.Value)
// 		t.Logf("target(%f,%f,%f)",
// 			expected.Hue, expected.Saturation, expected.Value)
// 		t.Errorf("actual(%f,%f,%f)",
// 			actual.Hue, actual.Saturation, actual.Value)
// 	}
// }

// func TestYamlHSV(t *testing.T) {
// 	var (
// 		err     error
// 		buffer  []byte
// 		want    HSV
// 		got     HSV
// 		hsv_set = make([]HSV, len(test_colors))
// 	)

// 	for i, rgba := range test_colors {
// 		want.FromRGB(rgba)
// 		hsv_set[i] = want

// 		buffer, err = yaml.Marshal(&rgba)
// 		if err != nil {
// 			t.Fatalf(err.Error())
// 		}
// 		t.Logf("\n%s", string(buffer))

// 		got_rgba := want.ToRGB()
// 		if got_rgba != rgba {
// 			t.Logf("hsv(%f,%f,%f)",
// 				want.Hue, want.Saturation, want.Value)
// 			t.Logf("target(%d,%d,%d)",
// 				rgba.R, rgba.G, rgba.B)
// 			t.Errorf("actual(%d,%d,%d)",
// 				got_rgba.R, got_rgba.G, got_rgba.B)
// 		}

// 		buffer, err = yaml.Marshal(&want)
// 		if err != nil {
// 			t.Fatalf(err.Error())
// 		}

// 		err = yaml.Unmarshal(buffer, &got)
// 		if err != nil {
// 			t.Fatalf(err.Error())
// 		}

// 		if got != want {
// 			t.Logf("want(%f,%f,%f)",
// 				want.Hue, want.Saturation, want.Value)
// 			t.Fatalf("got(%f,%f,%f)",
// 				got.Hue, got.Saturation, got.Value)
// 		}

// 		buffer, err = json.Marshal(&want)
// 		if err != nil {
// 			t.Fatalf(err.Error())
// 		}

// 		err = json.Unmarshal(buffer, &got)
// 		if err != nil {
// 			t.Fatalf(err.Error())
// 		}

// 		if got != want {
// 			t.Logf("want(%f,%f,%f)",
// 				want.Hue, want.Saturation, want.Value)
// 			t.Fatalf("got(%f,%f,%f)",
// 				got.Hue, got.Saturation, got.Value)
// 		}
// 	}

// 	buffer, err = yaml.Marshal(hsv_set)
// 	if err != nil {
// 		t.Fatalf(err.Error())
// 	}
// 	t.Logf("\n%s", string(buffer))

// 	buffer, err = json.MarshalIndent(hsv_set, "", "  ")
// 	if err != nil {
// 		t.Fatalf(err.Error())
// 	}
// 	t.Logf("\n%s", string(buffer))
// }

func TestCalc(t *testing.T) {
	f := func(c color.NRGBA) {
		var (
			primary     int = int(max(c.R, c.G, c.B))
			color_range int = primary - int(min(c.R, c.G, c.B))
			hue         int = 0
			diff        int = 0
		)
		fmt.Println(primary, color_range)
		if color_range != 0 {
			switch uint8(primary) {
			case c.R:
				diff = int(c.G) - int(c.B)
				hue = 255*diff/color_range + B2I(c.B > c.G)*int(hue_limit)
				fmt.Println("red", diff, hue)
			case c.G:
				diff = int(c.B) - int(c.R)
				hue = 255*diff/color_range + int(hue_green)
				fmt.Println("green", diff, hue)
			case c.B:
				diff = int(c.R) - int(c.G)
				hue = 255*diff/color_range + int(hue_blue)
				fmt.Println("blue", diff, hue)
			}
		}
	}

	f(color.NRGBA{55, 10, 11, 255})
}
