package ui

import (
	"fmt"
	"glow-gui/data"
	"glow-gui/glow"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type EffectDialog struct {
	*dialog.CustomDialog
	title       binding.String
	model       *data.Model
	applyButton *widget.Button
}

func NewEffectDialog(window fyne.Window, model *data.Model) (ef *EffectDialog) {
	ef = &EffectDialog{
		model: model,
		title: binding.NewString(),
	}

	nameLabel := widget.NewLabel("Title")
	nameEntry := widget.NewEntryWithData(ef.title)
	nameEntry.Validator = validation.NewAllStrings(ef.validateFileName)

	sep := widget.NewSeparator()
	frm := container.New(layout.NewFormLayout(),
		nameLabel, nameEntry, sep, sep)

	ef.CustomDialog = dialog.NewCustomWithoutButtons("Create Effect", frm, window)
	ef.applyButton = widget.NewButtonWithIcon("Apply", theme.ConfirmIcon(), ef.apply)
	ef.applyButton.Disable()

	revertButton := NewButtonItem(
		widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
			ef.CustomDialog.Hide()
		}))

	ef.CustomDialog.SetButtons([]fyne.CanvasObject{revertButton, ef.applyButton})
	return ef
}

func (ef *EffectDialog) apply() {
	title, _ := ef.title.Get()
	frame := &glow.Frame{}
	frame.Interval = uint32(RateBounds.OnVal)
	err := ef.model.Store.CreateNewEffect(title, frame)
	if err != nil {
		fmt.Println(title)
		fmt.Println(err)
	}
	ef.CustomDialog.Hide()
}

func (ef *EffectDialog) validateFileName(s string) (err error) {
	err = ef.model.Store.ValidateNewEffectName(s)
	if err != nil {
		ef.applyButton.Disable()
		fmt.Println(err)
		return
	}

	ef.title.Set(s)
	ef.applyButton.Enable()
	return
}
