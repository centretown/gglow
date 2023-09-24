package store

import (
	"fmt"
	"glow-gui/glow"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"gopkg.in/yaml.v3"
)

const (
	scheme          = "file://"
	max_buffer_size = 4096
)

var readbuf []byte = make([]byte, max_buffer_size)

const (
	FramePath   = "/home/dave/src/glow-gui/res/frames/"
	DerivedPath = "/home/dave/src/glow-gui/res/frames/derived"
)

var (
	FrameURI   fyne.ListableURI
	DerivedURI fyne.ListableURI
)

var uri_lookup = make(map[string]fyne.URI)
var lookUpList []string

func LookUpList() []string {
	return lookUpList
}

func LookupURI(s string) (uri fyne.URI, err error) {
	uri = uri_lookup[s]
	if uri == nil {
		err = fmt.Errorf("LookupURI: %s not found", s)
		return
	}
	return
}

func Setup() (err error) {
	var (
		uri     fyne.URI
		canList bool
	)

	uri, err = storage.ParseURI(scheme + FramePath)
	if err != nil {
		return
	}

	canList, err = storage.CanList(uri)
	if err != nil {
		return
	}

	if !canList {
		err = fmt.Errorf("%s un-listable", FramePath)
		return
	}

	FrameURI, err = storage.ListerForURI(uri)
	if err != nil {
		return
	}

	uriList, err := FrameURI.List()
	if err != nil {
		return
	}

	count := 0
	lookUpList = make([]string, len(uriList))
	for _, uri := range uriList {
		canList, err := storage.CanList(uri)
		if err == nil && !canList {
			s := makeTitle(uri)
			uri_lookup[s] = uri
			lookUpList[count] = s
			count++
		}
	}

	lookUpList = lookUpList[:count]
	return
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

func FrameListURI() fyne.ListableURI {
	return FrameURI
}

func ReadFrame(rdr fyne.URIReadCloser, frame *glow.Frame) (err error) {
	var (
		count int
	)

	defer rdr.Close()

	b := readbuf

	count, err = rdr.Read(b)
	if err != nil {
		return
	}

	b = b[:count]
	err = yaml.Unmarshal(b, frame)
	if err != nil {
		return
	}

	return
}

func LoadFrameURI(uri fyne.URI, frame *glow.Frame) (err error) {
	var (
		rdr fyne.URIReadCloser
	)
	rdr, err = storage.Reader(uri)
	if err != nil {
		return
	}

	return ReadFrame(rdr, frame)
}

func LoadFrame(fname string, frame *glow.Frame) (err error) {
	var (
		uri fyne.URI
		rdr fyne.URIReadCloser
	)

	uri, err = storage.ParseURI(scheme + fname)
	if err != nil {
		return
	}

	rdr, err = storage.Reader(uri)
	if err != nil {
		return
	}

	return ReadFrame(rdr, frame)
}

func StoreFrame(fname string, frame *glow.Frame) (err error) {
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
