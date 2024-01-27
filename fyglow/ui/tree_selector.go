package ui

import (
	"fmt"
	"gglow/fyglow/effectio"
	"strings"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func NewTreeSelector(effect *effectio.EffectIo, data binding.BoolTree) *widget.Tree {
	var auto bool
	tree := NewEffectTree(data, CreateLabel, UpdateLabel(data))
	tree.OnSelected = OnTreeSelected(effect, data, &auto)

	effect.AddFolderListener(binding.NewDataListener(func() {
		auto = true
		tree.Select(effect.FolderName())
	}))

	effect.AddFrameListener(binding.NewDataListener(func() {
		auto = true
		tree.Select(effect.FolderName() + separator + effect.EffectName())
	}))
	return tree
}

func OnTreeSelected(effect *effectio.EffectIo, data binding.BoolTree, auto *bool) func(widget.TreeNodeID) {
	return func(uid widget.TreeNodeID) {
		if *auto {
			*auto = false
			return
		}
		segment := strings.Split(uid, separator)
		if len(segment) < 1 {
			return
		}
		fmt.Println("loadfolder")
		effect.LoadFolder(segment[0])
		if len(segment) < 2 {
			return
		}
		fmt.Println("loadeffect")
		effect.LoadEffect(segment[1])
	}
}
