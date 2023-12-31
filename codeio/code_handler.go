package codeio

import (
	"gglow/glow"
	"gglow/iohandler"
	"io/fs"
	"os"
	"path/filepath"
)

var _ iohandler.OutHandler = (*CodeHandler)(nil)

var emptyList []string

type CodeHandler struct {
	path           string
	folders        []*iohandler.EffectItems
	currentFolder  *iohandler.EffectItems
	currentEffects []*iohandler.EffectItem
}

func NewCodeHandler(path string) (*CodeHandler, error) {
	ch := &CodeHandler{
		path:    path,
		folders: make([]*iohandler.EffectItems, 0),
	}

	return ch, nil
}

func (ch *CodeHandler) Create(path string) (err error) {
	var info fs.FileInfo
	info, err = os.Stat(path)
	if err == nil && info.IsDir() {
		return
	}
	if err != nil {
		err = os.Mkdir(path, os.ModePerm)
	}
	return err
}

func (ch *CodeHandler) WriteEffect(title string, frame *glow.Frame) error {
	ch.currentFolder.AddItem(iohandler.NewEffectItem(title, frame))
	return nil
}

func (ch *CodeHandler) WriteFolder(title string) error {
	ch.currentEffects = make([]*iohandler.EffectItem, 0)
	ch.currentFolder = iohandler.NewFolderList(title, ch.currentEffects)
	ch.folders = append(ch.folders, ch.currentFolder)
	return nil
}

func (ch *CodeHandler) SetFolder(key string) ([]string, error) {
	return emptyList, nil
}

func (ch *CodeHandler) OnExit() (err error) {
	return ch.process()
}

func (ch *CodeHandler) process() (err error) {
	generate := func(gen iohandler.Generator, fileName string) (err error) {
		err = gen.Open(filepath.Join(ch.path, fileName))
		if err != nil {
			return
		}
		defer gen.Close()
		err = gen.Write(ch.folders)
		return
	}

	err = generate(NewHeaderGenerator(), "catalog.h")
	if err == nil {
		err = generate(NewSourceGenerator(), "catalog.cpp")
	}
	if err == nil {
		err = generate(NewEffectGenerator(), "catalog_effects.yml")
	}
	return
}
