package ui

import (
	"glow-gui/control"
	"glow-gui/glow"
	"glow-gui/resources"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ColorPatchEditor struct {
	*dialog.CustomDialog
	model  *control.Model
	window fyne.Window

	source      *ColorPatch
	patch       *ColorPatch
	applyButton *widget.Button

	hue          binding.Float
	saturation   binding.Float
	value        binding.Float
	unused       binding.Bool
	removeButton *widget.Button
}

func NewColorPatchEditor(source *ColorPatch,
	model *control.Model,
	window fyne.Window) *ColorPatchEditor {

	pe := &ColorPatchEditor{
		source: source,
		patch:  NewColorPatchWithColor(source.GetHSVColor(), model, func() {}),
		model:  model,
		window: window,

		hue:        binding.NewFloat(),
		saturation: binding.NewFloat(),
		value:      binding.NewFloat(),
		unused:     binding.NewBool(),
	}

	pe.patch.Editing = true

	hueLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(pe.hue, "%3.0f"))
	hueSlider := NewButtonSlide(pe.hue, HueBounds)
	hueBox := container.NewBorder(nil, nil, widget.NewLabel("H"), hueLabel,
		container.NewBorder(nil, nil, nil, nil, hueSlider.Container))

	saturationLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(pe.saturation, "%3.0f"))
	saturationSlider := NewButtonSlide(pe.saturation, SaturationBounds)
	saturationBox := container.NewBorder(nil, nil, widget.NewLabel("S"), saturationLabel,
		saturationSlider.Container)

	valueLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(pe.value, "%3.0f"))
	valueSlider := NewButtonSlide(pe.value, ValueBounds)
	valueBox := container.NewBorder(nil, nil, widget.NewLabel("V"), valueLabel,
		valueSlider.Container)

	pickerButton := widget.NewButtonWithIcon(resources.PickerLabel.String(), theme.MoreHorizontalIcon(),
		pe.selectColorPicker(pe.patch))
	pe.setFields()
	pe.removeButton = widget.NewButtonWithIcon(resources.CutLabel.String(), theme.ContentCutIcon(),
		pe.remove)

	pe.hue.AddListener(binding.NewDataListener(pe.setHue))
	pe.saturation.AddListener(binding.NewDataListener(pe.setSaturation))
	pe.value.AddListener(binding.NewDataListener(pe.setValue))

	revertButton := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		pe.CustomDialog.Hide()
	})

	pe.applyButton = widget.NewButtonWithIcon("Apply", theme.ConfirmIcon(), pe.apply)
	vbox := container.NewVBox(
		pe.patch,
		hueBox,
		saturationBox,
		valueBox,
		widget.NewSeparator(), pickerButton)

	pe.CustomDialog = dialog.NewCustomWithoutButtons("", vbox, window)
	pe.CustomDialog.SetButtons([]fyne.CanvasObject{revertButton, pe.applyButton, pe.removeButton})
	return pe
}

func (pe *ColorPatchEditor) setFields() {
	pe.unused.Set(pe.patch.unused)
	pe.hue.Set(float64(pe.patch.colorHSV.Hue))
	pe.saturation.Set(float64(pe.patch.colorHSV.Saturation) * 100)
	pe.value.Set(float64(pe.patch.colorHSV.Value * 100))
}

func (pe *ColorPatchEditor) setHue() {
	h, _ := pe.hue.Get()
	hsv := pe.patch.GetHSVColor()
	hsv.Hue = float32(h)
	pe.setColor(hsv)
}
func (pe *ColorPatchEditor) setSaturation() {
	s, _ := pe.saturation.Get()
	hsv := pe.patch.GetHSVColor()
	hsv.Saturation = float32(s / 100)
	pe.setColor(hsv)
}
func (pe *ColorPatchEditor) setValue() {
	v, _ := pe.value.Get()
	hsv := pe.patch.GetHSVColor()
	hsv.Value = float32(v / 100)
	pe.setColor(hsv)
}

func (pe *ColorPatchEditor) setColor(hsv glow.HSV) {
	pe.patch.SetHSVColor(hsv)
	pe.model.SetDirty()
}

func (pe *ColorPatchEditor) remove() {
	pe.patch.EditCut()
	pe.apply()
}

func (pe *ColorPatchEditor) apply() {
	pe.source.CopyPatch(pe.patch)
	pe.CustomDialog.Hide()
}

func (le *ColorPatchEditor) selectColorPicker(patch *ColorPatch) func() {
	return func() {
		picker := dialog.NewColorPicker("Color Picker", "color", func(c color.Color) {
			if c != patch.GetColor() {
				patch.SetColor(c)
				le.model.SetDirty()
			}
		}, le.window)
		picker.Advanced = true
		picker.SetColor(patch.GetColor())
		picker.Show()
	}
}
