package ui

import (
	"gglow/fyglow/effectio"
	"gglow/text"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type ExportDialog struct {
	*dialog.CustomDialog
	window      fyne.Window
	effect      *effectio.EffectIo
	data        binding.BoolTree
	tree        *widget.Tree
	options     *widget.CheckGroup
	listener    binding.DataListener
	applyStatus binding.Bool
	labels      []string
}

func NewExportDialog(effect *effectio.EffectIo, window fyne.Window) *ExportDialog {
	xd := &ExportDialog{
		effect:      effect,
		window:      window,
		applyStatus: binding.NewBool(),
		labels:      make([]string, 0),
	}

	xd.options = widget.NewCheckGroup(
		driverLabels(), func(s []string) {
			xd.labels = s
			xd.listen()
		})
	xd.options.Horizontal = true

	xd.data = BuildBoolTree(effect)
	xd.listener = binding.NewDataListener(xd.listen)
	xd.tree = NewEffectTree(xd.data, xd.listener, createCheck, updateCheck(xd.data, xd.listener))

	lay := container.NewBorder(xd.options, nil, nil, nil, xd.tree)
	xd.CustomDialog = dialog.NewCustom(text.ExportLabel.String(), "", lay, window)

	cancelBtn := widget.NewButton(text.CancelLabel.String(), xd.cancel)
	applyBtn := widget.NewButton(text.NextLabel.String(), xd.apply)
	applyBtn.Disable()
	xd.CustomDialog.SetButtons([]fyne.CanvasObject{cancelBtn, applyBtn})

	xd.applyStatus.AddListener(binding.NewDataListener(func() {
		ok, _ := xd.applyStatus.Get()
		if ok {
			applyBtn.Enable()
		} else {
			applyBtn.Disable()
		}
	}))
	return xd
}

func (xd *ExportDialog) Start() {
	xd.CustomDialog.Show()
}

func (xd *ExportDialog) cancel() {
	xd.CustomDialog.Hide()
}

func (xd *ExportDialog) apply() {
	xd.CustomDialog.Hide()
	var path string
	dlg := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil || uri == nil {
			if err != nil {
				fyne.LogError("ShowFolderOpen", err)
			}
			return
		}
		path = uri.Path()
		drivers := driversFromLabels(xd.labels)
		act := BuildAction(xd.data, xd.effect, drivers, path)
		ConfirmActionDialog(act, xd.window)
	}, xd.window)

	dlg.Resize(xd.window.Canvas().Size())
	dlg.SetConfirmText(text.SelectLabel.String())
	dlg.Show()
}

func (xd *ExportDialog) listen() {
	var status bool
	statusPtr := &status
	statusSet := func() {
		xd.applyStatus.Set(*statusPtr)
	}
	defer statusSet()

	if len(xd.labels) < 1 {
		return
	}

	_, vals, err := xd.data.Get()
	if err != nil {
		fyne.LogError("ExportDialog.listen", err)
		return
	}

	for _, v := range vals {
		if v {
			status = true
			return
		}
	}
}
