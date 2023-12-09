package fields

import (
	"glow-gui/glow"

	"fyne.io/fyne/v2/data/binding"
)

type Model interface {
	KeyList() []string
	SummaryList() []string

	SetActive()
	EffectName() string
	GetFrame() *glow.Frame
	GetCurrentLayer() *glow.Layer
	SetCurrentLayer(i int)
	LayerIndex() int

	RefreshKeys(title string) []string
	ValidateNewFolderName(title string) (err error)
	ValidateNewEffectName(title string) (err error)

	AddFrameListener(listener binding.DataListener)
	AddLayerListener(listener binding.DataListener)
	AddChangeListener(listener binding.DataListener)

	SetChanged()
	HasChanged() bool

	IsFolder(title string) bool
	CreateNewEffect(title string, frame *glow.Frame) (err error)
	CreateNewFolder(title string) (err error)

	OnApply(f func(*glow.Frame))
	Apply()

	WriteEffect() (err error)
	ReadEffect(title string) (err error)

	CanUndo() bool
	UndoEffect()
}
