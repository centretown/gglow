package fileio

import (
	"fmt"
	"glow-gui/effects"
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
	effectName       string
	Frame            binding.Untyped
	Layer            binding.Untyped
	LayerSummaryList binding.StringList
	layerIndex       int
	Current          fyne.ListableURI
	keyList          []string
	FolderList       []string

	uriMap      map[string]fyne.URI
	route       []string
	stack       *Stack
	preferences fyne.Preferences
	rootPath    string

	// history *History
	// CanUndo binding.Bool
	// CanSave binding.Bool
	backup     *glow.Frame
	HasBackup  bool
	hasChanged binding.Bool
	isActive   bool

	saveActions []func(*glow.Frame)
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
		keyList:          make([]string, 0),
		FolderList:       make([]string, 0),
		Frame:            binding.NewUntyped(),
		Layer:            binding.NewUntyped(),
		LayerSummaryList: binding.NewStringList(),
		hasChanged:       binding.NewBool(),
		// CanUndo:    binding.NewBool(),
		// CanSave:    binding.NewBool(),

		preferences: preferences,
		uriMap:      make(map[string]fyne.URI),
		stack:       NewStack(rootURI),
		rootPath:    rootPath,
		backup:      &glow.Frame{},
		// history:     NewHistory(),
		saveActions: make([]func(*glow.Frame), 0),
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

	st.AddChangeListener(binding.NewDataListener(func() {
		if st.HasChanged() {
			f := st.GetFrame()
			fmt.Println("hasChanged makeBackup", f.Interval)
			st.makeBackup(true)
			// st.setFrame(&frame, st.LayerIndex)
		}
	}))

	return st
}

func (st *Store) SummaryList() []string {
	l, _ := st.LayerSummaryList.Get()
	return l
}

func (st *Store) LayerIndex() int {
	return st.layerIndex
}

func (st *Store) SetActive() {
	st.isActive = true
}

func (st *Store) setFrame(frame *glow.Frame, layerIndex int) {
	st.Frame.Set(frame)
	st.SetCurrentLayer(layerIndex)
	st.hasChanged.Set(false)
}

func (st *Store) GetFrame() *glow.Frame {
	f, _ := st.Frame.Get()
	return f.(*glow.Frame)
}

func (st *Store) SetCurrentLayer(i int) {
	frame := st.GetFrame()
	var layer *glow.Layer
	if i < len(frame.Layers) {
		st.layerIndex = i
		layer = &frame.Layers[i]
	} else {
		st.layerIndex = 0
		layer = &glow.Layer{}
	}
	st.Layer.Set(layer)
}

func (st *Store) GetCurrentLayer() *glow.Layer {
	layer, _ := st.Layer.Get()
	return layer.(*glow.Layer)
}

func (st *Store) Undo(title string) {

	if st.HasBackup {
		frame := st.backup
		st.setFrame(frame, st.layerIndex)
		st.makeBackup(false)
	}

	fyne.LogError("UndoEffect", fmt.Errorf("nothing to undo"))
}

func (st *Store) SetChanged() {
	if !st.isActive {
		return
	}
	st.hasChanged.Set(true)
}

func (st *Store) AddFrameListener(listener binding.DataListener) {
	st.Frame.AddListener(listener)
}

func (st *Store) AddLayerListener(listener binding.DataListener) {
	st.Layer.AddListener(listener)
}

func (st *Store) AddChangeListener(listener binding.DataListener) {
	st.hasChanged.AddListener(listener)
}

