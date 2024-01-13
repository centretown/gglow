package history

import (
	"fmt"
	"gglow/glow"
	"gglow/text"
	"strings"
)

type HistoryItem struct {
	Route []string
	Title string
	List  []*glow.Frame
}

func NewHistoryItem(route []string, title string) *HistoryItem {
	return &HistoryItem{
		Route: route,
		Title: title,
		List:  make([]*glow.Frame, 0),
	}
}

func (hi *HistoryItem) Path() string {
	return makePath(hi.Route, hi.Title)
}

func (hi *HistoryItem) HasHistory() bool {
	return len(hi.List) > 0
}

func (hi *HistoryItem) Push(source *glow.Frame) error {
	frame, err := glow.FrameDeepCopy(source)
	if err != nil {
		return err
	}
	hi.List = append(hi.List, frame)
	length := len(hi.List)
	fmt.Println("push", length)
	return nil
}

func (hi *HistoryItem) Pop() (*glow.Frame, error) {
	length := len(hi.List)
	fmt.Println("pop", length)
	if length < 1 {
		return nil,
			fmt.Errorf("%s: %s", makePath(hi.Route, hi.Title),
				text.MsgNotFound.String())
	}
	length--
	frame := hi.List[length]
	hi.List = hi.List[:length]
	return frame, nil
}

func makePath(route []string, title string) string {
	bld := &strings.Builder{}
	bld.Grow(guessLength + len(route)*guessLength)
	for _, s := range route {
		bld.WriteString(s)
		bld.WriteRune('/')
	}
	bld.WriteString(title)
	return bld.String()
}
