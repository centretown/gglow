package ui

import (
	"glow-gui/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type FrameTools struct {
	newFrame     *ButtonItem
	saveFrame    *ButtonItem
	deleteFrame  *ButtonItem
	createDialog *EffectDialog
	// createDialog *dialog.FileDialog
}

func NewFrameTools(model *data.Model, window fyne.Window) *FrameTools {
	ft := &FrameTools{}

	ft.createDialog = NewEffectDialog(window, model)
	// ft.createDialog = dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
	// }, window)
	// ft.createDialog.SetLocation(model.Store.Current)

	ft.saveFrame = NewButtonItem(
		widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), ft.save))
	ft.newFrame = NewButtonItem(
		widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
			ft.createDialog.Show()
		}))
	ft.deleteFrame = NewButtonItem(
		widget.NewButtonWithIcon("", theme.DeleteIcon(), ft.delete))
	return ft
}

func (ft *FrameTools) Items() (items []widget.ToolbarItem) {
	items = []widget.ToolbarItem{
		ft.newFrame,
		ft.saveFrame,
		ft.deleteFrame,
	}
	return
}

func (ft *FrameTools) save() {
}

func (ft *FrameTools) delete() {
}
