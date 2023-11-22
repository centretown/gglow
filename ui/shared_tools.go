package ui

import (
	"glow-gui/data"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type SharedTools struct {
	*widget.Toolbar
	saveButton  *ButtonItem
	undoButton  *ButtonItem
	model       *data.Model
	apply_funcs []func() `json:"-"`
}

func NewSharedTools(model *data.Model) *SharedTools {
	tl := &SharedTools{
		Toolbar: widget.NewToolbar(),
		model:   model,
	}

	tl.saveButton = NewButtonItem(
		widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), tl.apply))
	// tl.saveButton.Disable()
	tl.undoButton = NewButtonItem(
		widget.NewButtonWithIcon("", theme.ContentUndoIcon(), tl.undo))
	// tl.undoButton.Disable()

	tl.model.AddDirtyListener(binding.NewDataListener(func() {
		if tl.model.IsDirty() {
			tl.saveButton.Enable()
		} else {
			tl.saveButton.Disable()
		}
	}))

	tl.model.AddUndoListener(binding.NewDataListener(func() {
		if tl.model.CanUndo() {
			tl.undoButton.Enable()
		} else {
			tl.undoButton.Disable()
		}
	}))

	tl.AddItems(tl.saveButton, tl.undoButton)
	return tl
}

func (tl *SharedTools) AddItems(items ...widget.ToolbarItem) {
	tl.Toolbar.Items = append(tl.Toolbar.Items, items...)
}

func (tl *SharedTools) AddApply(f func()) {
	tl.apply_funcs = append(tl.apply_funcs, f)
}

func (tl *SharedTools) apply() {
	for _, f := range tl.apply_funcs {
		f()
	}

	tl.model.WriteEffect()
}

func (tl *SharedTools) undo() {
	tl.model.UndoEffect()
}
