package effects

import "glow-gui/glow"

type IoHandler interface {
	CreateNewDatabase() error
	FolderName() string
	EffectName() string
	ReadEffect(title string) (*glow.Frame, error)
	IsFolder(key string) bool
	KeyList() []string
	Refresh() ([]string, error)
	RefreshFolder(key string) ([]string, error)
	WriteEffect(title string, frame *glow.Frame) error
	WriteFolder(title string) error
	CreateNewEffect(title string, frame *glow.Frame) error
	CreateNewFolder(title string) error
	ValidateNewFolderName(title string) error
	ValidateNewEffectName(title string) error
	OnExit()
}
