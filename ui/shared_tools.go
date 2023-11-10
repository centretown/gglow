package ui

import (
	"glow-gui/data"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type SharedTools struct {
	*widget.Toolbar
	applyButton  *ButtonItem
	revertButton *ButtonItem
	isDirty      binding.Bool
	apply_funcs  []func()
	revert_funcs []func()
}

func NewSharedTools(model *data.Model, isDirty binding.Bool) *SharedTools {
	tl := &SharedTools{
		Toolbar: widget.NewToolbar(),
		isDirty: isDirty,
	}
	tl.applyButton = NewButtonItem(
		widget.NewButtonWithIcon("", theme.ConfirmIcon(), tl.apply))
	tl.applyButton.Disable()
	tl.revertButton = NewButtonItem(
		widget.NewButtonWithIcon("", theme.ContentUndoIcon(), tl.revert))
	tl.revertButton.Disable()

	isDirty.AddListener(binding.NewDataListener(func() {
		b, _ := isDirty.Get()
		if b {
			tl.applyButton.Enable()
			tl.revertButton.Enable()
		} else {
			tl.applyButton.Disable()
			tl.revertButton.Disable()
		}
	}))

	tl.AddItems(tl.applyButton, tl.revertButton)
	return tl
}

func (tl *SharedTools) AddItems(items ...widget.ToolbarItem) {
	tl.Toolbar.Items = append(tl.Toolbar.Items, items...)
}

func (tl *SharedTools) AddApply(f func()) {
	tl.apply_funcs = append(tl.apply_funcs, f)
}

func (tl *SharedTools) AddRevert(f func()) {
	tl.revert_funcs = append(tl.revert_funcs, f)
}

func (tl *SharedTools) apply() {
	for _, f := range tl.apply_funcs {
		f()
	}
	tl.isDirty.Set(false)
}

func (tl *SharedTools) revert() {
	for _, f := range tl.revert_funcs {
		f()
	}
	tl.isDirty.Set(false)
}
