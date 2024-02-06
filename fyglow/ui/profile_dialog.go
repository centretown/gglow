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

const MaxRGBLights int = 2500

type ProfileDialog struct {
	*dialog.CustomDialog
	preferences fyne.Preferences
	sourceStrip binding.Untyped
	columns     binding.Int
	rows        binding.Int
	background  color.Color

	boundsRow *IntEntryBounds
	boundsCol *IntEntryBounds
}

func NewProfileDialog(parent fyne.Window, p fyne.Preferences,
	sourceStrip binding.Untyped, background color.Color) *ProfileDialog {
	ll := &ProfileDialog{
		preferences: p,
		sourceStrip: sourceStrip,
		columns:     binding.NewInt(),
		rows:        binding.NewInt(),
		background:  background,

		boundsRow: &IntEntryBounds{1, 50, 0, 0},
		boundsCol: &IntEntryBounds{1, 50, 0, 0},
	}

	ll.columns.Set(p.Int(settings.StripColumns.String()))
	columnsEntry := NewRangeIntBox(ll.columns, ll.boundsCol)
	colItem := widget.NewFormItem(text.ColumnsLabel.String(), columnsEntry.Container)

	ll.rows.Set(p.Int(settings.StripRows.String()))
	rowsEntry := NewRangeIntBox(ll.rows, ll.boundsRow)
	rowsItem := widget.NewFormItem(text.RowsLabel.String(), rowsEntry.Container)

	frm := widget.NewForm(colItem, rowsItem)
	ll.CustomDialog = dialog.NewCustomWithoutButtons(text.GridLayoutLabel.String(),
		frm, parent)
	confirm := widget.NewButton(text.ApplyLabel.String(), ll.confirm)
	revert := widget.NewButton(text.CancelLabel.String(), ll.revert)
	ll.CustomDialog.SetButtons([]fyne.CanvasObject{revert, confirm})

	ll.columns.AddListener(binding.NewDataListener(func() {
		c, _ := ll.columns.Get()
		rowsEntry.SetMax(MaxRGBLights / c)
	}))
	ll.rows.AddListener(binding.NewDataListener(func() {
		r, _ := ll.rows.Get()
		columnsEntry.SetMax(MaxRGBLights / r)
	}))
	return ll
}

func (ll *ProfileDialog) confirm() {
	ll.Hide()
	columns, _ := ll.columns.Get()
	rows, _ := ll.rows.Get()
	ll.preferences.SetInt(settings.StripColumns.String(), columns)
	ll.preferences.SetInt(settings.StripRows.String(), rows)
	ll.sourceStrip.Set(NewLightStrip(columns*rows, rows, ll.background))
}

func (ll *ProfileDialog) revert() {
	ll.Hide()
	columns := ll.preferences.Int(settings.StripColumns.String())
	rows := ll.preferences.Int(settings.StripRows.String())
	ll.columns.Set(columns)
	ll.rows.Set(rows)
}
