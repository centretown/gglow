package ui

import (
	"gglow/action"
	"gglow/fyglow/effectio"
	"gglow/iohandler"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func BuildAction(data binding.BoolTree, effect *effectio.EffectIo, drivers []string, path string) (act *action.Action) {
	act = action.NewAction()
	act.Method = "clone"
	act.Input = effect.Accessor

	for _, driver := range drivers {
		output := &iohandler.Accessor{
			Driver:   driver,
			Path:     path,
			Database: path,
		}
		act.Outputs = append(act.Outputs, output)
	}

	folders := data.ChildIDs(binding.DataTreeRootID)
	for _, id := range folders {
		if filter, ok := SelectActionFilter(data, id); ok {
			act.FilterItems = append(act.FilterItems, filter)
		}
	}

	return
}

func SelectActionFilter(data binding.BoolTree, folder string) (item *action.FilterItem, selected bool) {
	effects := make([]string, 0)
	selected, _ = data.GetValue(folder)
	for _, id := range data.ChildIDs(folder) {
		if val, _ := data.GetValue(id); val {
			effects = append(effects, strings.TrimPrefix(id, folder+effectio.PathSeparator))
		}
	}

	if len(effects) > 0 {
		selected = true
	}

	if selected {
		item = &action.FilterItem{Folder: folder, Effects: effects}
	}
	return
}

func ConfirmView(act *action.Action) fyne.CanvasObject {
	rich := widget.NewRichTextWithText(act.NewActionView())
	scroll := container.NewScroll(rich)
	return scroll
}

func WrapVertical(top, bottom fyne.CanvasObject) fyne.CanvasObject {
	return container.NewBorder(
		container.NewBorder(nil, widget.NewSeparator(), nil, nil, top),
		nil, nil, nil, bottom)
}

func NewEffectTree(data binding.DataTree,
	create func(branch bool) fyne.CanvasObject,
	update func(widget.TreeNodeID, bool, fyne.CanvasObject)) *widget.Tree {
	return widget.NewTree(data.ChildIDs,
		IsBranch(data),
		create,
		update)
}

func NewEffectTreeWithListener(data binding.DataTree,
	listener binding.DataListener,
	create func(branch bool) fyne.CanvasObject,
	update func(widget.TreeNodeID, bool, fyne.CanvasObject)) *widget.Tree {
	return widget.NewTree(data.ChildIDs,
		IsBranch(data),
		create,
		update)
}

func IsBranch(data binding.DataTree) func(string) bool {
	return func(id string) bool {
		return len(data.ChildIDs(id)) > 0
	}
}

func CreateCheck(branch bool) fyne.CanvasObject {
	return widget.NewCheck("NewCheck template", func(b bool) {})
}

func UpdateCheck(data binding.DataTree,
	listener binding.DataListener) func(widget.TreeNodeID, bool, fyne.CanvasObject) {

	return func(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
		check := o.(*widget.Check)
		check.SetText(id)
		i, _ := data.GetItem(id)
		bb := i.(binding.Bool)
		check.Bind(bb)
		bb.AddListener(listener)
	}
}

func CreateLabel(branch bool) fyne.CanvasObject {
	return widget.NewLabel("NewLabel template")
}

func UpdateLabel(data binding.DataTree) func(widget.TreeNodeID, bool, fyne.CanvasObject) {
	return func(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
		label := o.(*widget.Label)
		label.SetText(id)
	}
}
