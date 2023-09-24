package ui

import (
	"fmt"
	"glow-gui/glow"
	"glow-gui/res"
	"glow-gui/store"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	widgetx "fyne.io/x/fyne/widget"
)

type Ui struct {
	frame *glow.Frame

	mainBox *fyne.Container
	strip   *LightStrip
	title   *widget.Label
	toolBox *widget.Toolbar
	list    *widget.List

	stop       chan int
	isSpinning bool

	window fyne.Window
	app    fyne.App

	appImage  *canvas.Image
	frameIcon *widget.Icon
	layerIcon *widget.Icon
}

func NewUi(app fyne.App, window fyne.Window) *Ui {
	gui := &Ui{
		frame:  &glow.Frame{},
		window: window,
		app:    app,
	}
	window.SetContent(gui.buildContent())
	return gui
}

func (ui *Ui) stopSpinner() {
	if ui.isSpinning {
		ui.stop <- 0
		ui.isSpinning = false
	}
}

func (ui *Ui) startSpinner() {
	ui.stopSpinner()

	frame := *ui.frame

	spin := func() {
		for {
			select {
			case <-ui.stop:
				return
			default:
				frame.Spin(ui.strip)
				time.Sleep(time.Duration(ui.frame.Interval) * 4 * time.Millisecond)
			}
		}
	}
	ui.isSpinning = true
	go spin()
}

func (ui *Ui) buildContent() *fyne.Container {
	ui.appImage = canvas.NewImageFromFile(res.GooseNoirImage.String())
	ui.frameIcon = widget.NewIcon(theme.DocumentIcon())
	ui.layerIcon = widget.NewIcon(theme.GridIcon())
	// title := widget.NewRichTextWithText(res.ChooseEffectLabel.String())
	ui.title = widget.NewLabel(res.LightEffectLabel.String())
	ui.title.Alignment = fyne.TextAlignCenter
	ui.title.TextStyle = fyne.TextStyle{Bold: true}

	ui.strip = NewLightStrip(50, 5, 16)
	ui.toolBox = ui.toolbar()
	tc := container.New(layout.NewCenterLayout(), ui.toolBox)

	ui.list = ui.NewLayerList()
	top := container.NewVBox(ui.title, ui.strip, tc)

	ui.mainBox = container.NewBorder(top, nil, nil, nil, ui.list)
	ui.stop = make(chan int)

	return ui.mainBox
}

func (ui *Ui) SetTitle(title string) {
	ui.title.SetText(title)
}

func (ui *Ui) SetWindowTitle(title string) {
	ui.window.SetTitle(title)
}

func (ui *Ui) SetFrame(frame *glow.Frame) {
	ui.frame = frame
	frame.Setup(uint16(ui.strip.length),
		uint16(ui.strip.rows),
		uint32(ui.strip.interval))
}

func (ui *Ui) toolbar() (t *widget.Toolbar) {
	var items []widget.ToolbarItem

	choose := NewToolbarSelect(ui)
	create := widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {})
	sep := widget.NewToolbarSeparator()
	play := widget.NewToolbarAction(theme.MediaPlayIcon(), func() {
		ui.startSpinner()
	})
	stop := widget.NewToolbarAction(theme.MediaStopIcon(), func() {
		ui.stopSpinner()
	})
	upload := widget.NewToolbarAction(theme.UploadIcon(), func() {})
	settings := widget.NewToolbarAction(theme.SettingsIcon(), func() {})

	items = append(items,
		choose,
		create,
		sep,
		play,
		stop,
		sep,
		upload,
		sep,
		settings)

	t = widget.NewToolbar(items...)

	return
}

func (ui *Ui) OnChangeFrame(s string) {
	uri, err := store.LookupURI(s)
	if err != nil {
		return
	}
	frame := &glow.Frame{}
	store.LoadFrameURI(uri, frame)
	ui.SetFrame(frame)
	ui.SetWindowTitle(res.WindowTitle + " - " + uri.Name())
	ui.SetTitle(res.LightEffectLabel.String() + " - " + s)
	ui.list.Refresh()
}

func (ui *Ui) NewLayerList() *widget.List {
	list := widget.NewList(
		// Length
		func() int {
			return len(ui.frame.Layers)
		},
		// CreateItem
		func() fyne.CanvasObject {
			icon := ui.layerIcon
			return &fyne.Container{
				Layout:  layout.NewBorderLayout(nil, nil, icon, nil),
				Objects: []fyne.CanvasObject{icon, widget.NewLabel("Template Object")}}
		},
		// UpdateItem
		func(id widget.ListItemID, item fyne.CanvasObject) {
			layer := &ui.frame.Layers[id]
			buf := fmt.Sprintf("#%d l:%d r:%d s:%d b%d e%d", id,
				layer.Length, layer.Rows, layer.Scan, layer.Begin, layer.End)
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(buf)
		})

	return list
}

func (ui *Ui) LightInfo() *fyne.Container {
	format := func(id res.LabelID, min, max float64,
		boundf *float64) (slider *widget.Slider,
		label *widget.Label, entry *widgetx.NumericalEntry) {

		boundValue := binding.BindFloat(boundf)
		slider = widget.NewSliderWithData(min, max, boundValue)

		entry = widgetx.NewNumericalEntry()

		// widget.NewEntryWithData(boundString)
		label = widget.NewLabel(id.String())
		return
	}

	grid := container.NewGridWithColumns(3)
	slider, label, entry := format(res.LengthLabel, 10, 100, &ui.strip.length)
	grid.Add(label)
	grid.Add(slider)
	grid.Add(entry)

	slider, label, entry = format(res.RowsLabel, 1, 10, &ui.strip.rows)
	grid.Add(label)
	grid.Add(slider)
	grid.Add(entry)

	slider, label, entry = format(res.IntervalLabel, 1, 1000, &ui.strip.interval)
	grid.Add(label)
	grid.Add(slider)
	grid.Add(entry)

	return grid
}

func ColorDialog(window fyne.Window) *fyne.Container {
	hueLabel := widget.NewLabel(res.HueLabel.String())
	hueEntry := widget.NewEntry()
	hueEntry.SetPlaceHolder(res.HueLabel.PlaceHolder())

	label2 := widget.NewLabel(res.SaturationLabel.String())
	value2 := widget.NewEntry()
	value2.SetPlaceHolder(res.SaturationLabel.PlaceHolder())

	button1 := widget.NewButton("Choose Color...", func() {
		picker := dialog.NewColorPicker("Color Picker", "color", func(c color.Color) {

		}, window)
		picker.Advanced = true
		picker.Show()
	})
	grid := container.New(layout.NewVBoxLayout(),
		hueLabel, hueEntry, label2, value2, button1)
	return grid
}
