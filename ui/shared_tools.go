package ui

import (
	"glow-gui/control"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type SharedTools struct {
	*widget.Toolbar
	saveButton *ButtonItem
	undoButton *ButtonItem
	model      *control.Model
}

func NewSharedTools(model *control.Model) *SharedTools {
	tl := &SharedTools{
		Toolbar: widget.NewToolbar(),
		model:   model,
	}

	tl.saveButton = NewButtonItem(
		widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), tl.save))
	tl.undoButton = NewButtonItem(
		widget.NewButtonWithIcon("", theme.ContentUndoIcon(), tl.undo))

	tl.model.AddChangeListener(binding.NewDataListener(func() {
		if tl.model.HasChanged() {
			tl.saveButton.Enable()
			tl.undoButton.Enable()
		} else {
			tl.undoButton.Disable()
			tl.saveButton.Disable()
		}
	}))

	// tl.model.AddUndoListener(binding.NewDataListener(func() {
	// 	if tl.model.CanUndo() {
	// 	} else {
	// 	}
	// }))

	tl.AddItems(tl.saveButton, tl.undoButton)
	return tl
}

func (tl *SharedTools) AddItems(items ...widget.ToolbarItem) {
	tl.Toolbar.Items = append(tl.Toolbar.Items, items...)
}

func (tl *SharedTools) save() {
	tl.model.WriteEffect()
}

func (tl *SharedTools) undo() {
	tl.model.UndoEffect()
}
