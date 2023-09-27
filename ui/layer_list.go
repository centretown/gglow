package ui

import (
	"fmt"
	"glow-gui/glow"
	"glow-gui/res"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type LayerItem int

const (
	LayerIcon LayerItem = iota
	HueItem
	ScanItem
	BeginItem
	EndItem
	LAYER_ITEM_COUNT
)

type LayerList struct {
	*widget.List
	frame        *glow.Frame
	numberFormat string
}

func (ll *LayerList) createLayerItem() fyne.CanvasObject {
	s := fmt.Sprintf("%4d", 0)

	createItem := func(labelID res.LabelID, iconID res.AppIconID, fn func()) *fyne.Container {
		head := widget.NewLabel(labelID.String())
		head.Alignment = fyne.TextAlignCenter
		btn := widget.NewButtonWithIcon("", res.AppIconResource(iconID), fn)
		data := widget.NewLabel(s)
		hbox := container.NewHBox(btn, data)
		vbox := container.NewVBox(head, hbox)
		return vbox
	}

	ctr := container.NewHBox(
		res.NewAppIcon(res.LayerIcon),
		createItem(res.HueShiftLabel, res.HueShiftIcon, func() {}),
		createItem(res.ScanLabel, res.ScanIcon, func() {}),
		createItem(res.BeginLabel, res.BeginIcon, func() {}),
		createItem(res.EndLabel, res.EndIcon, func() {}),
	)
	// widget.NewButtonWithIcon(res.HueShiftLabel.String(),
	// 	res.AppIconResource(res.HueShiftIcon), func() {}),
	// widget.NewLabel(s),

	// widget.NewButtonWithIcon(res.ScanLabel.String(),
	// 	res.AppIconResource(res.ScanIcon), func() {}),
	// widget.NewLabel(s),

	// widget.NewButtonWithIcon(res.BeginLabel.String(),
	// 	res.AppIconResource(res.BeginIcon), func() {}),
	// widget.NewLabel(s),

	// widget.NewButtonWithIcon(res.EndLabel.String(),
	// 	res.AppIconResource(res.EndIcon), func() {}),
	// widget.NewLabel(s))
	return ctr
}

func (ll *LayerList) updateLayerItem(id int, item fyne.CanvasObject) {
	layer := ll.frame.Layers[id]
	cntr := item.(*fyne.Container)

	setItem := func(id LayerItem, v interface{}) {
		vbox := cntr.Objects[id].(*fyne.Container)
		hbox := vbox.Objects[1].(*fyne.Container)
		label := hbox.Objects[1].(*widget.Label)
		label.SetText(fmt.Sprintf(ll.numberFormat, v))
	}
	setItem(HueItem, layer.HueShift)
	setItem(ScanItem, layer.Scan)
	setItem(BeginItem, layer.Begin)
	setItem(EndItem, layer.End)
}

func NewLayerList(frame *glow.Frame) *widget.List {
	ll := &LayerList{
		frame:        frame,
		numberFormat: "%3d",
	}

	ll.List = widget.NewList(
		func() int {
			return len(frame.Layers)
		},
		ll.createLayerItem,
		ll.updateLayerItem)

	return ll.List
}
