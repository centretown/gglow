package ui

import (
	"glow-gui/resources"

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
}

func NewLightStripLayout(parent fyne.Window, p fyne.Preferences,
	sourceStrip binding.Untyped) *LightStripLayout {
	ll := &LightStripLayout{
		preferences: p,
		sourceStrip: sourceStrip,
		columns:     binding.NewInt(),
		rows:        binding.NewInt(),
	}

	ll.columns.Set(p.Int(resources.StripColumns.String()))
	columnsEntry := NewRangeIntBox(ll.columns, &IntEntryBounds{1, 50, 0, 0})
	colItem := widget.NewFormItem(resources.ColumnsLabel.String(), columnsEntry.Container)

	ll.rows.Set(p.Int(resources.StripRows.String()))
	rowsEntry := NewRangeIntBox(ll.rows, &IntEntryBounds{1, 20, 0, 0})
	rowsItem := widget.NewFormItem(resources.RowsLabel.String(), rowsEntry.Container)

	frm := widget.NewForm(colItem, rowsItem)
	ll.CustomDialog = dialog.NewCustomWithoutButtons("form", frm, parent)
	confirm := widget.NewButton(resources.ApplyLabel.String(), ll.confirm)
	revert := widget.NewButton(resources.RevertLabel.String(), ll.revert)
	ll.CustomDialog.SetButtons([]fyne.CanvasObject{revert, confirm})
	return ll
}

func (ll *LightStripLayout) confirm() {
	ll.Hide()
	columns, _ := ll.columns.Get()
	rows, _ := ll.rows.Get()
	ll.preferences.SetInt(resources.StripColumns.String(), columns)
	ll.preferences.SetInt(resources.StripRows.String(), rows)
	ll.sourceStrip.Set(NewLightStrip(columns*rows, rows))
}

func (ll *LightStripLayout) revert() {
	ll.Hide()
	columns := ll.preferences.Int(resources.StripColumns.String())
	rows := ll.preferences.Int(resources.StripRows.String())
	ll.columns.Set(columns)
	ll.rows.Set(rows)
}
