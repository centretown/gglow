package iohandler

import (
	"gglow/glow"
)

type EffectIoHandler interface {
	IoHandler
	FolderName() string
	EffectName() string
	GetFrame() *glow.Frame
	ListCurrentFolder() []string
	SummaryList() []string

	SetActive()
	GetCurrentLayer() *glow.Layer
	SetCurrentLayer(i int)
	LayerIndex() int

	AddFrameListener(listener interface{})
	AddLayerListener(listener interface{})
	AddChangeListener(listener interface{})

	SetChanged()
	HasChanged() bool

	LoadFolder(title string) []string

	OnApply(f func(*glow.Frame))
	Apply()

	SaveEffect() (err error)
	LoadEffect(title string) (err error)

	CanUndo() bool
	UndoEffect()
}
