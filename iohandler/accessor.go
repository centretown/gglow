package iohandler

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Accessor struct {
	Driver   string
	Path     string
	User     string
	Password string
	Host     string
	Port     string
	Database string
	Folder   string
	Effect   string
}

type AccessorView struct {
	Driver   string
	Path     string
	Host     string
	Port     string
	Database string
}

func NewAccessor() *Accessor {
	cf := &Accessor{}
	return cf
}

func NewAccessorView(accessor *Accessor) *AccessorView {
	av := &AccessorView{
		Driver:   accessor.Driver,
		Path:     accessor.Path,
		Host:     accessor.Host,
		Port:     accessor.Port,
		Database: accessor.Database,
	}
	return av
}

func LoadAccessor(path string) (accessor *Accessor, err error) {
	var rdr *os.File
	rdr, err = os.Open(path)
	if err != nil {
		return
	}
	defer rdr.Close()

	var buf []byte
	buf, err = io.ReadAll(rdr)
	if err != nil {
		return
	}

	accessor = &Accessor{}
	err = yaml.Unmarshal(buf, accessor)
	return
}

func SaveAccessor(path string, accessor *Accessor) (err error) {
	var buf []byte
	buf, err = yaml.Marshal(accessor)
	if err != nil {
		return
	}
	err = os.WriteFile(path, buf, os.ModePerm)
	return
}
