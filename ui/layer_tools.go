package ui

import (
	"glow-gui/data"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type LayerTools struct {
	*widget.Toolbar
	// editor      *LayerEditor
	model       *data.Model
	applyLayer  *ButtonItem
	revertLayer *ButtonItem
	insertLayer *ButtonItem
	removeLayer *ButtonItem
}

func NewLayerTools(model *data.Model) *LayerTools {
	lt := &LayerTools{
		model: model,
		// editor: editor,
	}

	lt.applyLayer = NewButtonItem(
		widget.NewButtonWithIcon("", theme.ConfirmIcon(), lt.apply))
	lt.revertLayer = NewButtonItem(
		widget.NewButtonWithIcon("", theme.CancelIcon(), lt.revert))
	lt.insertLayer = NewButtonItem(
		widget.NewButtonWithIcon("", theme.ContentAddIcon(), lt.remove))
	lt.removeLayer = NewButtonItem(
		widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), lt.remove))

	lt.Toolbar = widget.NewToolbar(
		lt.applyLayer,
		lt.revertLayer,
		lt.insertLayer,
		lt.removeLayer,
	)
	lt.Hide()
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
