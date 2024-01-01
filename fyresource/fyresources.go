package fyresource

import (
	"gglow/glow"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

var AppID = "com.centretown.glow.preferences"
var WindowSize = fyne.Size{Width: 480, Height: 480}

type ImageID uint16

const (
	DarkGanderImage ImageID = iota
)

var imagePath = []string{
	"resources/dark-gander.png",
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

type AppIconID int

const (
	FrameIcon AppIconID = iota
	LayerIcon
	HueShiftIcon
	ScanIcon
	BeginIcon
	EndIcon
	EffectsIcon
	SpeedIcon
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
	"speed.svg",
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

var gridIconPath = "resources/icons/"

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

func LoadIcons(theme string) (err error) {
	var res fyne.Resource

	for i := 0; i < len(gridIconFiles); i++ {
		res, err = fyne.LoadResourceFromPath(path.Join(gridIconPath, theme, gridIconFiles[i]))
		// res, err = fyne.LoadResourceFromPath(gridIconPath + theme + "/" + gridIconFiles[i])
		if err != nil {
			return
		}
		gridResoures[i] = res
	}

	for i := 0; i < len(appIconFiles); i++ {
		res, err = fyne.LoadResourceFromPath(gridIconPath + theme + "/" + appIconFiles[i])
		if err != nil {
			return
		}
		appResoures[i] = res
	}

	return
}
