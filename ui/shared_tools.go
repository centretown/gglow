package ui

import (
	"glow-gui/effects"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type SharedTools struct {
	*widget.Toolbar
	saveButton  *ButtonItem
	applyButton *ButtonItem
	undoButton  *ButtonItem
	effect      effects.Effect
}

func NewSharedTools(effect effects.Effect) *SharedTools {
	tl := &SharedTools{
		Toolbar: widget.NewToolbar(),
		effect:  effect,
	}

	tl.saveButton = NewButtonItem(
		widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), tl.save))
	tl.undoButton = NewButtonItem(
		widget.NewButtonWithIcon("", theme.ContentUndoIcon(), tl.undo))
	tl.applyButton = NewButtonItem(
		widget.NewButtonWithIcon("", theme.ConfirmIcon(), tl.apply))

	tl.effect.AddChangeListener(binding.NewDataListener(func() {
		if tl.effect.HasChanged() {
			tl.saveButton.Enable()
			tl.undoButton.Enable()
			tl.applyButton.Enable()
			return
		}
		if !tl.effect.CanUndo() {
			tl.undoButton.Disable()
		}
		tl.saveButton.Disable()
		tl.applyButton.Disable()
	}))

	// tl.effect.AddUndoListener(binding.NewDataListener(func() {
	// 	if tl.effect.CanUndo() {
	// 	} else {
	// 	}
	// }))

	tl.AddItems(tl.saveButton, tl.applyButton, tl.undoButton)
	return tl
}

func (tl *SharedTools) AddItems(items ...widget.ToolbarItem) {
	tl.Toolbar.Items = append(tl.Toolbar.Items, items...)
}

func (tl *SharedTools) save() {
	tl.effect.WriteEffect()
}

func (tl *SharedTools) apply() {
	tl.effect.Apply()
}

func (tl *SharedTools) undo() {
	tl.effect.UndoEffect()
}
