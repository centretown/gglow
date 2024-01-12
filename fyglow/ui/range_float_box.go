package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type RangeFloatBox struct {
	*fyne.Container
	Decrease *widget.Button
	Increase *widget.Button
	Entry    *RangeEntryFloat
}

func NewRangeFloatBox(field binding.Float, bounds *FloatEntryBounds) *RangeFloatBox {
	buttonCheck := func(val float64) func() {
		return IncrementFloat(val, field, bounds)
	}

	rb := &RangeFloatBox{
		Decrease: widget.NewButtonWithIcon("", theme.MoveDownIcon(), buttonCheck(-1)),
		Entry:    NewRangeEntryFloat(field, bounds),
		Increase: widget.NewButtonWithIcon("", theme.MoveUpIcon(), buttonCheck(1)),
	}

	rb.Container = container.NewHBox(rb.Decrease, rb.Entry, rb.Increase)
	return rb
}

func (rb *RangeFloatBox) Disable() {
	rb.Decrease.Disable()
	rb.Increase.Disable()
	rb.Entry.Disable()
}

func (rb *RangeFloatBox) Enable() {
	rb.Decrease.Enable()
	rb.Increase.Enable()
	rb.Entry.Enable()
}
