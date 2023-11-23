package store

import (
	"fmt"
	"glow-gui/glow"
	"glow-gui/resources"
	"strings"
	"time"
)

const (
	guessLength = 16
)

type HistoryItem struct {
	List   []*glow.Frame
	cursor int
}

func NewHistoryItem() *HistoryItem {
	return &HistoryItem{
		List: make([]*glow.Frame, 0),
	}
}

type History struct {
	TimeStamp    time.Time
	HistoryItems map[string]*HistoryItem
}

func NewHistory() *History {
	h := &History{
		TimeStamp:    time.Now(),
		HistoryItems: make(map[string]*HistoryItem),
	}
	return h
}

func (h *History) makePath(route []string, title string) string {
	bld := &strings.Builder{}
	bld.Grow(guessLength + len(route)*guessLength)
	for _, s := range route {
		bld.WriteString(s)
		bld.WriteRune('/')
	}
	bld.WriteString(title)
	return bld.String()
}

func (h *History) Add(route []string, title string, source *glow.Frame) error {
	path := h.makePath(route, title)
	item, ok := h.HistoryItems[path]
	if !ok {
		item = NewHistoryItem()
		h.HistoryItems[path] = item
	}

	frame, err := glow.FrameDeepCopy(source)
	if err != nil {
		return err
	}

	// last := len(item.List) - 1
	// if item.cursor < last {
	// 	item.List = item.List[:item.cursor+1]
	// }

	item.List = append(item.List, frame)
	item.cursor = len(item.List) - 1
	fmt.Println("HistoryAdd", title, item.cursor)
	return nil
}

func (h *History) HasPrevious(route []string, title string) bool {
	item, ok := h.HistoryItems[h.makePath(route, title)]
	var length int
	if ok {
		length = len(item.List)
		ok = length > 0 && item.cursor >= 0
		fmt.Print(item.cursor, " ")
	}
	fmt.Println(ok, length)
	return ok
}

func (h *History) Previous(route []string, title string) (frame *glow.Frame, err error) {
	path := h.makePath(route, title)
	item, ok := h.HistoryItems[path]
	if !ok {
		err = fmt.Errorf("%s: %s", path, resources.MsgNotFound.String())
		return
	}

	if len(item.List) < 1 || item.cursor < 0 {
		err = fmt.Errorf("%s: %s", path, resources.MsgListEmpty.String())
		return
	}

	frame = item.List[item.cursor]
	item.cursor--
	return
}

func (h *History) Dump() {
	// const dump_file = "history_dump.json"
}
