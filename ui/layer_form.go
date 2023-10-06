package ui

import (
	"glow-gui/data"
	"glow-gui/res"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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

	var (
		label fyne.CanvasObject
		tab   *container.TabItem
	)

	// first tab
	frm := lf.createHueTab()
	tab = container.NewTabItem(res.HueLabel.String(), frm)
	lf.AppTabs.Append(tab)

	// second tab
	frm = lf.createScanTab()
	tab = container.NewTabItem(res.ScanLabel.String(), frm)
	lf.AppTabs.Append(tab)

	// grid tab

	frm = container.New(layout.NewFormLayout())
	label = lf.createLabel(res.OriginLabel, res.HueShiftIcon)
	sel := widget.NewSelect(res.OriginLabels, func(s string) {})
	sel.OnChanged = func(s string) {
		model.Fields.Origin.Set(sel.SelectedIndex())
	}
	model.Fields.Origin.AddListener(binding.NewDataListener(func() {
		index, _ := model.Fields.Origin.Get()
		sel.SetSelectedIndex(index)
	}))

	frm.Objects = append(frm.Objects, label, sel)

	label = lf.createLabel(res.OrientationLabel, res.HueShiftIcon)
	sel = widget.NewSelect(res.OrientationLabels, func(s string) {})
	sel.OnChanged = func(s string) {
		model.Fields.Orientation.Set(sel.SelectedIndex())
	}
	model.Fields.Orientation.AddListener(binding.NewDataListener(func() {
		index, _ := model.Fields.Orientation.Get()
		sel.SetSelectedIndex(index)
	}))

	frm.Objects = append(frm.Objects, label, sel)

	tab = container.NewTabItem(res.GridLabel.String(), frm)
	lf.AppTabs.Append(tab)

	// colors tab
	frm = container.New(layout.NewFormLayout())
	tab = container.NewTabItem(res.ChromaLabel.String(), frm)
	lf.AppTabs.Append(tab)

	return lf
}

func (lf *LayerForm) createLabel(labelID res.LabelID, iconID res.AppIconID) fyne.CanvasObject {
	icon := widget.NewIcon(res.AppIconResource(iconID))
	label := widget.NewLabel(labelID.String())
	hbox := container.NewHBox(icon, label)
	return hbox
}

func (lf *LayerForm) createHueTab() *fyne.Container {
	frm := container.New(layout.NewFormLayout())
	label := lf.createLabel(res.HueShiftLabel, res.HueShiftIcon)
	slider := widget.NewSliderWithData(-10, 10, lf.model.Fields.HueShift)
	dataLabel := widget.NewLabelWithData(
		binding.FloatToStringWithFormat(lf.model.Fields.HueShift, "%.0f"))
	box := container.NewBorder(nil, nil, dataLabel, nil, slider)
	checkLabel := widget.NewLabel(res.DynamicLabel.String())
	check := widget.NewCheck("",
		func(b bool) {
			if b {
				label.Show()
				box.Show()
			} else {
				label.Hide()
				box.Hide()
			}
		})

	lf.model.Fields.HueShift.AddListener(binding.NewDataListener(func() {
		f, _ := lf.model.Fields.HueShift.Get()
		if f == 0 {
			check.SetChecked(false)
		} else {
			check.SetChecked(true)
		}

	}))
	frm.Objects = append(frm.Objects, checkLabel, check, label, box)
	return frm
}

func (lf *LayerForm) createScanTab() *fyne.Container {
	frm := container.New(layout.NewFormLayout())
	label := lf.createLabel(res.ScanLengthLabel, res.ScanIcon)
	slider := widget.NewSliderWithData(-10, 10, lf.model.Fields.Scan)
	dataLabel := widget.NewLabelWithData(
		binding.FloatToStringWithFormat(lf.model.Fields.Scan, "%.0f"))
	box := container.NewBorder(nil, nil, dataLabel, nil, slider)
	checkLabel := widget.NewLabel(res.ScannerLabel.String())
	check := widget.NewCheck("",
		func(b bool) {
			if b {
				label.Show()
				box.Show()
			} else {
				label.Hide()
				box.Hide()
			}
		})
	f, _ := lf.model.Fields.Scan.Get()
	if f == 0 {
		check.SetChecked(false)
		label.Hide()
		box.Hide()
	}
	frm.Objects = append(frm.Objects, checkLabel, check, label, box)
	return frm
}
