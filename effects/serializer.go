package effects

import (
	"encoding/json"
	"glow-gui/glow"
	"strings"

	"fyne.io/fyne/v2"
	"gopkg.in/yaml.v3"
)

type Serializer interface {
	Scan(buffer []byte, frame *glow.Frame) (err error)
	Format(frame *glow.Frame) (buffer []byte, err error)
	MakeFileName(title string) string
}

type YamlSerializer struct {
}

func (yml *YamlSerializer) Scan(buffer []byte, frame *glow.Frame) (err error) {
	err = yaml.Unmarshal(buffer, frame)
	return
}

func (yml *YamlSerializer) Format(frame *glow.Frame) (buffer []byte, err error) {
	buffer, err = yaml.Marshal(frame)
	return
}

func (yml *YamlSerializer) MakeFileName(title string) string {
	s := strings.ReplaceAll(title, " ", "_")
	s += ".yaml"
	return s
}

type JsonSerializer struct {
}

func (jsn *JsonSerializer) Scan(buffer []byte, frame *glow.Frame) (err error) {
	err = json.Unmarshal(buffer, frame)
	return
}

func (jsn *JsonSerializer) Format(frame *glow.Frame) (buffer []byte, err error) {
	buffer, err = json.Marshal(frame)
	return
}

func (jsn *JsonSerializer) MakeFileName(title string) string {
	s := strings.ReplaceAll(title, " ", "_")
	s += ".json"
	return s
}

func UriSerializer(uri fyne.URI) Serializer {
	switch uri.Extension() {
	case ".yaml":
		return &YamlSerializer{}
	}
	return &JsonSerializer{}
}
