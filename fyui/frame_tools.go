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
	effect.AddChangeListener(binding.NewDataListener(func() {
		if effect.HasChanged() {
			saveButton.Enable()
			return
		}
		saveButton.Disable()
	}))

	createDialog := NewEffectDialog(effect, window)
	ft.Toolbar.Append(NewButtonItem(
		widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
			createDialog.Start()
		})))

	folderDialog := NewFolderDialog(effect, window)
	ft.Toolbar.Append(NewButtonItem(
		widget.NewButtonWithIcon("", theme.FolderNewIcon(), func() {
			folderDialog.Start()
		})))

	ft.Toolbar.Append(NewButtonItem(
		widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {})))

	return ft
}
