package ui

import (
	"glow-gui/glow"
	"glow-gui/res"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

// type LayerListData struct {
// 	binding.StringList
// 	frame *glow.Frame
// }

// func NewLayerListData(frame *glow.Frame) *LayerListData {

// 	list := binding.NewStringList()

// 	for _, l := range frame.Layers {
// 		list.Append(Summarize(&l))
// 	}

// 	return &LayerListData{
// 		StringList: list,
// 		frame:      frame,
// 	}

// }

// func (lld *LayerListData) GetItem(i int) (binding.DataItem, error) {
// 	return lld.list.
// 	// return binding.String())
// }

// func (lld *LayerListData) Length() int {
// 	return len(lld.frame.Layers)
// }

type LayerList struct {
	binding.DataItem
	*widget.List
	window fyne.Window
	data   binding.UntypedList
}

func NewLayerList(window fyne.Window, frame *glow.Frame) *LayerList {
	ll := &LayerList{
		window: window,
	}

	ll.data = binding.NewUntypedList()
	ll.SetFrame(frame)

	ll.List = widget.NewListWithData(ll.data,
		ll.createLayerItem,
		ll.updateLayerItem)

	ll.DataItem = ll.data
	return ll
}

func (ll *LayerList) createLayerItem() fyne.CanvasObject {
	return container.NewHBox(
		widget.NewButtonWithIcon("",
			res.AppIconResource(res.LayerIcon), func() {}),
		widget.NewLabel("template"))
}

func (ll *LayerList) updateLayerItem(i binding.DataItem, box fyne.CanvasObject) {
	c := box.(*fyne.Container)
	x, err := i.(binding.Untyped).Get()
	if err != nil {
		fyne.LogError("Error getting data item", err)
	}

	layer := x.(*glow.Layer)
	s := binding.NewString()
	s.Set(Summarize(layer))

	l := c.Objects[1].(*widget.Label)
	l.Bind(s)

	b := c.Objects[0].(*widget.Button)
	b.OnTapped = func() {
		f := NewLayerForm(ll.window, layer)
		f.Show()
	}
}

func (ll *LayerList) SetFrame(frame *glow.Frame) {
	list := make([]interface{}, 0, len(frame.Layers))
	for i := range frame.Layers {
		list = append(list, &frame.Layers[i])
	}
	ll.data.Set(list)
}
