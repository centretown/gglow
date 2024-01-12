package iohandler

import "gglow/glow"

const Dots = ".."

type OutHandler interface {
	Create(name string) error
	WriteEffect(title string, frame *glow.Frame) error
	WriteFolder(title string) error
	SetCurrentFolder(key string) ([]string, error)
	OnExit() error
}

type InHandler interface {
	FolderName() string
	EffectName() string
	ReadEffect(title string) (*glow.Frame, error)
	ReadFolder(string) ([]string, error)

	ValidateNewFolderName(title string) error
	ValidateNewEffectName(title string) error
	IsFolder(key string) bool
	IsRootFolder() bool

	ListCurrent() []string
	SetCurrentFolder(key string) ([]string, error)
	SetRootCurrent() ([]string, error)

	OnExit() error
}

type IoHandler interface {
	InHandler
	OutHandler
	CreateNewEffect(title string, frame *glow.Frame) error
	CreateNewFolder(title string) error
}
