package ui

import (
	"fmt"
	"gglow/fyglow/effectio"
	"gglow/resources"

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
	effect      *effectio.EffectIo
	applyButton *widget.Button
}

func NewFolderDialog(effect *effectio.EffectIo, window fyne.Window) *FolderDialog {
	fd := &FolderDialog{
		effect: effect,
		title:  binding.NewString(),
	}

	nameLabel := widget.NewLabel(resources.FolderLabel.String())
	nameEntry := widget.NewEntryWithData(fd.title)
	nameEntry.Validator = validation.NewAllStrings(fd.validateFolderName)

	sep := widget.NewSeparator()
	frm := container.New(layout.NewFormLayout(),
		nameLabel, nameEntry, sep, sep)

	fd.CustomDialog = dialog.NewCustomWithoutButtons(resources.AddFolderLabel.String(),
		frm, window)
	fd.applyButton = widget.NewButtonWithIcon(resources.ApplyLabel.String(),
		theme.ConfirmIcon(), fd.Apply)
	fd.applyButton.Disable()

	revertButton := NewButtonItem(
		widget.NewButtonWithIcon(resources.CancelLabel.String(),
			theme.CancelIcon(), func() {
				fd.CustomDialog.Hide()
			}))
	fd.CustomDialog.SetButtons([]fyne.CanvasObject{revertButton, fd.applyButton})
	return fd
}

func (fd *FolderDialog) Cancel() {
	fd.CustomDialog.Hide()
}

func (fd *FolderDialog) Apply() {
	fd.CustomDialog.Hide()
	title, _ := fd.title.Get()
	fmt.Println("FolderDialog.Apply", title)
	err := fd.effect.AddFolder(title)
	if err != nil {
		fyne.LogError(title, err)
	}
}

func (fd *FolderDialog) validateFolderName(s string) error {
	err := fd.effect.ValidateNewFolderName(s)
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
