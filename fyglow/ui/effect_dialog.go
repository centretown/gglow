package ui

import (
	"gglow/fyglow/effectio"
	"gglow/glow"
	"gglow/text"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/validation"
)

type EffectDialog struct {
	*SimpleDialog
}

func NewEffectDialog(effect *effectio.EffectIo, window fyne.Window) (ef *EffectDialog) {
	ef = &EffectDialog{}

	ef.SimpleDialog = NewSimpleDialog(effect, window,
		text.EffectLabel.String(), text.EffectLabel.String())

	ef.NameEntry.Validator = validation.NewAllStrings(ef.validateFileName)

	ef.Apply = func() {
		frame := glow.NewFrame()
		frame.Interval = uint32(RateBounds.OnVal)
		name, _ := ef.Name.Get()
		err := ef.effect.AddEffect(name, frame)
		if err != nil {
			fyne.LogError(name, err)
		}
	}

	return ef
}

func (ef *EffectDialog) validateFileName(s string) error {
	err := ef.effect.ValidateNewEffectName(s)
	if err != nil {
		ef.ApplyButton.Disable()
		return err
	}

	ef.Name.Set(s)
	ef.ApplyButton.Enable()
	return err
}
