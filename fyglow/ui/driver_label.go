package ui

import (
	"gglow/iohandler"
	"gglow/text"
)

type driverKey struct {
	label  string
	driver string
}

var driverKeys = []driverKey{
	{text.CodeLabel.String(), iohandler.DRIVER_CODE},
	{text.DataLabel.String(), iohandler.DRIVER_SQLLITE3},
}

func DriverLabels() []string {
	labels := make([]string, 0, len(driverKeys))
	for _, k := range driverKeys {
		labels = append(labels, k.label)
	}
	return labels
}

func DriverFromLabel(label string) string {
	for _, k := range driverKeys {
		if k.label == label {
			return k.driver
		}
	}
	return ""
}
