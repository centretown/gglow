package control

import (
	"fmt"
	"glow-gui/glow"
	"glow-gui/store"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

type Manager struct {
	store            *store.Store
	Frame            binding.Untyped
	Layer            binding.Untyped
	LayerSummaryList binding.StringList
	isDirty          binding.Bool
	canUndo          binding.Bool
	WindowHasContent bool
}

func NewController(store *store.Store) *Manager {
	m := &Manager{
		store:            store,
		Frame:            store.Frame,
		Layer:            store.Layer,
		LayerSummaryList: binding.NewStringList(),
		isDirty:          store.IsDirty,
		canUndo:          store.HasPrevious,
	}

	m.Frame.AddListener(binding.NewDataListener(m.onChangeFrame))
	return m
}

func (m *Manager) EffectName() string {
	return m.store.EffectName
}

func (m *Manager) IsFolder(title string) bool {
	return m.store.IsFolder(title)
}

func (m *Manager) CreateNewEffect(title string, frame *glow.Frame) (err error) {
	return m.store.CreateNewEffect(title, frame)
}

func (m *Manager) CreateNewFolder(title string) (err error) {
	return m.store.CreateNewFolder(title)
}

func (m *Manager) RefreshKeys(title string) []string {
	m.store.RefreshKeys(title)
	return m.KeyList()
}

func (m *Manager) ValidateNewFolderName(title string) (err error) {
	return m.store.ValidateNewFolderName(title)
}

func (m *Manager) ValidateNewEffectName(title string) (err error) {
	return m.store.ValidateNewEffectName(title)
}

func (m *Manager) KeyList() []string {
	return m.store.KeyList
}

func (m *Manager) UpdateHistory() error {
	return m.store.UpdateHistory()
}

func (m *Manager) WriteEffect() (err error) {
	frame := m.GetFrame()
	err = m.store.WriteEffect(m.EffectName(), frame)
	current := *frame
	m.Frame.Set(&current)
	return
}

func (m *Manager) ReadEffect(title string) (err error) {
	err = m.store.ReadEffect(title)
	if err != nil {
		fyne.LogError("ReadEffect", err)
		return
	}

	// fmt.Println("Model ReadEffect", title)
	return
}

func (m *Manager) AddFrameListener(listener binding.DataListener) {
	m.Frame.AddListener(listener)
}

func (m *Manager) AddDirtyListener(listener binding.DataListener) {
	m.isDirty.AddListener(listener)
}

func (m *Manager) SetDirty() {
	if m.WindowHasContent {
		m.store.SetDirty(true)
	}
}

func (m *Manager) IsDirty() bool {
	return m.store.GetDirty()
}

func (m *Manager) AddUndoListener(listener binding.DataListener) {
	m.canUndo.AddListener(listener)
}

func (m *Manager) CanUndo() bool {
	b, _ := m.canUndo.Get()
	return b
}

func (m *Manager) UndoEffect() {
	if !m.store.CanUndo(m.EffectName()) {
		fyne.LogError("UndoEffect",
			fmt.Errorf("%s has no history", m.EffectName()))
		return
	}

	state, err := m.store.Undo(m.EffectName())
	if err != nil {
		fyne.LogError("UndoEffect", err)
		return
	}

	fmt.Println("UndoEffect succeeded", state.Frame.Interval, m.EffectName())
}

func (m *Manager) onChangeFrame() {
	frame := m.GetFrame()
	summaries := make([]string, 0, len(frame.Layers))
	for i, layer := range frame.Layers {
		summaries = append(summaries, Summarize(&layer, i+1))
	}
	m.LayerSummaryList.Set(summaries)
}

func (m *Manager) GetFrame() *glow.Frame {
	return m.store.GetFrame()
}

func (m *Manager) GetCurrentLayer() *glow.Layer {
	return m.store.GetCurrentLayer()
}

func (m *Manager) SetCurrentLayer(i int) {
	m.store.SetCurrentLayer(i)
}

func (m *Manager) LayerIndex() int {
	return m.store.LayerIndex
}
