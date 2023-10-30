package ui

import (
	"glow-gui/data"
	"glow-gui/glow"
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
	model   *data.Model
	layer   *glow.Layer
	fields  *data.Fields
	window  fyne.Window
	isDirty binding.Bool

	bDynamic  bool
	bScanner  bool
	bOverride bool

	selectOrigin      *widget.Select
	selectOrientation *widget.Select

	checkScan *widget.Check
	checkHue  *widget.Check
	checkRate *widget.Check

	scanBox *RangeIntBox
	hueBox  *RangeIntBox
	rateBox *RangeIntBox

	dropDownOffset float32
	applyButton    *widget.Button
	revertButton   *widget.Button
	buttonBox      *fyne.Container
	rateBounds     *IntEntryBounds
	hueBounds      *IntEntryBounds
	scanBounds     *IntEntryBounds
}

func NewLayerEditor(model *data.Model, window fyne.Window) *LayerEditor {
	le := &LayerEditor{
		Container: container.NewHBox(),
		window:    window,

		model: model,
		layer: model.GetCurrentLayer(),

		fields:  data.NewFields(),
		isDirty: binding.NewBool(),

		rateBounds: &IntEntryBounds{MinVal: 16, MaxVal: 360, OnVal: 48, OffVal: 0},
		hueBounds:  &IntEntryBounds{MinVal: -10, MaxVal: 10, OnVal: 1, OffVal: 0},
		scanBounds: &IntEntryBounds{MinVal: 1, MaxVal: 10, OnVal: 1, OffVal: 0},

		selectOrigin:      widget.NewSelect(resources.OriginLabels, func(s string) {}),
		selectOrientation: widget.NewSelect(resources.OrientationLabels, func(s string) {}),

		applyButton:  widget.NewButton(resources.ApplyLabel.String(), func() {}),
		revertButton: widget.NewButton(resources.RevertLabel.String(), func() {}),

		dropDownOffset: theme.CaptionTextSize() + 2*(theme.InnerPadding()+theme.Padding()+theme.LineSpacing()),
	}

	le.applyButton.Disable()
	le.revertButton.Disable()
	le.buttonBox = container.NewCenter(container.NewHBox(le.revertButton, le.applyButton))

	le.buildEntries()
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

func (le *LayerEditor) buildEntries() {
	le.hueBox = NewRangeIntBox(le.fields.HueShift, le.hueBounds)
	le.fields.HueShift.AddListener(binding.NewDataListener(func() {
		shift, _ := le.fields.HueShift.Get()
		le.isDirty.Set(int16(shift) != le.layer.HueShift)
	}))
	le.checkHue = widget.NewCheck("", checkRangeBox(le.hueBox, le.fields.HueShift))

	le.scanBox = NewRangeIntBox(le.fields.Scan, le.scanBounds)
	le.fields.Scan.AddListener(binding.NewDataListener(func() {
		scan, _ := le.fields.Scan.Get()
		le.isDirty.Set(uint16(scan) != le.layer.Scan)
	}))
	le.checkScan = widget.NewCheck("", checkRangeBox(le.scanBox, le.fields.Scan))

	le.rateBox = NewRangeIntBox(le.fields.Rate, le.rateBounds)
	le.fields.Rate.AddListener(binding.NewDataListener(func() {
		rate, _ := le.fields.Rate.Get()
		le.isDirty.Set(uint32(rate) != le.layer.Rate)
	}))
	le.checkRate = widget.NewCheck("", checkRangeBox(le.rateBox, le.fields.Rate))
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
				le.showDropDown(makeDropDown(le.newRateForm()))),
		}
}

func (le *LayerEditor) showDropDown(dropDown *widget.PopUp) func() {
	f := func() {
		le.isDirty.Set(le.fields.IsDirty(le.model.GetCurrentLayer()))
		le.applyButton.OnTapped = le.apply(dropDown)
		le.revertButton.OnTapped = le.revert()
		dropDown.Resize(fyne.Size{Width: minStripWidth - 2*theme.Padding(),
			Height: minStripHeight})
		offset := fyne.NewDelta(theme.Padding(), le.dropDownOffset)
		dropDown.Move(le.Container.Position().Add(offset))
		dropDown.Show()
	}
	return f
}

func (le *LayerEditor) newRateForm() *fyne.Container {
	label := widget.NewLabel(resources.RateLabel.String())
	checkLabel := widget.NewLabel(resources.OverrideLabel.String())
	return container.New(layout.NewFormLayout(),
		checkLabel, le.checkRate, label, le.rateBox.Container)
}

func (le *LayerEditor) newHueForm() *fyne.Container {
	label := widget.NewLabel(resources.HueShiftLabel.String())
	checkLabel := widget.NewLabel(resources.DynamicLabel.String())
	return container.New(layout.NewFormLayout(),
		checkLabel, le.checkHue, label, le.hueBox.Container)
}

func (le *LayerEditor) newScanForm() *fyne.Container {
	label := widget.NewLabel(resources.LengthLabel.String())
	checkLabel := widget.NewLabel(resources.ScanLabel.String())
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

	return container.New(layout.NewFormLayout(),
		labelOrigin, le.selectOrigin,
		labelOrientation, le.selectOrientation)
}

func (le *LayerEditor) setFields() {
	le.layer = le.model.GetCurrentLayer()
	le.fields.FromLayer(le.layer)

	le.selectOrigin.SetSelectedIndex(int(le.layer.Grid.Origin))
	le.selectOrientation.SetSelectedIndex(int(le.layer.Grid.Orientation))

	le.bDynamic = (le.layer.HueShift != int16(le.hueBounds.OffVal))
	le.checkHue.SetChecked(le.bDynamic)
	le.hueBox.Enable(le.bDynamic)

	le.bScanner = (le.layer.Scan != uint16(le.scanBounds.OffVal))
	le.checkScan.SetChecked(le.bScanner)
	le.scanBox.Enable(le.bScanner)

	le.bOverride = (le.layer.Rate != uint32(le.rateBounds.OffVal))
	le.checkRate.SetChecked(le.bOverride)
	le.rateBox.Enable(le.bOverride)
}

func (le *LayerEditor) apply(dropDown *widget.PopUp) func() {
	return func() {
		dropDown.Hide()
	}
}

func (le *LayerEditor) revert() func() {
	return func() {
		le.setFields()
		le.isDirty.Set(false)
	}
}
