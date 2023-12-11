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
	FileName(title string) string
}

type YamlSerializer struct {
}

func (yml *YamlSerializer) Scan(buffer []byte, frame *glow.Frame) error {
	return yaml.Unmarshal(buffer, frame)
}

func (yml *YamlSerializer) Format(frame *glow.Frame) ([]byte, error) {
	return yaml.Marshal(frame)
}

func (yml *YamlSerializer) FileName(title string) string {
	return strings.ReplaceAll(title, " ", "_") + ".yaml"
}

type JsonSerializer struct {
}

func (jsn *JsonSerializer) Scan(buffer []byte, frame *glow.Frame) error {
	return json.Unmarshal(buffer, frame)
}

func (jsn *JsonSerializer) Format(frame *glow.Frame) ([]byte, error) {
	return json.Marshal(frame)
}

func (jsn *JsonSerializer) FileName(title string) string {
	return strings.ReplaceAll(title, " ", "_") + ".json"
}

func UriSerializer(uri fyne.URI) Serializer {
	switch uri.Extension() {
	case ".yaml":
		return &YamlSerializer{}
	}
	return &JsonSerializer{}
}
