package ui

import (
	"gglow/fyglow/effectio"
	"gglow/fyglow/resource"
	"gglow/text"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	MenuFolders = iota
	MenuFolderAdd
	MenuFolderRemove
	MenuEffects
	MenuEffectSave
	MenuEffectAdd
	MenuEffectRemove
	MenuLayers
	MenuLayerAdd
	MenuLayerInsert
	MenuLayerRemove
	MenuFileshare
	MenuQuit
	MENU_ITEM_COUNT
)

var MenuItems [MENU_ITEM_COUNT]*fyne.MenuItem

func BuildMenu(effect *effectio.EffectIo, window fyne.Window) (menu *fyne.Menu) {
	menu = fyne.NewMenu("")
	addFolder := NewFolderDialog(effect, window)
	addEffect := NewEffectDialog(effect, window)
	expWizard := NewExportWizard(effect, window)
	MenuItems = [MENU_ITEM_COUNT]*fyne.MenuItem{
		{
			Label:  text.FolderLabel.String(),
			Icon:   theme.FolderIcon(),
			Action: func() {},
		},
		{
			Label:    text.NewLabel.String(),
			Icon:     theme.FolderNewIcon(),
			Shortcut: ShiftCtrlN,
			Action:   addFolder.Start,
		},
		{
			Label:  text.RemoveLabel.String(),
			Icon:   theme.DeleteIcon(),
			Action: func() {},
		},

		{
			Label:  text.EffectsLabel.String(),
			Icon:   resource.IconEffect(),
			Action: func() {},
		},
		{
			Label:    text.SaveLabel.String(),
			Icon:     theme.DocumentSaveIcon(),
			Shortcut: CtrlS,
			Action:   effect.SaveEffect,
		},
		{
			Label:    text.NewLabel.String(),
			Icon:     resource.IconFrameAdd(),
			Shortcut: CtrlN,
			Action:   addEffect.Start,
		},
		{
			Label:  text.RemoveLabel.String(),
			Icon:   resource.IconFrameRemove(),
			Action: func() {},
		},

		{
			Label:  text.LayersLabel.String(),
			Icon:   resource.IconLayer(),
			Action: func() {},
		},
		{
			Label:    text.NewLabel.String(),
			Icon:     theme.ContentAddIcon(),
			Shortcut: CtrlN,
			Action:   effect.AddLayer,
		},
		{
			Label:  text.InsertLabel.String(),
			Icon:   theme.MoreVerticalIcon(),
			Action: effect.InsertLayer,
		},
		{
			Label:  text.RemoveLabel.String(),
			Icon:   theme.ContentRemoveIcon(),
			Action: effect.RemoveLayer,
		},

		{
			Label:    text.ExportLabel.String(),
			Icon:     resource.IconFileShare(),
			Shortcut: CtrlE,
			Action:   expWizard.Start,
		},
		{
			Label:    text.QuitLabel.String(),
			Icon:     resource.IconExit(),
			Shortcut: CtrlQ,
			Action:   func() { os.Exit(0) },
		},
	}

	MenuItems[MenuFolders].ChildMenu = &fyne.Menu{
		Items: []*fyne.MenuItem{
			MenuItems[MenuFolderAdd],
			MenuItems[MenuFolderRemove]},
	}

	MenuItems[MenuEffects].ChildMenu = &fyne.Menu{
		Items: []*fyne.MenuItem{
			MenuItems[MenuEffectSave],
			MenuItems[MenuEffectAdd],
			MenuItems[MenuEffectRemove]},
	}

	MenuItems[MenuLayers].ChildMenu = &fyne.Menu{
		Items: []*fyne.MenuItem{
			MenuItems[MenuLayerAdd],
			MenuItems[MenuLayerInsert],
			MenuItems[MenuLayerRemove]},
	}

	menu.Items = append(menu.Items,
		MenuItems[MenuFolders],
		&fyne.MenuItem{IsSeparator: true},
		MenuItems[MenuEffects],
		&fyne.MenuItem{IsSeparator: true},
		MenuItems[MenuLayers],
		&fyne.MenuItem{IsSeparator: true},
		MenuItems[MenuFileshare],
		&fyne.MenuItem{IsSeparator: true},
		MenuItems[MenuQuit],
	)

	for _, item := range MenuItems {
		if item.Shortcut != nil {
			AddGlobalShortCut(window, &GlobalShortCut{
				Shortcut: item.Shortcut.(*desktop.CustomShortcut),
				Action:   item.Action})
		}
	}
	return menu
}

func NewButtonItemFromMenu(id int) (bi *ButtonItem) {
	return NewButtonItem(
		widget.NewButtonWithIcon("", MenuItems[id].Icon,
			MenuItems[id].Action))
}
