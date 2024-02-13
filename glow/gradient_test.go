package glow

import (
	"fmt"
	"image"
	"image/color"
	"testing"

	"github.com/disintegration/imaging"
)

var grad_colors = [][]color.NRGBA{
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
		color.NRGBA{255, 0, 0, 255},
		color.NRGBA{0, 255, 255, 255},
		color.NRGBA{0, 255, 0, 255},
		color.NRGBA{0, 0, 255, 255},
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
			// ci, ri := 5, 0
			dst := image.NewNRGBA(grad_rectangles[ri])
			var lg = NewLinearGradient(TopLeft, Vertical, grad_colors[ci])
			lg.Draw(dst)
			title = fmt.Sprintf("%s/gradient-%02x-%02x.png", folder, ri, ci)
			err := imaging.Save(dst, title)
			if err != nil {
				t.Fatal(err)
			}
		}
	}

}
