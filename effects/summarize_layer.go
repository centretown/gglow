package effects

import (
	"fmt"
	"glow-gui/glow"
	"glow-gui/resources"
	"strings"
)

func SummarizeLayer(layer *glow.Layer, index int) string {
	bldr := strings.Builder{}
	bldr.Grow(80)

	space := " "
	bldr.WriteString(fmt.Sprintf("%d: ", index))

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
		bldr.WriteString(resources.ScanLabel.String() + space)
	}

	if layer.Begin != 0 {
		bldr.WriteString(fmt.Sprintf("%d%%", layer.Begin))
	}

	return bldr.String()
}
