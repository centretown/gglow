package effects

import (
	"fmt"
	"glow-gui/glow"
	"glow-gui/glowio"
	"glow-gui/settings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

var _ Effect = (*EffectIo)(nil)

const Dots = ".."

func defaultFrame() (frame *glow.Frame) {
	frame = &glow.Frame{}
	frame.Interval = 48
	var layer glow.Layer
	layer.Chroma.Colors = append(layer.Chroma.Colors,
		glow.HSV{Hue: 0, Saturation: 0, Value: 100})
	frame.Layers = append(frame.Layers, layer)
	return
}

type EffectIo struct {
	glowio.IoHandler
	effectName       string
	folderName       string
	Frame            binding.Untyped
	Layer            binding.Untyped
	LayerSummaryList binding.StringList
	layerIndex       int

	preferences fyne.Preferences

	// history *History
	// CanUndo binding.Bool
	// CanSave binding.Bool
	backup     *glow.Frame
	HasBackup  bool
	hasChanged binding.Bool
	isActive   bool

	saveActions []func(*glow.Frame)
	config      *settings.Configuration
	// changeDetected bool
}

func NewEffectIo(io glowio.IoHandler, preferences fyne.Preferences, config *settings.Configuration) *EffectIo {

	eff := &EffectIo{
		IoHandler:        io,
		Frame:            binding.NewUntyped(),
		Layer:            binding.NewUntyped(),
		LayerSummaryList: binding.NewStringList(),
		hasChanged:       binding.NewBool(),
		// CanUndo:    binding.NewBool(),
		// CanSave:    binding.NewBool(),

		preferences: preferences,
		backup:      &glow.Frame{},
		// history:     NewHistory(),
		saveActions: make([]func(*glow.Frame), 0),
		config:      config,
	}
	// fmt.Println(config.Method, config.Path, config.Folder, config.Effect)

	folder := config.Folder
	effect := config.Effect

	// folder := preferences.StringWithFallback(settings.EffectFolder.String(), "")
	if folder != "" {
		eff.RefreshFolder(folder)
	}
	// effect := preferences.StringWithFallback(settings.Effect.String(), "")
	if len(effect) > 0 {
		eff.LoadEffect(effect)
	} else {
		eff.setFrame(defaultFrame(), 0)
	}

	eff.AddChangeListener(binding.NewDataListener(func() {
		if eff.HasChanged() {
			// f := eff.GetFrame()
			// fmt.Println("hasChanged makeBackup", f.Interval)
			eff.makeBackup(true)
		}
	}))

	eff.AddFrameListener(binding.NewDataListener(eff.onChangeFrame))
	return eff
}

func (eff *EffectIo) LoadFolder(folder string) []string {
	keys, _ := eff.IoHandler.RefreshFolder(folder)
	return keys
}

func (st *EffectIo) SummaryList() []string {
	l, _ := st.LayerSummaryList.Get()
	return l
}

func (st *EffectIo) LayerIndex() int {
	return st.layerIndex
}

func (st *EffectIo) SetActive() {
	st.isActive = true
}

func (st *EffectIo) setFrame(frame *glow.Frame, layerIndex int) {
	st.Frame.Set(frame)
	st.SetCurrentLayer(layerIndex)
	st.hasChanged.Set(false)
}

func (st *EffectIo) GetFrame() *glow.Frame {
	f, _ := st.Frame.Get()
	return f.(*glow.Frame)
}

func (st *EffectIo) SetCurrentLayer(i int) {
	frame := st.GetFrame()
	var layer *glow.Layer
	if i < len(frame.Layers) {
		st.layerIndex = i
		layer = &frame.Layers[i]
	} else {
		st.layerIndex = 0
		layer = &glow.Layer{}
	}
	st.Layer.Set(layer)
}

func (st *EffectIo) GetCurrentLayer() *glow.Layer {
	layer, _ := st.Layer.Get()
	return layer.(*glow.Layer)
}

func (st *EffectIo) Undo(title string) {

	if st.HasBackup {
		frame := st.backup
		st.setFrame(frame, st.layerIndex)
		st.makeBackup(false)
	}

	fyne.LogError("UndoEffect", fmt.Errorf("nothing to undo"))
}

func (st *EffectIo) SetChanged() {
	if !st.isActive {
		return
	}
	st.hasChanged.Set(true)
}

func (st *EffectIo) AddFrameListener(listener binding.DataListener) {
	st.Frame.AddListener(listener)
}

func (st *EffectIo) AddLayerListener(listener binding.DataListener) {
	st.Layer.AddListener(listener)
}

func (st *EffectIo) AddChangeListener(listener binding.DataListener) {
	st.hasChanged.AddListener(listener)
}

func (m *EffectIo) HasChanged() bool {
	b, _ := m.hasChanged.Get()
	return b
}

func (eff *EffectIo) OnExit() {
	eff.IoHandler.OnExit()
	eff.preferences.SetStringList(eff.config.Driver,
		[]string{eff.config.Path, eff.folderName, eff.effectName})
}

func (eff *EffectIo) LoadEffect(title string) error {
	frame, err := eff.IoHandler.ReadEffect(title)
	if err != nil {
		return err
	}

	eff.effectName = title
	eff.folderName = eff.IoHandler.FolderName()

	eff.setFrame(frame, 0)
	eff.makeBackup(false)
	eff.hasChanged.Set(false)
	return nil
}

func (st *EffectIo) SaveEffect() error {
	title, frame := st.EffectName(), st.GetFrame()
	err := ValidateEffectName(title)
	if err != nil {
		return err
	}

	err = st.IoHandler.WriteEffect(title, frame)
	if err != nil {
		return err
	}

	st.hasChanged.Set(false)
	return err
}

func (st *EffectIo) makeBackup(b bool) {
	st.HasBackup = b
	if b {
		st.backup, _ = glow.FrameDeepCopy(st.GetFrame())
	} else {
		st.backup = &glow.Frame{}
	}
}

func (st *EffectIo) Apply() {
	frame := st.GetFrame()

	for _, saveAction := range st.saveActions {
		saveAction(frame)
	}

	current := *frame
	st.Frame.Set(&current)
}

func (st *EffectIo) OnApply(f func(*glow.Frame)) {
	st.saveActions = append(st.saveActions, f)
}

func (st *EffectIo) EffectName() string {
	return st.effectName
}

func (st *EffectIo) UndoEffect() {
	st.Undo(st.effectName)
	st.Apply()
}

func (st *EffectIo) CanUndo() bool {
	return st.HasBackup
}

func (st *EffectIo) onChangeFrame() {
	frame := st.GetFrame()
	summaries := make([]string, 0, len(frame.Layers))
	for i, layer := range frame.Layers {
		summaries = append(summaries, SummarizeLayer(&layer, i+1))
	}
	st.LayerSummaryList.Set(summaries)
}
