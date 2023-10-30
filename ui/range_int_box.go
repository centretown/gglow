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
	Bounds       *EntryBoundsInt
	Decrease     *widget.Button
	Increase     *widget.Button
	Entry        *RangeEntryInt
	Field        binding.Int
	initialValue int
}

func NewRangeIntBox(field binding.Int, bounds *EntryBoundsInt) *RangeIntBox {

	rb := &RangeIntBox{
		Bounds: bounds,
		Entry:  NewRangeEntryInt(field, bounds),
		Field:  field,
	}

	rb.Decrease = widget.NewButtonWithIcon("", theme.MoveDownIcon(), rb.buttonCheck(-1))
	rb.Increase = widget.NewButtonWithIcon("", theme.MoveUpIcon(), rb.buttonCheck(1))
	rb.Container = container.NewHBox(rb.Decrease, rb.Entry, rb.Increase)
	return rb
}

func (rb *RangeIntBox) buttonCheck(inc int) func() {
	return func() {
		f, _ := rb.Field.Get()
		f += inc
		if f >= rb.Bounds.MinVal && f <= rb.Bounds.MaxVal {
			rb.Field.Set(f)
		}
	}
}

func (rb *RangeIntBox) Enable(b bool) {
	if b {
		rb.Decrease.Enable()
		rb.Increase.Enable()
		rb.Entry.Enable()
		return
	}
	rb.Decrease.Disable()
	rb.Increase.Disable()
	rb.Entry.Disable()
}
