package ui

import (
	"glow-gui/data"
	"glow-gui/resources"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type LayerForm struct {
	*fyne.Container
	model     *data.Model
	isDynamic binding.Bool
	isScanner binding.Bool
	window    fyne.Window
}

func NewLayerForm(model *data.Model, window fyne.Window) *LayerForm {
	lf := &LayerForm{
		Container: container.NewHBox(),
		model:     model,
		isDynamic: binding.NewBool(),
		isScanner: binding.NewBool(),
		window:    window,
	}

	// theme.
	showForm := func(btn *widget.Button, popup *widget.PopUp) func() {
		f := func() {
			popup.Resize(fyne.Size{Width: minStripWidth - 2*theme.Padding(),
				Height: minStripHeight})
			offset := fyne.NewDelta(theme.Padding(), btn.Size().Height+2*theme.Padding())
			popup.Move(lf.Container.Position().Add(offset))
			popup.Show()
		}
		return f
	}

	doNothingOnTapped := func() {}

	huePopup := lf.newHueDropDown()
	hueButton := widget.NewButton(resources.HueLabel.String(), doNothingOnTapped)
	hueButton.OnTapped = showForm(hueButton, huePopup)

	scanForm := lf.newScanDropDown()
	scanButton := widget.NewButton(resources.ScanLabel.String(), doNothingOnTapped)
	scanButton.OnTapped = showForm(scanButton, scanForm)

	gridForm := lf.newGridDropDown()
	gridButton := widget.NewButton(resources.GridLabel.String(), doNothingOnTapped)
	gridButton.OnTapped = showForm(gridButton, gridForm)

	// colorsForm := container.New(layout.NewFormLayout())
	colorsButton := widget.NewButton(resources.ChromaLabel.String(), doNothingOnTapped)
	// colorsButton.OnTapped = showForm(colorsButton, colorsForm)

	lf.Container.Objects = []fyne.CanvasObject{gridButton, scanButton, colorsButton, hueButton}

	lf.model.Layer.AddListener(binding.NewDataListener(func() {
		f, _ := lf.model.Fields.HueShift.Get()
		lf.isDynamic.Set(f != 0)
		f, _ = lf.model.Fields.Scan.Get()
		lf.isScanner.Set(f > 0)
	}))

	return lf
}

func (lf *LayerForm) newHueDropDown() *widget.PopUp {
	checkLabel, check := NewLabelCheck(resources.DynamicLabel.String(), lf.isDynamic)
	bounds := EntryBounds{MinVal: -10, MaxVal: 10, OnVal: 1, OffVal: 0}
	le := NewDisabledEntry(resources.HueShiftLabel.String(), lf.isDynamic,
		lf.model.Fields.HueShift, &bounds)
	return lf.makePopup(checkLabel, check, le.Label, le.InBox)
}

func (lf *LayerForm) newScanDropDown() *widget.PopUp {
	checkLabel, check := NewLabelCheck(resources.ScanLabel.String(), lf.isScanner)
	bounds := EntryBounds{MinVal: 0, MaxVal: 10, OnVal: 1, OffVal: 0}
	le := NewDisabledEntry(resources.LengthLabel.String(), lf.isScanner,
		lf.model.Fields.Scan, &bounds)
	return lf.makePopup(checkLabel, check, le.Label, le.InBox)
}

func (lf *LayerForm) newGridDropDown() *widget.PopUp {

	labelOrigin := widget.NewLabel(resources.OriginLabel.String())
	selectOrigin := widget.NewSelect(resources.OriginLabels, func(s string) {})
	selectOrigin.OnChanged = func(s string) {
		lf.model.Fields.Origin.Set(selectOrigin.SelectedIndex())
	}

	labelOrientation := widget.NewLabel(resources.OrientationLabel.String())
	selectOrientation := widget.NewSelect(resources.OrientationLabels, func(s string) {})
	selectOrientation.OnChanged = func(s string) {
		lf.model.Fields.Orientation.Set(selectOrientation.SelectedIndex())
	}

	lf.model.Layer.AddListener(binding.NewDataListener(func() {
		origin, _ := lf.model.Fields.Origin.Get()
		selectOrigin.SetSelectedIndex(origin)
		orientation, _ := lf.model.Fields.Orientation.Get()
		selectOrientation.SetSelectedIndex(orientation)
	}))

	return lf.makePopup(labelOrigin, selectOrigin,
		labelOrientation, selectOrientation)
}

func (lf *LayerForm) makePopup(objects ...fyne.CanvasObject) *widget.PopUp {
	frm := container.New(layout.NewFormLayout())
	frm.Objects = append(frm.Objects, objects...)
	return widget.NewPopUp(frm, lf.window.Canvas())
}
