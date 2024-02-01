package ui

import (
	"fmt"
	"gglow/fyglow/effectio"
	"gglow/fyglow/resource"
	"gglow/text"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type LayerTools struct {
	*widget.Toolbar
	effect *effectio.EffectIo
}

func NewLayerTools(effect *effectio.EffectIo, window fyne.Window, menu *fyne.Menu) *LayerTools {
	lt := &LayerTools{
		Toolbar: widget.NewToolbar(),
		effect:  effect,
	}

	lt.Toolbar.Append(NewButtonItem(widget.NewButtonWithIcon("",
		theme.ContentAddIcon(), lt.addLayer)))
	lt.Toolbar.Append(NewButtonItem(widget.NewButtonWithIcon("",
		theme.MoreVerticalIcon(), lt.insertLayer)))
	lt.Toolbar.Append(NewButtonItem(widget.NewButtonWithIcon("",
		theme.ContentRemoveIcon(), lt.removeLayer)))

	AddGlobalShortCut(window,
		&GlobalShortCut{Shortcut: CtrlL, Action: lt.addLayer})

	itemNew := &fyne.MenuItem{Label: text.NewLabel.String(),
		Icon: theme.ContentAddIcon(), Action: lt.addLayer, Shortcut: CtrlL}
	itemInsert := &fyne.MenuItem{Label: text.InsertLabel.String(),
		Icon: theme.MoreVerticalIcon(), Action: lt.insertLayer}
	itemRemove := &fyne.MenuItem{Label: text.RemoveLabel.String(),
		Icon: theme.ContentRemoveIcon(), Action: lt.removeLayer}
	itemLayer := &fyne.MenuItem{
		Icon:      resource.IconLayer(),
		Label:     text.LayersLabel.String(),
		ChildMenu: &fyne.Menu{Label: "", Items: []*fyne.MenuItem{itemNew, itemInsert, itemRemove}}}
	menu.Items = append(menu.Items, &fyne.MenuItem{IsSeparator: true}, itemLayer)
	return lt
}

func (lt *LayerTools) insertLayer() {
	lt.effect.InsertLayer()
	fmt.Println("insert")
}

func (lt *LayerTools) addLayer() {
	lt.effect.AddLayer()
}

func (lt *LayerTools) removeLayer() {
	lt.effect.RemoveLayer()
}
