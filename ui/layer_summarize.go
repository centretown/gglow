package ui

import (
	"fmt"
	"glow-gui/glow"
	"glow-gui/res"
	"strings"
)

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

	if layer.Begin != 0 || layer.End != 100 {
		bldr.WriteString(fmt.Sprintf("%d%%",
			layer.End-layer.Begin))
	}

	return bldr.String()
}
