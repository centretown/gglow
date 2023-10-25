package ui

import (
	"fmt"
	"glow-gui/resources"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type LightStripSpeed struct {
	*dialog.CustomDialog
	preferences fyne.Preferences
	interval    binding.Float
}

func NewLightStripSpeed(parent fyne.Window, p fyne.Preferences,
	sourceStrip binding.Untyped) *LightStripSpeed {
	fs := &LightStripSpeed{
		preferences: p,
		interval:    binding.NewFloat(),
	}

	fs.interval.Set(p.Float(resources.StripInterval.String()))
	intervalEntry := NewRangeEntry(fs.interval, &EntryBounds{16, 360, 0, 0})
	intervalItem := widget.NewFormItem(resources.IntervalLabel.String(), intervalEntry)

	frm := widget.NewForm(intervalItem)
	fs.CustomDialog = dialog.NewCustomWithoutButtons("form", frm, parent)

	confirm := widget.NewButton(resources.ConfirmLabel.String(), func() {
		if err := frm.Validate(); err != nil {
			fmt.Println("confirm", err)
			return
		}
		fs.Hide()
		// update new values
		stripInterval, _ := fs.interval.Get()
		fs.preferences.SetFloat(resources.StripInterval.String(), stripInterval)

		// strip := NewLightStrip(stripColumns*stripRows, stripRows)
		// sourceStrip.Set(strip)
	})

	cancel := widget.NewButton(resources.CancelLabel.String(), func() {
		fs.Hide()

		interval := fs.preferences.Float(resources.StripInterval.String())
		fs.interval.Set(interval)
		intervalEntry.SetText(fmt.Sprintf("%.0f", interval))
	})

	fs.CustomDialog.SetButtons([]fyne.CanvasObject{cancel, confirm})
	return fs
}
