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

	// hue shift tab
	frm := lf.createHueTab()
	tab := container.NewTabItem(res.HueLabel.String(), frm)
	lf.AppTabs.Append(tab)

	// scan tab
	frm = lf.createScanTab()
	tab = container.NewTabItem(res.ScanLabel.String(), frm)
	lf.AppTabs.Append(tab)

	// grid tab
	frm = lf.createGridTab()
	tab = container.NewTabItem(res.GridLabel.String(), frm)
	lf.AppTabs.Append(tab)

	// colors tab
	frm = container.New(layout.NewFormLayout())
	tab = container.NewTabItem(res.ChromaLabel.String(), frm)
	lf.AppTabs.Append(tab)

	return lf
}

func (lf *LayerForm) createLabel(labelID res.LabelID, iconID res.AppIconID) fyne.CanvasObject {
	icon := container.NewPadded(widget.NewIcon(res.AppIconResource(iconID)))
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
	box := container.NewBorder(nil, nil,
		dataLabel, nil, slider)
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
		check.SetChecked(f != 0)
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
	lf.model.Fields.Scan.AddListener(binding.NewDataListener(func() {
		f, _ := lf.model.Fields.Scan.Get()
		check.SetChecked((f > 0))
	}))
	frm.Objects = append(frm.Objects, checkLabel, check, label, box)
	return frm
}

func (lf *LayerForm) createGridTab() *fyne.Container {
	frm := container.New(layout.NewFormLayout())
	labelOrigin := lf.createLabel(res.OriginLabel, res.HueShiftIcon)
	selectOrigin := widget.NewSelect(res.OriginLabels, func(s string) {})
	selectOrigin.OnChanged = func(s string) {
		lf.model.Fields.Origin.Set(selectOrigin.SelectedIndex())
	}
	lf.model.Fields.Origin.AddListener(binding.NewDataListener(func() {
		index, _ := lf.model.Fields.Origin.Get()
		selectOrigin.SetSelectedIndex(index)
	}))

	frm.Objects = append(frm.Objects, labelOrigin, selectOrigin)

	labelOrientation := lf.createLabel(res.OrientationLabel, res.HueShiftIcon)
	selectOrientation := widget.NewSelect(res.OrientationLabels, func(s string) {})
	selectOrientation.OnChanged = func(s string) {
		lf.model.Fields.Orientation.Set(selectOrientation.SelectedIndex())
	}
	lf.model.Fields.Orientation.AddListener(binding.NewDataListener(func() {
		index, _ := lf.model.Fields.Orientation.Get()
		selectOrientation.SetSelectedIndex(index)
	}))
	frm.Objects = append(frm.Objects, labelOrientation, selectOrientation)
	return frm
}
