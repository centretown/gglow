package store

// import (
// 	"glow-gui/glow"
// 	"testing"

// 	"fyne.io/fyne/v2/storage"
// 	"fyne.io/fyne/v2/test"
// 	"gopkg.in/yaml.v3"
// )

// var files = []string{
// 	"AAA_Spotlight.yaml",
// 	"Black_and_White.yaml",
// 	"Rainbow_Diagonal.yaml",
// 	"Rainbow_Horizontal.yaml",
// 	"Rainbow_Vertical.yaml",
// 	"Scan_Complementary.yaml",
// 	"Scan_Double.yaml",
// 	"Scan_Gradient.yaml",
// 	"Split_in_Three.yaml",
// 	"Split_in_Two.yaml",
// }

// func TestSetup(t *testing.T) {
// 	test.NewApp()
// 	err := Setup()

// 	if err != nil {
// 		t.Fatalf(err.Error())
// 	}

// 	uri := Current
// 	if uri == nil {
// 		t.Fatalf("nil FrameURI")
// 	}

// 	canList, err := storage.CanList(uri)
// 	if err != nil {
// 		t.Fatalf(err.Error())
// 	}

// 	if canList == false {
// 		t.Fatalf("can't list %v", uri.Path())
// 	}
// }

// func TestLoadFrame(t *testing.T) {
// 	test.NewApp()

// 	for i := range files {
// 		var frame glow.Frame
// 		err := loadFrame(ExamplesPath+files[i], &frame)
// 		if err != nil {
// 			t.Fatalf(err.Error())
// 		}

// 		var b []byte
// 		frame.Setup(100, 5)
// 		b, err = yaml.Marshal(&frame)
// 		if err != nil {
// 			t.Fatalf(err.Error())
// 		}

// 		t.Logf("%s:", files[i])
// 		t.Logf("%s", string(b))
// 	}

// }

// func TestStoreFrame(t *testing.T) {
// 	test.NewApp()

// 	var (
// 		frame glow.Frame
// 		fname = "empty_frame.yaml"
// 	)

// 	err := StoreFrame(BasePath+fname, &frame)
// 	if err != nil {
// 		t.Fatalf(err.Error())
// 	}

// 	err = loadFrame(BasePath+fname, &frame)
// 	if err != nil {
// 		t.Fatalf(err.Error())
// 	}

// 	var b []byte

// 	b, err = yaml.Marshal(&frame)
// 	if err != nil {
// 		t.Fatalf(err.Error())
// 	}

// 	t.Logf("%s:", fname)
// 	t.Logf("\n%s", string(b))

// }
