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

type RangeEntryInt struct {
	widget.Entry
}

func NewRangeEntryInt(field binding.Int, bounds *IntEntryBounds) *RangeEntryInt {
	e := &RangeEntryInt{}
	binder := binding.IntToString(field)
	e.Bind(binder)

	validateRange := func(s string) error {
		val, _ := strconv.ParseInt(s, 10, 64)
		if int(val) >= bounds.MinVal && int(val) <= bounds.MaxVal {
			return nil
		}
		return fmt.Errorf("%d not in range %d-%d",
			val, bounds.MinVal, bounds.MaxVal)
	}

	e.Validator = validation.NewAllStrings(validateRange)
	e.ExtendBaseWidget(e)
	return e
}

func (e *RangeEntryInt) TypedRune(r rune) {
	if (r >= '0' && r <= '9') || r == ',' {
		e.Entry.TypedRune(r)
	}
}

func (e *RangeEntryInt) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	content := paste.Clipboard.Content()
	if _, err := strconv.ParseInt(content, 10, 64); err == nil {
		e.Entry.TypedShortcut(shortcut)
	}
}

func (e *RangeEntryInt) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}
