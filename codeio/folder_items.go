package codeio

import "gglow/glow"

type EffectItem struct {
	Title    string
	Constant string
	Frame    *glow.Frame
}

func NewEffectItem(title string, frame *glow.Frame) *EffectItem {
	ei := &EffectItem{
		Title: title,
		Frame: frame,
	}
	return ei
}

type FolderList struct {
	Title string
	List  []*EffectItem
}

func NewFolderList(title string, list []*EffectItem) *FolderList {
	fl := &FolderList{
		Title: title,
		List:  list,
	}
	return fl
}

func (fl *FolderList) AddItem(item *EffectItem) {
	fl.List = append(fl.List, item)
}
