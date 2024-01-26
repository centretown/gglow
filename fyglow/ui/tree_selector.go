package ui

import (
	"gglow/fyglow/effectio"
	"strings"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func NewTreeSelector(effect *effectio.EffectIo) *widget.Tree {
	data := BuildBoolTree(effect)
	tree := NewEffectTree(data, CreateLabel, UpdateLabel(data))
	tree.OnSelected = OnTreeSelected(effect, data)
	return tree
}

func OnTreeSelected(effect *effectio.EffectIo, data binding.BoolTree) func(widget.TreeNodeID) {
	return func(uid widget.TreeNodeID) {
		segment := strings.Split(uid, separator)
		if len(segment) < 1 {
			return
		}
		effect.LoadFolder(segment[0])
		if len(segment) < 2 {
			return
		}
		effect.LoadEffect(segment[1])
	}
}
