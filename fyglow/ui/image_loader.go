package ui

import (
	"gglow/fyglow/effectio"
	"gglow/glow"
	"gglow/pic"
	"gglow/text"
	"image"
	"image/color"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/disintegration/imaging"
)

const (
	ViewWidth     = 400
	ViewHeight    = 400
	DesiredWidth  = 720
	DesiredHeight = 720
)

type ImageLoader struct {
	*dialog.CustomDialog
	window        fyne.Window
	effect        *effectio.EffectIo
	layer         *glow.Layer
	picture       *image.NRGBA
	view          *canvas.Image
	pickLabel     *widget.Label
	pickDialog    *dialog.FileDialog
	path          string
	height, width int
	filter        pic.ResampleItem
}

func NewImageLoader(effect *effectio.EffectIo, window fyne.Window) *ImageLoader {

	ld := &ImageLoader{
		window: window,
		effect: effect,
		height: ViewWidth,
		width:  ViewHeight,
		filter: pic.Box,
	}

	ld.pickDialog = imagePicker(ld.updatePath, window)
	ld.pickLabel = widget.NewLabel(ld.path)
	pickButton := widget.NewButton("Load...", ld.pickClick)

	ld.view = canvas.NewImageFromImage(imaging.New(ld.width, ld.height, color.Black))
	ld.view.Resize(DesiredSize(float32(ld.width), float32(ld.height),
		window.Canvas().Size()))

	sel := widget.NewSelect(pic.ResampleList, nil)
	sel.SetSelectedIndex(int(ld.filter))
	sel.OnChanged = func(string) { ld.filter = pic.ResampleItem(sel.SelectedIndex()) }

	content := container.NewBorder(container.NewCenter(ld.pickLabel), nil,
		sel, nil, ld.view)
	ld.CustomDialog = dialog.NewCustom(text.ImageLoad.String(), "", content, window)

	applyButton := widget.NewButton(text.ApplyLabel.String(), ld.apply)
	cancelButton := widget.NewButton(text.CancelLabel.String(), func() {
		ld.CustomDialog.Hide()
	})
	ld.SetButtons([]fyne.CanvasObject{cancelButton, pickButton, applyButton})
	return ld
}

func (ld *ImageLoader) pickClick() {
	var (
		listUri fyne.ListableURI
	)
	dir := filepath.Dir(ld.path)
	uri, err := storage.ParseURI("file://" + dir)
	if err == nil {
		listUri, err = storage.ListerForURI(uri)
		if err == nil {
			ld.pickDialog.SetLocation(listUri)
		}
	}
	if err != nil {
		fyne.LogError("LOAD IMAGE", err)
	}
	ld.pickDialog.Resize(DesiredSize(DesiredWidth, DesiredHeight,
		ld.window.Canvas().Size()))
	ld.pickDialog.Show()

}

func (ld *ImageLoader) updatePath(path string, refresh bool) {
	var err error
	ld.path = path
	ld.pickLabel.SetText(path)
	if len(path) > 1 {
		ld.picture, err = pic.ResamplePath(path, ld.height, ld.width, ld.filter.Filter())
		if err != nil {
			fyne.LogError("Image Loader", err)
		}
	}
	if ld.picture == nil {
		ld.picture = imaging.New(ld.width, ld.height, color.Black)
	}

	ld.view.Image = ld.picture
	if refresh {
		ld.view.Refresh()
	}
}

func DesiredSize(w, h float32, max fyne.Size) fyne.Size {
	dz := fyne.NewSize(w, h)
	if dz.Height > max.Height {
		dz.Height = max.Height
	}
	if dz.Width > max.Width {
		dz.Width = max.Width
	}
	return dz
}

func (ld *ImageLoader) Start() {
	ld.layer = ld.effect.GetCurrentLayer()
	ld.updatePath(ld.layer.ImageName, false)
	ld.view.Image = ld.picture
	dz := DesiredSize(DesiredWidth, DesiredHeight, ld.window.Canvas().Size())
	ld.Resize(dz)
	ld.Show()
}

func imagePicker(update func(string, bool), window fyne.Window) *dialog.FileDialog {
	dlg := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
		if err != nil || uc == nil {
			// fyne.LogError("Image loader", err)
			return
		}
		uc.Close()
		update(uc.URI().Path(), true)
	}, window)
	return dlg
}

func (ld *ImageLoader) apply() {
	ld.CustomDialog.Hide()
	layer := ld.effect.GetCurrentLayer()
	layer.ImageName = ld.path
	ld.effect.SetChanged()
}
