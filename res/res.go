package res

import (
	"glow-gui/glow"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

var AppID = "com.centretown.glow.preferences"
var WindowSize = fyne.Size{Width: 480, Height: 480}

type ImageID uint16

const (
	GooseNoirImage ImageID = iota
)

var imagePath = []string{
	"res/dark-gander.png",
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

const (
	StripLength   float64 = 36
	StripRows     float64 = 4
	StripInterval float64 = 64
)

func (id ContentID) String() string {
	return contentLabels[id]
}

func ContentLabels() []string {
	return contentLabels
}

type LabelID uint16

const (
	GlowLabel LabelID = iota
	LengthLabel
	RowsLabel
	IntervalLabel
	LayersLabel
	GridLabel
	ChromaLabel
	HueShiftLabel
	ScanLengthLabel
	BeginLabel
	EndLabel
	OriginLabel
	OrientationLabel
	ColorsLabel
	HueLabel
	ScanLabel
	SaturationLabel
	ValueLabel
	EffectsLabel
	ScannerLabel
	BackDropLabel
	DynamicLabel
	StaticLabel
	GradientLabel
	ReversedLabel
)

var entryLabels = []string{
	"Glow",
	"Length", "Rows", "Interval", "Layers",
	"Grid", "Colors", "Hue Shift", "Scan Length",
	"Begin", "End",
	"Origin", "Orientation",
	"Colors", "Hue", "Scan",
	"Saturation", "Value",
	"Effects",
	"Scanner", "Backdrop",
	"Dynamic", "Static",
	"Gradient", "Reversed",
}

func (id LabelID) String() string {
	return entryLabels[id]
}

func (id LabelID) PlaceHolder() string {
	return strings.ToLower(entryLabels[id])
}

type OrientationID glow.Orientation

var OrientationLabels = []string{
	"Level",
	"Upright",
	"Tilted",
}

func (id OrientationID) String() string {
	return OrientationLabels[id]
}

func (id OrientationID) PlaceHolder() string {
	return strings.ToLower(OrientationLabels[id])
}

type OriginID glow.Origin

var OriginLabels = []string{
	"Top Left",
	"Top Right",
	"Bottom Left",
	"Bottom Right",
}

func (id OriginID) String() string {
	return OriginLabels[id]
}

func (id OriginID) PlaceHolder() string {
	return strings.ToLower(OriginLabels[id])
}

type AppIconID int

const (
	FrameIcon AppIconID = iota
	LayerIcon
	HueShiftIcon
	ScanIcon
	BeginIcon
	EndIcon
	EffectsIcon
	APP_ICON_COUNT
)

func NewAppIcon(i AppIconID) (w *widget.Icon) {
	if int(i) >= len(appResoures) {
		i = 0
	}
	w = widget.NewIcon(appResoures[i])
	return
}

var appIconFiles = []string{
	"frame.svg",
	"layer.svg",
	"hue_shift.svg",
	"scan.svg",
	"begin.svg",
	"end.svg",
	"effect.svg",
}

var appResoures = make([]fyne.Resource, int(APP_ICON_COUNT))

func AppIconResource(i AppIconID) fyne.Resource {
	return appResoures[i]
}

const (
	GridBottomLeftHorizontal uint16 = iota
	GridBottomLeftVertical
	GridBottomLeftDiagonal

	GridBottomRightHorizontal
	GridBottomRightVertical
	GridBottomRightDiagonal

	GridTopLeftHorizontal
	GridTopRightVertical
	GridTopLeftDiagonal

	GridTopLeftVertical
	GridTopRightDiagonal
	GridTopRightHorizontal
	GRID_ICON_COUNT
)

var gridIconFiles = []string{
	"top_left_horizontal.svg",
	"top_left_vertical.svg",
	"top_left_diagonal.svg",

	"top_right_horizontal.svg",
	"top_right_vertical.svg",
	"top_right_diagonal.svg",

	"bottom_left_horizontal.svg",
	"bottom_left_vertical.svg",
	"bottom_left_diagonal.svg",

	"bottom_right_horizontal.svg",
	"bottom_right_vertical.svg",
	"bottom_right_diagonal.svg",
}

var gridResoures = make([]fyne.Resource,
	int(glow.ORIGIN_COUNT)*int(glow.ORIENTATION_COUNT))

func NewGridIcon(origin, orientation int) (w *widget.Icon) {
	i := origin*int(glow.ORIGIN_COUNT-1) + orientation
	if i >= len(gridResoures) {
		i = 0
	}
	w = widget.NewIcon(gridResoures[i])
	return
}

func LoadGridIcons(theme string) (err error) {
	var res fyne.Resource

	for i := 0; i < len(gridIconFiles); i++ {
		res, err = fyne.LoadResourceFromPath("res/icons/" + theme + "/" + gridIconFiles[i])
		if err != nil {
			return
		}
		gridResoures[i] = res
	}

	for i := 0; i < len(appIconFiles); i++ {
		res, err = fyne.LoadResourceFromPath("res/icons/" + theme + "/" + appIconFiles[i])
		if err != nil {
			return
		}
		appResoures[i] = res
	}

	return
}
