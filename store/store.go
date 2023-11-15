package store

import (
	"fmt"
	"glow-gui/glow"
	"glow-gui/resources"
	"glow-gui/settings"
	"io"
	"os"
	"strings"
	"unicode"

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

func ValidateEffectName(title string) error {
	title = strings.TrimSpace(title)

	if len(title) < 1 {
		return fmt.Errorf(resources.MsgRequired.String())
	}

	if title == "NULL" {
		return fmt.Errorf(resources.MsgRequired.String())
	}

	for i, c := range title {
		if i == 0 && !unicode.IsUpper(c) {
			return fmt.Errorf(resources.MsgFirstUpper.String())
		}
		if !(c == ' ' || unicode.IsLetter(c) || unicode.IsDigit(c)) {
			return fmt.Errorf(resources.MsgAlphaNumeric.String())
		}
	}
	return nil
}

func (store *Store) IsDuplicate(title string) error {
	_, found := store.uriMap[title]
	if found {
		return fmt.Errorf(resources.MsgDuplicate.String())
	}
	return nil
}

func ValidateFolderName(title string) error {
	title = strings.TrimSpace(title)

	if len(title) < 1 {
		return fmt.Errorf(resources.MsgRequired.String())
	}

	if title == "NULL" {
		return fmt.Errorf(resources.MsgRequired.String())
	}

	for _, c := range title {
		if !(c == '_' || unicode.IsLetter(c) || unicode.IsDigit(c)) {
			return fmt.Errorf(resources.MsgAlphaNumeric.String())
		}
	}
	return nil
}

func (store *Store) ValidateNewFolderName(title string) error {
	err := ValidateFolderName(title)
	if err != nil {
		return err
	}

	err = store.IsDuplicate(title)
	return err
}

func (store *Store) ValidateNewEffectName(title string) error {
	err := ValidateEffectName(title)
	if err != nil {
		return err
	}
	err = store.IsDuplicate(title)
	return err
}

func (store *Store) CreateNewEffect(title string, frame *glow.Frame) error {
	err := store.IsDuplicate(title)
	if err != nil {
		return err
	}
	return store.WriteEffect(title, frame)
}

func (store *Store) CreateNewFolder(title string) error {
	err := store.IsDuplicate(title)
	if err != nil {
		return err
	}
	return store.WriteFolder(title)
}

func (store *Store) WriteFolder(title string) error {
	title = strings.TrimSpace(title)
	err := ValidateFolderName(title)
	if err != nil {
		return err
	}

	path := scheme + store.Current.Path() + "/" + title
	fmt.Println(path)
	uri, err := storage.ParseURI(path)
	if err != nil {
		return err
	}

	err = storage.CreateListable(uri)
	return err
}

func (store *Store) WriteEffect(title string, frame *glow.Frame) error {
	title = strings.TrimSpace(title)

	err := ValidateEffectName(title)
	if err != nil {
		return err
	}

	path := scheme + store.Current.Path() + "/" + MakeFileName(title)
	fmt.Println(path)
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

	n, err := wrt.Write(buf)
	if err != nil {
		return err
	}

	store.makeLookupList()

	fmt.Println("bytes written", n)
	return nil
}
