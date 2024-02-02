package ui

import (
	"gglow/fyglow/effectio"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func NewFrameToolbar(effect *effectio.EffectIo) (toolbar *widget.Toolbar) {
	toolbar = widget.NewToolbar()
	saveButton := NewButtonItemFromMenu(MenuEffectSave)
	toolbar.Append(saveButton)
	toolbar.Append(NewButtonItemFromMenu(MenuEffectAdd))
	toolbar.Append(NewButtonItemFromMenu(MenuEffectRemove))
	toolbar.Append(NewButtonItemFromMenu(MenuFolderAdd))
	toolbar.Append(NewButtonItemFromMenu(MenuFolderRemove))
	effect.AddChangeListener(binding.NewDataListener(func() {
		if effect.HasChanged() {
			saveButton.Enable()
		} else {
			saveButton.Disable()
		}
	}))

	// AddGlobalShortCut(window,
	// 	&GlobalShortCut{Shortcut: CtrlS,
	// 		Action: effectItems[MenuEffectSave].Action})
	// AddGlobalShortCut(window,
	// 	&GlobalShortCut{Shortcut: CtrlN, Action: effectAdd})
	// AddGlobalShortCut(window,
	// 	&GlobalShortCut{Shortcut: ShiftCtrlN, Action: addFolder.Start})

	return toolbar
}
