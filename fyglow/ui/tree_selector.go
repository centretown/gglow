package ui

import (
	"gglow/fyglow/effectio"

	"fyne.io/fyne/v2/widget"
)

func NewTreeSelector(effect *effectio.EffectIo) *widget.Tree {
	data := effect.TreeData()
	tree := NewEffectTree(data, CreateLabel, UpdateLabel(data))
	tree.OnSelected = OnTreeSelected(effect, tree)

	// effect.AddFolderListener(binding.NewDataListener(func() {
	// 	tree.Select(effect.FolderName())
	// }))

	// effect.AddFrameListener(binding.NewDataListener(func() {
	// 	tree.Select(effect.FolderName() + effectio.PathSeparator + effect.EffectName())
	// }))
	// tree.OnUnselected = func(uid widget.TreeNodeID) {
	// 	if tree.IsBranch(uid) && uid != effect.FolderName() {
	// 		tree.CloseBranch(uid)
	// 	}
	// }
	return tree
}

func OnTreeSelected(effect *effectio.EffectIo, tree *widget.Tree) func(widget.TreeNodeID) {
	return func(uid widget.TreeNodeID) {
		if tree.IsBranch(uid) {
			tree.OpenBranch(uid)
			return
		}

		effect.Select(uid)
	}
}
