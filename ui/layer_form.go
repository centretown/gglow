package ui

import (
	"glow-gui/glow"
	"glow-gui/res"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type LayerForm struct {
	*dialog.CustomDialog
	window fyne.Window
	layer  *glow.Layer
}

func (lf *LayerForm) createForm() fyne.CanvasObject {

	createLabel := func(labelID res.LabelID, iconID res.AppIconID) fyne.CanvasObject {
		icon := widget.NewIcon(res.AppIconResource(iconID))
		label := widget.NewLabel(labelID.String())
		hbox := container.NewHBox(icon, label)
		return hbox
	}

	var (
		label, data fyne.CanvasObject
		frm         *fyne.Container
		tabs        *container.AppTabs = container.NewAppTabs()
		tab         *container.TabItem
	)

	frm = container.New(layout.NewFormLayout())

	// first tab
	label = createLabel(res.HueShiftLabel, res.HueShiftIcon)
	data = widget.NewLabel(strconv.Itoa(int(lf.layer.HueShift)))
	frm.Objects = append(frm.Objects, label, data)

	label = createLabel(res.ScanLabel, res.ScanIcon)
	data = widget.NewLabel(strconv.Itoa(int(lf.layer.Scan)))
	frm.Objects = append(frm.Objects, label, data)

	label = createLabel(res.BeginLabel, res.BeginIcon)
	data = widget.NewLabel(strconv.Itoa(int(lf.layer.Begin)))
	frm.Objects = append(frm.Objects, label, data)

	label = createLabel(res.EndLabel, res.EndIcon)
	data = widget.NewLabel(strconv.Itoa(int(lf.layer.End)))
	frm.Objects = append(frm.Objects, label, data)
	tab = container.NewTabItem("Layout", frm)
	tabs.Append(tab)

	// grid tab
	frm = container.New(layout.NewFormLayout())
	label = createLabel(res.OriginLabel, res.HueShiftIcon)
	sel := widget.NewSelect(res.OriginLabels, func(s string) {})
	sel.SetSelectedIndex(int(lf.layer.Grid.Origin))
	frm.Objects = append(frm.Objects, label, sel)

	label = createLabel(res.OrientationLabel, res.HueShiftIcon)
	sel = widget.NewSelect(res.OrientationLabels, func(s string) {})
	frm.Objects = append(frm.Objects, label, sel)
	sel.SetSelectedIndex(int(lf.layer.Grid.Orientation))
	tab = container.NewTabItem(res.GridLabel.String(), frm)
	tabs.Append(tab)

	// colors tab
	frm = container.New(layout.NewFormLayout())
	tab = container.NewTabItem(res.ChromaLabel.String(), frm)
	tabs.Append(tab)

	return tabs
}

func NewLayerForm(window fyne.Window, layer *glow.Layer) *dialog.CustomDialog {
	lf := &LayerForm{
		window: window,
		layer:  layer,
	}

	title := window.Title()
	dismiss := "Cancel"

	lf.CustomDialog = dialog.NewCustom(
		title,
		dismiss,
		lf.createForm(),
		window)

	return lf.CustomDialog
}

func ColorDialog(window fyne.Window) *fyne.Container {
	hueLabel := widget.NewLabel(res.HueLabel.String())
	hueEntry := widget.NewEntry()
	hueEntry.SetPlaceHolder(res.HueLabel.PlaceHolder())

	label2 := widget.NewLabel(res.SaturationLabel.String())
	value2 := widget.NewEntry()
	value2.SetPlaceHolder(res.SaturationLabel.PlaceHolder())

	button1 := widget.NewButton("Choose Color...", func() {
		picker := dialog.NewColorPicker("Color Picker", "color", func(c color.Color) {

		}, window)
		picker.Advanced = true
		picker.Show()
	})
	grid := container.New(layout.NewVBoxLayout(),
		hueLabel, hueEntry, label2, value2, button1)
	return grid
}
