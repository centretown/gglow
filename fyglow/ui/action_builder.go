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
	"gopkg.in/yaml.v3"
)

const separator = "/"

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
		if filter, ok := SelectBoolFilter(data, id); ok {
			act.FilterItems = append(act.FilterItems, filter)
		}
	}

	return
}

func BuildBoolTree(effect *effectio.EffectIo) binding.BoolTree {
	var data binding.BoolTree = binding.NewBoolTree()
	folders, _ := effect.ListFolder(iohandler.DOTS)
	for _, folder := range folders {
		data.Append(binding.DataTreeRootID, folder.Key, false)
	}

	for _, folder := range folders {
		ls, _ := effect.ListFolder(folder.Key)
		for _, l := range ls {
			if l.Key != iohandler.DOTS {
				val := l.Value + separator + l.Key
				data.Append(folder.Key, val, false)
			}
		}
	}

	return data
}

func ConfirmObject(act *action.Action) fyne.CanvasObject {
	buf, _ := yaml.Marshal(act)
	seg := string(buf)
	rich := widget.NewRichTextWithText(seg)
	scroll := container.NewScroll(rich)
	return scroll
}

func ShowActionResultsObject(act *action.Action) fyne.CanvasObject {
	buf, _ := yaml.Marshal(act)
	seg := string(buf)
	rich := widget.NewRichTextWithText(seg)
	scroll := container.NewScroll(rich)
	return scroll
}

func WrapVertical(top, bottom fyne.CanvasObject) fyne.CanvasObject {
	return container.NewBorder(
		container.NewBorder(nil, widget.NewSeparator(), nil, nil, top),
		nil, nil, nil, bottom)
}

func SelectBoolFilter(data binding.BoolTree, folder string) (item *action.FilterItem, selected bool) {
	effects := make([]string, 0)
	selected, _ = data.GetValue(folder)
	for _, id := range data.ChildIDs(folder) {
		if val, _ := data.GetValue(id); val {
			effects = append(effects, strings.TrimPrefix(id, folder+separator))
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
