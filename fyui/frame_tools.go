package fyui

import (
	"gglow/fyio"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type FrameTools struct {
	*widget.Toolbar
}

func NewFrameTools(effect *fyio.EffectIo, window fyne.Window) *FrameTools {
	ft := &FrameTools{
		Toolbar: widget.NewToolbar(),
	}

	saveButton := NewButtonItem(
		widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func() {
			effect.SaveEffect()
		}))
	ft.Toolbar.Append(saveButton)

	addEffect := NewEffectDialog(effect, window)
	effectButton := NewButtonItem(
		widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
			addEffect.Start()
		}))
	ft.Toolbar.Append(effectButton)

	effect.AddFolderListener(binding.NewDataListener(func() {
		if effect.IsRootFolder() {
			effectButton.Disable()
		} else {
			effectButton.Enable()
		}
	}))

	ft.Toolbar.Append(NewButtonItem(
		widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {})))

	addFolder := NewFolderDialog(effect, window)
	ft.Toolbar.Append(NewButtonItem(
		widget.NewButtonWithIcon("", theme.FolderNewIcon(), func() {
			addFolder.Start()
		})))

	ft.Toolbar.Append(NewButtonItem(
		widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {})))

	effect.AddChangeListener(binding.NewDataListener(func() {
		if effect.HasChanged() {
			saveButton.Enable()
		} else {
			saveButton.Disable()
		}
	}))

	return ft
}
