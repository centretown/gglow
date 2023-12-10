package effects

import "glow-gui/glow"

type IoHandler interface {
	WriteEffect(title string, frame *glow.Frame) error
	CreateNewEffect(title string, frame *glow.Frame) error
	ReadEffect(title string) (*glow.Frame, error)
	CreateNewFolder(title string) error
	WriteFolder(title string) error
	ValidateNewFolderName(title string) error
	ValidateNewEffectName(title string) error
	IsFolder(key string) bool
	KeyList() []string
	RefreshKeys(key string) ([]string, error)
	OnExit()
}
