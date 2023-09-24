package res

import (
	"glow-gui/glow"
	"strings"

	"fyne.io/fyne/v2"
)

var AppID = "com.centretown.glow.preferences"
var WindowTitle = "Light Effects"
var WindowSize = fyne.Size{Width: 600, Height: 400}

type ImageID uint16

const (
	GooseNoirImage ImageID = iota
)

var imagePath = []string{
	"res/gangsta-goose.png",
}

func (id ImageID) String() string {
	return imagePath[id]
}

var AppImage fyne.Resource

func (id ImageID) Load() (res fyne.Resource, err error) {
	res, err = fyne.LoadResourceFromPath(id.String())
	AppImage = res
	return
}

type ContentID uint16

const (
	FrameContent ContentID = iota
	LayerContent
	ChromaContent
	ContentCount
)

var contentLabels = []string{"Frames", "Layers", "Colors"}

func (id ContentID) String() string {
	return contentLabels[id]
}

func ContentLabels() []string {
	return contentLabels
}

type LabelID uint16

const (
	LengthLabel LabelID = iota
	RowsLabel
	IntervalLabel
	LayersLabel
	GridLabel
	ChromaLabel
	HueShiftLabel
	ScanLabel
	BeginLabel
	EndLabel
	OriginLabel
	OrientationLabel
	ColorsLabel
	HueLabel
	SaturationLabel
	ValueLabel
	ChooseEffectLabel
	LightEffectLabel
)

var entryLabels = []string{
	"Length", "Rows", "Interval", "Layers",
	"Grid", "Chroma", "Hue Shift", "Scan Length",
	"Begin At", "End At",
	"Origin", "Orientation",
	"Colors", "Hue", "Saturation", "Value",
	"Choose Effect", "Light Effect",
}

func (id LabelID) String() string {
	return entryLabels[id]
}

func (id LabelID) PlaceHolder() string {
	return strings.ToLower(entryLabels[id])
}

type OrientationID glow.Orientation

var orientationLabels = []string{
	"Horizontal",
	"Vertical",
	"Diagonal",
}

func (id OrientationID) String() string {
	return orientationLabels[id]
}

func (id OrientationID) PlaceHolder() string {
	return strings.ToLower(orientationLabels[id])
}

type OriginID glow.Origin

var originLabels = []string{
	"Top Left",
	"Top Right",
	"Bottom Left",
	"Bottom Right",
}

func (id OriginID) String() string {
	return originLabels[id]
}

func (id OriginID) PlaceHolder() string {
	return strings.ToLower(originLabels[id])
}
