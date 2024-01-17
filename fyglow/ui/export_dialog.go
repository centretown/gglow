package ui

import (
	"fmt"
	"gglow/fyglow/effectio"
	"gglow/iohandler"
	"gglow/text"
	"strings"

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
		ls := xd.data.ChildIDs("")
		xd.showData(ls, 1)
	})
	return xd
}

func (xd *ExportDialog) showData(ids []string, depth int) {
	for _, id := range ids {
		val, _ := xd.data.GetValue(id)
		fmt.Println(strings.Repeat("\t", depth), id, val)
		xd.showData(xd.data.ChildIDs(id), depth+1)
	}
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

func (xd *ExportDialog) buildData() binding.BoolTree {
	data := binding.NewBoolTree()
	folders, _ := xd.effect.ListFolder(iohandler.DOTS)
	for _, folder := range folders {
		data.Append("", folder.Key, false)
	}

	for _, folder := range folders {
		ls, _ := xd.effect.ListFolder(folder.Key)
		for _, l := range ls {
			if l.Key != iohandler.DOTS {
				val := l.Value + "/" + l.Key
				data.Append(folder.Key, val, false)
			}
		}
	}
	return data
}
func (xd *ExportDialog) update(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
	check := o.(*widget.Check)
	check.SetText(id)
	item, _ := xd.data.GetItem(id)
	check.Bind(item.(binding.Bool))
}

func (xd *ExportDialog) buildTree() *widget.Tree {
	xd.data = xd.buildData()
	return widget.NewTree(
		xd.data.ChildIDs,
		xd.isBranch,
		xd.create,
		xd.update)

}
