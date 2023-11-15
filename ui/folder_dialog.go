package ui

import (
	"glow-gui/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type FolderDialog struct {
	*dialog.CustomDialog
	title       binding.String
	model       *data.Model
	applyButton *widget.Button
}

func NewFolderDialog(window fyne.Window, model *data.Model) *FolderDialog {
	fd := &FolderDialog{
		model: model,
		title: binding.NewString(),
	}

	nameLabel := widget.NewLabel("Folder")
	nameEntry := widget.NewEntryWithData(fd.title)
	nameEntry.Validator = validation.NewAllStrings(fd.validateFolderName)

	sep := widget.NewSeparator()
	frm := container.New(layout.NewFormLayout(),
		nameLabel, nameEntry, sep, sep)

	fd.CustomDialog = dialog.NewCustomWithoutButtons("Add Folder", frm, window)
	fd.applyButton = widget.NewButtonWithIcon("Apply", theme.ConfirmIcon(), fd.apply)
	fd.applyButton.Disable()

	revertButton := NewButtonItem(
		widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
			fd.CustomDialog.Hide()
		}))
	fd.CustomDialog.SetButtons([]fyne.CanvasObject{revertButton, fd.applyButton})
	return fd
}

func (fd *FolderDialog) apply() {
	title, _ := fd.title.Get()
	err := fd.model.Store.CreateNewFolder(title)
	if err != nil {
		fyne.LogError(title, err)
	}
	fd.CustomDialog.Hide()
}

func (fd *FolderDialog) validateFolderName(s string) error {
	err := fd.model.Store.ValidateNewFolderName(s)
	if err != nil {
		fd.applyButton.Disable()
		return err
	}
	fd.applyButton.Enable()
	return nil
}

func (fd *FolderDialog) Start() {
	fd.title.Set("")
	fd.CustomDialog.Show()
}
