package effects

import (
	"glow-gui/glow"

	"fyne.io/fyne/v2/data/binding"
)

type Effect interface {
	FolderName() string
	EffectName() string
	GetFrame() *glow.Frame
	KeyList() []string
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

	ValidateNewFolderName(title string) (err error)
	ValidateNewEffectName(title string) (err error)

	IsFolder(title string) bool
	CreateNewEffect(title string, frame *glow.Frame) (err error)
	CreateNewFolder(title string) (err error)

	RefreshFolder(title string) []string

	OnApply(f func(*glow.Frame))
	Apply()

	WriteEffect() (err error)
	ReadEffect(title string) (err error)

	CanUndo() bool
	UndoEffect()
}
