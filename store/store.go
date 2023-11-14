package store

import (
	"fmt"
	"glow-gui/glow"
	"glow-gui/resources"
	"glow-gui/settings"
	"io"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"gopkg.in/yaml.v3"
)

const (
	scheme            = "file://"
	DefaultEffectPath = "/home/dave/src/glow-gui/storage/effects/"
	ExamplesPath      = "/home/dave/src/glow-gui/storage/effects/examples/"
	dots              = ".."
)

type Store struct {
	Current    fyne.ListableURI
	Derived    fyne.ListableURI
	KeyList    []string
	FolderList []string

	uriMap      map[string]fyne.URI
	route       []string
	stack       *Stack
	preferences fyne.Preferences
	rootPath    string
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

	store := &Store{
		preferences: preferences,
		uriMap:      make(map[string]fyne.URI),
		stack:       NewStack(rootURI),
		rootPath:    rootPath,
		KeyList:     make([]string, 0),
		FolderList:  make([]string, 0),
	}

	store.stack.Push(rootURI)
	store.makeLookupList()
	store.route = preferences.StringListWithFallback(settings.EffectRoute.String(),
		[]string{MakeTitle(rootURI)})

	for _, s := range store.route[1:] {
		store.RefreshKeys(s)
	}

	return store
}

func (store *Store) OnExit() {
	store.preferences.SetStringList(settings.EffectRoute.String(), store.route)
	store.preferences.SetString(settings.EffectPath.String(), store.rootPath)
}

func (store *Store) IsFolder(key string) bool {
	uri, ok := store.uriMap[key]
	if ok {
		ok, _ = storage.CanList(uri)
	}
	return ok
}

func (store *Store) RefreshKeys(key string) {

	uri, ok := store.uriMap[key]
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
		store.stack.Pop()
	} else {
		store.stack.Push(listable)
	}

	store.makeLookupList()
}

func (store *Store) makeLookupList() (err error) {

	store.uriMap = make(map[string]fyne.URI)
	currentUri, isRoot := store.stack.Current()
	if !isRoot {
		store.uriMap[dots] = store.Current
	}
	store.Current = currentUri

	uriList, err := store.Current.List()
	store.KeyList = make([]string, 0, len(uriList)+1)
	store.FolderList = make([]string, 0)
	if !isRoot {
		store.KeyList = append(store.KeyList, dots)
		store.FolderList = append(store.FolderList, dots)
	}

	for _, uri := range uriList {
		title := MakeTitle(uri)
		store.uriMap[title] = uri
		store.KeyList = append(store.KeyList, title)
		isList, _ := storage.CanList(uri)
		if isList {
			store.FolderList = append(store.FolderList, title)
		}
	}
	return
}

func (store *Store) LoadFrame(key string, frame *glow.Frame) error {

	uri, ok := store.uriMap[key]
	if !ok {
		err := fmt.Errorf(resources.MsgGetEffectLookup.Format(key))
		fyne.LogError("LoadFrame", err)
		return err
	}

	rdr, err := storage.Reader(uri)
	if err != nil {
		fyne.LogError(resources.MsgGetFrame.Format(key), err)
		return err
	}
	defer rdr.Close()

	buffer, err := io.ReadAll(rdr)
	if err != nil {
		fyne.LogError(resources.MsgGetFrame.Format(key), err)
		return err
	}

	err = yaml.Unmarshal(buffer, frame)
	if err != nil {
		fyne.LogError(resources.MsgGetFrame.Format(key), err)
		return err
	}

	store.route = store.stack.Route()
	return nil
}

func (store *Store) StoreFrame(fname string, frame *glow.Frame) (err error) {
	var (
		uri   fyne.URI
		wrt   fyne.URIWriteCloser
		buf   []byte
		count int
	)

	buf, err = yaml.Marshal(frame)
	if err != nil {
		return
	}

	uri, err = storage.ParseURI(scheme + fname)
	if err != nil {
		return
	}

	wrt, err = storage.Writer(uri)
	if err != nil {
		return
	}
	defer wrt.Close()

	count, err = wrt.Write(buf)
	if err != nil {
		return
	}

	if count == 0 {
		err = fmt.Errorf("StoreFrame: zero bytes written")
	}

	return
}

func (store *Store) ValidateNewEffectName(s string) (err error) {
	s = strings.TrimSpace(s)

	if len(s) < 1 {
		err = fmt.Errorf(resources.MsgRequired.String())
		return
	}

	isAlpha := func(c rune) bool {
		return (c >= 'A' && c <= 'Z') ||
			(c >= 'a' && c <= 'z')
	}

	isNumeric := func(c rune) bool {
		return (c >= '0' && c <= '9')
	}

	isAlphaNumeric := func(c rune) bool {
		return isAlpha(c) || isNumeric(c)
	}

	for i, c := range s {
		if i == 0 && !isAlpha(c) {
			err = fmt.Errorf(resources.MsgFirstAlpha.String())
			return
		}

		if c != ' ' && !isAlphaNumeric(c) {
			err = fmt.Errorf(resources.MsgAlphaNumeric.String())
			return
		}
	}

	_, ok := store.uriMap[s]
	if ok {
		err = fmt.Errorf(resources.MsgDuplicate.String())
		return
	}

	return
}

func (store *Store) CreateNewEffect(title string, frame *glow.Frame) (err error) {
	err = store.ValidateNewEffectName(title)
	if err != nil {
		return
	}

	path := scheme + store.Current.Path() + "/" + MakeFileName(title)
	fmt.Println(path)
	uri, err := storage.ParseURI(path)
	if err != nil {
		return
	}

	writable, err := storage.CanWrite(uri)
	if err != nil {
		return
	}

	if !writable {
		err = fmt.Errorf(resources.MsgNotWritable.String())
		return
	}

	wrt, err := storage.Writer(uri)
	if err != nil {
		return
	}
	defer wrt.Close()

	buf, err := yaml.Marshal(frame)
	n, err := wrt.Write(buf)
	if err != nil {
		return
	}

	store.makeLookupList()

	fmt.Println("bytes written", n)
	return
}
