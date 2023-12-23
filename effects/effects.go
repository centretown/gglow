package effects

import (
	"gglow/glow"
	"gglow/iohandler"

	"fyne.io/fyne/v2/data/binding"
)

type Effect interface {
	iohandler.IoHandler
	FolderName() string
	EffectName() string
	GetFrame() *glow.Frame
	ListCurrentFolder() []string
	SummaryList() []string

	SetActive()
	GetCurrentLayer() *glow.Layer
	SetCurrentLayer(i int)
	LayerIndex() int

	AddFrameListener(listener binding.DataListener)
	AddLayerListener(listener binding.DataListener)
	AddChangeListener(listener binding.DataListener)

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
