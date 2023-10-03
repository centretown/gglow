package ui

import (
	"glow-gui/data"
	"glow-gui/glow"
	"glow-gui/res"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type LayerForm struct {
	*container.AppTabs
	model      *data.Model
	changeView func()
}

func NewLayerForm(model *data.Model, changeView func()) *LayerForm {
	lf := &LayerForm{
		AppTabs:    container.NewAppTabs(),
		model:      model,
		changeView: changeView,
	}

	createLabel :=
		func(labelID res.LabelID, iconID res.AppIconID) fyne.CanvasObject {
			icon := widget.NewIcon(res.AppIconResource(iconID))
			label := widget.NewLabel(labelID.String())
			hbox := container.NewHBox(icon, label)
			return hbox
		}

	var (
		label, data fyne.CanvasObject
		frm         *fyne.Container
		tab         *container.TabItem
	)

	frm = container.New(layout.NewFormLayout())

	layer := &glow.Layer{}

	// first tab
	label = createLabel(res.HueShiftLabel, res.HueShiftIcon)
	data = widget.NewLabel(strconv.Itoa(int(layer.HueShift)))
	frm.Objects = append(frm.Objects, label, data)
	tab = container.NewTabItem("Hue", frm)
	lf.AppTabs.Append(tab)

	// second
	frm = container.New(layout.NewFormLayout())
	label = createLabel(res.ScanLabel, res.ScanIcon)
	data = widget.NewLabel(strconv.Itoa(int(layer.Scan)))
	frm.Objects = append(frm.Objects, label, data)
	tab = container.NewTabItem("Scan", frm)
	lf.AppTabs.Append(tab)

	// label = createLabel(res.BeginLabel, res.BeginIcon)
	// data = widget.NewLabel(strconv.Itoa(int(layer.Begin)))
	// frm.Objects = append(frm.Objects, label, data)

	// label = createLabel(res.EndLabel, res.EndIcon)
	// data = widget.NewLabel(strconv.Itoa(int(layer.End)))
	// frm.Objects = append(frm.Objects, label, data)
	// grid tab
	frm = container.New(layout.NewFormLayout())
	label = createLabel(res.OriginLabel, res.HueShiftIcon)
	sel := widget.NewSelect(res.OriginLabels, func(s string) {})
	sel.SetSelectedIndex(int(layer.Grid.Origin))
	frm.Objects = append(frm.Objects, label, sel)

	label = createLabel(res.OrientationLabel, res.HueShiftIcon)
	sel = widget.NewSelect(res.OrientationLabels, func(s string) {})
	frm.Objects = append(frm.Objects, label, sel)
	sel.SetSelectedIndex(int(layer.Grid.Orientation))
	tab = container.NewTabItem(res.GridLabel.String(), frm)
	lf.AppTabs.Append(tab)

	// colors tab
	frm = container.New(layout.NewFormLayout())
	tab = container.NewTabItem(res.ChromaLabel.String(), frm)
	lf.AppTabs.Append(tab)

	return lf
}
