package ui

import (
	"glow-gui/data"
	"glow-gui/resources"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type LayerForm struct {
	*container.AppTabs
	model     *data.Model
	isDynamic binding.Bool
	isScanner binding.Bool
}

func NewLayerForm(model *data.Model) *LayerForm {
	lf := &LayerForm{
		AppTabs:   container.NewAppTabs(),
		model:     model,
		isDynamic: binding.NewBool(),
		isScanner: binding.NewBool(),
	}

	// hue shift tab
	frm := lf.createHueTab()
	tab := container.NewTabItem(resources.HueLabel.String(), frm)
	lf.AppTabs.Append(tab)

	// scan tab
	frm = lf.createScanTab()
	tab = container.NewTabItem(resources.ScanLabel.String(), frm)
	lf.AppTabs.Append(tab)

	// grid tab
	frm = lf.createGridTab()
	tab = container.NewTabItem(resources.GridLabel.String(), frm)
	lf.AppTabs.Append(tab)

	// colors tab
	frm = container.New(layout.NewFormLayout())
	tab = container.NewTabItem(resources.ChromaLabel.String(), frm)
	lf.AppTabs.Append(tab)

	lf.model.Layer.AddListener(binding.NewDataListener(func() {
		f, _ := lf.model.Fields.HueShift.Get()
		lf.isDynamic.Set(f != 0)
		f, _ = lf.model.Fields.Scan.Get()
		lf.isScanner.Set(f > 0)
	}))

	return lf
}

func (lf *LayerForm) createLabel(labelID resources.LabelID, iconID resources.AppIconID) fyne.CanvasObject {
	icon := container.NewPadded(widget.NewIcon(resources.AppIconResource(iconID)))
	label := widget.NewLabel(labelID.String())
	hbox := container.NewHBox(icon, label)
	return hbox
}

func (lf *LayerForm) createCheckSlide(field binding.Float,
	label fyne.CanvasObject, checkLabelID resources.LabelID,
	isChecked binding.Bool) *fyne.Container {

	frm := container.New(layout.NewFormLayout())
	slider := widget.NewSliderWithData(1, 10, field)
	dataLabel := widget.NewLabelWithData(
		binding.FloatToStringWithFormat(field, "%.0f"))
	box := container.NewBorder(nil, nil,
		dataLabel, nil, slider)

	box.Hide()
	label.Hide()
	checkLabel := widget.NewLabel(checkLabelID.String())

	isChecked.AddListener(binding.NewDataListener(func() {
		b, _ := isChecked.Get()
		if b {
			label.Show()
			box.Show()
		} else {
			label.Hide()
			box.Hide()
		}
	}))

	check := widget.NewCheckWithData("", isChecked)
	frm.Objects = append(frm.Objects, checkLabel, check, label, box)
	return frm
}

func (lf *LayerForm) createHueTab() *fyne.Container {
	label := lf.createLabel(resources.HueShiftLabel, resources.HueShiftIcon)
	frm := lf.createCheckSlide(lf.model.Fields.HueShift, label, resources.DynamicLabel, lf.isDynamic)
	return frm
}

func (lf *LayerForm) createScanTab() *fyne.Container {
	label := lf.createLabel(resources.ScanLengthLabel, resources.ScanIcon)
	frm := lf.createCheckSlide(lf.model.Fields.Scan, label, resources.ScannerLabel, lf.isScanner)
	return frm
}

func (lf *LayerForm) createGridTab() *fyne.Container {
	frm := container.New(layout.NewFormLayout())
	labelOrigin := lf.createLabel(resources.OriginLabel, resources.HueShiftIcon)
	selectOrigin := widget.NewSelect(resources.OriginLabels, func(s string) {})
	selectOrigin.OnChanged = func(s string) {
		lf.model.Fields.Origin.Set(selectOrigin.SelectedIndex())
	}
	lf.model.Fields.Origin.AddListener(binding.NewDataListener(func() {
		index, _ := lf.model.Fields.Origin.Get()
		selectOrigin.SetSelectedIndex(index)
	}))

	frm.Objects = append(frm.Objects, labelOrigin, selectOrigin)

	labelOrientation := lf.createLabel(resources.OrientationLabel, resources.HueShiftIcon)
	selectOrientation := widget.NewSelect(resources.OrientationLabels, func(s string) {})
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
