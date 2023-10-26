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

type LayerEditor struct {
	*fyne.Container
	model  *data.Model
	fields *data.Fields
	// frame  *glow.Frame
	isDynamic  binding.Bool
	isScanner  binding.Bool
	isOverride binding.Bool
	window     fyne.Window

	selectOrigin      *widget.Select
	selectOrientation *widget.Select
}

func NewLayerEditor(model *data.Model, window fyne.Window) *LayerEditor {
	lf := &LayerEditor{
		Container:  container.NewHBox(),
		model:      model,
		fields:     data.NewFields(),
		isDynamic:  binding.NewBool(),
		isScanner:  binding.NewBool(),
		isOverride: binding.NewBool(),

		selectOrigin:      widget.NewSelect(resources.OriginLabels, func(s string) {}),
		selectOrientation: widget.NewSelect(resources.OrientationLabels, func(s string) {}),

		window: window,
	}

	buttonHeight := theme.CaptionTextSize() + 2*(theme.InnerPadding()+theme.Padding()+theme.LineSpacing())
	showForm := func(popup *widget.PopUp) func() {
		f := func() {
			popup.Resize(fyne.Size{Width: minStripWidth - 2*theme.Padding(),
				Height: minStripHeight})
			offset := fyne.NewDelta(theme.Padding(), buttonHeight)
			popup.Move(lf.Container.Position().Add(offset))
			popup.Show()
		}
		return f
	}

	hueDropDown := lf.newHueDropDown()
	hueButton := widget.NewButton(resources.HueLabel.String(), showForm(hueDropDown))

	scanDropDown := lf.newScanDropDown()
	scanButton := widget.NewButton(resources.ScanLabel.String(), showForm(scanDropDown))

	gridDropDown := lf.newGridDropDown()
	gridButton := widget.NewButton(resources.GridLabel.String(), showForm(gridDropDown))

	colorsButton := widget.NewButton(resources.ChromaLabel.String(), func() {})

	rateDropDown := lf.newRateDropDown()
	rateButton := widget.NewButton(resources.RateLabel.String(), showForm(rateDropDown))

	lf.Container.Objects =
		[]fyne.CanvasObject{
			gridButton,
			scanButton,
			colorsButton,
			hueButton,
			rateButton}

	// lf.model.Frame.AddListener(binding.NewDataListener(func() {
	// 	lf.frame, _ = glow.FrameDeepCopy(lf.model.GetFrame())
	// }))

	lf.model.Layer.AddListener(binding.NewDataListener(lf.setFields))

	return lf
}

func (lf *LayerEditor) setFields() {
	lf.fields.FromLayer(lf.model.GetCurrentLayer())

	origin, _ := lf.fields.Origin.Get()
	lf.selectOrigin.SetSelectedIndex(origin)
	orientation, _ := lf.fields.Orientation.Get()
	lf.selectOrientation.SetSelectedIndex(orientation)

	i, _ := lf.fields.HueShift.Get()
	lf.isDynamic.Set(i != 0)
	i, _ = lf.fields.Scan.Get()
	lf.isScanner.Set(i > 0)
	i, _ = lf.fields.Rate.Get()
	lf.isOverride.Set(i > 0)
}

func (lf *LayerEditor) newRateDropDown() *widget.PopUp {
	checkLabel, check := NewLabelCheck(resources.Override.String(), lf.isOverride)
	bounds := EntryBoundsInt{MinVal: 16, MaxVal: 360, OnVal: 48, OffVal: 0}
	le := NewDisabledEntry(resources.RateLabel.String(), lf.isOverride,
		lf.fields.Rate, &bounds)
	return lf.makeDropDown(checkLabel, check, le.Label, le.RangeBox.Container)
}

func (lf *LayerEditor) newHueDropDown() *widget.PopUp {
	checkLabel, check := NewLabelCheck(resources.DynamicLabel.String(), lf.isDynamic)
	bounds := EntryBoundsInt{MinVal: -10, MaxVal: 10, OnVal: 1, OffVal: 0}
	le := NewDisabledEntry(resources.HueShiftLabel.String(), lf.isDynamic,
		lf.fields.HueShift, &bounds)
	return lf.makeDropDown(checkLabel, check, le.Label, le.RangeBox.Container)
}

func (lf *LayerEditor) newScanDropDown() *widget.PopUp {
	checkLabel, check := NewLabelCheck(resources.ScanLabel.String(), lf.isScanner)
	bounds := EntryBoundsInt{MinVal: 1, MaxVal: 10, OnVal: 1, OffVal: 0}
	le := NewDisabledEntry(resources.LengthLabel.String(), lf.isScanner,
		lf.fields.Scan, &bounds)
	return lf.makeDropDown(checkLabel, check, le.Label, le.RangeBox.Container)
}

func (lf *LayerEditor) newGridDropDown() *widget.PopUp {

	labelOrigin := widget.NewLabel(resources.OriginLabel.String())
	lf.selectOrigin.OnChanged = func(s string) {
		lf.fields.Origin.Set(lf.selectOrigin.SelectedIndex())
	}

	labelOrientation := widget.NewLabel(resources.OrientationLabel.String())
	lf.selectOrientation.OnChanged = func(s string) {
		lf.fields.Orientation.Set(lf.selectOrientation.SelectedIndex())
	}

	return lf.makeDropDown(labelOrigin, lf.selectOrigin,
		labelOrientation, lf.selectOrientation)
}

func (lf *LayerEditor) makeDropDown(objects ...fyne.CanvasObject) *widget.PopUp {
	vbox := container.NewVBox()
	dropDown := widget.NewPopUp(vbox, lf.window.Canvas())

	frm := container.New(layout.NewFormLayout())
	frm.Objects = append(frm.Objects, objects...)
	applyButton := widget.NewButton(resources.ApplyLabel.String(), lf.apply(dropDown))
	revertButton := widget.NewButton(resources.RevertLabel.String(), lf.revert(dropDown))
	hbox := container.NewCenter(container.NewHBox(revertButton, applyButton))
	vbox.Objects = append(vbox.Objects, frm, widget.NewSeparator(), hbox)
	return dropDown
}

func (lf *LayerEditor) apply(dropDown *widget.PopUp) func() {
	return func() {
		dropDown.Hide()
		// lf.fields.ToLayer()
	}
}

func (lf *LayerEditor) revert(dropDown *widget.PopUp) func() {
	return func() {
		// dropDown.Hide()
		lf.setFields()
	}
}
