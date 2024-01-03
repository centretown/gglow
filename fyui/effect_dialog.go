package fyui

import (
	"gglow/glow"
	"gglow/iohandler"

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
	title  binding.String
	effect iohandler.EffectIoHandler

	applyButton *widget.Button
}

func NewEffectDialog(effect iohandler.EffectIoHandler, window fyne.Window) (ef *EffectDialog) {
	ef = &EffectDialog{
		effect: effect,
		title:  binding.NewString(),
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
	ef.CustomDialog.Hide()
	title, _ := ef.title.Get()
	frame := glow.NewFrame()
	frame.Interval = uint32(RateBounds.OnVal)
	err := ef.effect.CreateNewEffect(title, frame)
	if err != nil {
		fyne.LogError(title, err)
	}

	//refresh the current folder
	ef.effect.LoadFolder(ef.effect.FolderName())

	err = ef.effect.LoadEffect(title)
	if err != nil {
		fyne.LogError(title, err)
	}
}

func (ef *EffectDialog) validateFileName(s string) error {
	err := ef.effect.ValidateNewEffectName(s)
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
