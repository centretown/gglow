package effectio

import (
	"fmt"
	"gglow/glow"
	"gglow/text"
	"strings"
)

func SummarizeLayer(layer *glow.Layer, index int) string {
	bldr := strings.Builder{}
	bldr.Grow(80)

	space := " "
	bldr.WriteString(fmt.Sprintf("%d: ", index))

	if layer.HueShift != 0 {
		bldr.WriteString(text.DynamicLabel.String() + space)
	} else {
		bldr.WriteString(text.StaticLabel.String() + space)
	}

	bldr.WriteString(text.OrientationID(
		layer.Grid.Orientation).String() + space)

	if len(layer.Chroma.Colors) > 1 {
		bldr.WriteString(text.GradientLabel.String() + space)
	}

	if layer.Scan > 0 {
		bldr.WriteString(text.ScanLabel.String() + space)
	}

	if layer.Begin != 0 {
		bldr.WriteString(fmt.Sprintf("%d%%", layer.Begin))
	}

	return bldr.String()
}
