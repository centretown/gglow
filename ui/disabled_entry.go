package ui

import (
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type DisabledEntry struct {
	Label    *widget.Label
	RangeBox *RangeIntBox
}

func NewDisabledEntry(labelText string, enabled binding.Bool,
	field binding.Int, bounds *EntryBoundsInt) (le *DisabledEntry) {

	le = &DisabledEntry{}
	le.Label = widget.NewLabel(labelText)
	le.RangeBox = NewRangeIntBox(field, bounds)

	enabled.AddListener(binding.NewDataListener(func() {
		b, _ := enabled.Get()
		if b {
			f, _ := field.Get()
			if f == 0 {
				field.Set(bounds.OnVal)
			}
			le.Enable()
		} else {
			le.Disable()
			field.Set(bounds.OffVal)
		}
	}))
	return
}

func (le *DisabledEntry) Disable() {
	le.RangeBox.Disable()
}

func (le *DisabledEntry) Enable() {
	le.RangeBox.Enable()
}
