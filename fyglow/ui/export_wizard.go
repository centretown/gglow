package ui

import (
	"gglow/action"
	"gglow/fyglow/effectio"
	"gglow/text"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const (
	STEP_FOLDER = iota
	STEP_PROFILE
	STEP_CONFIRM
	STEP_REVIEW
	STEP_COUNT
	STEP_LAST = STEP_COUNT - 1
)

type ExportWizard struct {
	*dialog.CustomDialog
	window fyne.Window

	data     binding.BoolTree
	path     binding.String
	step     binding.Int
	listener binding.DataListener
	driver   string

	tree        *widget.Tree
	tabs        *container.AppTabs
	tabList     []*container.TabItem
	reviewLabel *widget.Label
	nextButton  *widget.Button

	act *action.Action
}

func NewExportWizard(effect *effectio.EffectIo, window fyne.Window) *ExportWizard {
	exz := &ExportWizard{
		window:  window,
		data:    effect.BuildTreeData(),
		step:    binding.NewInt(),
		path:    binding.NewString(),
		tabList: make([]*container.TabItem, STEP_COUNT),
		act:     action.NewAction(),
	}

	exz.listener = binding.NewDataListener(exz.listen)
	exz.tree = NewEffectTreeWithListener(exz.data, exz.listener,
		CreateCheck,
		UpdateCheck(exz.data, exz.listener))

	exz.tabList[STEP_FOLDER] = exz.destinationTab()
	exz.path.AddListener(binding.NewDataListener(func() {
		exz.listen()
	}))

	exz.tabList[STEP_PROFILE] = exz.profileTab()

	exz.reviewLabel = widget.NewLabel(text.ReviewLabel.String())
	exz.tabList[STEP_REVIEW] = container.NewTabItem(text.ReviewLabel.String(),
		WrapVertical(exz.reviewLabel, ConfirmView(exz.act)))

	confirm := widget.NewLabel(text.ConfirmLabel.String())
	exz.tabList[STEP_CONFIRM] = container.NewTabItem(text.ConfirmLabel.String(),
		WrapVertical(confirm, ConfirmView(exz.act)))

	exz.tabs = container.NewAppTabs(exz.tabList...)
	exz.tabs.OnSelected = exz.onSelected(effect, confirm)

	cancelBtn := widget.NewButton(text.CancelLabel.String(), exz.cancel)
	exz.nextButton = widget.NewButton(text.NextLabel.String(), exz.next)
	exz.nextButton.Disable()

	exz.CustomDialog = dialog.NewCustom(text.ExportLabel.String(),
		text.CancelLabel.String(), exz.tabs, window)
	exz.CustomDialog.SetButtons([]fyne.CanvasObject{cancelBtn, exz.nextButton})

	exz.step.AddListener(binding.NewDataListener(exz.updateTabItems))

	return exz
}

func (exz *ExportWizard) onSelected(effect *effectio.EffectIo, confirm fyne.CanvasObject) func(ti *container.TabItem) {
	return func(ti *container.TabItem) {
		current := exz.tabs.SelectedIndex()
		switch current {
		case STEP_CONFIRM:
			path, _ := exz.path.Get()
			exz.act = BuildAction(exz.data, effect, []string{exz.driver}, path)
			exz.tabList[STEP_CONFIRM].Content =
				WrapVertical(confirm, ConfirmView(exz.act))
			exz.nextButton.SetText(text.ProceedLabel.String())
			exz.nextButton.Importance = widget.HighImportance
			exz.nextButton.Refresh()
			return
		default:
			exz.nextButton.SetText(text.NextLabel.String())
			exz.nextButton.Importance = widget.LowImportance
		}
		exz.setButtons()
	}
}

func (exz *ExportWizard) profileTab() *container.TabItem {
	profile := widget.NewSelect([]string{"light strip", "grid 8x8", "grid 4x9"},
		func(s string) {})
	profile.SetSelectedIndex(0)
	return container.NewTabItem(text.ProfileLabel.String(),
		WrapVertical(profile, exz.tree))
}

func (exz *ExportWizard) destinationTab() *container.TabItem {
	options := widget.NewRadioGroup(
		DriverLabels(), func(s string) { exz.driver = DriverFromLabel(s) })
	options.Horizontal = true
	options.SetSelected(DriverLabels()[0])

	folderSelector := NewFolderSelector(exz.path, exz.window)
	pathLabel := widget.NewLabelWithData(exz.path)
	destinationBox := container.NewHBox(options, folderSelector, pathLabel)
	return container.NewTabItem(text.DestinationLabel.String(),
		WrapVertical(destinationBox, exz.tree))
}

func (exz *ExportWizard) onProceed() {
	err := exz.act.Process()
	if err != nil {
		fyne.LogError("ConfirmAction", err)
	}
	if exz.act.HasErrors() {
		exz.reviewLabel.SetText(text.ActionError.String())
	} else {
		exz.reviewLabel.SetText(text.ActionSuccess.String())
	}
	exz.tabList[STEP_REVIEW].Content =
		WrapVertical(exz.reviewLabel, ConfirmView(exz.act))
	exz.tabs.EnableIndex(STEP_REVIEW)
	exz.tabs.SelectIndex(STEP_REVIEW)
}

func (exz *ExportWizard) cancel() {
	exz.CustomDialog.Hide()
}

func (exz *ExportWizard) updateTabItems() {
	step, _ := exz.step.Get()
	for i := range exz.tabs.Items {
		if i <= step {
			exz.tabs.EnableIndex(i)
		} else {
			exz.tabs.DisableIndex(i)
		}
	}
	exz.setButtons()
}

func (exz *ExportWizard) next() {
	current := exz.tabs.SelectedIndex()
	step, _ := exz.step.Get()
	if current < step {
		if current < STEP_LAST {
			exz.tabs.SelectIndex(current + 1)
		}
	} else if current == STEP_CONFIRM {
		exz.onProceed()
	}
}

func (exz *ExportWizard) setButtons() {
	step, _ := exz.step.Get()
	current := exz.tabs.SelectedIndex()
	if current < step {
		exz.nextButton.Enable()
	} else {
		exz.nextButton.Disable()
	}
	exz.nextButton.Refresh()
}

func (exz *ExportWizard) listen() {
	step, _ := exz.step.Get()
	stepPtr := &step
	stepSet := func() {
		exz.step.Set(*stepPtr)
	}
	defer stepSet()

	path, _ := exz.path.Get()
	if len(path) < 1 {
		step = STEP_FOLDER
		return
	}

	step = STEP_PROFILE
	_, vals, err := exz.data.Get()
	if err != nil {
		fyne.LogError("ExportDialog.listen", err)
		return
	}

	for _, v := range vals {
		if v {
			step = STEP_CONFIRM
			return
		}
	}
}
