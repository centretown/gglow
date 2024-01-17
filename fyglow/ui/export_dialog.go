package ui

import (
	"gglow/fyglow/effectio"
	"gglow/text"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type ExportDialog struct {
	*dialog.CustomDialog
	effect *effectio.EffectIo
	data   binding.BoolTree
	tree   *widget.Tree
}

func NewExportDialog(effect *effectio.EffectIo, window fyne.Window) *ExportDialog {
	xd := &ExportDialog{
		effect: effect,
	}
	xd.tree = xd.buildTree()
	xd.CustomDialog = dialog.NewCustom("Export", text.CancelLabel.String(), xd.tree, window)

	xd.CustomDialog.SetOnClosed(func() {
		action := BuildAction(xd.data, effect)
		err := action.Process()
		if err != nil {
			fyne.LogError("MAKECODE", err)
		}
	})
	return xd
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
