package data

import (
	"fmt"
	"glow-gui/glow"
	"glow-gui/resources"
	"strings"
)

func Summarize(layer *glow.Layer) string {
	bldr := strings.Builder{}
	bldr.Grow(80)

	space := " "

	if layer.HueShift != 0 {
		bldr.WriteString(resources.DynamicLabel.String() + space)
	} else {
		bldr.WriteString(resources.StaticLabel.String() + space)
	}

	bldr.WriteString(resources.OrientationID(
		layer.Grid.Orientation).String() + space)

	if len(layer.Chroma.Colors) > 1 {
		bldr.WriteString(resources.GradientLabel.String() + space)
	}

	if layer.Scan > 0 {
		bldr.WriteString(resources.ScannerLabel.String() + space)
	}

	if layer.Begin != 0 || layer.End != 100 {
		bldr.WriteString(fmt.Sprintf("%d%%",
			layer.End-layer.Begin))
	}

	return bldr.String()
}
