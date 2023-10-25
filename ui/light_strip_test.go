package ui

import (
	"glow-gui/glow"
	"os"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
	"gopkg.in/yaml.v3"
)

func TestLightStrip(t *testing.T) {
	app := test.NewApp()
	w := app.NewWindow("Circle")

	strip := NewLightStrip(50, 5)
	if len(strip.lights) != int(strip.length) {
		t.Fatalf("LightStrip got %d want %.0f",
			len(strip.lights), strip.length)
	}

	w.SetContent(strip)
	w.Resize(fyne.NewSize(600, 600))

	showStripDimensions(t, strip)

	t.Logf("background W:%.2f H:%.2f X:%.2f Y:%.2f",
		strip.background.Size().Width, strip.background.Size().Height,
		strip.background.Position().X, strip.background.Position().Y)

	t.Logf("base W:%.2f H:%.2f X:%.2f Y:%.2f",
		strip.BaseWidget.Size().Width, strip.BaseWidget.Size().Height,
		strip.BaseWidget.Position().X, strip.BaseWidget.Position().Y)

	t.Logf("strip W:%.2f H:%.2f X:%.2f Y:%.2f",
		strip.Size().Width, strip.Size().Height,
		strip.Position().X, strip.Position().Y)

	t.Logf("strip X:%.2f Y:%.2f",
		strip.MinSize().Width, strip.MinSize().Height)

	w.ShowAndRun()

	t.Logf("rend X:%.2f Y:%.2f",
		strip.MinSize().Width, strip.MinSize().Height)
}

func showStripDimensions(t *testing.T, strip *LightStrip) {
	for i, c := range strip.lights {
		t.Logf("circle:%d W:%.2f H:%.2f X:%.2f Y:%.2f", i,
			c.Size().Width, c.Size().Height,
			c.Position().X, c.Position().Y)
	}
}

const (
	scheme    = "file://"
	framePath = "/home/dave/src/glow-gui/res/frames/"
)

var files = []string{
	"spotlight.yaml",
	"split.yaml",
	"split_in_two.yaml",
	"black_and_white.yaml",
	"complementary_scan.yaml",
	"double_scan.yaml",
	"gradient_scan.yaml",
}

func TestYamlRead(t *testing.T) {
	var buf []byte = make([]byte, 2048)

	frames := make([]glow.Frame, len(files))
	for i := range files {
		readOs(t, framePath+files[i], &frames[i], buf)
	}

}

func readOs(t *testing.T,
	fname string, frame *glow.Frame, buf []byte) {

	f, err := os.Open(fname)
	if err != nil {
		t.Fatal(err.Error())
	}

	b := buf

	count, err := f.Read(b)
	if err != nil {
		t.Fatal(err.Error())
	}
	f.Close()

	t.Logf("%s %d bytes read\n%s", fname, count, string(b))

	b = b[:count]
	err = yaml.Unmarshal(b, frame)
	if err != nil {
		t.Fatalf(err.Error())
	}

	frame.Setup(36, 4)
	frame.SetInterval(16)

	t.Logf("%v", frame)
}

func readStorage(t *testing.T,
	fname string, frame *glow.Frame, buf []byte) {

	uri, err := storage.ParseURI(scheme + fname)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("name:%s, scheme:%s mimetype=%s extension=%s",
		uri.Name(),
		uri.Scheme(),
		uri.MimeType(),
		uri.Extension())

	f, err := storage.Reader(uri)
	if err != nil {
		t.Fatal(err.Error())
	}

	b := buf

	count, err := f.Read(b)
	if err != nil {
		t.Fatal(err.Error())
	}
	f.Close()

	t.Logf("%d bytes read\n%s", count, string(b))

	b = b[:count]
	err = yaml.Unmarshal(b, frame)
	if err != nil {
		t.Fatalf(err.Error())
	}

	frame.Setup(36, 4)
	frame.SetInterval(16)

	t.Logf("%v", frame)

}

func TestStorage(t *testing.T) {
	test.NewApp()
	var buf []byte = make([]byte, 2048)
	frames := make([]glow.Frame, len(files))
	for i := range files {
		readStorage(t, framePath+files[i], &frames[i], buf)
	}
}
