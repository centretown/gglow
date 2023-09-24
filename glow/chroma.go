package glow

import (
	"fmt"
	"image/color"
)

type Chroma struct {
	Length      uint16 `yaml:"length" json:"length"`
	HueShift    int16  `yaml:"hue_shift" json:"hue_shift"`
	Colors      []HSV  `yaml:"colors" json:"colors"`
	segmentSize uint16
	quick_color color.RGBA
}

func (chroma *Chroma) Setup(length uint16,
	source HSV, target HSV, hueShift int16) error {
	chroma.Length = length
	chroma.HueShift = hueShift
	chroma.Colors = append(chroma.Colors, source, target)
	return chroma.Validate()
}

func (chroma *Chroma) SetupLength(length uint16, hueShift int16) error {
	chroma.Length = length
	chroma.HueShift = hueShift
	return chroma.Validate()
}

func (chroma *Chroma) Validate() error {
	if chroma.Length == 0 {
		return fmt.Errorf("Chroma.Setup zero length")
	}

	if len(chroma.Colors) == 0 {
		chroma.Colors = append(chroma.Colors, HSV{0, 0, 1})
	}

	chroma.quick_color = chroma.Colors[0].ToRGB()
	size := uint16(len(chroma.Colors) - 1)
	if size < 2 {
		chroma.segmentSize = chroma.Length
	} else {
		chroma.segmentSize = chroma.Length / size
	}
	return nil
}

func (chroma *Chroma) Map(index uint16) color.RGBA {
	size := uint16(len(chroma.Colors))
	if size < 2 || index == 0 {
		return chroma.quick_color
	}
	colorIndex := index / chroma.segmentSize
	offset := index % chroma.segmentSize
	first := chroma.Colors[colorIndex]
	last := chroma.Colors[colorIndex+1]
	result := first.ToGradient(last, offset, chroma.segmentSize)
	return result.ToRGB()
}

func (chroma *Chroma) UpdateColors() {
	if chroma.HueShift == 0 {
		return
	}

	var hsv *HSV
	for i := range chroma.Colors {
		hsv = &chroma.Colors[i]

		hsv.Hue += float32(chroma.HueShift)

		if hsv.Hue >= HueMax {
			hsv.Hue = 360 - hsv.Hue
		} else if hsv.Hue < 0 {
			hsv.Hue = 360 + hsv.Hue
		}

		// fmt.Println(i, "hue update", hsv.Hue, chroma.HueShift)
	}

	chroma.quick_color = chroma.Colors[0].ToRGB()
}

func (chroma *Chroma) AddColors(hsv ...HSV) {
	chroma.Colors = append(chroma.Colors, hsv...)
}
