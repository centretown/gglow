package storageio

import (
	"fmt"
	"glow-gui/effects"
	"glow-gui/glow"
	"glow-gui/resources"
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

var _ effects.IoHandler = (*StorageHandler)(nil)

const (
	scheme = "file://"
)

type StorageHandler struct {
	RootURI  fyne.ListableURI
	Current  fyne.ListableURI
	uriMap   map[string]fyne.URI
	rootPath string
	keyList  []string
	folder   string
	title    string

	serializer effects.Serializer
}

func NewStorageHandler(rootPath string) (*StorageHandler, error) {
	path := scheme + rootPath

	uri, err := storage.ParseURI(path)
	if err != nil {
		fyne.LogError(resources.MsgParseEffectPath.Format(path), err)
		return nil, err
	}

	rootURI, err := storage.ListerForURI(uri)
	if err != nil {
		fyne.LogError(resources.MsgPathNotFolder.Format(path), err)
		return nil, err
	}

	fh := &StorageHandler{
		RootURI:    rootURI,
		uriMap:     make(map[string]fyne.URI),
		rootPath:   rootPath,
		keyList:    make([]string, 0),
		serializer: &effects.JsonSerializer{},
	}

	return fh, nil
}

func (fh *StorageHandler) Refresh() ([]string, error) {
	fh.Current = fh.RootURI
	fh.folder = ".."
	err := fh.makeLookupList()
	return fh.keyList, err
}

func (fh *StorageHandler) IsFolder(key string) bool {
	uri, ok := fh.uriMap[key]
	if ok {
		ok, _ = storage.CanList(uri)
	}
	return ok
}

func (fh *StorageHandler) RefreshFolder(folder string) ([]string, error) {
	if folder == effects.Dots {
		return fh.Refresh()
	}

	uri, ok := fh.uriMap[folder]
	if !ok {
		err := fmt.Errorf(resources.MsgGetEffectLookup.Format(folder))
		fyne.LogError("RefreshFolder", err)
		return fh.keyList, err
	}

	listable, err := storage.ListerForURI(uri)
	if err != nil {
		err = fmt.Errorf(resources.MsgPathNotFolder.Format(folder))
		fyne.LogError("RefreshFolder", err)
		return fh.keyList, err
	}

	fh.Current = listable
	fh.folder = folder
	err = fh.makeLookupList()
	return fh.keyList, err
}

func (fh *StorageHandler) KeyList() []string {
	return fh.keyList
}

func (fh *StorageHandler) makeLookupList() (err error) {

	fh.uriMap = make(map[string]fyne.URI)
	currentUri := fh.Current
	isRoot := currentUri == fh.RootURI
	if !isRoot {
		fh.uriMap[effects.Dots] = fh.Current
	}
	fh.Current = currentUri

	var uriList []fyne.URI
	uriList, err = fh.Current.List()
	if err != nil {
		return
	}

	fh.keyList = make([]string, 0, len(uriList)+1)
	if !isRoot {
		fh.keyList = append(fh.keyList, effects.Dots)
	}

	for _, uri := range uriList {
		title := MakeTitle(uri)
		fh.uriMap[title] = uri
		fh.keyList = append(fh.keyList, title)
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

	path := scheme + fh.Current.Path() + "/" + fh.serializer.FileName(title)
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

	fh.folder = fh.Current.Name()
	fh.title = title
	return frame, err
}

func (fh *StorageHandler) FolderName() string {
	return fh.folder
}

func (fh *StorageHandler) EffectName() string {
	return fh.title
}

func (fh *StorageHandler) ValidateNewFolderName(title string) error {
	err := effects.ValidateFolderName(title)
	if err != nil {
		return err
	}

	err = fh.isDuplicate(title)
	return err
}

func (fh *StorageHandler) ValidateNewEffectName(title string) error {
	err := effects.ValidateEffectName(title)
	if err != nil {
		return err
	}
	err = fh.isDuplicate(title)
	return err
}

func (fh *StorageHandler) OnExit() {
}

func (fh *StorageHandler) CreateNewDatabase() error {
	return nil
}
