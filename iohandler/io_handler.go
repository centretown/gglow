package iohandler

import "gglow/glow"

const dots = ".."

func IsFolder(item string) bool {
	return item == dots
}

func AsFolder() string {
	return dots
}

const DRIVER_MYSQL = "mysql"
const DRIVER_SQLLITE3 = "sqlite3"
const DRIVER_POSTGRES = "postgres"
const DRIVER_CODE = "code"

type KeyValue struct {
	Key   string
	Value string
}

type OutHandler interface {
	Create(name string) error
	CreateEffect(folder, title string, frame *glow.Frame) error
	UpdateEffect(folder, title string, frame *glow.Frame) error
	CreateFolder(folder string) error
	OnExit() error
}

type InHandler interface {
	ReadEffect(folder, title string) (*glow.Frame, error)
	ListEffects(string) ([]string, error)
	ListKeys(string) ([]KeyValue, error)
	ListFolders() ([]string, error)
	OnExit() error
}

type IoHandler interface {
	InHandler
	OutHandler
}
