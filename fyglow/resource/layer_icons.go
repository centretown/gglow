package resource

// import (
// 	"fyne.io/fyne/v2"
// )

// var icon0 = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
// <svg
//    height="24px"
//    viewBox="0 0 24 24"
//    width="24px"
// `

// // fill="#ffffff"
// // <defs
// // id="defs10" />
// var icon1 = `
//    version="1.1"
//    id="svg6"
//    xmlns="http://www.w3.org/2000/svg"
//    xmlns:svg="http://www.w3.org/2000/svg">
//   <path d="M0 0h24v24H0V0z" fill="none" />
//   <path `
// var icon2 = ` />
// </svg>
// `

// // var icons = []string{
// // 	`d="M4 6H2v14c0 1.1.9 2 2 2h14v-2H4V6zm16-4H8c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zm0 14H8V4h12v12z"`,
// // 	`d="M4 6H2v14c0 1.1.9 2 2 2h14v-2H4V6zm16-4H8c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zm0 14H8V4h12v12zm-7-2h2v-3h3V9h-3V6h-2v3h-3v2h3z"`,
// // 	`d="M4 6H2v14c0 1.1.9 2 2 2h14v-2H4V6zm16-4H8c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zm0 14H8V4h12v12zM18 11H10v-2h8v2z"`,
// // }

// var header1 = `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 0 24 24" width="24px" `

// // var header3 = `fill="#ffffff">`
// var header2 = `<path d="M0 0h24v24H0V0z" fill="none" />`

// var icons = []string{
// 	`<path
//     d="M4 6H2v14c0 1.1.9 2 2 2h14v-2H4V6zm16-4H8c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zm0 14H8V4h12v12z" />
// </svg>`,
// 	`<path
//     d="M4 6H2v14c0 1.1.9 2 2 2h14v-2H4V6zm16-4H8c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zm0 14H8V4h12v12zm-7-2h2v-3h3V9h-3V6h-2v3h-3v2h3z" />
// </svg>`,
// 	`<path
//     d="M4 6H2v14c0 1.1.9 2 2 2h14v-2H4V6zm16-4H8c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zm0 14H8V4h12v12zM18 11H10v-2h8v2z" />
// </svg>`,
// 	`<path
//      id="path4-3"
//      d="m 17.82584,5.3380582 -1.41,-1.41 -6,6 6,5.9999998 1.41,-1.41 -4.58,-4.5899998 z M 4,6 H 2 v 14 c 0,1.1 0.9,2 2,2 H 18 V 20 H 4 Z M 20,2 H 8 C 6.9,2 6,2.9 6,4 v 12 c 0,1.1 0.9,2 2,2 h 12 c 1.1,0 2,-0.9 2,-2 V 4 C 22,2.9 21.1,2 20,2 Z m 0,14 H 8 V 4 h 12 z" />
// </svg>
// `,
// }

// var iconNames = []string{
// 	"layer",
// 	"layerAdd",
// 	"layerRemove",
// 	"layerInsert",
// }

// func makeSVG(i uint, fill string) []byte {
// 	return []byte(header1 + fill + header2 + icons[i])
// }

// func LoadIcons(theme string) {
// 	var fill string
// 	if theme == "light" {
// 		fill = `fill="#000000">`
// 	} else {
// 		fill = `fill="#ffffff">`
// 	}

// 	for i := uint(0); i < IconCount; i++ {
// 		// buf = makeSVG(i)
// 		iconResources[i] = fyne.NewStaticResource(iconNames[i],
// 			makeSVG(i, fill))
// 	}
// }

// func GetLayerIcon(i uint) fyne.Resource {
// 	if i >= IconCount {
// 		i = 0
// 	}
// 	return iconResources[i]
// }
