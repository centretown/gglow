package ui

import (
	"glow-gui/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type FrameTools struct {
	frameMenu *ButtonItem

	newFrame     *ButtonItem
	newFolder    *ButtonItem
	saveFrame    *ButtonItem
	deleteFrame  *ButtonItem
	createDialog *EffectDialog
	folderDialog *FolderDialog

	toolBar *widget.Toolbar
	popUp   *widget.PopUp
	window  fyne.Window
}

func NewFrameTools(model *data.Model, window fyne.Window) *FrameTools {
	ft := &FrameTools{
		window: window,
	}

	ft.createDialog = NewEffectDialog(window, model)
	ft.folderDialog = NewFolderDialog(window, model)

	ft.frameMenu = NewButtonItem(
		widget.NewButtonWithIcon("", theme.FileIcon(), ft.menu))

	ft.saveFrame = NewButtonItem(
		widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), ft.save))
	ft.newFolder = NewButtonItem(
		widget.NewButtonWithIcon("", theme.FolderNewIcon(), func() {
			ft.popUp.Hide()
			ft.folderDialog.Start()
		}))
	ft.newFrame = NewButtonItem(
		widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
			ft.popUp.Hide()
			ft.createDialog.Start()
		}))
	ft.deleteFrame = NewButtonItem(
		widget.NewButtonWithIcon("", theme.DeleteIcon(), ft.delete))
	ft.toolBar = widget.NewToolbar(
		ft.newFolder,
		ft.newFrame,
		ft.deleteFrame,
	)

	ft.popUp = widget.NewPopUp(ft.toolBar, window.Canvas())

	return ft
}

func (ft *FrameTools) Items() (items []widget.ToolbarItem) {
	items = []widget.ToolbarItem{
		ft.saveFrame,
		ft.frameMenu,
		// ft.newFrame,
		// ft.saveFrame,
		// ft.deleteFrame,
	}
	return
}

func (ft *FrameTools) save() {
	ft.popUp.Hide()
}

func (ft *FrameTools) delete() {
	ft.popUp.Hide()
}

func (ft *FrameTools) menu() {
	pos := fyne.Position{
		X: -ft.popUp.MinSize().Width/2 + ft.frameMenu.Button.MinSize().Width/2,
		Y: -ft.popUp.MinSize().Height,
	}

	ft.popUp.ShowAtRelativePosition(pos, ft.frameMenu.Button)
}
