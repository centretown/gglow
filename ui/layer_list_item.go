package ui

import (
	"fmt"
	"glow-gui/glow"
	"glow-gui/res"
	"strings"
)

// import (
// 	"fmt"
// 	"glow-gui/glow"
// 	"glow-gui/res"
// 	"strings"

// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/data/binding"
// 	"fyne.io/fyne/v2/widget"
// )

// type LayerListItemData struct {
// 	binding.String
// 	layer *glow.Layer
// }

// func NewLayerListItemData(layer *glow.Layer) (ld *LayerListItemData) {
// 	ld = &LayerListItemData{
// 		String: binding.NewString(),
// 		layer:  layer,
// 	}
// 	return
// }

// type LayerListItem struct {
// 	fyne.CanvasObject
// 	binding.DataItem
// 	Label *widget.Label
// }

// func NewLayerListItem(layer *glow.Layer) *LayerListItem {
// 	item := &LayerListItem{
// 		CanvasObject: hbox,
// 		DataItem:     NewLayerListItemData(layer),
// 		Label:        label,
// 	}

// 	itemData := NewLayerListItemData(layer)

// 	hbox := container.NewHBox()
// 	hbox.Add(widget.NewButtonWithIcon("",
// 		res.AppIconResource(res.LayerIcon), func() {}))
// 	label := widget.NewLabelWithData(itemData)
// 	hbox.Add(label)

// 	return item
// }

// // func (ll *LayerListItem) Get() (string, error) {
// // 	return Summarize(ll.layer), nil
// // }

// // func (ll *LayerListItem) Set(string) error {
// // 	return nil
// // }

func Summarize(layer *glow.Layer) string {
	bldr := strings.Builder{}
	bldr.Grow(80)

	space := " "

	if layer.HueShift != 0 {
		bldr.WriteString(res.DynamicLabel.String() + space)
	} else {
		bldr.WriteString(res.StaticLabel.String() + space)
	}

	bldr.WriteString(res.OrientationID(
		layer.Grid.Orientation).String() + space)

	if len(layer.Chroma.Colors) > 1 {
		bldr.WriteString(res.GradientLabel.String() + space)
	}

	if layer.Scan > 0 {
		bldr.WriteString(res.ScannerLabel.String() + space)
	}
	//  else {
	// 	bldr.WriteString(res.BackDropLabel.String() + space)
	// }

	if layer.Begin != 0 || layer.End != 100 {

		bldr.WriteString(fmt.Sprintf("%d%%",
			layer.End-layer.Begin))
	}

	return bldr.String()
}
