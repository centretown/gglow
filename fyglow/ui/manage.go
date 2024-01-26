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

type Manager struct {
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

func NewManager(effect *effectio.EffectIo, window fyne.Window) *Manager {
	mgr := &Manager{
		window:  window,
		data:    BuildBoolTree(effect),
		step:    binding.NewInt(),
		path:    binding.NewString(),
		tabList: make([]*container.TabItem, STEP_COUNT),
		act:     action.NewAction(),
	}

	mgr.listener = binding.NewDataListener(mgr.listen)
	mgr.tree = NewEffectTreeWithListener(mgr.data, mgr.listener,
		CreateCheck,
		UpdateCheck(mgr.data, mgr.listener))

	mgr.tabList[STEP_FOLDER] = mgr.destinationTab()
	mgr.path.AddListener(binding.NewDataListener(func() {
		mgr.listen()
	}))

	mgr.tabList[STEP_PROFILE] = mgr.profileTab()

	mgr.reviewLabel = widget.NewLabel(text.ReviewLabel.String())
	mgr.tabList[STEP_REVIEW] = container.NewTabItem(text.ReviewLabel.String(),
		WrapVertical(mgr.reviewLabel, ConfirmView(mgr.act)))

	confirm := widget.NewLabel(text.ConfirmLabel.String())
	mgr.tabList[STEP_CONFIRM] = container.NewTabItem(text.ConfirmLabel.String(),
		WrapVertical(confirm, ConfirmView(mgr.act)))

	mgr.tabs = container.NewAppTabs(mgr.tabList...)
	mgr.tabs.OnSelected = mgr.onSelected(effect, confirm)

	cancelBtn := widget.NewButton(text.CancelLabel.String(), mgr.cancel)
	mgr.nextButton = widget.NewButton(text.NextLabel.String(), mgr.next)
	mgr.nextButton.Disable()

	mgr.CustomDialog = dialog.NewCustom(text.ManageLabel.String(),
		text.CancelLabel.String(), mgr.tabs, window)
	mgr.CustomDialog.SetButtons([]fyne.CanvasObject{cancelBtn, mgr.nextButton})

	mgr.step.AddListener(binding.NewDataListener(mgr.updateTabItems))

	return mgr
}

func (mgr *Manager) onSelected(effect *effectio.EffectIo, confirm fyne.CanvasObject) func(ti *container.TabItem) {
	return func(ti *container.TabItem) {
		current := mgr.tabs.SelectedIndex()
		switch current {
		case STEP_CONFIRM:
			path, _ := mgr.path.Get()
			mgr.act = BuildAction(mgr.data, effect, []string{mgr.driver}, path)
			mgr.tabList[STEP_CONFIRM].Content =
				WrapVertical(confirm, ConfirmView(mgr.act))
			mgr.nextButton.SetText(text.ProceedLabel.String())
			mgr.nextButton.Importance = widget.HighImportance
			mgr.nextButton.Refresh()
			return
		default:
			mgr.nextButton.SetText(text.NextLabel.String())
			mgr.nextButton.Importance = widget.LowImportance
		}
		mgr.setButtons()
	}
}

func (mgr *Manager) profileTab() *container.TabItem {
	profile := widget.NewSelect([]string{"light strip", "grid 8x8", "grid 4x9"},
		func(s string) {})
	profile.SetSelectedIndex(0)
	return container.NewTabItem(text.ProfileLabel.String(),
		WrapVertical(profile, mgr.tree))
}

func (mgr *Manager) destinationTab() *container.TabItem {
	options := widget.NewRadioGroup(
		driverLabels(), func(s string) { mgr.driver = driverFromLabel(s) })
	options.Horizontal = true
	options.SetSelected(driverLabels()[0])

	folderSelector := NewFolderSelector(mgr.path, mgr.window)
	pathLabel := widget.NewLabelWithData(mgr.path)
	destinationBox := container.NewHBox(options, folderSelector, pathLabel)
	return container.NewTabItem(text.DestinationLabel.String(),
		WrapVertical(destinationBox, mgr.tree))
}

func (mgr *Manager) onProceed() {
	err := mgr.act.Process()
	if err != nil {
		fyne.LogError("ConfirmAction", err)
	}
	if mgr.act.HasErrors() {
		mgr.reviewLabel.SetText(text.ActionError.String())
	} else {
		mgr.reviewLabel.SetText(text.ActionSuccess.String())
	}
	mgr.tabList[STEP_REVIEW].Content =
		WrapVertical(mgr.reviewLabel, ConfirmView(mgr.act))
	mgr.tabs.EnableIndex(STEP_REVIEW)
	mgr.tabs.SelectIndex(STEP_REVIEW)
}

func (mgr *Manager) cancel() {
	mgr.CustomDialog.Hide()
}

func (mgr *Manager) updateTabItems() {
	step, _ := mgr.step.Get()
	for i := range mgr.tabs.Items {
		if i <= step {
			mgr.tabs.EnableIndex(i)
		} else {
			mgr.tabs.DisableIndex(i)
		}
	}
	mgr.setButtons()
}

func (mgr *Manager) next() {
	current := mgr.tabs.SelectedIndex()
	step, _ := mgr.step.Get()
	if current < step {
		if current < STEP_LAST {
			mgr.tabs.SelectIndex(current + 1)
		}
	} else if current == STEP_CONFIRM {
		mgr.onProceed()
	}
}

func (mgr *Manager) setButtons() {
	step, _ := mgr.step.Get()
	current := mgr.tabs.SelectedIndex()
	if current < step {
		mgr.nextButton.Enable()
	} else {
		mgr.nextButton.Disable()
	}
	mgr.nextButton.Refresh()
}

func (mgr *Manager) listen() {
	step, _ := mgr.step.Get()
	stepPtr := &step
	stepSet := func() {
		mgr.step.Set(*stepPtr)
	}
	defer stepSet()

	path, _ := mgr.path.Get()
	if len(path) < 1 {
		step = STEP_FOLDER
		return
	}

	step = STEP_PROFILE
	_, vals, err := mgr.data.Get()
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
