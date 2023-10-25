package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type DisabledEntry struct {
	Label    *widget.Label
	Decrease *widget.Button
	Increase *widget.Button
	Entry    *RangeEntry
	InBox    *fyne.Container
}

func NewDisabledEntry(labelText string, enabled binding.Bool,
	field binding.Float, bounds *EntryBounds) (le *DisabledEntry) {

	le = &DisabledEntry{}
	le.Label = widget.NewLabel(labelText)
	le.Entry = NewRangeEntry(field, bounds)

	buttonCheck := func(val float64) func() {
		return func() {
			f, _ := field.Get()
			f += val
			if f >= bounds.MinVal && f <= bounds.MaxVal {
				field.Set(f)
			}
		}
	}

	le.Decrease = widget.NewButtonWithIcon("", theme.MoveDownIcon(), buttonCheck(-1))
	le.Increase = widget.NewButtonWithIcon("", theme.MoveUpIcon(), buttonCheck(1))
	le.InBox = container.NewHBox(le.Decrease, le.Entry, le.Increase)

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
	le.Decrease.Disable()
	le.Increase.Disable()
	le.Entry.Disable()
}

func (le *DisabledEntry) Enable() {
	le.Decrease.Enable()
	le.Increase.Enable()
	le.Entry.Enable()
}
