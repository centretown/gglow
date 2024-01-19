package ui

import (
	"gglow/fyglow/effectio"
	"gglow/iohandler"
	"gglow/text"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var driverLabels = []string{text.CodeLabel.String(), text.DataLabel.String()}
var availableDrivers = []string{iohandler.DRIVER_CODE, iohandler.DRIVER_SQLLITE3}

func getDrivers(sel []string) (list []string) {
	list = make([]string, 0, len(sel))
	for _, s := range sel {
		for i, label := range driverLabels {
			if s == label {
				list = append(list, availableDrivers[i])
			}
		}
	}
	return
}

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

	xd.listener = binding.NewDataListener(xd.listen)
	xd.options = widget.NewCheckGroup(
		driverLabels, func(s []string) {
			xd.labels = s
			xd.listen()
		})
	xd.options.Horizontal = true

	xd.tree = xd.buildTree()
	lay := container.NewBorder(xd.options, nil, nil, nil, xd.tree)
	xd.CustomDialog = dialog.NewCustom(text.ExportLabel.String(), "", lay, window)

	cancelBtn := widget.NewButton(text.CancelLabel.String(), xd.cancel)
	applyBtn := widget.NewButton(text.NextLabel.String(), xd.apply)
	applyBtn.Disable()
	xd.SetButtons([]fyne.CanvasObject{cancelBtn, applyBtn})

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
		drivers := getDrivers(xd.labels)
		act := BuildAction(xd.data, xd.effect, drivers, path)
		err = act.Process()
		if err != nil {
			fyne.LogError("Export Code", err)
		}

		ShowActionResults(act, xd.window)

	}, xd.window)

	dlg.Resize(xd.window.Canvas().Size())
	dlg.SetConfirmText(text.SelectLabel.String())
	dlg.Show()
}

func (xd *ExportDialog) cancel() {
	xd.CustomDialog.Hide()
}

func (xd *ExportDialog) Start() {
	xd.CustomDialog.Show()
}

func (xd *ExportDialog) isBranch(id widget.TreeNodeID) bool {
	children := xd.data.ChildIDs(id)
	return len(children) > 0
}

func (xd *ExportDialog) create(branch bool) fyne.CanvasObject {
	return widget.NewCheck("NewCheck template", func(b bool) {})
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

func (xd *ExportDialog) update(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
	check := o.(*widget.Check)
	check.SetText(id)
	i, _ := xd.data.GetItem(id)
	bb := i.(binding.Bool)
	check.Bind(bb)
	bb.AddListener(xd.listener)
}

func (xd *ExportDialog) buildTree() (wtr *widget.Tree) {
	xd.data = BuildBoolTree(xd.effect)
	wtr = widget.NewTree(
		xd.data.ChildIDs,
		xd.isBranch,
		xd.create,
		xd.update)

	return
}
