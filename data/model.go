package data

import (
	"fmt"
	"glow-gui/glow"
	"glow-gui/store"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

type Model struct {
	store            *store.Store
	Frame            binding.Untyped
	LayerIndex       int
	Layer            binding.Untyped
	LayerSummaryList binding.StringList
	isDirty          binding.Bool
	canUndo          binding.Bool
	WindowHasContent bool
}

func NewModel(store *store.Store) *Model {
	m := &Model{
		store:            store,
		Frame:            store.Frame,
		LayerSummaryList: binding.NewStringList(),
		Layer:            binding.NewUntyped(),
		isDirty:          store.IsDirty,
		canUndo:          store.HasPrevious,
	}

	m.Layer.Set(m.GetFrame().Layers[0])
	m.Frame.AddListener(binding.NewDataListener(m.onChangeFrame))
	return m
}

func (m *Model) EffectName() string {
	return m.store.EffectName
}

func (m *Model) IsFolder(title string) bool {
	return m.store.IsFolder(title)
}

func (m *Model) CreateNewEffect(title string, frame *glow.Frame) (err error) {
	return m.store.CreateNewEffect(title, frame)
}

func (m *Model) CreateNewFolder(title string) (err error) {
	return m.store.CreateNewFolder(title)
}

func (m *Model) RefreshKeys(title string) []string {
	m.store.RefreshKeys(title)
	return m.KeyList()
}

func (m *Model) ValidateNewFolderName(title string) (err error) {
	return m.store.ValidateNewFolderName(title)
}

func (m *Model) ValidateNewEffectName(title string) (err error) {
	return m.store.ValidateNewEffectName(title)
}

func (m *Model) KeyList() []string {
	return m.store.KeyList
}

func (m *Model) UpdateHistory() error {
	return m.store.UpdateHistory()
}

func (m *Model) WriteEffect() (err error) {
	frame := m.GetFrame()
	err = m.store.WriteEffect(m.EffectName(), frame)
	current := *frame
	m.Frame.Set(&current)
	return
}

func (m *Model) ReadEffect(title string) (err error) {
	err = m.store.ReadEffect(title)
	if err != nil {
		fyne.LogError("ReadEffect", err)
		return
	}

	m.LayerIndex = 0
	m.Layer.Set(m.GetFrame().Layers[m.LayerIndex])
	// fmt.Println("Model ReadEffect", title)
	return
}

func (m *Model) AddFrameListener(listener binding.DataListener) {
	m.Frame.AddListener(listener)
}

func (m *Model) AddDirtyListener(listener binding.DataListener) {
	m.isDirty.AddListener(listener)
}

func (m *Model) SetDirty() {
	if m.WindowHasContent {
		m.store.SetDirty(true)
	}
}

func (m *Model) IsDirty() bool {
	return m.store.GetDirty()
}

func (m *Model) AddUndoListener(listener binding.DataListener) {
	m.canUndo.AddListener(listener)
}

func (m *Model) CanUndo() bool {
	b, _ := m.canUndo.Get()
	return b
}

func (m *Model) UndoEffect() {
	if !m.store.CanUndo(m.EffectName()) {
		fyne.LogError("UndoEffect",
			fmt.Errorf("%s has no history", m.EffectName()))
		return
	}

	frame, err := m.store.Undo(m.EffectName())
	if err != nil {
		fyne.LogError("UndoEffect", err)
		return
	}

	m.Layer.Set(frame.Layers[m.LayerIndex])

	// fmt.Println(m.EffectName)
	// err = m.Frame.Set(frame)
	// if err != nil {
	// 	fyne.LogError("UndoEffect", err)
	// 	return
	// }
	fmt.Println("UndoEffect", frame.Interval, m.EffectName())
}

func (m *Model) onChangeFrame() {
	frame := m.GetFrame()

	summaries := make([]string, 0, len(frame.Layers))
	for i, layer := range frame.Layers {
		summaries = append(summaries, Summarize(&layer, i+1))
	}
	m.LayerSummaryList.Set(summaries)
	m.SetCurrentLayer(m.LayerIndex)
}

func (m *Model) GetFrame() *glow.Frame {
	return m.store.GetFrame()
}

func (m *Model) GetCurrentLayer() *glow.Layer {
	layer, _ := m.Layer.Get()
	return layer.(*glow.Layer)
}

func (m *Model) SetCurrentLayer(i int) {
	frame := m.GetFrame()
	var layer *glow.Layer
	if i < len(frame.Layers) {
		m.LayerIndex = i
		layer = &frame.Layers[i]
	} else {
		m.LayerIndex = 0
		layer = &glow.Layer{}
	}
	m.Layer.Set(layer)
}
