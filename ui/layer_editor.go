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
	window fyne.Window
	// frame  *glow.Frame
	isDirty binding.Bool
	// isDynamic  binding.Bool
	// isScanner  binding.Bool
	// isOverride binding.Bool

	bDynamic  bool
	bScanner  bool
	bOverride bool

	selectOrigin      *widget.Select
	selectOrientation *widget.Select
	checkScan         *widget.Check
	checkHue          *widget.Check
	checkRate         *widget.Check

	scanBox *RangeIntBox
	hueBox  *RangeIntBox
	rateBox *RangeIntBox

	offsetHeight float32
	applyButton  *widget.Button
	revertButton *widget.Button
	buttonBox    *fyne.Container
	rateBounds   *EntryBoundsInt
	hueBounds    *EntryBoundsInt
	scanBounds   *EntryBoundsInt
}

func NewLayerEditor(model *data.Model, window fyne.Window) *LayerEditor {
	le := &LayerEditor{
		Container: container.NewHBox(),
		window:    window,

		model:   model,
		fields:  data.NewFields(),
		isDirty: binding.NewBool(),
		// isDynamic:  binding.NewBool(),
		// isScanner:  binding.NewBool(),
		// isOverride: binding.NewBool(),

		rateBounds: &EntryBoundsInt{MinVal: 16, MaxVal: 360, OnVal: 48, OffVal: 0},
		hueBounds:  &EntryBoundsInt{MinVal: 1, MaxVal: 10, OnVal: 1, OffVal: 0},
		scanBounds: &EntryBoundsInt{MinVal: 1, MaxVal: 10, OnVal: 1, OffVal: 0},

		selectOrigin:      widget.NewSelect(resources.OriginLabels, func(s string) {}),
		selectOrientation: widget.NewSelect(resources.OrientationLabels, func(s string) {}),

		applyButton:  widget.NewButton(resources.ApplyLabel.String(), func() {}),
		revertButton: widget.NewButton(resources.RevertLabel.String(), func() {}),
		offsetHeight: theme.CaptionTextSize() + 2*(theme.InnerPadding()+theme.Padding()+theme.LineSpacing()),
	}

	le.applyButton.Disable()
	le.revertButton.Disable()
	le.buttonBox = container.NewCenter(container.NewHBox(le.revertButton, le.applyButton))

	le.hueBox = NewRangeIntBox(le.fields.HueShift, le.hueBounds)
	le.scanBox = NewRangeIntBox(le.fields.Scan, le.scanBounds)
	le.rateBox = NewRangeIntBox(le.fields.Rate, le.rateBounds)

	le.checkScan = widget.NewCheck("", checkBox(le.scanBox, le.fields.Scan))
	le.checkHue = widget.NewCheck("", checkBox(le.hueBox, le.fields.HueShift))
	le.checkRate = widget.NewCheck("", checkBox(le.rateBox, le.fields.Rate))

	le.buildButtons()

	le.model.Layer.AddListener(binding.NewDataListener(le.setFields))
	le.isDirty.AddListener(binding.NewDataListener(func() {
		b, _ := le.isDirty.Get()
		if b {
			le.applyButton.Enable()
			le.revertButton.Enable()
		} else {
			le.applyButton.Disable()
			le.revertButton.Disable()
		}
	}))
	return le
}

func checkBox(rangeBox *RangeIntBox, field binding.Int) func(bool) {
	return func(b bool) {
		if b {
			i, _ := field.Get()
			if i == rangeBox.Bounds.OffVal {
				field.Set(rangeBox.Bounds.OnVal)
			}
			rangeBox.Enable(true)
		} else {
			field.Set(rangeBox.Bounds.OffVal)
			rangeBox.Enable(false)
		}
	}
}

func (le *LayerEditor) IsDirty() bool {
	b, _ := le.isDirty.Get()
	return b
}

func (le *LayerEditor) buildButtons() {

	sep := widget.NewSeparator()
	makeDropDown := func(form *fyne.Container) *widget.PopUp {
		dropDown := widget.NewPopUp(container.NewVBox(form, sep, le.buttonBox), le.window.Canvas())
		return dropDown
	}

	le.Container.Objects =
		[]fyne.CanvasObject{
			widget.NewButton(resources.GridLabel.String(),
				le.showDropDown(makeDropDown(le.newGridForm()))),

			widget.NewButton(resources.ScanLabel.String(),
				le.showDropDown(makeDropDown(le.newScanForm()))),

			widget.NewButton(resources.ChromaLabel.String(), func() {}),

			widget.NewButton(resources.HueLabel.String(),
				le.showDropDown(makeDropDown(le.newHueForm()))),

			widget.NewButton(resources.RateLabel.String(),
				le.showDropDown(makeDropDown(le.newRateForm())))}

}

