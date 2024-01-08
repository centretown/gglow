package fyui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

type GlobalShortCut struct {
	Shortcut *desktop.CustomShortcut
	Apply    func()
	Enabled  func() bool
}
type DialogShortCut struct {
	Apply   func()
	Enabled func() bool
}

var CtrlS = &desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl}
var CtrlN = &desktop.CustomShortcut{KeyName: fyne.KeyN, Modifier: fyne.KeyModifierControl}
var CtrlQ = &desktop.CustomShortcut{KeyName: fyne.KeyQ, Modifier: fyne.KeyModifierControl}

var Esc = &desktop.CustomShortcut{KeyName: fyne.KeyEscape}

func AddGlobalShortCut(window fyne.Window, sc *GlobalShortCut) {
	window.Canvas().AddShortcut(sc.Shortcut, func(shortcut fyne.Shortcut) {
		if sc.Enabled() {
			fmt.Println(sc.Shortcut.KeyName, sc.Shortcut.Modifier)
			sc.Apply()
		}
	})
}

var shortCutMap map[fyne.KeyName]*DialogShortCut = make(map[fyne.KeyName]*DialogShortCut)

func ProcessDialogShortcuts(ke *fyne.KeyEvent) {
	fmt.Println("ProcessDialogShortcuts", ke.Name)
	dsc, ok := shortCutMap[ke.Name]
	if ok {
		if dsc.Enabled() {
			dsc.Apply()
		}
	}
}

func InitializeDialogShortcuts(window fyne.Window) {
	window.Canvas().SetOnTypedKey(ProcessDialogShortcuts)
}

func AddDialogShortcuts(key fyne.KeyName, dsc *DialogShortCut) {
	shortCutMap[key] = dsc
}

func ResetDialogShortcuts() {
	shortCutMap = make(map[fyne.KeyName]*DialogShortCut)
}
