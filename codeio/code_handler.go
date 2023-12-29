package codeio

import (
	"fmt"
	"gglow/glow"
	"gglow/iohandler"
	"io/fs"
	"os"
)

var _ iohandler.OutHandler = (*CodeHandler)(nil)

var emptyList []string

type CodeHandler struct {
	path   string
	folder string
	title  string

	folders        []*FolderList
	currentFolder  *FolderList
	currentEffects []*EffectItem
}

func NewCodeHandler(path string) (*CodeHandler, error) {
	ch := &CodeHandler{
		path:    path,
		folders: make([]*FolderList, 0),
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
	ch.title = title
	ch.currentFolder.AddItem(NewEffectItem(title, frame))
	return nil
}

func (ch *CodeHandler) WriteFolder(title string) error {
	ch.folder = title
	ch.currentFolder = NewFolderList(title, ch.currentEffects)
	ch.folders = append(ch.folders, ch.currentFolder)
	return nil
}

func (ch *CodeHandler) FolderName() string {
	return ch.folder
}

func (ch *CodeHandler) SetFolder(key string) ([]string, error) {
	ch.folder = key
	return emptyList, nil
}

func (ch *CodeHandler) OnExit() {

	err := os.Chdir(ch.path)
	if err != nil {
		fmt.Println(err)
		return
	}

	header := NewHeaderGenerator()
	err = header.Open("catalog.h")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer header.Close()
	err = header.Write(ch.folders)
	if err != nil {
		fmt.Println(err)
		return
	}

	source := NewSourceGenerator()
	err = source.Open("catalog.cpp")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer source.Close()
	err = source.Write(ch.folders)
	if err != nil {
		fmt.Println(err)
		return
	}

	esphome := NewEffectGenerator()
	err = esphome.Open("grid_effects.yml")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer esphome.Close()
	err = esphome.Write(ch.folders)
	if err != nil {
		fmt.Println(err)
		return
	}
}
