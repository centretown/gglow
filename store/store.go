package store

import (
	"fmt"
	"glow-gui/glow"
	"glow-gui/res"
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"gopkg.in/yaml.v3"
)

const (
	scheme          = "file://"
	max_buffer_size = 4096
)

// var readbuf []byte = make([]byte, max_buffer_size)

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

	path := scheme + FramePath

	formatMessage := func(id res.MessageID, msg string, err error) error {
		id.Log(path, err)
		return fmt.Errorf("%s %s %s", id, msg, err.Error())
	}

	uri, err = storage.ParseURI(path)
	if err != nil {
		err = formatMessage(res.MsgParseEffectPath, path, err)
		return
	}

	canList, err = storage.CanList(uri)
	if err != nil {
		err = formatMessage(res.MsgNoAccess, path, err)
		return
	}

	if !canList {
		err = formatMessage(res.MsgPathNotFolder, path, err)
		return
	}

	FrameURI, err = storage.ListerForURI(uri)
	if err != nil {
		err = formatMessage(res.MsgNoList, path, err)
		return
	}

	uriList, err := FrameURI.List()
	if err != nil {
		err = formatMessage(res.MsgNoList, path, err)
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

func readFrame(rdr fyne.URIReadCloser, frame *glow.Frame) (err error) {
	defer rdr.Close()

	var b []byte

	b, err = io.ReadAll(rdr)
	if err != nil {
		return
	}

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

	return readFrame(rdr, frame)
}

func loadFrame(fname string, frame *glow.Frame) (err error) {
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

	return readFrame(rdr, frame)
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
