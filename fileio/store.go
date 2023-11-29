package fileio

import (
	"fmt"
	"glow-gui/glow"
	"glow-gui/resources"
	"glow-gui/settings"
	"io"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/storage"
	"gopkg.in/yaml.v3"
)

const (
	scheme            = "file://"
	DefaultEffectPath = "/home/dave/src/glow-gui/cabinet/effects/"
	ExamplesPath      = "/home/dave/src/glow-gui/cabinet/examples/"
	dots              = ".."
)

func defaultFrame() (frame *glow.Frame) {
	frame = &glow.Frame{}
	frame.Interval = 48
	frame.Layers = append(frame.Layers, glow.Layer{})
	return
}

type Store struct {
	EffectName string
	Frame      binding.Untyped
	Layer      binding.Untyped
	LayerIndex int
	Current    fyne.ListableURI
	KeyList    []string
	FolderList []string

	uriMap      map[string]fyne.URI
	route       []string
	stack       *Stack
	preferences fyne.Preferences
	rootPath    string

	// history *History
	// CanUndo binding.Bool
	// CanSave binding.Bool
	IsDirty binding.Bool

	// changeDetected bool
}

func NewStore(preferences fyne.Preferences) *Store {

	rootPath := preferences.StringWithFallback(settings.EffectPath.String(),
		DefaultEffectPath)
	path := scheme + rootPath

	uri, err := storage.ParseURI(path)
	if err != nil {
		fyne.LogError(resources.MsgParseEffectPath.Format(path), err)
		os.Exit(1)
	}

	rootURI, err := storage.ListerForURI(uri)
	if err != nil {
		fyne.LogError(resources.MsgPathNotFolder.Format(path), err)
		os.Exit(1)
	}

	st := &Store{
		KeyList:    make([]string, 0),
		FolderList: make([]string, 0),
		Frame:      binding.NewUntyped(),
		Layer:      binding.NewUntyped(),
		IsDirty:    binding.NewBool(),
		// CanUndo:    binding.NewBool(),
		// CanSave:    binding.NewBool(),

		preferences: preferences,
		uriMap:      make(map[string]fyne.URI),
		stack:       NewStack(rootURI),
		rootPath:    rootPath,
		// history:     NewHistory(),
	}

	st.stack.Push(rootURI)
	st.makeLookupList()
	st.route = preferences.StringListWithFallback(settings.EffectRoute.String(),
		[]string{MakeTitle(rootURI)})

	for _, s := range st.route[1:] {
		st.RefreshKeys(s)
	}

	effect := preferences.StringWithFallback(settings.Effect.String(), "")
	if len(effect) > 0 {
		st.ReadEffect(effect)
	} else {
		st.setFrame(defaultFrame(), 0)
	}

	return st
}

func (st *Store) setFrame(frame *glow.Frame, layerIndex int) {
	st.Frame.Set(frame)
	st.SetCurrentLayer(layerIndex)
	st.IsDirty.Set(false)
}

func (st *Store) GetFrame() *glow.Frame {
	f, _ := st.Frame.Get()
	return f.(*glow.Frame)
}

func (st *Store) SetCurrentLayer(i int) {
	frame := st.GetFrame()
	var layer *glow.Layer
	if i < len(frame.Layers) {
		st.LayerIndex = i
		layer = &frame.Layers[i]
	} else {
		st.LayerIndex = 0
		layer = &glow.Layer{}
	}
	st.Layer.Set(layer)
}

func (st *Store) GetCurrentLayer() *glow.Layer {
	layer, _ := st.Layer.Get()
	return layer.(*glow.Layer)
}

func (st *Store) Undo(title string) {

	if !st.GetDirty() {
		fyne.LogError("UndoEffect", fmt.Errorf("nothing to undo"))
		return
	}

	frame := *st.GetFrame()
	st.setFrame(&frame, st.LayerIndex)
}

func (st *Store) SetDirty(dirty bool) {
	st.IsDirty.Set(dirty)
}

func (m *Store) GetDirty() bool {
	b, _ := m.IsDirty.Get()
	return b
}

func (st *Store) OnExit() {
	st.preferences.SetStringList(settings.EffectRoute.String(), st.route)
	st.preferences.SetString(settings.EffectPath.String(), st.rootPath)
}

func (st *Store) IsFolder(key string) bool {
	uri, ok := st.uriMap[key]
	if ok {
		ok, _ = storage.CanList(uri)
	}
	return ok
}

