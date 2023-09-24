package store

import (
	"glow-gui/glow"
	"testing"

	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
	"gopkg.in/yaml.v3"
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

func TestSetup(t *testing.T) {
	test.NewApp()
	err := Setup()

	if err != nil {
		t.Fatalf(err.Error())
	}

	uri := FrameURI
	if uri == nil {
		t.Fatalf("nil FrameURI")
	}

	canList, err := storage.CanList(uri)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if canList == false {
		t.Fatalf("can't list %v", uri.Path())
	}
}

func TestLoadFrame(t *testing.T) {
	test.NewApp()

	for i := range files {
		var frame glow.Frame
		err := LoadFrame(FramePath+files[i], &frame)
		if err != nil {
			t.Fatalf(err.Error())
		}

		var b []byte
		frame.Setup(100, 5, 32)
		b, err = yaml.Marshal(&frame)
		if err != nil {
			t.Fatalf(err.Error())
		}

		t.Logf("%s:", files[i])
		t.Logf("%s", string(b))
	}

}

func TestStoreFrame(t *testing.T) {
	test.NewApp()

	var (
		frame glow.Frame
		fname = "empty_frame.yaml"
	)

	err := StoreFrame(DerivedPath+fname, &frame)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = LoadFrame(DerivedPath+fname, &frame)
	if err != nil {
		t.Fatalf(err.Error())
	}

	var b []byte

	b, err = yaml.Marshal(&frame)
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Logf("%s:", fname)
	t.Logf("\n%s", string(b))

}
