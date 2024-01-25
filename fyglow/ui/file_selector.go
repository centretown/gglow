package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func NewFolderSelector(path binding.String, window fyne.Window) fyne.CanvasObject {
	btn := widget.NewButton("Select Folder", func() {
		dlg := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				if err != nil {
					fyne.LogError("ShowFolderOpen", err)
				}
				return
			}
			path.Set(uri.Path())
		}, window)
		dlg.Show()
	})
	return btn
}
