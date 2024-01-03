package fyio

import (
	"fmt"
	"gglow/glow"
	"gglow/iohandler"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

var _ iohandler.EffectIoHandler = (*EffectIo)(nil)

const Dots = ".."

type EffectIo struct {
	iohandler.IoHandler
	config           *iohandler.Accessor
	Frame            binding.Untyped
	Layer            binding.Untyped
	LayerSummaryList binding.StringList
	effectName       string
	folderName       string
	layerIndex       int
	hasChanged       binding.Bool
	isActive         bool
	saveActions      []func(*glow.Frame)
}

func NewEffectIo(io iohandler.IoHandler, preferences fyne.Preferences, config *iohandler.Accessor) *EffectIo {

	eff := &EffectIo{
		IoHandler:        io,
		Frame:            binding.NewUntyped(),
		Layer:            binding.NewUntyped(),
		LayerSummaryList: binding.NewStringList(),
		hasChanged:       binding.NewBool(),
		saveActions:      make([]func(*glow.Frame), 0),
		config:           config,
	}

	folder := config.Folder
	effect := config.Effect

	if folder != "" {
		eff.SetFolder(folder)
	}
	if len(effect) > 0 {
		eff.LoadEffect(effect)
	} else {
		eff.setFrame(glow.NewFrame(), 0)
	}

	eff.AddFrameListener(binding.NewDataListener(eff.onChangeFrame))
	return eff
}

func (eff *EffectIo) LoadFolder(folder string) []string {
	keys, _ := eff.IoHandler.SetFolder(folder)
	return keys
}

func (eff *EffectIo) SummaryList() []string {
	l, _ := eff.LayerSummaryList.Get()
	return l
}

func (eff *EffectIo) LayerIndex() int {
	return eff.layerIndex
}

func (eff *EffectIo) SetActive() {
	eff.isActive = true
}

func (eff *EffectIo) setFrame(frame *glow.Frame, layerIndex int) {
	eff.Frame.Set(frame)
	eff.SetCurrentLayer(layerIndex)
	eff.SetUnchanged()
}

func (eff *EffectIo) GetFrame() *glow.Frame {
	f, _ := eff.Frame.Get()
	return f.(*glow.Frame)
}

func (eff *EffectIo) SetCurrentLayer(i int) {
	frame := eff.GetFrame()
	var layer *glow.Layer
	if i < len(frame.Layers) {
		eff.layerIndex = i
		layer = &frame.Layers[i]
	} else {
		eff.layerIndex = 0
		layer = &glow.Layer{}
	}
	eff.Layer.Set(layer)
}

func (eff *EffectIo) GetCurrentLayer() *glow.Layer {
	layer, _ := eff.Layer.Get()
	return layer.(*glow.Layer)
}

func (eff *EffectIo) AddFrameListener(listener interface{}) {
	eff.Frame.AddListener(listener.(binding.DataListener))
}

func (eff *EffectIo) AddLayer(*glow.Layer) (err error)  { return }
func (eff *EffectIo) AddEffect(*glow.Frame) (err error) { return }

func (eff *EffectIo) AddLayerListener(listener interface{}) {
	eff.Layer.AddListener(listener.(binding.DataListener))
}

func (eff *EffectIo) SetChanged() {
	if !eff.isActive {
		return
	}
	eff.hasChanged.Set(true)
}

func (eff *EffectIo) SetUnchanged() {
	eff.hasChanged.Set(false)
}

func (eff *EffectIo) AddChangeListener(listener interface{}) {
	eff.hasChanged.AddListener(listener.(binding.DataListener))
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
	//apply changes
	frame := eff.GetFrame()
	for _, saveAction := range eff.saveActions {
		saveAction(frame)
	}
	current := *frame
	eff.Frame.Set(&current)

	title, frame := eff.EffectName(), eff.GetFrame()
	err := ValidateEffectName(title)
	if err != nil {
		return err
	}

	err = eff.IoHandler.WriteEffect(title, frame)
	if err != nil {
		return err
	}

	eff.SetUnchanged()
	return err
}

func (eff *EffectIo) OnSave(f func(*glow.Frame)) {
	eff.saveActions = append(eff.saveActions, f)
}

func (eff *EffectIo) EffectName() string {
	return eff.effectName
}

func (eff *EffectIo) onChangeFrame() {
	fmt.Printf("onChangeFrame")
	frame := eff.GetFrame()
	summaries := make([]string, 0, len(frame.Layers))
	for i, layer := range frame.Layers {
		summaries = append(summaries, SummarizeLayer(&layer, i+1))
	}
	eff.LayerSummaryList.Set(summaries)
}
