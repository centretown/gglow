package glow

import "image/color"

const (
	byte_limit        uint16 = 0xff
	hue_segment_count uint16 = 6
	hue_limit         uint16 = byte_limit * hue_segment_count
	hue_segment       uint16 = hue_limit / hue_segment_count
	hue_red           uint16 = 0
	hue_yellow        uint16 = hue_segment
	hue_green         uint16 = hue_limit * 2 / hue_segment_count
	hue_cyan          uint16 = hue_limit * 3 / hue_segment_count
	hue_blue          uint16 = hue_limit * 4 / hue_segment_count
	hue_magenta       uint16 = hue_limit * 5 / hue_segment_count
)

type HSVColor struct {
	Hue        uint16 `yaml:"hue" json:"hue"`
	Saturation uint8  `yaml:"saturation" json:"saturation"`
	Value      uint8  `yaml:"value" json:"value"`
}

func (hsv *HSVColor) FromRGB(c color.NRGBA) {
	var (
		primary     int = int(max(c.R, c.G, c.B))
		color_range int = primary - int(min(c.R, c.G, c.B))
		hue         int = 0
		diff        int = 0
	)

	if color_range != 0 {
		if c.R == uint8(primary) {
			diff = int(c.G) - int(c.B)
			hue = 255*diff/color_range + B2I(diff < 0)*int(hue_limit)
		} else if c.G == uint8(primary) {
			diff = int(c.B) - int(c.R)
			hue = 255*diff/color_range + int(hue_green)
		} else {
			diff = int(c.R) - int(c.G)
			hue = 255*diff/color_range + int(hue_blue)
		}
	}

	hsv.Hue = uint16(hue)
	hsv.Value = uint8(primary)

	if primary == 0 {
		hsv.Saturation = 0
	} else {
		hsv.Saturation = uint8((color_range * 255) / primary)
	}

}

var hueFuncs = [6]func(hsv *HSVColor) color.NRGBA{
	// red to yellow
	func(hsv *HSVColor) color.NRGBA {
		return color.NRGBA{
			A: 255,
			R: uint8(byte_limit),
			G: uint8(hsv.Hue)}
	},
	// yellow to green
	func(hsv *HSVColor) color.NRGBA {
		return color.NRGBA{
			A: 255,
			R: uint8(hue_green - hsv.Hue),
			G: uint8(byte_limit)}
	},
	// green to cyan
	func(hsv *HSVColor) color.NRGBA {
		return color.NRGBA{
			A: 255,
			G: uint8(byte_limit),
			B: uint8(hsv.Hue - hue_green)}
	},
	// cyan to blue
	func(hsv *HSVColor) color.NRGBA {
		return color.NRGBA{
			A: 255,
			G: uint8(hue_blue - hsv.Hue),
			B: uint8(byte_limit)}
	},
	// blue to magenta
	func(hsv *HSVColor) color.NRGBA {
		return color.NRGBA{
			A: 255,
			B: uint8(byte_limit),
			R: uint8(hsv.Hue - hue_blue)}
	},
	// magenta to red
	func(hsv *HSVColor) color.NRGBA {
		return color.NRGBA{
			A: 255,
			B: uint8(hue_limit - hsv.Hue),
			R: uint8(byte_limit)}
	},
}

func (hsv *HSVColor) ToRGB() (c color.NRGBA) {
	index := hsv.Hue / hue_segment // 0 - 5
	c = hueFuncs[index](hsv)

	var saturation_multiplier uint16 = 1 + uint16(hsv.Saturation)
	var saturation_added uint16 = hue_segment - uint16(hsv.Saturation)
	var value_multiplier uint16 = 1 + uint16(hsv.Value)

	var color_result uint16 = (uint16(c.R)*saturation_multiplier)>>8 +
		saturation_added
	c.R = uint8((color_result * value_multiplier) >> 8)

	color_result = (uint16(c.G)*saturation_multiplier)>>8 +
		saturation_added
	c.G = uint8(color_result * value_multiplier >> 8)

	color_result = (uint16(c.B)*saturation_multiplier)>>8 +
		saturation_added
	c.B = uint8(color_result * value_multiplier >> 8)

	return c
}

// func (hsv *HSVColor) ToGradient(target HSVColor, index, length uint16, reverse bool) HSVColor {
// 	var dh int = (int(hsv.Hue)-int(target.Hue))*B2I(reverse) +
// 		(int(target.Hue)-int(hsv.Hue))*B2I(!reverse)
// 	hue *= index
// 	hue /= int(length)

// 	var ds int = int(hsv.Saturation) - int(target.Saturation)*dir
// 	var dv int = int(hsv.Value) - int(target.Value)*dir

// 	target.Hue += uint16(B2I(hsv.Hue > target.Hue)) * hue_limit
// 	var gradient_hue uint16 = hsv.Hue + ((target.Hue-hsv.Hue)*index)/length
// 	var gradient_saturation int16 = uint16(hsv.Saturation) +
// 		(uint16(target.Saturation)-uint16(hsv.Saturation))*index/length
// }
