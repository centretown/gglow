package fyui

import (
	"gglow/fyio"
	"strings"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type EffectSelect struct {
	*widget.SelectEntry
	selection binding.String
	effect    *fyio.EffectIo
	lookup    map[string]string
	options   []string
}

func NewEffectSelect(effect *fyio.EffectIo) *EffectSelect {
	fs := &EffectSelect{
		selection: binding.NewString(),
		effect:    effect,
		options:   make([]string, 0),
		lookup:    make(map[string]string),
	}

	fs.SelectEntry = widget.NewSelectEntry(effect.ListCurrent())
	fs.selection.Set(effect.EffectName())
	fs.Entry.Bind(fs.selection)
	fs.SelectEntry.OnChanged = fs.onChange

	effect.AddFrameListener(binding.NewDataListener(func() {
		fs.selection.Set(effect.EffectName())
	}))

	effect.AddFolderListener(binding.NewDataListener(func() {
		fs.selection.Set(effect.FolderName())
		fs.options = fs.effect.ListCurrent()
		fs.SelectEntry.SetOptions(fs.options)
		fs.buildLookup()
	}))
	return fs
}

func (fs *EffectSelect) buildLookup() {
	fs.lookup = make(map[string]string)
	for _, s := range fs.options {
		fs.lookup[strings.ToLower(s)] = s
	}
}

func (fs *EffectSelect) onChange(title string) {
	title, complete := fs.Parse(title)
	if !complete {
		return
	}
	if fs.effect.IsFolder(title) {
		fs.options = fs.effect.LoadFolder(title)
		fs.buildLookup()
		fs.SelectEntry.SetOptions(fs.options)
	} else {
		fs.effect.LoadEffect(title)
	}
}

func (fs *EffectSelect) Parse(title string) (result string, complete bool) {
	result = title
	length := len(result)
	if length < 1 {
		fs.SelectEntry.SetOptions(fs.options)
		return
	}

	search := strings.ToLower(result)
	result, complete = fs.lookup[search]
	if complete {
		fs.SelectEntry.SetOptions(fs.options)
		return
	}

	ls := make([]string, 0)
	for _, s := range fs.options {
		if strings.Contains(strings.ToLower(s), search) {
			ls = append(ls, s)
		}
	}
	fs.SelectEntry.SetOptions(ls)
	return
}
