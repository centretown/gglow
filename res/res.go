package res

import (
	"glow-gui/glow"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

var AppID = "com.centretown.glow.preferences"
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
	GlowEffectsLabel
)

var entryLabels = []string{
	"Length", "Rows", "Interval", "Layers",
	"Grid", "Chroma", "Shift", "Scan",
	"Begin", "End",
	"Origin", "Orientation",
	"Colors", "Hue", "Saturation", "Value",
	"pick an effect...", "Glow Effects",
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
