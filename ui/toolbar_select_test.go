package ui

import (
	"glow-gui/res"
	"glow-gui/store"
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestToolBarSelect(t *testing.T) {
	app := test.NewApp()
	window := app.NewWindow(res.WindowTitle)
	store.Setup()

	gui := NewUi(app, window)
	tool := NewToolbarSelect(gui)

	if tool == nil {
		t.Fatalf("nil tool")
	}

	t.Logf("%v", tool.Chooser.Options)
}
