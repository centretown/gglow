package resource

import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestIcons(t *testing.T) {
	test.NewApp()

	LoadIcons("dark")
	fill := `fill="#ffffff">`

	for i := uint(0); i < ICON_COUNT; i++ {
		s := makeSVG(i, fill)
		fmt.Println(i+1, string(s))
	}

	res := Icon(IconLayerID)
	btn := widget.NewButtonWithIcon("", res, func() {})
	btn.CreateRenderer()

	res = Icon(IconLayerAddID)
	btn = widget.NewButtonWithIcon("", res, func() {})
	btn.CreateRenderer()

	res = Icon(IconLayerRemoveID)
	btn = widget.NewButtonWithIcon("", res, func() {})
	btn.CreateRenderer()

	res = Icon(IconLayerInsertID)
	btn = widget.NewButtonWithIcon("", res, func() {})
	btn.CreateRenderer()
}