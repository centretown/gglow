package store

import (
	"testing"

	"fyne.io/fyne/v2/test"
)

var titles = []string{
	"AAA_Spotlight",
	"Black_and_White",
	"Rainbow_Diagonal",
	"Rainbow_Horizontal",
	"Rainbow_Vertical",
	"Scan_Complementary",
	"Scan_Double",
	"Scan_Gradient",
	"Split_in_Three",
	"Split_in_Two",
}

func TestStoreRead(t *testing.T) {
	app := test.NewApp()
	store := NewStore(app.Preferences())

	t.Log(store.Current.Name(), store.Current.Path())
	t.Log(store.KeyList)

	for _, title := range store.KeyList {
		err := store.ReadEffect(title)
		if err != nil && !store.IsFolder(title) {
			t.Fatalf(err.Error())
		}
		t.Log(title, store.Current.Name(), store.EffectName)
	}
}

func TestStoreUndo(t *testing.T) {
	app := test.NewApp()
	store := NewStore(app.Preferences())
	t.Log(store.Current.Name(), store.EffectName)

	if len(store.KeyList) < 1 {
		t.Fatal("No Keys")
	}
	title := store.KeyList[0]

	err := store.ReadEffect(title)
	if err != nil {
		t.Fatalf(err.Error())
	}

	frame := store.GetFrame()
	t.Log("before", title, store.Current.Name(), store.EffectName, frame.Interval)

	interval := frame.Interval
	store.UpdateHistory()
	frame.Interval++
	store.SetDirty(true)

	t.Log("new interval after update", frame.Interval)

	frame, err = store.Undo(title)
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Log("after undo 1", frame.Interval)
	if interval != frame.Interval {
		t.Fatal("undo level 1")
	}

	frame.Interval += 2
	store.SetDirty(true)
	frame, err = store.Undo(title)
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Log("after undo 0", frame.Interval)
	if interval != frame.Interval {
		t.Fatal("undo level 0")
	}

}

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
