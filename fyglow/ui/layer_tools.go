package ui

import (
	"fmt"
	"gglow/fyglow/effectio"
	"gglow/fyglow/resource"
	"gglow/text"

	"fyne.io/fyne/v2"
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
		resource.IconLayerAdd(), lt.addLayer)))
	lt.Toolbar.Append(NewButtonItem(widget.NewButtonWithIcon("",
		resource.IconLayerInsert(), lt.insertLayer)))
	lt.Toolbar.Append(NewButtonItem(widget.NewButtonWithIcon("",
		resource.IconLayerRemove(), lt.removeLayer)))

	AddGlobalShortCut(window,
		&GlobalShortCut{Shortcut: CtrlL, Action: lt.addLayer})

	itemNew := &fyne.MenuItem{Label: text.NewLabel.String(),
		Icon: resource.IconLayerAdd(), Action: lt.addLayer, Shortcut: CtrlL}
	itemInsert := &fyne.MenuItem{Label: text.InsertLabel.String(),
		Icon: resource.IconLayerInsert(), Action: lt.insertLayer}
	itemRemove := &fyne.MenuItem{Label: text.RemoveLabel.String(),
		Icon: resource.IconLayerRemove(), Action: lt.removeLayer}
	itemLayer := &fyne.MenuItem{
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
