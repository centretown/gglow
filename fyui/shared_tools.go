package fyui

// import (
// 	"gglow/fyio"

// 	"fyne.io/fyne/v2/data/binding"
// 	"fyne.io/fyne/v2/theme"
// 	"fyne.io/fyne/v2/widget"
// )

// type SharedTools struct {
// 	*widget.Toolbar
// 	saveButton *ButtonItem
// 	effect     *fyio.EffectIo
// }

// func NewSharedTools(effect *fyio.EffectIo) *SharedTools {
// 	tl := &SharedTools{
// 		Toolbar: widget.NewToolbar(),
// 		effect:  effect,
// 	}

// 	tl.saveButton = NewButtonItem(
// 		widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func() {
// 			tl.effect.SaveEffect()
// 		}))

// 	tl.effect.AddChangeListener(binding.NewDataListener(func() {
// 		if tl.effect.HasChanged() {
// 			tl.saveButton.Enable()
// 			return
// 		}
// 		tl.saveButton.Disable()
// 	}))

// 	tl.AddItems(tl.saveButton)
// 	return tl
// }

// // func (tl *SharedTools) AddItems(items ...widget.ToolbarItem) {
// // 	tl.Toolbar.Items = append(tl.Toolbar.Items, items...)
// // }
