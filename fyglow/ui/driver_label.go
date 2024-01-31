package ui

import (
	"gglow/iohandler"
	"gglow/text"
)

var driverMap map[string]string = map[string]string{
	text.CodeLabel.String(): iohandler.DRIVER_CODE,
	text.DataLabel.String(): iohandler.DRIVER_SQLLITE3,
}

func DriverLabels() []string {
	labels := make([]string, 0, len(driverMap))
	for k := range driverMap {
		labels = append(labels, k)
	}
	return labels
}

func DriverFromLabel(label string) (driver string) {
	return driverMap[label]
}
