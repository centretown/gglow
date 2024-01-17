package ui

import "fyne.io/fyne/v2/data/binding"

type ExItem struct {
	Folder   string
	Effect   string
	Selected bool
}

type ExportItem struct {
	binding.DataItem
	binding.String
}
