package ui

import (
	"fyne.io/fyne/v2/widget"
)

func NewLayerToolbar() *widget.Toolbar {
	toolbar := widget.NewToolbar()
	toolbar.Append(NewButtonItemFromMenu(MenuLayerAdd))
	toolbar.Append(NewButtonItemFromMenu(MenuLayerInsert))
	toolbar.Append(NewButtonItemFromMenu(MenuLayerRemove))
	// AddGlobalShortCut(window,
	// 	&GlobalShortCut{Shortcut: CtrlL, Action: lt.addLayer})
	return toolbar
}
