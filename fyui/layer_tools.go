package fyui

import (
	"fmt"
	"gglow/fyio"
	"gglow/resources"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type LayerTools struct {
	*widget.Toolbar
}

func NewLayerTools(effect *fyio.EffectIo, window fyne.Window, menu *fyne.Menu) *LayerTools {
	lt := &LayerTools{
		Toolbar: widget.NewToolbar(),
	}

	applyNew := func() { fmt.Println("Apply New Layer") }
	applyRemove := func() { fmt.Println("Remove Layer") }

	lt.Toolbar.Append(NewButtonItem(
		widget.NewButtonWithIcon("", theme.ContentAddIcon(), applyNew)))
	lt.Toolbar.Append(NewButtonItem(
		widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), applyRemove)))

	AddGlobalShortCut(window,
		&GlobalShortCut{Shortcut: CtrlL, Action: applyNew})

	itemNew := &fyne.MenuItem{Label: resources.NewLabel.String(),
		Icon:     theme.ContentAddIcon(),
		Action:   applyNew,
		Shortcut: CtrlL}
	itemRemove := &fyne.MenuItem{Label: resources.RemoveLabel.String(),
		Icon:   theme.ContentRemoveIcon(),
		Action: applyRemove}
	itemLayer := &fyne.MenuItem{
		Label:     resources.LayersLabel.String(),
		ChildMenu: &fyne.Menu{Label: "", Items: []*fyne.MenuItem{itemNew, itemRemove}}}
	menu.Items = append(menu.Items, &fyne.MenuItem{IsSeparator: true}, itemLayer)
	return lt
}
