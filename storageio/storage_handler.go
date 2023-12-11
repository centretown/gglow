package storageio

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
	"fyne.io/fyne/v2/storage"
)

const (
	scheme            = "file://"
	defaultEffectPath = "/home/dave/src/glow-gui/cabinet/effects/"
	examplesPath      = "/home/dave/src/glow-gui/cabinet/examples/"
	dots              = ".."
)

type StorageHandler struct {
	Current     fyne.ListableURI
	FolderList  []string
	stack       *URIStack
	uriMap      map[string]fyne.URI
	rootPath    string
	keyList     []string
	route       []string
	preferences fyne.Preferences
	serializer  effects.Serializer
}

func NewStorageHandler(preferences fyne.Preferences) *StorageHandler {
	rootPath := preferences.StringWithFallback(settings.EffectPath.String(),
		defaultEffectPath)
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

	fh := &StorageHandler{
		preferences: preferences,
		FolderList:  make([]string, 0),
		stack:       NewStack(rootURI),
		uriMap:      make(map[string]fyne.URI),
		rootPath:    rootPath,
		keyList:     make([]string, 0),
		serializer:  &effects.JsonSerializer{},
	}

	fh.stack.Push(rootURI)
	fh.makeLookupList()
	fh.route = preferences.StringListWithFallback(settings.EffectRoute.String(),
		[]string{MakeTitle(rootURI)})

	for _, s := range fh.route[1:] {
		fh.RefreshKeys(s)
	}

	return fh
}

func (fh *StorageHandler) IsFolder(key string) bool {
	uri, ok := fh.uriMap[key]
	if ok {
		ok, _ = storage.CanList(uri)
	}
	return ok
}

func (fh *StorageHandler) RefreshKeys(key string) ([]string, error) {

	uri, ok := fh.uriMap[key]
	if !ok {
		err := fmt.Errorf(resources.MsgGetEffectLookup.Format(key))
		fyne.LogError("RefreshKeys", err)
		return fh.keyList, err
	}

	listable, err := storage.ListerForURI(uri)
	if err != nil {
		fyne.LogError(resources.MsgPathNotFolder.Format(key), err)
		return fh.keyList, err
	}

	if key == dots {
		fh.stack.Pop()
	} else {
		fh.stack.Push(listable)
	}

	fh.makeLookupList()
	return fh.keyList, nil
}

func (fh *StorageHandler) KeyList() []string {
	return fh.keyList
}

func (fh *StorageHandler) makeLookupList() (err error) {

	fh.uriMap = make(map[string]fyne.URI)
	currentUri, isRoot := fh.stack.Current()
	if !isRoot {
		fh.uriMap[dots] = fh.Current
	}
	fh.Current = currentUri

	uriList, err := fh.Current.List()
	fh.keyList = make([]string, 0, len(uriList)+1)
	fh.FolderList = make([]string, 0)
	if !isRoot {
		fh.keyList = append(fh.keyList, dots)
		fh.FolderList = append(fh.FolderList, dots)
	}

	for _, uri := range uriList {
		title := MakeTitle(uri)
		fh.uriMap[title] = uri
		fh.keyList = append(fh.keyList, title)
		isList, _ := storage.CanList(uri)
		if isList {
			fh.FolderList = append(fh.FolderList, title)
		}
	}
	return
}

func (fh *StorageHandler) isDuplicate(title string) error {
	_, found := fh.uriMap[title]
	if found {
		return fmt.Errorf(resources.MsgDuplicate.String())
	}
	return nil
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

func (fh *StorageHandler) CreateNewEffect(title string, frame *glow.Frame) error {
	err := fh.isDuplicate(title)
	if err != nil {
		return err
	}
	return fh.WriteEffect(title, frame)
}

func (fh *StorageHandler) CreateNewFolder(title string) error {
	err := fh.isDuplicate(title)
	if err != nil {
		return err
	}
	return fh.WriteFolder(title)
}

func (fh *StorageHandler) WriteFolder(title string) error {
	title = strings.TrimSpace(title)
	err := effects.ValidateFolderName(title)
	if err != nil {
		return err
	}

	path := scheme + fh.Current.Path() + "/" + title
	uri, err := storage.ParseURI(path)
	if err != nil {
		return err
	}

	err = storage.CreateListable(uri)
	return err
}

func (fh *StorageHandler) WriteEffect(title string, frame *glow.Frame) error {
	title = strings.TrimSpace(title)

	path := scheme + fh.Current.Path() + "/" +
		fh.serializer.MakeFileName(title)
	fmt.Println(path)
	// path := scheme + fh.Current.Path() + "/" + MakeFileName(title)
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

	buf, err := fh.serializer.Format(frame)

	// buf, err := json.Marshal(frame)
	if err != nil {
		return err
	}

	_, err = wrt.Write(buf)
	if err != nil {
		return err
	}

	fh.makeLookupList()
	return err
}

func (fh *StorageHandler) ReadEffect(title string) (*glow.Frame, error) {

	frame := &glow.Frame{}

	uri, ok := fh.uriMap[title]
	if !ok {
		err := fmt.Errorf(resources.MsgGetEffectLookup.Format(title))
		fyne.LogError("ReadEffect", err)
		return frame, err
	}

	rdr, err := storage.Reader(uri)
	if err != nil {
		fyne.LogError(resources.MsgGetFrame.Format(title), err)
		return frame, err
	}
	defer rdr.Close()

	buffer, err := io.ReadAll(rdr)
	if err != nil {
		fyne.LogError(resources.MsgGetFrame.Format(title), err)
		return frame, err
	}

	serializer := effects.UriSerializer(uri)
	err = serializer.Scan(buffer, frame)
	if err != nil {
		fyne.LogError(resources.MsgGetFrame.Format(title), err)
		return frame, err
	}

	fh.route = fh.stack.Route()
	return frame, err
}

func (fh *StorageHandler) ValidateNewFolderName(title string) error {
	err := effects.ValidateFolderName(title)
	if err != nil {
		return err
	}

	err = fh.isDuplicate(title)
	return err
}

func (st *StorageHandler) ValidateNewEffectName(title string) error {
	err := effects.ValidateEffectName(title)
	if err != nil {
		return err
	}
	err = st.isDuplicate(title)
	return err
}

func (fh *StorageHandler) OnExit() {
	fh.preferences.SetStringList(settings.EffectRoute.String(), fh.route)
	fh.preferences.SetString(settings.EffectPath.String(), fh.rootPath)
}
