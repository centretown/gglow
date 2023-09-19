package ui

import (
	"glow-gui/glow"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestLightStrip(t *testing.T) {

	frame := glow.Frame{Length: 36, Rows: 4}
	strip := NewLightStrip(&frame)
	if len(strip.objects) != int(frame.Length) {
		t.Fatalf("LightStrip got %d want %d",
			len(strip.objects), frame.Length)
	}

	app := test.NewApp()
	w := app.NewWindow("Circle")
	w.SetContent(strip)
	w.Resize(fyne.NewSize(600, 600))

	r := strip.CreateRenderer()
	t.Logf("base W:%.2f H:%.2f X:%.2f Y:%.2f",
		strip.background.Size().Width, strip.background.Size().Height,
		strip.background.Position().X, strip.background.Position().Y)

	t.Logf("base W:%.2f H:%.2f X:%.2f Y:%.2f",
		strip.DisableableWidget.Size().Width, strip.DisableableWidget.Size().Height,
		strip.DisableableWidget.Position().X, strip.DisableableWidget.Position().Y)

	showStripDimensions(t, strip)
	r.Layout(strip.DisableableWidget.Size())
	showStripDimensions(t, strip)
	w.ShowAndRun()
}

func showStripDimensions(t *testing.T, strip *LightStrip) {
	for i, c := range strip.objects {
		t.Logf("circle:%d W:%.2f H:%.2f X:%.2f Y:%.2f", i,
			c.Size().Width, c.Size().Height,
			c.Position().X, c.Position().Y)
	}

}
