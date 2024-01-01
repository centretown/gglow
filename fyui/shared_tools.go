package fyui

import (
	"gglow/iohandler"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type SharedTools struct {
	*widget.Toolbar
	saveButton *ButtonItem
	effect     iohandler.EffectIoHandler
}

func NewSharedTools(effect iohandler.EffectIoHandler) *SharedTools {
	tl := &SharedTools{
		Toolbar: widget.NewToolbar(),
		effect:  effect,
	}

	tl.saveButton = NewButtonItem(
		widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), tl.save))

	tl.effect.AddChangeListener(binding.NewDataListener(func() {
		if tl.effect.HasChanged() {
			tl.saveButton.Enable()
			return
		}
		tl.saveButton.Disable()
	}))

	tl.AddItems(tl.saveButton)
	return tl
}

func (tl *SharedTools) AddItems(items ...widget.ToolbarItem) {
	tl.Toolbar.Items = append(tl.Toolbar.Items, items...)
}

func (tl *SharedTools) save() {
	tl.effect.SaveEffect()
}
