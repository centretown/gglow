package storageio

import (
	"fmt"
	"gglow/fyio"
	"gglow/glow"
	"gglow/iohandler"
	"gglow/resources"
	"io"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

// var _ io.IoHandler = (*StorageHandler)(nil)

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

	serializer iohandler.Serializer
}

func makeFolder(folder string) (err error) {
	var info os.FileInfo
	info, err = os.Stat(folder)
	if err != nil {
		err = os.Mkdir(folder, os.ModeDir|os.ModePerm)
	} else if !info.IsDir() {
		err = fmt.Errorf("%s exists but is not a directory", folder)
	}
	return
}

func RootURI(rootPath string) (fyne.ListableURI, error) {
	path := scheme + rootPath

	uri, err := storage.ParseURI(path)
	if err != nil {
		fyne.LogError(resources.MsgParseEffectPath.Format(path), err)
		return nil, err
	}

	rootURI, err := storage.ListerForURI(uri)
	if err != nil {
		err = makeFolder(rootPath)
		if err == nil {
			return RootURI(rootPath)
		}
	}
	return rootURI, err
}

func NewStorageHandler(rootPath string) (*StorageHandler, error) {
	rootURI, err := RootURI(rootPath)
	if err != nil {
		return nil, err
	}

	fh := &StorageHandler{
		RootURI:    rootURI,
		uriMap:     make(map[string]fyne.URI),
		rootPath:   rootPath,
		keyList:    make([]string, 0),
		serializer: &iohandler.JsonSerializer{},
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
	if folder == fyio.Dots {
		return fh.Refresh()
	}

	uri, ok := fh.uriMap[folder]
	if !ok {
		path := fh.rootPath + "/" + folder
		err := makeFolder(path)
		if err != nil {
			return fh.keyList, err
		}
		path = scheme + path
		uri, err = storage.ParseURI(path)
		if err != nil {
			fyne.LogError("RefreshFolder", err)
			return fh.keyList, err
		}
		fh.uriMap[folder] = uri
	}

	listable, err := storage.ListerForURI(uri)
	if err != nil {
		fyne.LogError("ListerForURI", err)
		err = fmt.Errorf(resources.MsgPathNotFolder.Format(folder))
		fyne.LogError("ListerForURI", err)
		return fh.keyList, err
	}

	fh.Current = listable
	fh.folder = folder
	err = fh.makeLookupList()
	return fh.keyList, err
}

func (fh *StorageHandler) ListCurrentFolder() []string {
	return fh.keyList
}

func (fh *StorageHandler) makeLookupList() (err error) {

	fh.uriMap = make(map[string]fyne.URI)
	currentUri := fh.Current
	isRoot := currentUri == fh.RootURI
	if !isRoot {
		fh.uriMap[fyio.Dots] = fh.Current
	}
	fh.Current = currentUri

	var uriList []fyne.URI
	uriList, err = fh.Current.List()
	if err != nil {
		return
	}

	fh.keyList = make([]string, 0, len(uriList)+1)
	if !isRoot {
		fh.keyList = append(fh.keyList, fyio.Dots)
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

func (fh *StorageHandler) WriteFolder(folder string) error {
	fh.folder = folder
	fmt.Println(folder, "folder")
	folder = strings.TrimSpace(folder)
	err := fyio.ValidateFolderName(folder)
	if err != nil {
		return err
	}

	path := scheme + fh.Current.Path() + "/" + folder
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

	serializer := iohandler.UriSerializer(uri.Extension())
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
	err := fyio.ValidateFolderName(title)
	if err != nil {
		return err
	}

	err = fh.isDuplicate(title)
	return err
}

func (fh *StorageHandler) ValidateNewEffectName(title string) error {
	err := fyio.ValidateEffectName(title)
	if err != nil {
		return err
	}
	err = fh.isDuplicate(title)
	return err
}

func (fh *StorageHandler) OnExit() {
}

func (fh *StorageHandler) CreateNewDatabase(rootPath string) (err error) {
	return
}
