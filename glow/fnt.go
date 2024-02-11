package glow

import (
	"image"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var fontData = goregular.TTF

func DrawText(title string, rect image.Rectangle, chroma *Chroma) (*image.NRGBA, error) {
	fnt, err := opentype.Parse(fontData)
	if err != nil {
		log.Fatal(err)
	}

	fontsize, height := rect.Max.Y, rect.Max.Y
	face, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size:    float64(fontsize),
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatal(err)
	}

	drawer := font.Drawer{Face: face}
	width := drawer.MeasureString(title).Ceil()
	width += width / 20
	rect.Max.X = width

	dst := image.NewNRGBA(rect)
	src := image.Transparent
	hsv := HSV{Hue: 0, Saturation: 1, Value: 1}

	if len(chroma.Colors) > 0 {
		hsv = chroma.Colors[0]
	}
	c := hsv.ToRGB()

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			dst.SetNRGBA(x, y, c)
		}
	}
	drawer.Dst = dst
	drawer.Src = src
	drawer.Dot = fixed.P(0, height)
	drawer.DrawString(title)

	return dst, nil
}
