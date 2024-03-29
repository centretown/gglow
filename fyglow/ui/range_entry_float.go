package ui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/widget"
)

type RangeEntryFloat struct {
	widget.Entry
}

func NewRangeEntryFloat(field binding.Float, bounds *FloatEntryBounds) *RangeEntryFloat {
	e := &RangeEntryFloat{}
	binder := binding.FloatToStringWithFormat(field, "%4.0f")
	e.Bind(binder)

	validateRange := func(s string) error {
		val, _ := strconv.ParseFloat(s, 64)
		if val >= bounds.MinVal && val <= bounds.MaxVal {
			return nil
		}
		return fmt.Errorf("%.0f not in range %.0f-%.0f",
			val, bounds.MinVal, bounds.MaxVal)
	}

	e.Validator = validation.NewAllStrings(validateRange)
	e.ExtendBaseWidget(e)
	return e
}

func (e *RangeEntryFloat) TypedRune(r rune) {
	if (r >= '0' && r <= '9') || r == '.' || r == ',' {
		e.Entry.TypedRune(r)
	}
}

func (e *RangeEntryFloat) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	content := paste.Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err == nil {
		e.Entry.TypedShortcut(shortcut)
	}
}

func (e *RangeEntryFloat) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}
