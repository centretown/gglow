package store

import (
	"strings"

	"fyne.io/fyne/v2"
)

func MakeTitle(uri fyne.URI) (s string) {
	s = uri.Name()
	i := strings.Index(s, uri.Extension())
	if i > 0 {
		s = s[:i]
	}
	s = strings.ReplaceAll(s, "_", " ")
	return
}

func MakeFileName(title string) string {
	s := strings.ReplaceAll(title, " ", "_")
	s += ".yaml"
	return s
}
