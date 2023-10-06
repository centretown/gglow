package ui

import (
	"glow-gui/data"
	"glow-gui/res"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type LayerList struct {
	*widget.List
	model      *data.Model
	changeView func()
}

func NewLayerList(model *data.Model, changeView func()) *LayerList {
	ls := &LayerList{
		model:      model,
		changeView: changeView,
	}

	ls.List = widget.NewListWithData(model.LayerSummaryList,
		ls.createLayerItem,
		ls.updateLayerItem)

	ls.List.OnSelected = func(id widget.ListItemID) {
		ls.model.SetCurrentLayer(id)
	}

	return ls
}

const (
	ButtonPos = 0
	LabelPos  = 1
)

func (ls *LayerList) createLayerItem() fyne.CanvasObject {
	return container.NewHBox(
		widget.NewButtonWithIcon("",
			res.AppIconResource(res.LayerIcon), ls.changeView),
		widget.NewLabel("template"),
	)
}

func (ls *LayerList) updateLayerItem(item binding.DataItem,
	canvasObj fyne.CanvasObject) {

	s := item.(binding.String)
	text, _ := s.Get()
	box := canvasObj.(*fyne.Container)
	label := box.Objects[LabelPos].(*widget.Label)
	label.SetText(text)

	// layer := item.Get()
	// box := canvasObj.(*fyne.Container)
	// unTypedLayer, err := item.(binding.Untyped).Get()
	// if err != nil {
	// 	res.MsgGetLayer.Log("layerList item", err)
	// 	panic(err.Error()) //PANIC
	// }

	// layer := unTypedLayer.(*glow.Layer)
	// binder := binding.NewString()
	// binder.Set(Summarize(layer))

	// label.Bind(binder)

	// button := box.Objects[ButtonPos].(*widget.Button)
	// button.OnTapped = func() {
	// 	// form := NewLayerDialog(ls.window, layer)
	// 	// form.Show()
	// }
}
