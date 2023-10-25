package ui

import (
	"fmt"
	"glow-gui/resources"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type LightStripLayout struct {
	*dialog.CustomDialog
	preferences fyne.Preferences
	columns     binding.Float
	rows        binding.Float
}

func NewLightStripLayout(parent fyne.Window, p fyne.Preferences,
	sourceStrip binding.Untyped) *LightStripLayout {
	fe := &LightStripLayout{
		preferences: p,
		columns:     binding.NewFloat(),
		rows:        binding.NewFloat(),
	}

	fe.columns.Set(p.Float(resources.StripColumns.String()))
	columnsEntry := NewRangeEntry(fe.columns, &EntryBounds{1, 50, 0, 0})
	colItem := widget.NewFormItem(resources.ColumnsLabel.String(), columnsEntry)
	// colItem := createRangeEntryItem(resources.ColumnsLabel.String(), fe.columns, 1, 50)

	fe.rows.Set(p.Float(resources.StripRows.String()))
	rowsEntry := NewRangeEntry(fe.rows, &EntryBounds{1, 20, 0, 0})
	rowsItem := widget.NewFormItem(resources.RowsLabel.String(), rowsEntry)

	frm := widget.NewForm(colItem, rowsItem)
	fe.CustomDialog = dialog.NewCustomWithoutButtons("form", frm, parent)

	confirm := widget.NewButton(resources.ConfirmLabel.String(), func() {
		if err := frm.Validate(); err != nil {
			fmt.Println("confirm", err)
			return
		}
		fe.Hide()
		// update new values
		stripColumns, _ := fe.columns.Get()
		fe.preferences.SetFloat(resources.StripColumns.String(), stripColumns)
		stripRows, _ := fe.rows.Get()
		fe.preferences.SetFloat(resources.StripRows.String(), stripRows)
		strip := NewLightStrip(stripColumns*stripRows, stripRows)
		sourceStrip.Set(strip)
	})

	cancel := widget.NewButton(resources.CancelLabel.String(), func() {
		fe.Hide()
		// restore old values
		columns := fe.preferences.Float(resources.StripColumns.String())
		fe.columns.Set(columns)
		columnsEntry.SetText(fmt.Sprintf("%.0f", columns))

		rows := fe.preferences.Float(resources.StripRows.String())
		fe.rows.Set(rows)
		rowsEntry.SetText(fmt.Sprintf("%.0f", rows))
	})

	fe.CustomDialog.SetButtons([]fyne.CanvasObject{cancel, confirm})
	return fe
}
