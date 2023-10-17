package ui

import (
	"glow-gui/resources"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type LightStripSettings struct {
	*dialog.CustomDialog
	preferences fyne.Preferences
	// sourceStrip binding.Untyped
	columns  binding.Float
	rows     binding.Float
	interval binding.Float
}

func NewLightStripSettings(parent fyne.Window, p fyne.Preferences,
	sourceStrip binding.Untyped) *LightStripSettings {
	fe := &LightStripSettings{
		preferences: p,
		columns:     binding.NewFloat(),
		rows:        binding.NewFloat(),
		interval:    binding.NewFloat(),
	}

	frm := container.New(layout.NewFormLayout())
	fe.columns.Set(p.Float(resources.StripColumns.String()))
	label := widget.NewLabel(resources.ColumnsLabel.String())
	colSlide := createSlide(fe.columns, 1, 50)
	frm.Objects = append(frm.Objects, label, colSlide)

	fe.rows.Set(p.Float(resources.StripRows.String()))
	label = widget.NewLabel(resources.RowsLabel.String())
	rowSlide := createSlide(fe.rows, 1, 20)
	frm.Objects = append(frm.Objects, label, rowSlide)

	fe.interval.Set(p.Float(resources.StripInterval.String()))
	label = widget.NewLabel(resources.IntervalLabel.String())
	intervalSlide := createSlide(fe.interval, 16, 360)
	frm.Objects = append(frm.Objects, label, intervalSlide)

	fe.CustomDialog = dialog.NewCustomWithoutButtons("form", frm, parent)

	confirm := widget.NewButton(resources.ConfirmLabel.String(), func() {
		// update new values
		stripColumns, _ := fe.columns.Get()
		fe.preferences.SetFloat(resources.StripColumns.String(), stripColumns)
		stripRows, _ := fe.rows.Get()
		fe.preferences.SetFloat(resources.StripRows.String(), stripRows)
		stripInterval, _ := fe.interval.Get()
		fe.preferences.SetFloat(resources.StripInterval.String(), stripInterval)
		fe.Hide()
		strip := NewLightStrip(stripColumns*stripRows, stripRows, stripInterval)
		sourceStrip.Set(strip)
	})

	cancel := widget.NewButton(resources.CancelLabel.String(), func() {
		// restore old values
		fe.columns.Set(fe.preferences.Float(resources.StripColumns.String()))
		fe.rows.Set(fe.preferences.Float(resources.StripRows.String()))
		fe.interval.Set(fe.preferences.Float(resources.StripInterval.String()))
		fe.Hide()
	})
	fe.CustomDialog.SetButtons([]fyne.CanvasObject{cancel, confirm})
	return fe
}
