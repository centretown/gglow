package fyui

import (
	"fmt"
	"gglow/fyio"
	"strings"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type EffectSelect struct {
	*widget.SelectEntry
	selection binding.String
	// folderName binding.String
	// effectName binding.String
	effect  *fyio.EffectIo
	lookup  map[string]string
	options []string
	// auto       bool
}

func NewEffectSelect(effect *fyio.EffectIo) *widget.SelectEntry {
	fs := &EffectSelect{
		// folderName: binding.NewString(),
		// effectName: binding.NewString(),
		selection: binding.NewString(),
		effect:    effect,
		options:   make([]string, 0),
		lookup:    make(map[string]string),
	}
	// fs.selection = binding.NewSprintf("%s/%s", fs.folderName, fs.effectName)
	// fs.Select = widget.NewSelect(effect.ListCurrentFolder(),
	// 	fs.onChange)

	fs.SelectEntry = widget.NewSelectEntry(effect.ListCurrentFolder())
	fs.selection.Set(effect.EffectName())
	fs.Entry.Bind(fs.selection)
	fs.SelectEntry.OnChanged = fs.onChange

	effect.AddFrameListener(binding.NewDataListener(func() {
		fmt.Println("AddFrameListener")
		fs.selection.Set(effect.EffectName())
		// fs.effectName.Set(effect.EffectName())
		// selected := fs.SelectEntry.Selected
		// if selected != effect.EffectName() {
		// fs.auto = true
		// fs.SelectEntry.SetSelected(effect.EffectName())
		// }
	}))

	effect.AddFolderListener(binding.NewDataListener(func() {
		fmt.Println("AddFolderListener")
		fs.selection.Set(effect.FolderName())
		fs.options = fs.effect.ListCurrentFolder()
		fs.SelectEntry.SetOptions(fs.options)
		fs.buildLookup()
		// fs.selection.Set(effect.FolderName() + "/" + effect.EffectName())
	}))
	return fs.SelectEntry
}

func (fs *EffectSelect) buildLookup() {
	fs.lookup = make(map[string]string)
	for _, s := range fs.options {
		fs.lookup[strings.ToLower(s)] = s
	}
}

func (fs *EffectSelect) onChange(title string) {
	title, ok := fs.Parse(title)
	if !ok {
		return
	}
	fs.selection.Set(title)
	fmt.Println("onChange")

	if fs.effect.IsFolder(title) {
		fs.options = fs.effect.LoadFolder(title)
		fs.buildLookup()
		fs.SelectEntry.SetOptions(fs.options)
	} else {
		fs.effect.LoadEffect(title)
	}
}

func (fs *EffectSelect) Parse(title string) (string, bool) {
	length := len(title)
	if length < 1 {
		fs.SelectEntry.SetOptions(fs.options)
		return title, false
	}

	search := strings.ToLower(title)
	if result, ok := fs.lookup[search]; ok {
		fmt.Println("result", result)
		return result, ok
	}

	ls := make([]string, 0)
	for _, s := range fs.options {
		if strings.HasPrefix(strings.ToLower(s), search) {
			ls = append(ls, s)
		}
	}
	fs.SelectEntry.SetOptions(ls)
	return title, false
}
