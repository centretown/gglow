package ui

import (
	"glow-gui/data"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type LayerTools struct {
	*widget.Toolbar
	// editor      *LayerEditor
	model        *data.Model
	ApplyButton  *ButtonItem
	RevertButton *ButtonItem
	InsertButton *ButtonItem
	RemoveButton *ButtonItem
}

func NewLayerTools(model *data.Model) *LayerTools {
	lt := &LayerTools{
		model: model,
		// editor: editor,
	}

	lt.ApplyButton = NewButtonItem(
		widget.NewButtonWithIcon("", theme.ConfirmIcon(), lt.apply))
	lt.RevertButton = NewButtonItem(
		widget.NewButtonWithIcon("", theme.CancelIcon(), lt.revert))
	lt.InsertButton = NewButtonItem(
		widget.NewButtonWithIcon("", theme.ContentAddIcon(), lt.remove))
	lt.RemoveButton = NewButtonItem(
		widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), lt.remove))

	lt.Toolbar = widget.NewToolbar(
		lt.ApplyButton,
		lt.RevertButton,
		lt.InsertButton,
		lt.RemoveButton,
	)
	return lt
}

// func (lt *LayerTools) Show() {
// 	lt.Toolbar.Show()
// }

// func (lt *LayerTools) Hide() {
// 	lt.Toolbar.Hide()
// }

func (lt *LayerTools) apply() {
}

func (lt *LayerTools) revert() {
}

func (lt *LayerTools) add() {
}

func (lt *LayerTools) remove() {
}
