package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type ButtonItem struct {
	*widget.Button
}

func NewButtonItem(button *widget.Button) (bi *ButtonItem) {
	button.Importance = widget.LowImportance
	bi = &ButtonItem{
		Button: button,
	}
	return
}

func (bi *ButtonItem) ToolbarObject() fyne.CanvasObject {
	return bi.Button
}

func (bi *ButtonItem) Enable() {
	bi.Button.Enable()
}
func (bi *ButtonItem) Disable() {
	bi.Button.Disable()
}

type LabelItem struct {
	*widget.Label
}

func NewLabelItem(label *widget.Label) *LabelItem {
	lb := &LabelItem{
		Label: label,
	}
	return lb
}

func (lb *LabelItem) ToolbarObject() fyne.CanvasObject {
	return lb.Label
}
