package resource

import (
	"fyne.io/fyne/v2"
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
	IconFrameID uint = iota
	IconFrameAddID
	IconFrameRemoveID
	IconLayerInsertID
	IconEffectID
	IconFileShareID
	IconExitID
	ICON_COUNT
)

var appResoures = make([]fyne.Resource, int(ICON_COUNT))

func IconLayer() fyne.Resource       { return appResoures[IconFrameID] }
func IconFrameAdd() fyne.Resource    { return appResoures[IconFrameAddID] }
func IconFrameRemove() fyne.Resource { return appResoures[IconFrameRemoveID] }
func IconLayerInsert() fyne.Resource { return appResoures[IconLayerInsertID] }
func IconEffect() fyne.Resource      { return appResoures[IconEffectID] }
func IconFileShare() fyne.Resource   { return appResoures[IconFileShareID] }
func IconExit() fyne.Resource        { return appResoures[IconExitID] }
func Icon(i uint) fyne.Resource {
	if i >= ICON_COUNT {
		i = 0
	}
	return appResoures[i]
}

var header1 = `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 0 24 24" width="24px" `
var header2 = `<path d="M0 0h24v24H0V0z" fill="none" />`
var icons = []string{
	`<path
    d="M4 6H2v14c0 1.1.9 2 2 2h14v-2H4V6zm16-4H8c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zm0 14H8V4h12v12z" />
</svg>`,
	`<path
    d="M4 6H2v14c0 1.1.9 2 2 2h14v-2H4V6zm16-4H8c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zm0 14H8V4h12v12zm-7-2h2v-3h3V9h-3V6h-2v3h-3v2h3z" />
</svg>`,
	`<path
    d="M4 6H2v14c0 1.1.9 2 2 2h14v-2H4V6zm16-4H8c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zm0 14H8V4h12v12zM18 11H10v-2h8v2z" />
</svg>`,
	`<path
     d="m 17.82584,5.3380582 -1.41,-1.41 -6,6 6,5.9999998 1.41,-1.41 -4.58,-4.5899998 z M 4,6 H 2 v 14 c 0,1.1 0.9,2 2,2 H 18 V 20 H 4 Z M 20,2 H 8 C 6.9,2 6,2.9 6,4 v 12 c 0,1.1 0.9,2 2,2 h 12 c 1.1,0 2,-0.9 2,-2 V 4 C 22,2.9 21.1,2 20,2 Z m 0,14 H 8 V 4 h 12 z" />
</svg>`,
	`<g><rect fill="none" height="24" width="24" x="0"/></g><g><g><polygon points="19,9 20.25,6.25 23,5 20.25,3.75 19,1 17.75,3.75 15,5 17.75,6.25"/><polygon points="19,15 17.75,17.75 15,19 17.75,20.25 19,23 20.25,20.25 23,19 20.25,17.75"/><path d="M11.5,9.5L9,4L6.5,9.5L1,12l5.5,2.5L9,20l2.5-5.5L17,12L11.5,9.5z M9.99,12.99L9,15.17l-0.99-2.18L5.83,12l2.18-0.99 L9,8.83l0.99,2.18L12.17,12L9.99,12.99z"/></g></g>
</svg>`,
	`<g><path d="M0,0h24v24H0V0z" fill="none"/></g><g><g><path d="M6,5H4v16c0,1.1,0.9,2,2,2h10v-2H6V5z"/><path d="M18,1h-8C8.9,1,8,1.9,8,3v14c0,1.1,0.9,2,2,2h8c1.1,0,2-0.9,2-2V3C20,1.9,19.1,1,18,1z M18,17h-8v-1h8V17z M18,14h-8V6h8 V14z M18,4h-8V3h8V4z"/><path d="M12.5,10.25h1.63l-0.69,0.69L14.5,12L17,9.5L14.5,7l-1.06,1.06l0.69,0.69H12c-0.55,0-1,0.45-1,1V12h1.5V10.25z"/></g></g>
</svg>`,
	`<path d="M10.09 15.59L11.5 17l5-5-5-5-1.41 1.41L12.67 11H3v2h9.67l-2.58 2.59zM19 3H5c-1.11 0-2 .9-2 2v4h2V5h14v14H5v-4H3v4c0 1.1.89 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2z"/>
</svg>`,
}

var iconNames = []string{
	"layer",
	"layerAdd",
	"layerRemove",
	"layerInsert",
	"effect",
	"fileShare",
	"exit",
}

func makeSVG(i uint, fill string) []byte {
	return []byte(header1 + fill + header2 + icons[i])
}

func LoadIcons(theme string) {
	var fill string
	if theme == "light" {
		fill = `fill="#000000">`
	} else {
		fill = `fill="#ffffff">`
	}

	for i := uint(0); i < ICON_COUNT; i++ {
		appResoures[i] = fyne.NewStaticResource(iconNames[i],
			makeSVG(i, fill))
	}
}
