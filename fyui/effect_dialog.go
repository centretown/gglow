package fyui

import (
	"gglow/fyio"
	"gglow/glow"
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

type EffectDialog struct {
	*dialog.CustomDialog
	title  binding.String
	effect *fyio.EffectIo

	applyButton *widget.Button
}

func NewEffectDialog(effect *fyio.EffectIo, window fyne.Window) (ef *EffectDialog) {
	ef = &EffectDialog{
		effect: effect,
		title:  binding.NewString(),
	}

	nameLabel := widget.NewLabel(resources.TitleLabel.String())
	nameEntry := widget.NewEntryWithData(ef.title)
	nameEntry.Validator = validation.NewAllStrings(ef.validateFileName)

	sep := widget.NewSeparator()
	frm := container.New(layout.NewFormLayout(),
		nameLabel, nameEntry, sep, sep)

	ef.CustomDialog = dialog.NewCustomWithoutButtons(resources.AddEffectLabel.String(),
		frm, window)
	ef.applyButton = widget.NewButtonWithIcon(resources.ApplyLabel.String(),
		theme.ConfirmIcon(), ef.apply)
	ef.applyButton.Disable()

	revertButton := NewButtonItem(
		widget.NewButtonWithIcon(resources.CancelLabel.String(),
			theme.CancelIcon(), func() {
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
	err := ef.effect.AddEffect(title, frame)
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
