package store

import (
	"encoding/json"
	"fmt"
	"glow-gui/glow"
	"testing"

	"fyne.io/fyne/v2/test"
)

func dumpItem(t *testing.T, item *glow.Frame) string {
	b, err := json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestHistory(t *testing.T) {
	app := test.NewApp()

	current := &glow.Frame{}
	current.SetInterval(48)
	current.Setup(100, 10)

	store := NewStore(app.Preferences())
	history := store.history

	title := "effect"

	err := history.Add(store.route, title,
		glow.NewGlowState(current, 0))
	if err != nil {
		t.Fatal(err)
	}

	current.SetInterval(300)
	current.Setup(200, 10)

	err = history.Add(store.route, title,
		glow.NewGlowState(current, 0))
	if err != nil {
		t.Fatal(err)
	}

	state, err := history.Previous(store.route, title)
	if err != nil {
		t.Fatal(err)
	}

	frame := state.Frame

	if frame.Interval != 300 {
		t.Fatal(fmt.Errorf("frame.Interval != 300"))
	}
	if frame.Length != 200 {
		t.Fatal(fmt.Errorf("frame.Length != 200"))
	}
	t.Log(dumpItem(t, frame))

	state, err = history.Previous(store.route, title)
	if err != nil {
		t.Fatal(err)
	}
	frame = state.Frame

	if frame.Interval != 48 {
		t.Fatal(fmt.Errorf("frame.Interval != 48"))
	}
	if frame.Length != 100 {
		t.Fatal(fmt.Errorf("frame.Length != 100"))
	}
	t.Log(dumpItem(t, frame))
}
