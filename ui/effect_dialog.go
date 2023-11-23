package ui

import (
	"glow-gui/control"
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
	model       *control.Manager
	applyButton *widget.Button
}

func NewEffectDialog(model *control.Manager, window fyne.Window) (ef *EffectDialog) {
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
	err := ef.model.CreateNewEffect(title, frame)
	if err != nil {
		fyne.LogError(title, err)
	}
	ef.CustomDialog.Hide()
}

func (ef *EffectDialog) validateFileName(s string) error {
	err := ef.model.ValidateNewEffectName(s)
	if err != nil {
		ef.applyButton.Disable()
		return err
	}

	ef.title.Set(s)
	ef.applyButton.Enable()
	return err
}

func (ef *EffectDialog) Start() {
	ef.title.Set("")
	ef.CustomDialog.Show()
}
