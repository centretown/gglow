package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

type GlobalShortCut struct {
	Shortcut *desktop.CustomShortcut
	Action   func()
}

var CtrlS = &desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl}
var CtrlN = &desktop.CustomShortcut{KeyName: fyne.KeyN, Modifier: fyne.KeyModifierControl}
var CtrlL = &desktop.CustomShortcut{KeyName: fyne.KeyL, Modifier: fyne.KeyModifierControl}
var ShiftCtrlN = &desktop.CustomShortcut{KeyName: fyne.KeyN, Modifier: fyne.KeyModifierShift | fyne.KeyModifierControl}
var CtrlQ = &desktop.CustomShortcut{KeyName: fyne.KeyQ, Modifier: fyne.KeyModifierControl}
var CtrlE = &desktop.CustomShortcut{KeyName: fyne.KeyE, Modifier: fyne.KeyModifierControl}
var CtrlM = &desktop.CustomShortcut{KeyName: fyne.KeyM, Modifier: fyne.KeyModifierControl}

func AddGlobalShortCut(window fyne.Window, cut *GlobalShortCut) {
	window.Canvas().AddShortcut(cut.Shortcut, func(shortcut fyne.Shortcut) {
		cut.Action()
	})
}
