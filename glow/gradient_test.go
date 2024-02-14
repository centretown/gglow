package glow

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/disintegration/imaging"
)

var grad_colors = [][]color.NRGBA{
	{},
	{
		color.NRGBA{255, 255, 255, 255},
	},
	{
		color.NRGBA{255, 0, 0, 255},
		color.NRGBA{0, 255, 255, 255},
	},
	{
		color.NRGBA{255, 0, 0, 255},
		color.NRGBA{0, 255, 255, 255},
		color.NRGBA{0, 255, 0, 255},
	},
	{
		color.NRGBA{21, 21, 0, 255},
		color.NRGBA{21, 119, 33, 255},
		color.NRGBA{0, 64, 0, 255},
		color.NRGBA{0, 31, 0, 255},
	},
	{
		color.NRGBA{255, 0, 0, 255},
		color.NRGBA{0, 255, 255, 255},
		color.NRGBA{0, 255, 0, 255},
		color.NRGBA{0, 0, 255, 255},
		color.NRGBA{255, 0, 255, 255},
	},
	{
		color.NRGBA{255, 0, 0, 255},
		color.NRGBA{0, 255, 255, 255},
		color.NRGBA{0, 255, 0, 255},
		color.NRGBA{0, 0, 255, 255},
		color.NRGBA{255, 0, 255, 255},
		color.NRGBA{255, 0, 0, 255},
		color.NRGBA{0, 255, 255, 255},
	},
	{
		color.NRGBA{255, 0, 0, 255},
		color.NRGBA{0, 255, 255, 255},
		color.NRGBA{0, 255, 0, 255},
		color.NRGBA{0, 0, 255, 255},
		color.NRGBA{255, 0, 255, 255},
		color.NRGBA{255, 0, 0, 255},
		color.NRGBA{0, 255, 255, 255},
		color.NRGBA{0, 0, 255, 255},
	},
}

var grad_rectangles = []image.Rectangle{
	{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: 433, Y: 139},
	},
	{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: 219, Y: 337},
	},
}

func showDeltaAndColors(delta *Delta) {
	for i, d := range delta.segments {
		fmt.Printf("%04d - Begin=%04d End=%04d (R:%04d,G:%04d,B:%04d,A:%04d)\n",
			i, d.Begin, d.End, d.R, d.G, d.B, d.A)
		fmt.Printf("%02x:%02x:%02x:%02x Color\n", d.C.R, d.C.G, d.C.B, d.C.A)
		for x := d.Begin; x <= d.End; x++ {
			c := d.Point(x)
			fmt.Printf("%02x:%02x:%02x ", c.R, c.G, c.B)
			if x%8 == 7 {
				fmt.Println()
			}
		}
		fmt.Println()
		fmt.Println()
	}
	fmt.Println()
}

func showDelta(delta *Delta) {
	for i, d := range delta.segments {
		fmt.Printf("%04d - Begin=%04d End=%04d (R:%04d,G:%04d,B:%04d,A:%04d)\n",
			i, d.Begin, d.End, d.R, d.G, d.B, d.A)
	}
	fmt.Println()
}

func TestDeltaPoints(t *testing.T) {
	// for j := range 1 {
	// 	for i := range 4 {
	ci := 1
	ri := 0
	rect := grad_rectangles[ri]
	delta := NewDelta(grad_colors[ci], rect.Dx())
	fmt.Printf("Rect %02d Colors %02d\n", ri, ci)
	showDeltaAndColors(delta)
	// 	}
	// }
}

func TestDeltas(t *testing.T) {
	for j := range grad_rectangles {
		for i := range grad_colors {
			rect := grad_rectangles[j]
			deltas := NewDelta(grad_colors[i], rect.Dx())
			showDelta(deltas)
		}
	}
}

func TestLinearGradient(t *testing.T) {

	var (
		folder = "test_pics"
		title  = ""
	)

	for ci := range grad_colors {
		for ri := range grad_rectangles {
			dst := image.NewNRGBA(grad_rectangles[ri])
			var lg = NewLinearGradient(TopRight, Diagonal, grad_colors[ci])
			lg.Draw(dst)
			title = fmt.Sprintf("%s/gradient-%02x-%02x.png", folder, ri, ci)
			err := imaging.Save(dst, title)
			if err != nil {
				t.Fatal(err)
			}
		}
	}

}

func TestAngle2Divisors(t *testing.T) {
	degrees := []float64{0, 11.25, 22.5, 45, 67.5, 90, 180, 360}
	for i, degree := range degrees {
		radians := 2 * degree * math.Pi / 360
		fmt.Printf("%02d\tdeg: %8.3f\trad: %5.3f\tsin: %5.3f\tcos: %5.3f\n",
			i, degree, radians, math.Sin(radians), math.Cos(radians))
	}

}
