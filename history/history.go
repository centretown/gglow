package history

import (
	"fmt"
	"gglow/glow"
	"gglow/text"
	"time"
)

const (
	guessLength = 16
)

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

func (h *History) Find(path string) (*HistoryItem, bool) {
	item, ok := h.HistoryItems[path]
	return item, ok
}

func (h *History) Add(route []string, title string, source *glow.Frame) error {
	path := makePath(route, title)
	item, ok := h.Find(path)
	if !ok {
		item = NewHistoryItem(route, title)
		h.HistoryItems[path] = item
	}
	return item.Push(source)
}

func (h *History) HasHistory(route []string, title string) bool {
	item, ok := h.Find(makePath(route, title))
	if ok {
		ok = item.HasHistory()
	}
	return ok
}

func (h *History) RestorePrevious(route []string, title string) (source *glow.Frame, err error) {
	path := makePath(route, title)
	item, ok := h.Find(path)
	if !ok {
		err = fmt.Errorf("%s: %s", path, text.MsgNotFound.String())
		return
	}

	if !item.HasHistory() {
		err = fmt.Errorf("%s: %s", path, text.MsgListEmpty.String())
		return

	}

	return item.Pop()
}

func (h *History) Dump() {
	for k, v := range h.HistoryItems {
		fmt.Println(k)
		fmt.Println(v.Route, v.Title, len(v.List))
	}
}