func (st *Store) RefreshKeys(key string) {

	uri, ok := st.uriMap[key]
	if !ok {
		err := fmt.Errorf(resources.MsgGetEffectLookup.Format(key))
		fyne.LogError("RefreshKeys", err)
		return
	}

	listable, err := storage.ListerForURI(uri)
	if err != nil {
		fyne.LogError(resources.MsgPathNotFolder.Format(key), err)
		return
	}

	if key == dots {
		st.stack.Pop()
	} else {
		st.stack.Push(listable)
	}

	st.makeLookupList()
}

func (st *Store) makeLookupList() (err error) {

	st.uriMap = make(map[string]fyne.URI)
	currentUri, isRoot := st.stack.Current()
	if !isRoot {
		st.uriMap[dots] = st.Current
	}
	st.Current = currentUri

	uriList, err := st.Current.List()
	st.KeyList = make([]string, 0, len(uriList)+1)
	st.FolderList = make([]string, 0)
	if !isRoot {
		st.KeyList = append(st.KeyList, dots)
		st.FolderList = append(st.FolderList, dots)
	}

	for _, uri := range uriList {
		title := MakeTitle(uri)
		st.uriMap[title] = uri
		st.KeyList = append(st.KeyList, title)
		isList, _ := storage.CanList(uri)
		if isList {
			st.FolderList = append(st.FolderList, title)
		}
	}
	return
}

func (st *Store) IsDuplicate(title string) error {
	_, found := st.uriMap[title]
	if found {
		return fmt.Errorf(resources.MsgDuplicate.String())
	}
	return nil
}

func (st *Store) CreateNewEffect(title string, frame *glow.Frame) error {
	err := st.IsDuplicate(title)
	if err != nil {
		return err
	}
	return st.WriteEffect(title, frame)
}

func (st *Store) CreateNewFolder(title string) error {
	err := st.IsDuplicate(title)
	if err != nil {
		return err
	}
	return st.WriteFolder(title)
}

func (st *Store) WriteFolder(title string) error {
	title = strings.TrimSpace(title)
	err := ValidateFolderName(title)
	if err != nil {
		return err
	}

	path := scheme + st.Current.Path() + "/" + title
	uri, err := storage.ParseURI(path)
	if err != nil {
		return err
	}

	err = storage.CreateListable(uri)
	return err
}

func (st *Store) ReadEffect(title string) error {

	frame := &glow.Frame{}

	uri, ok := st.uriMap[title]
	if !ok {
		err := fmt.Errorf(resources.MsgGetEffectLookup.Format(title))
		fyne.LogError("ReadEffect", err)
		return err
	}

	rdr, err := storage.Reader(uri)
	if err != nil {
		fyne.LogError(resources.MsgGetFrame.Format(title), err)
		return err
	}
	defer rdr.Close()

	buffer, err := io.ReadAll(rdr)
	if err != nil {
		fyne.LogError(resources.MsgGetFrame.Format(title), err)
		return err
	}

	err = yaml.Unmarshal(buffer, frame)
	if err != nil {
		fyne.LogError(resources.MsgGetFrame.Format(title), err)
		return err
	}

	st.route = st.stack.Route()
	st.EffectName = title
	st.setFrame(frame, 0)
	return nil
}

func (st *Store) WriteEffect(title string, frame *glow.Frame) error {
	title = strings.TrimSpace(title)

	err := ValidateEffectName(title)
	if err != nil {
		return err
	}

	path := scheme + st.Current.Path() + "/" + MakeFileName(title)
	uri, err := storage.ParseURI(path)
	if err != nil {
		return err
	}

	writable, err := storage.CanWrite(uri)
	if err != nil {
		return err
	}

	if !writable {
		err = fmt.Errorf(resources.MsgNotWritable.String())
		return err
	}

	wrt, err := storage.Writer(uri)
	if err != nil {
		return err
	}
	defer wrt.Close()

	buf, err := yaml.Marshal(frame)
	if err != nil {
		return err
	}

	_, err = wrt.Write(buf)
	if err != nil {
		return err
	}

	st.makeLookupList()
	st.IsDirty.Set(false)
	return err
}

func (st *Store) ValidateNewFolderName(title string) error {
	err := ValidateFolderName(title)
	if err != nil {
		return err
	}

	err = st.IsDuplicate(title)
	return err
}

func (st *Store) ValidateNewEffectName(title string) error {
	err := ValidateEffectName(title)
	if err != nil {
		return err
	}
	err = st.IsDuplicate(title)
	return err
}
