package ui

import (
	"gglow/settings"
	"gglow/text"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type LightStripLayout struct {
	*dialog.CustomDialog
	preferences fyne.Preferences
	sourceStrip binding.Untyped
	columns     binding.Int
	rows        binding.Int
	background  color.Color
}

func NewLightStripLayout(parent fyne.Window, p fyne.Preferences,
	sourceStrip binding.Untyped, background color.Color) *LightStripLayout {
	ll := &LightStripLayout{
		preferences: p,
		sourceStrip: sourceStrip,
		columns:     binding.NewInt(),
		rows:        binding.NewInt(),
		background:  background,
	}

	ll.columns.Set(p.Int(settings.StripColumns.String()))
	columnsEntry := NewRangeIntBox(ll.columns, &IntEntryBounds{1, 50, 0, 0})
	colItem := widget.NewFormItem(text.ColumnsLabel.String(), columnsEntry.Container)

	ll.rows.Set(p.Int(settings.StripRows.String()))
	rowsEntry := NewRangeIntBox(ll.rows, &IntEntryBounds{1, 20, 0, 0})
	rowsItem := widget.NewFormItem(text.RowsLabel.String(), rowsEntry.Container)

	frm := widget.NewForm(colItem, rowsItem)
	ll.CustomDialog = dialog.NewCustomWithoutButtons(text.GridLayoutLabel.String(),
		frm, parent)
	confirm := widget.NewButton(text.ApplyLabel.String(), ll.confirm)
	revert := widget.NewButton(text.CancelLabel.String(), ll.revert)
	ll.CustomDialog.SetButtons([]fyne.CanvasObject{revert, confirm})
	return ll
}

func (ll *LightStripLayout) confirm() {
	ll.Hide()
	columns, _ := ll.columns.Get()
	rows, _ := ll.rows.Get()
	ll.preferences.SetInt(settings.StripColumns.String(), columns)
	ll.preferences.SetInt(settings.StripRows.String(), rows)
	ll.sourceStrip.Set(NewLightStrip(columns*rows, rows, ll.background))
}

func (ll *LightStripLayout) revert() {
	ll.Hide()
	columns := ll.preferences.Int(settings.StripColumns.String())
	rows := ll.preferences.Int(settings.StripRows.String())
	ll.columns.Set(columns)
	ll.rows.Set(rows)
}
