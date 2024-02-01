package ui

import (
	"fmt"
	"gglow/fyglow/effectio"
	"gglow/fyglow/resource"
	"gglow/text"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type FrameTools struct {
	*widget.Toolbar
}

func NewFrameTools(effect *effectio.EffectIo, window fyne.Window,
	menu *fyne.Menu) *FrameTools {

	ft := &FrameTools{
		Toolbar: widget.NewToolbar(),
	}

	addFolder := NewFolderDialog(effect, window)
	addEffect := NewEffectDialog(effect, window)

	effectSave := func() {
		if effect.HasChanged() {
			effect.SaveEffect()
		}
	}

	effectAdd := func() {
		addEffect.Start()
	}

	folderAdd := func() {
		addFolder.Start()
	}

	saveButton := NewButtonItem(
		widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), effectSave))
	ft.Toolbar.Append(saveButton)

	effectButton := NewButtonItem(
		widget.NewButtonWithIcon("", resource.IconFrameAdd(), effectAdd))
	ft.Toolbar.Append(effectButton)

	ft.Toolbar.Append(NewButtonItem(
		widget.NewButtonWithIcon("", resource.IconFrameRemove(), func() {})))

	ft.Toolbar.Append(NewButtonItem(
		widget.NewButtonWithIcon("", theme.FolderNewIcon(), folderAdd)))

	ft.Toolbar.Append(NewButtonItem(
		widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {})))

	// effect.AddFolderListener(binding.NewDataListener(func() {
	// 	if effect.IsRootFolder() {
	// 		effectButton.Disable()
	// 	} else {
	// 		effectButton.Enable()
	// 	}
	// }))

	effect.AddChangeListener(binding.NewDataListener(func() {
		if effect.HasChanged() {
			saveButton.Enable()
		} else {
			saveButton.Disable()
		}
	}))

	AddGlobalShortCut(window,
		&GlobalShortCut{Shortcut: CtrlS, Action: effectSave})
	AddGlobalShortCut(window,
		&GlobalShortCut{Shortcut: CtrlN, Action: effectAdd})
	AddGlobalShortCut(window,
		&GlobalShortCut{Shortcut: ShiftCtrlN, Action: addFolder.Start})

	itemNewFolder := &fyne.MenuItem{Label: text.NewLabel.String(),
		Icon:   theme.FolderNewIcon(),
		Action: addFolder.Start, Shortcut: ShiftCtrlN,
	}
	itemTrash := &fyne.MenuItem{Label: text.TrashLabel.String(),
		Icon:   theme.DeleteIcon(),
		Action: func() { fmt.Println("Trash Folder") },
	}
	itemFolders := &fyne.MenuItem{
		Icon:      theme.FolderIcon(),
		Label:     text.FolderLabel.String(),
		ChildMenu: &fyne.Menu{Label: "", Items: []*fyne.MenuItem{itemNewFolder, itemTrash}},
	}
	menu.Items = append(menu.Items, &fyne.MenuItem{IsSeparator: true},
		itemFolders)

	itemSave := &fyne.MenuItem{Label: text.SaveLabel.String(),
		Icon:     theme.DocumentSaveIcon(),
		Action:   effectSave,
		Shortcut: CtrlS,
	}
	itemNew := &fyne.MenuItem{Label: text.NewLabel.String(),
		Icon:     resource.IconFrameAdd(),
		Action:   effectAdd,
		Shortcut: CtrlN,
	}
	itemRemove := &fyne.MenuItem{Label: text.RemoveLabel.String(),
		Icon:   resource.IconFrameRemove(),
		Action: func() { fmt.Println("Remove Effect") },
	}
	itemEffects := &fyne.MenuItem{
		Icon:  resource.IconEffect(),
		Label: text.EffectsLabel.String(),
		ChildMenu: &fyne.Menu{Label: "",
			Items: []*fyne.MenuItem{itemSave, itemNew, itemRemove}},
	}

	menu.Items = append(menu.Items,
		&fyne.MenuItem{IsSeparator: true}, itemEffects)

	return ft
}
