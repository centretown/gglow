package ui

import (
	"fmt"
	"gglow/fyglow/effectio"
	"gglow/iohandler"
	"gglow/text"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type ExportDialog struct {
	*dialog.CustomDialog
	effect  *effectio.EffectIo
	folders []string
	items   map[string][]string
	checks  map[string]bool
	tree    *widget.Tree
}

func NewExportDialogB(effect *effectio.EffectIo, window fyne.Window) *ExportDialog {
	xd := &ExportDialog{
		effect:  effect,
		checks:  make(map[string]bool),
		folders: make([]string, 0),
		items:   make(map[string][]string),
	}
	xd.tree = xd.buildTreeFromBool()
	xd.CustomDialog = dialog.NewCustom("Export", text.CancelLabel.String(), xd.tree, window)
	return xd
}

func NewExportDialog(effect *effectio.EffectIo, window fyne.Window) *ExportDialog {
	xd := &ExportDialog{
		effect:  effect,
		checks:  make(map[string]bool),
		folders: make([]string, 0),
		items:   make(map[string][]string),
	}
	xd.tree = xd.buildTree()
	xd.CustomDialog = dialog.NewCustom("Export", text.CancelLabel.String(), xd.tree, window)
	return xd
}

func (xd *ExportDialog) Start() {
	xd.CustomDialog.Show()
}

func (xd *ExportDialog) buildItems() {
	xd.folders = make([]string, 0)
	xd.items = make(map[string][]string)
	xd.checks = make(map[string]bool)
	ls, _ := xd.effect.ListFolder(iohandler.DOTS)
	for _, l := range ls {
		xd.folders = append(xd.folders, l.Key)
		xd.checks[l.Key] = false
	}

	for _, folder := range xd.folders {
		ls, _ := xd.effect.ListFolder(folder)
		itemList := make([]string, 0)
		for _, l := range ls {
			if l.Key != iohandler.DOTS {
				val := l.Value + "/" + l.Key
				itemList = append(itemList, val)
				xd.checks[val] = false
			}
		}

		xd.items[folder] = itemList
	}
}

func (xd *ExportDialog) buildTree() *widget.Tree {
	xd.buildItems()
	return widget.NewTree(
		xd.childUIDs,
		xd.isBranch,
		xd.create,
		xd.update)

}

func (xd *ExportDialog) childUIDs(id widget.TreeNodeID) []widget.TreeNodeID {
	if id == "" || id == ".." {
		return xd.folders
	}
	item, ok := xd.items[id]
	if ok {
		return item
	}
	return []string{}
}

func (xd *ExportDialog) isBranch(id widget.TreeNodeID) bool {
	if id == "" || id == ".." {
		return true
	}
	_, ok := xd.items[id]
	return ok
}

func (xd *ExportDialog) create(branch bool) fyne.CanvasObject {
	fmt.Println("CREATE", branch)
	return widget.NewCheck("Leaf template", func(b bool) {})
}

func (xd *ExportDialog) update(id string, branch bool, o fyne.CanvasObject) {
	check := o.(*widget.Check)
	check.SetText(id)
	check.OnChanged = func(b bool) {
		xd.checks[id] = b
	}
}

func (xd *ExportDialog) buildBoolTree() binding.BoolTree {
	tr := binding.NewBoolTree()
	folders, _ := xd.effect.ListFolder(iohandler.DOTS)
	fmt.Println(folders)
	for _, folder := range folders {
		tr.Append("", folder.Key, false)
		fmt.Println(folder.Key, "folder")
	}

	for _, folder := range folders {
		ls, _ := xd.effect.ListFolder(folder.Key)
		for _, l := range ls {
			val := l.Value + "/" + l.Key
			tr.Append(folder.Key, val, false)
			fmt.Println(val, "effect")
		}
	}
	return tr
}

func (xd *ExportDialog) buildTreeFromBool() *widget.Tree {
	data := xd.buildBoolTree()
	return widget.NewTree(
		data.ChildIDs,
		func(id widget.TreeNodeID) bool {
			children := data.ChildIDs(id)
			return len(children) > 0
		},
		// xd.isBranch,
		xd.create,
		func(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
			check := o.(*widget.Check)
			check.SetText(id)
			item, _ := data.GetItem(id)
			check.Bind(item.(binding.Bool))
		})

}

func (xd *ExportDialog) updateb(item binding.DataItem, branch bool, o fyne.CanvasObject) {
	var check *widget.Check = o.(*widget.Check)
	check.Bind(item.(binding.Bool))
	fmt.Println("UB", check)
}
