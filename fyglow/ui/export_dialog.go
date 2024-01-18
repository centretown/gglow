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
	window  fyne.Window
	effect  *effectio.EffectIo
	data    binding.BoolTree
	tree    *widget.Tree
	options *widget.CheckGroup
}

func NewExportDialog(effect *effectio.EffectIo, window fyne.Window) *ExportDialog {
	xd := &ExportDialog{
		effect: effect,
		window: window,
	}
	xd.options = widget.NewCheckGroup(
		[]string{"Code", "Data"}, func(s []string) {})
	xd.options.Horizontal = true
	xd.tree = xd.buildTree()
	lay := container.NewBorder(xd.options, nil, nil, nil, xd.tree)
	xd.CustomDialog = dialog.NewCustom("Export", "", lay, window)

	cancel := widget.NewButton(text.CancelLabel.String(), xd.cancel)
	apply := widget.NewButton(text.ApplyLabel.String(), xd.apply)
	xd.SetButtons([]fyne.CanvasObject{cancel, apply})
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
		act := BuildAction(xd.data, xd.effect, "code", path)
		err = act.Process()
		if err != nil {
			fyne.LogError("Export Code", err)
		}

		ShowActionResults(act, xd.window)

	}, xd.window)

	dlg.Resize(xd.window.Canvas().Size())
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

func (xd *ExportDialog) update(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
	check := o.(*widget.Check)
	check.SetText(id)
	i, _ := xd.data.GetItem(id)
	check.Bind(i.(binding.Bool))
}

func (xd *ExportDialog) buildTree() *widget.Tree {
	xd.data = BuildBoolTree(xd.effect)
	return widget.NewTree(
		xd.data.ChildIDs,
		xd.isBranch,
		xd.create,
		xd.update)

}
