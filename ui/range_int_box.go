package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type RangeIntBox struct {
	*fyne.Container
	Decrease *widget.Button
	Increase *widget.Button
	Entry    *RangeEntryInt
}

func NewRangeIntBox(field binding.Int, bounds *EntryBoundsInt) *RangeIntBox {
	buttonCheck := func(val int) func() {
		return func() {
			f, _ := field.Get()
			f += val
			if f >= bounds.MinVal && f <= bounds.MaxVal {
				field.Set(f)
			}
		}
	}

	rb := &RangeIntBox{
		Decrease: widget.NewButtonWithIcon("", theme.MoveDownIcon(), buttonCheck(-1)),
		Entry:    NewRangeEntryInt(field, bounds),
		Increase: widget.NewButtonWithIcon("", theme.MoveUpIcon(), buttonCheck(1)),
	}

	rb.Container = container.NewHBox(rb.Decrease, rb.Entry, rb.Increase)
	return rb
}

func (rb *RangeIntBox) Disable() {
	rb.Decrease.Disable()
	rb.Increase.Disable()
	rb.Entry.Disable()
}

func (rb *RangeIntBox) Enable() {
	rb.Decrease.Enable()
	rb.Increase.Enable()
	rb.Entry.Enable()
}
