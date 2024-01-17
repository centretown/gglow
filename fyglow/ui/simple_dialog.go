package ui

import (
	"gglow/fyglow/effectio"
	"gglow/text"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type SimpleDialog struct {
	*dialog.CustomDialog
	Name      binding.String
	NameEntry *widget.Entry
	NameLabel *widget.Label
	Apply     func()
	Cancel    func()

	window      fyne.Window
	effect      *effectio.EffectIo
	ApplyButton *widget.Button
}

func NewSimpleDialog(effect *effectio.EffectIo, window fyne.Window,
	title string, nameLabel string) *SimpleDialog {
	sd := &SimpleDialog{
		window: window,
		effect: effect,
		Name:   binding.NewString(),
	}
	sd.NameLabel = widget.NewLabel(nameLabel)
	sd.NameEntry = widget.NewEntryWithData(sd.Name)

	sep := widget.NewSeparator()
	frm := container.New(layout.NewFormLayout(),
		sd.NameLabel, sd.NameEntry, sep, sep)
	sd.CustomDialog = dialog.NewCustomWithoutButtons(title, frm, window)

	sd.ApplyButton = widget.NewButtonWithIcon(text.ApplyLabel.String(),
		theme.ConfirmIcon(), sd.apply)
	sd.ApplyButton.Disable()

	cancelButton := NewButtonItem(
		widget.NewButtonWithIcon(text.CancelLabel.String(),
			theme.CancelIcon(), sd.cancel))
	sd.CustomDialog.SetButtons([]fyne.CanvasObject{cancelButton, sd.ApplyButton})

	sd.NameEntry.OnSubmitted = func(s string) {
		if sd.ApplyButton.Disabled() {
			sd.cancel()
			return
		}
		sd.apply()
	}
	return sd
}

func (sd *SimpleDialog) cancel() {
	if sd.Cancel != nil {
		sd.Cancel()
	}
	sd.CustomDialog.Hide()
}

func (sd *SimpleDialog) apply() {
	if sd.Apply != nil {
		sd.Apply()
	}
	sd.CustomDialog.Hide()
}
func (sd *SimpleDialog) Start() {
	sd.Name.Set("")
	sd.CustomDialog.Show()
	sd.window.Canvas().Focus(sd.NameEntry)
}
