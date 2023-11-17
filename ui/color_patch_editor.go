package ui

import (
	"glow-gui/glow"
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
	source      *ColorPatch
	patch       *ColorPatch
	applyButton *widget.Button
	isDirty     binding.Bool
	window      fyne.Window

	hue          binding.Float
	saturation   binding.Float
	value        binding.Float
	disabled     binding.Bool
	disableCheck *widget.Check
}

func NewColorPatchEditor(source *ColorPatch,
	isDirty binding.Bool,
	window fyne.Window) *ColorPatchEditor {

	pe := &ColorPatchEditor{
		source:  source,
		patch:   NewColorPatchWithColor(source.GetHSVColor(), func() {}),
		isDirty: isDirty,
		window:  window,

		hue:        binding.NewFloat(),
		saturation: binding.NewFloat(),
		value:      binding.NewFloat(),
		disabled:   binding.NewBool(),
	}

	hueLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(pe.hue, "%.0f"))
	hueSlider := widget.NewSliderWithData(0, 360, pe.hue)
	hueBox := container.NewBorder(nil, nil, widget.NewLabel("H"), hueLabel, hueSlider)

	saturationLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(pe.saturation, "%.0f"))
	saturationSlider := widget.NewSliderWithData(0, 100, pe.saturation)
	saturationBox := container.NewBorder(nil, nil, widget.NewLabel("S"), saturationLabel, saturationSlider)

	valueLabel := widget.NewLabelWithData(binding.FloatToStringWithFormat(pe.value, "%.0f"))
	valueSlider := widget.NewSliderWithData(0, 100, pe.value)
	valueBox := container.NewBorder(nil, nil, widget.NewLabel("V"), valueLabel, valueSlider)

	// pe.patch.SetTapped(pe.selectColorPicker(pe.patch))
	pickerButton := widget.NewButtonWithIcon("Color Picker", theme.MoreHorizontalIcon(),
		pe.selectColorPicker(pe.patch))
	pe.disableCheck = widget.NewCheckWithData("Disable", pe.disabled)

	pe.setFields()

	pe.disabled.AddListener(binding.NewDataListener(
		func() {
			disable, _ := pe.disabled.Get()
			pe.patch.SetDisabled(disable)
			if disable {
				pickerButton.Disable()
				pe.isDirty.Set(true)
			} else {
				pe.patch.SetHSVColor(pe.source.GetHSVColor())
				pe.isDirty.Set(false)
			}
		}))
	pe.hue.AddListener(binding.NewDataListener(pe.setColor))
	pe.saturation.AddListener(binding.NewDataListener(pe.setColor))
	pe.value.AddListener(binding.NewDataListener(pe.setColor))

	revertButton := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		pe.CustomDialog.Hide()
	})

	pe.applyButton = widget.NewButtonWithIcon("Apply", theme.ConfirmIcon(), pe.apply)
	vbox := container.NewVBox(pe.disableCheck,
		pe.patch,
		hueBox,
		saturationBox,
		valueBox,
		pickerButton,
		widget.NewSeparator())

	pe.CustomDialog = dialog.NewCustomWithoutButtons("color", vbox, window)
	pe.CustomDialog.SetButtons([]fyne.CanvasObject{revertButton, pe.applyButton})
	return pe
}

func (pe *ColorPatchEditor) setFields() {
	pe.disabled.Set(pe.patch.disabled)
	pe.disableCheck.SetChecked(pe.patch.disabled)
	pe.hue.Set(float64(pe.patch.colorHSV.Hue))
	pe.saturation.Set(float64(pe.patch.colorHSV.Saturation) * 100)
	pe.value.Set(float64(pe.patch.colorHSV.Value * 100))
}

func (pe *ColorPatchEditor) setColor() {
	h, _ := pe.hue.Get()
	s, _ := pe.saturation.Get()
	v, _ := pe.value.Get()
	hsv := glow.HSV{
		Hue:        float32(h),
		Saturation: float32(s / 100),
		Value:      float32(v / 100)}
	pe.patch.SetHSVColor(hsv)
	pe.isDirty.Set(true)
}

func (pe *ColorPatchEditor) apply() {
	pe.source.Copy(pe.patch)
	pe.CustomDialog.Hide()
}

func (le *ColorPatchEditor) selectColorPicker(patch *ColorPatch) func() {
	return func() {
		picker := dialog.NewColorPicker("Color Picker", "color", func(c color.Color) {
			if c != patch.GetColor() {
				patch.SetColor(c)
				le.isDirty.Set(true)
			}
		}, le.window)
		picker.Advanced = true
		picker.SetColor(patch.GetColor())
		picker.Show()
	}
}
