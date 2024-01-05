package fyio

import (
	"fmt"
	"gglow/glow"
	"gglow/iohandler"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

type EffectIo struct {
	iohandler.IoHandler

	config *iohandler.Accessor
	frame  *glow.Frame
	layer  *glow.Layer

	folderName string
	effectName string
	layerIndex int

	folderWatch binding.Int
	frameWatch  binding.Int
	layerWatch  binding.Int
	hasChanged  binding.Bool
	summaryList []string

	isActive    bool
	saveActions []func(*glow.Frame)
}

func NewEffect(io iohandler.IoHandler, preferences fyne.Preferences, config *iohandler.Accessor) *EffectIo {

	eff := &EffectIo{
		IoHandler:   io,
		folderWatch: binding.NewInt(),
		frameWatch:  binding.NewInt(),
		layerWatch:  binding.NewInt(),

		hasChanged:  binding.NewBool(),
		saveActions: make([]func(*glow.Frame), 0),
		summaryList: make([]string, 0),
		config:      config,
	}

	eff.frame = glow.NewFrame()
	eff.layer = &eff.frame.Layers[0]

	folder := config.Folder
	effect := config.Effect

	if folder != "" {
		eff.SetFolder(folder)
	}
	if len(effect) > 0 {
		eff.LoadEffect(effect)
	}
	return eff
}

func (eff *EffectIo) SummaryList() []string {
	return eff.summaryList
}

func (eff *EffectIo) LayerIndex() int {
	return eff.layerIndex
}

func (eff *EffectIo) SetActive() {
	eff.isActive = true
}

func (eff *EffectIo) alert(x binding.Int) {
	a, _ := x.Get()
	x.Set(a + 1)
}

func (eff *EffectIo) alertFolder() {
	eff.alert(eff.folderWatch)
}

func (eff *EffectIo) LoadFolder(folder string) []string {
	ls, err := eff.IoHandler.SetFolder(folder)
	if err != nil {
		fyne.LogError("loadfolder", err)
		return ls
	}
	eff.alertFolder()
	return ls
}

func (eff *EffectIo) alertFrame() {
	eff.alert(eff.frameWatch)
}

func (eff *EffectIo) alertLayer() {
	eff.alert(eff.layerWatch)
}

func (eff *EffectIo) GetFrame() *glow.Frame {
	return eff.frame
}

func (eff *EffectIo) setFrame(frame *glow.Frame, layerIndex int) {
	eff.frame = frame

	summaries := make([]string, 0, len(frame.Layers))
	for i, layer := range frame.Layers {
		summaries = append(summaries, SummarizeLayer(&layer, i+1))
	}
	eff.summaryList = summaries

	eff.setLayer(layerIndex)
	eff.alertFrame()
	eff.alertLayer()
	eff.SetUnchanged()
}

func (eff *EffectIo) setLayer(index int) {
	lCount := len(eff.frame.Layers)
	if lCount == 0 {
		fyne.LogError("setLayer", fmt.Errorf("no Layers in frame"))
		os.Exit(1)
		return
	}

	if index >= lCount {
		index = lCount - 1
	}
	eff.layerIndex = index
	eff.layer = &eff.frame.Layers[index]
}

func (eff *EffectIo) SetCurrentLayer(i int) {
	eff.setLayer(i)
	eff.alertLayer()
}

func (eff *EffectIo) GetCurrentLayer() *glow.Layer {
	return eff.layer
}

func (eff *EffectIo) AddFolderListener(listener binding.DataListener) {
	eff.folderWatch.AddListener(listener)
}
func (eff *EffectIo) AddFrameListener(listener binding.DataListener) {
	eff.frameWatch.AddListener(listener)
}
func (eff *EffectIo) AddLayerListener(listener binding.DataListener) {
	eff.layerWatch.AddListener(listener)
}
func (eff *EffectIo) AddChangeListener(listener binding.DataListener) {
	eff.hasChanged.AddListener(listener)
}

// func (eff *EffectIo) AddLayer(*glow.Layer) (err error)  { return }
// func (eff *EffectIo) AddEffect(*glow.Frame) (err error) { return }

func (eff *EffectIo) SetChanged() {
	if !eff.isActive {
		return
	}
	eff.hasChanged.Set(true)
}

func (eff *EffectIo) SetUnchanged() {
	eff.hasChanged.Set(false)
}

func (eff *EffectIo) HasChanged() bool {
	b, _ := eff.hasChanged.Get()
	return b
}

func (eff *EffectIo) LoadEffect(title string) error {
	frame, err := eff.IoHandler.ReadEffect(title)
	if err != nil {
		return err
	}

	eff.effectName = title
	eff.folderName = eff.IoHandler.FolderName()

	eff.setFrame(frame, 0)
	eff.SetUnchanged()
	return nil
}

func (eff *EffectIo) SaveEffect() error {
	title := eff.EffectName()
	err := ValidateEffectName(title)
	if err != nil {
		return err
	}

	//apply changes
	for _, saveAction := range eff.saveActions {
		saveAction(eff.frame)
	}

	err = eff.IoHandler.WriteEffect(title, eff.frame)
	if err != nil {
		return err
	}

	eff.setFrame(eff.frame, eff.layerIndex)
	eff.SetUnchanged()
	return err
}

func (eff *EffectIo) OnSave(f func(*glow.Frame)) {
	eff.saveActions = append(eff.saveActions, f)
}

func (eff *EffectIo) EffectName() string {
	return eff.effectName
}
