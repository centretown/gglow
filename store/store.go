package store

import (
	"fmt"
	"glow-gui/glow"
	"glow-gui/resources"
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
	Current     fyne.ListableURI
	Derived     fyne.ListableURI
	uriMap      map[string]fyne.URI
	keyList     []string
	folderList  []string
	route       []string
	stack       *Stack
	preferences fyne.Preferences
	basePath    string
}

func NewStore(preferences fyne.Preferences) *Store {

	basePath := preferences.StringWithFallback(resources.EffectPath.String(),
		DefaultEffectPath)
	path := scheme + basePath

	uri, err := storage.ParseURI(path)
	if err != nil {
		fyne.LogError(resources.MsgParseEffectPath.Format(path), err)
		os.Exit(1)
	}

	listable, err := storage.ListerForURI(uri)
	if err != nil {
		fyne.LogError(resources.MsgPathNotFolder.Format(path), err)
		os.Exit(1)
	}

	store := &Store{
		uriMap:      make(map[string]fyne.URI),
		stack:       NewStack(listable),
		preferences: preferences,
		basePath:    basePath,
		keyList:     make([]string, 0),
		folderList:  make([]string, 0),
	}

	store.stack.Push(listable)
	store.makeLookupList()
	store.route = preferences.StringListWithFallback(resources.EffectRoute.String(),
		[]string{makeTitle(listable)})

	for _, s := range store.route[1:] {
		store.RefreshLookupList(s)
	}

	return store
}

func (store *Store) LookUpList() []string {
	return store.keyList
}

func (store *Store) OnExit() {
	store.preferences.SetStringList(resources.EffectRoute.String(), store.route)
	store.preferences.SetString(resources.EffectPath.String(), store.basePath)
}

// func formatMessage(id resources.MessageID,
// 	key string, err error) error {
// 	id.FormatMessage(key, err)
// 	return fmt.Errorf("%s %s %v", id, key, err)
// }

func (store *Store) IsFolder(key string) bool {
	uri, ok := store.uriMap[key]
	if ok {
		ok, _ = storage.CanList(uri)
	}
	return ok
}

func (store *Store) RefreshLookupList(key string) {

	uri, ok := store.uriMap[key]
	if !ok {
		fyne.LogError("RefreshLookupList", fmt.Errorf("%s not found", key))
		return
	}

	listable, err := storage.ListerForURI(uri)
	if err != nil {
		fyne.LogError("RefreshLookupList",
			fmt.Errorf("%s not listable %v", key, err))
		return
	}

	if key == dots {
		store.stack.Pop()
	} else {
		store.stack.Push(listable)
	}

	store.makeLookupList()
}

func (store *Store) FolderList() []string {
	return store.folderList
}

func (store *Store) makeLookupList() (err error) {

	store.uriMap = make(map[string]fyne.URI)
	currentUri, isRoot := store.stack.Current()
	if !isRoot {
		store.uriMap[dots] = store.Current
	}

	store.Current = currentUri
	uriList, err := store.Current.List()
	store.keyList = make([]string, 0, len(uriList)+1)
	store.folderList = make([]string, 0)
	if !isRoot {
		store.keyList = append(store.keyList, dots)
	}

	for _, uri := range uriList {
		title := makeTitle(uri)
		store.uriMap[title] = uri
		store.keyList = append(store.keyList, title)
		isList, _ := storage.CanList(uri)
		if isList {
			store.folderList = append(store.folderList, title)
		}
	}
	return
}

func (store *Store) LoadFrame(key string, frame *glow.Frame) error {

	uri, ok := store.uriMap[key]
	if !ok {
		err := fmt.Errorf("key not in URI map")
		fyne.LogError(resources.MsgGetEffectLookup.Format(key), err)
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

	// successfully loaded so update route
	store.route = make([]string, 0)
	for _, uri := range store.stack.Dump() {
		store.route = append(store.route, makeTitle(uri))
	}

	return nil
}

func (store *Store) StoreFrame(fname string, frame *glow.Frame) (err error) {
	var (
		uri fyne.URI
		wrt fyne.URIWriteCloser
		buf []byte
		// exists bool
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

func (store *Store) FrameListURI() fyne.ListableURI {
	return store.Current
}

func makeTitle(uri fyne.URI) (s string) {
	s = uri.Name()
	i := strings.Index(s, uri.Extension())
	if i > 0 {
		s = s[:i]
	}
	s = strings.ReplaceAll(s, "_", " ")
	return
}
