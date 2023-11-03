package ui

import (
	"glow-gui/data"
	"glow-gui/resources"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type FrameTools struct {
	*widget.Toolbar
	model        *data.Model
	frameRate    *ButtonItem
	saveFrame    *ButtonItem
	restoreFrame *ButtonItem
	newFrame     *ButtonItem
}

func NewFrameTools(model *data.Model) *FrameTools {
	ft := &FrameTools{
		model: model,
	}
	ft.frameRate = NewButtonItem(
		widget.NewButtonWithIcon("",
			resources.AppIconResource(resources.SpeedIcon), ft.speed))
	ft.saveFrame = NewButtonItem(
		widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), ft.save))
	ft.restoreFrame = NewButtonItem(
		widget.NewButtonWithIcon("", theme.MediaReplayIcon(), ft.reload))
	ft.newFrame = NewButtonItem(
		widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), ft.create))

	ft.Toolbar = widget.NewToolbar(
		ft.frameRate,
		ft.saveFrame,
		ft.restoreFrame,
		ft.newFrame,
	)
	return ft
}

func (ft *FrameTools) create() {
}
func (ft *FrameTools) speed() {
}
func (ft *FrameTools) reload() {
}
func (ft *FrameTools) save() {
}