func (m *Store) HasChanged() bool {
	b, _ := m.hasChanged.Get()
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

func (st *Store) KeyList() []string {
	return st.keyList
}

func (st *Store) RefreshKeys(key string) []string {

	uri, ok := st.uriMap[key]
	if !ok {
		err := fmt.Errorf(resources.MsgGetEffectLookup.Format(key))
		fyne.LogError("RefreshKeys", err)
		return st.keyList
	}

	listable, err := storage.ListerForURI(uri)
	if err != nil {
		fyne.LogError(resources.MsgPathNotFolder.Format(key), err)
		return st.keyList
	}

	if key == dots {
		st.stack.Pop()
	} else {
		st.stack.Push(listable)
	}

	st.makeLookupList()
	return st.keyList
}

func (st *Store) makeLookupList() (err error) {

	st.uriMap = make(map[string]fyne.URI)
	currentUri, isRoot := st.stack.Current()
	if !isRoot {
		st.uriMap[dots] = st.Current
	}
	st.Current = currentUri

	uriList, err := st.Current.List()
	st.keyList = make([]string, 0, len(uriList)+1)
	st.FolderList = make([]string, 0)
	if !isRoot {
		st.keyList = append(st.keyList, dots)
		st.FolderList = append(st.FolderList, dots)
	}

	for _, uri := range uriList {
		title := MakeTitle(uri)
		st.uriMap[title] = uri
		st.keyList = append(st.keyList, title)
		isList, _ := storage.CanList(uri)
		if isList {
			st.FolderList = append(st.FolderList, title)
		}
	}
	return
}

func (st *Store) isDuplicate(title string) error {
	_, found := st.uriMap[title]
	if found {
		return fmt.Errorf(resources.MsgDuplicate.String())
	}
	return nil
}

func (st *Store) CreateNewEffect(title string, frame *glow.Frame) error {
	err := st.isDuplicate(title)
	if err != nil {
		return err
	}
	return st.writeEffect(title, frame)
}

func (st *Store) CreateNewFolder(title string) error {
	err := st.isDuplicate(title)
	if err != nil {
		return err
	}
	return st.WriteFolder(title)
}

func (st *Store) WriteFolder(title string) error {
	title = strings.TrimSpace(title)
	err := effects.ValidateFolderName(title)
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
	st.effectName = title
	st.setFrame(frame, 0)
	st.makeBackup(false)
	st.hasChanged.Set(false)
	return nil
}

func (st *Store) WriteEffect() error {
	return st.writeEffect(st.EffectName(), st.GetFrame())
}

func (st *Store) writeEffect(title string, frame *glow.Frame) error {
	title = strings.TrimSpace(title)

	err := effects.ValidateEffectName(title)
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
	st.hasChanged.Set(false)
	return err
}

func (st *Store) ValidateNewFolderName(title string) error {
	err := effects.ValidateFolderName(title)
	if err != nil {
		return err
	}

	err = st.isDuplicate(title)
	return err
}

func (st *Store) ValidateNewEffectName(title string) error {
	err := effects.ValidateEffectName(title)
	if err != nil {
		return err
	}
	err = st.isDuplicate(title)
	return err
}

func (st *Store) makeBackup(b bool) {
	st.HasBackup = b
	if b {
		st.backup, _ = glow.FrameDeepCopy(st.GetFrame())
	} else {
		st.backup = &glow.Frame{}
	}
}

func (st *Store) Apply() {
	frame := st.GetFrame()

	for _, saveAction := range st.saveActions {
		saveAction(frame)
	}

	current := *frame
	st.Frame.Set(&current)
}

func (st *Store) OnApply(f func(*glow.Frame)) {
	st.saveActions = append(st.saveActions, f)
}

func (st *Store) EffectName() string {
	return st.effectName
}

func (st *Store) UndoEffect() {
	st.Undo(st.effectName)
	st.Apply()
}

func (st *Store) CanUndo() bool {
	return st.HasBackup
}

func (st *Store) onChangeFrame() {
	frame := st.GetFrame()
	summaries := make([]string, 0, len(frame.Layers))
	for i, layer := range frame.Layers {
		summaries = append(summaries, effects.SummarizeLayer(&layer, i+1))
	}
	st.LayerSummaryList.Set(summaries)
}

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
