package fyui

// import (
// 	"gglow/fyio"

// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/data/binding"
// 	"fyne.io/fyne/v2/widget"
// )

// type FolderTools struct {
// 	*fyne.Container
// 	folderName binding.String
// 	effectName binding.String
// }

// func NewFolderTools(effect *fyio.EffectIo, window fyne.Window, menu *fyne.Menu) *FolderTools {
// 	ft := &FolderTools{
// 		Container:  container.NewHBox(),
// 		folderName: binding.NewString(),
// 		effectName: binding.NewString(),
// 	}

// 	folderLabel := widget.NewLabelWithData(ft.folderName)
// 	effectLabel := widget.NewLabelWithData(ft.effectName)

// 	effect.AddFrameListener(binding.NewDataListener(func() {
// 		ft.folderName.Set(effect.FolderName())
// 		ft.effectName.Set(effect.EffectName())
// 	}))
// 	effect.AddFolderListener(binding.NewDataListener(func() {
// 		ft.folderName.Set(effect.FolderName())
// 	}))

// 	ft.Container.Objects = append(ft.Container.Objects, folderLabel, effectLabel)
// 	return ft
// }
