package ui

import (
	"glow-gui/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type FrameTools struct {
	saveFrame    *ButtonItem
	newFrame     *ButtonItem
	createDialog *CreateEffect
}

func NewFrameTools(model *data.Model, window fyne.Window) *FrameTools {
	ft := &FrameTools{}

	ft.createDialog = ft.newEffectDialog(window)

	ft.saveFrame = NewButtonItem(
		widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), ft.save))
	ft.newFrame = NewButtonItem(
		widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
			ft.createDialog.Show()
		}))
	return ft
}

func (ft *FrameTools) Items() (items []widget.ToolbarItem) {
	items = []widget.ToolbarItem{
		ft.saveFrame,
		ft.newFrame,
	}
	return
}

type CreateEffect struct {
	*dialog.CustomDialog
	name       binding.String
	profileAdd binding.Bool
	profile    binding.String
}

func (ft *FrameTools) newEffectDialog(window fyne.Window) (ce *CreateEffect) {
	ce = &CreateEffect{
		name:       binding.NewString(),
		profileAdd: binding.NewBool(),
		profile:    binding.NewString(),
	}

	nameLabel := widget.NewLabel("Name")
	nameEntry := widget.NewEntryWithData(ce.name)

	profileCheckLabel := widget.NewLabel("Add Profile")
	profileCheck := widget.NewCheckWithData("", ce.profileAdd)

	profileLabel := widget.NewLabel("Profile Name")
	profile := widget.NewEntryWithData(ce.profile)

	frm := container.New(layout.NewFormLayout(),
		nameLabel, nameEntry,
		profileCheckLabel, profileCheck,
		profileLabel, profile)

	ce.CustomDialog = dialog.NewCustomWithoutButtons("Create Effect", frm, window)

	applyButton := NewButtonItem(
		widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {

		}))

	revertButton := NewButtonItem(
		widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
			ce.CustomDialog.Hide()
		}))
	ce.CustomDialog.SetButtons([]fyne.CanvasObject{revertButton, applyButton})

	return ce
}

func (ft *FrameTools) restore() {
}

func (ft *FrameTools) save() {
}
