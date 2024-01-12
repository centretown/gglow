package fyui

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

func AddGlobalShortCut(window fyne.Window, sc *GlobalShortCut) {
	window.Canvas().AddShortcut(sc.Shortcut, func(shortcut fyne.Shortcut) {
		sc.Action()
	})
}
