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
	scheme       = "file://"
	BasePath     = "/home/dave/src/glow-gui/storage/effects/"
	ExamplesPath = "/home/dave/src/glow-gui/storage/effects/examples/"
	dots         = "..."
)

type Store struct {
	Current fyne.ListableURI
	Derived fyne.ListableURI
	uriMap  map[string]fyne.URI
	keyList []string
	stack   *Stack
}

func NewStore() *Store {
	path := scheme + BasePath
	uri, err := storage.ParseURI(path)

	if err != nil {
		formatMessage(resources.MsgParseEffectPath, path, err)
		os.Exit(1)
	}

	listable, err := storage.ListerForURI(uri)
	if err != nil {
		formatMessage(resources.MsgParseEffectPath, path, err)
		os.Exit(1)
	}

	store := &Store{
		uriMap: make(map[string]fyne.URI),
		stack:  NewStack(listable),
	}

	store.makeLookupList()
	return store
}

func (store *Store) LookUpList() []string {
	return store.keyList
}

func (store *Store) LookupURI(s string) (uri fyne.URI, err error) {
	uri = store.uriMap[s]
	if uri == nil {
		err = fmt.Errorf("LookupURI: %s not found", s)
		return
	}
	return
}

func formatMessage(id resources.MessageID,
	path string, err error) error {
	id.Log(path, err)
	return fmt.Errorf("%s %s %s", id, path, err.Error())
}

func (store *Store) IsFolder(path string) bool {
	uri, ok := store.uriMap[path]
	if ok {
		ok, _ = storage.CanList(uri)
	}
	return ok
}

func (store *Store) RefreshLookupList(key string) {

	uri, ok := store.uriMap[key]
	if !ok {
		fyne.LogError("", fmt.Errorf("%s not found", key))
		return
	}

	listable, err := storage.ListerForURI(uri)
	if err != nil {
		fyne.LogError("", fmt.Errorf("%s not listable", key))
		return
	}

	if key != "..." {
		store.stack.Push(listable)
	}
	store.makeLookupList()
}

func (store *Store) makeLookupList() (err error) {

	store.uriMap = make(map[string]fyne.URI)

	currentUri, isBase := store.stack.Pop()

	if !isBase {
		store.uriMap[dots] = store.Current
	}

	store.Current = currentUri
	uriList, err := store.Current.List()
	store.keyList = make([]string, 0, len(uriList)+1)
	if !isBase {
		store.keyList = append(store.keyList, dots)
	}

	for _, uri := range uriList {
		title := makeTitle(uri)
		store.uriMap[title] = uri
		store.keyList = append(store.keyList, title)
	}

	return
}

func (store *Store) LoadFrameURI(uri fyne.URI, frame *glow.Frame) (err error) {
	var (
		rdr fyne.URIReadCloser
	)
	rdr, err = storage.Reader(uri)
	if err != nil {
		return
	}

	return store.readFrame(rdr, frame)
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

func (store *Store) readFrame(rdr fyne.URIReadCloser, frame *glow.Frame) (err error) {
	defer rdr.Close()

	var b []byte

	b, err = io.ReadAll(rdr)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(b, frame)
	return
}

// func (store *Store) loadFrame(fname string, frame *glow.Frame) (err error) {
// 	var (
// 		uri fyne.URI
// 		rdr fyne.URIReadCloser
// 	)

// 	uri, err = storage.ParseURI(scheme + fname)
// 	if err != nil {
// 		return
// 	}

// 	rdr, err = storage.Reader(uri)
// 	if err != nil {
// 		return
// 	}

// 	return store.readFrame(rdr, frame)
// }