func (le *LayerEditor) setFields() {
	layer := le.model.GetCurrentLayer()
	le.fields.FromLayer(layer)

	le.selectOrigin.SetSelectedIndex(int(layer.Grid.Origin))
	le.selectOrientation.SetSelectedIndex(int(layer.Grid.Orientation))

	le.bDynamic = (layer.HueShift != int16(le.hueBounds.OffVal))
	le.checkHue.SetChecked(le.bDynamic)
	le.hueBox.Enable(le.bDynamic)

	le.bScanner = (layer.Scan != uint16(le.scanBounds.OffVal))
	le.checkScan.SetChecked(le.bScanner)
	le.scanBox.Enable(le.bScanner)

	le.bOverride = (layer.Rate != uint32(le.rateBounds.OffVal))
	le.checkRate.SetChecked(le.bOverride)
	le.rateBox.Enable(le.bOverride)
}

func (le *LayerEditor) showDropDown(dropDown *widget.PopUp) func() {
	f := func() {
		le.isDirty.Set(le.fields.IsDirty(le.model.GetCurrentLayer()))
		le.applyButton.OnTapped = le.apply(dropDown)
		le.revertButton.OnTapped = le.revert()
		dropDown.Resize(fyne.Size{Width: minStripWidth - 2*theme.Padding(),
			Height: minStripHeight})
		offset := fyne.NewDelta(theme.Padding(), le.offsetHeight)
		dropDown.Move(le.Container.Position().Add(offset))
		dropDown.Show()
	}
	return f
}

func addDirty[T int16 | uint16 | uint32](isDirty binding.Bool,
	field binding.Int, val T) {

	field.AddListener(binding.NewDataListener(func() {
		i, _ := field.Get()
		isDirty.Set(val != T(i))
	}))
}

func (le *LayerEditor) newRateForm() *fyne.Container {
	label := widget.NewLabel(resources.RateLabel.String())
	checkLabel := widget.NewLabel(resources.OverrideLabel.String())
	rate := le.model.GetCurrentLayer().Rate
	addDirty(le.isDirty, le.fields.Rate, rate)
	return container.New(layout.NewFormLayout(),
		checkLabel, le.checkRate, label, le.rateBox.Container)
}

func (le *LayerEditor) newHueForm() *fyne.Container {
	label := widget.NewLabel(resources.HueShiftLabel.String())
	checkLabel := widget.NewLabel(resources.DynamicLabel.String())
	shift := le.model.GetCurrentLayer().HueShift

	addDirty(le.isDirty, le.fields.HueShift, shift)
	return container.New(layout.NewFormLayout(),
		checkLabel, le.checkHue, label, le.hueBox.Container)
}

func (le *LayerEditor) newScanForm() *fyne.Container {
	label := widget.NewLabel(resources.LengthLabel.String())
	checkLabel := widget.NewLabel(resources.ScanLabel.String())
	scan := le.model.GetCurrentLayer().Scan

	addDirty(le.isDirty, le.fields.Scan, scan)
	return container.New(layout.NewFormLayout(),
		checkLabel, le.checkScan, label, le.scanBox.Container)
}

func (le *LayerEditor) newGridForm() *fyne.Container {
	labelOrigin := widget.NewLabel(resources.OriginLabel.String())
	le.selectOrigin.OnChanged = func(s string) {
		current := le.model.GetCurrentLayer().Grid.Origin
		selected := le.selectOrigin.SelectedIndex()
		if selected != int(current) {
			le.fields.Origin.Set(selected)
			le.isDirty.Set(true)
		}
	}
	labelOrientation := widget.NewLabel(resources.OrientationLabel.String())
	le.selectOrientation.OnChanged = func(s string) {
		current := le.model.GetCurrentLayer().Grid.Orientation
		selected := le.selectOrientation.SelectedIndex()
		if selected != int(current) {
			le.fields.Orientation.Set(selected)
			le.isDirty.Set(true)
		}
	}

	// addDirty(le.isDirty, le.fields.Origin, le.model.GetCurrentLayer().Grid.Origin)
	// addDirty(le.isDirty, le.fields.Orientation, le.model.GetCurrentLayer().Grid.Orientation)

	return container.New(layout.NewFormLayout(),
		labelOrigin, le.selectOrigin,
		labelOrientation, le.selectOrientation)
}

func (le *LayerEditor) apply(dropDown *widget.PopUp) func() {
	return func() {
		dropDown.Hide()

		// lf.fields.ToLayer()
	}
}

func (le *LayerEditor) revert() func() {
	return func() {
		// dropDown.Hide()
		le.setFields()
		le.isDirty.Set(false)
	}
}
