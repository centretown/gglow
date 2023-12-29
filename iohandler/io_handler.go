package iohandler

import "gglow/glow"

type OutHandler interface {
	Create(name string) error
	WriteEffect(title string, frame *glow.Frame) error
	WriteFolder(title string) error
	FolderName() string
	SetFolder(key string) ([]string, error)
	OnExit()
}

type InHandler interface {
	FolderName() string
	EffectName() string
	ReadEffect(title string) (*glow.Frame, error)
	SetFolder(key string) ([]string, error)
	ValidateNewFolderName(title string) error
	ValidateNewEffectName(title string) error
	IsFolder(key string) bool
	ListCurrentFolder() []string
	RootFolder() ([]string, error)
	OnExit()
}

type IoHandler interface {
	OutHandler
	CreateNewEffect(title string, frame *glow.Frame) error
	CreateNewFolder(title string) error
	InHandler
}
