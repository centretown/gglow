package ui

import (
	"gglow/iohandler"
	"gglow/text"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

var driverMap map[string]string = map[string]string{
	text.CodeLabel.String(): iohandler.DRIVER_CODE,
	text.DataLabel.String(): iohandler.DRIVER_SQLLITE3,
}

func driverLabels() []string {
	labels := make([]string, 0, len(driverMap))
	for k, _ := range driverMap {
		labels = append(labels, k)
	}
	return labels
}

func driverFromLabel(label string) (driver string) {
	return driverMap[label]
}

func driversFromLabels(options []string) (list []string) {
	list = make([]string, 0, len(options))
	for _, s := range options {
		v, ok := driverMap[s]
		if ok {
			list = append(list, v)
		}
	}
	return
}

func NewEffectSelector(data binding.BoolTree,
	listener binding.DataListener,
	top fyne.CanvasObject) fyne.CanvasObject {
	tree := NewEffectTreeWithListener(data, listener, CreateCheck, UpdateCheck(data, listener))
	return container.NewBorder(
		container.NewBorder(nil, widget.NewSeparator(), nil, nil, top),
		nil, nil, nil, tree)
}
