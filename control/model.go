package control

import (
	"fmt"
	"glow-gui/data"
	"glow-gui/glow"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

type Model struct {
	store            *data.Store
	Frame            binding.Untyped
	Layer            binding.Untyped
	LayerSummaryList binding.StringList
	isDirty          binding.Bool
	saveActions      []func() `json:"-"`
	WindowHasContent bool
}

func NewController(store *data.Store) *Model {
	m := &Model{
		store:            store,
		Frame:            store.Frame,
		Layer:            store.Layer,
		LayerSummaryList: binding.NewStringList(),
		isDirty:          store.IsDirty,
		saveActions:      make([]func(), 0),
	}

	m.Frame.AddListener(binding.NewDataListener(m.onChangeFrame))
	return m
}

func (m *Model) AddSaveAction(f func()) {
	m.saveActions = append(m.saveActions, f)
}

func (m *Model) EffectName() string {
	return m.store.EffectName
}

func (m *Model) GetFrame() *glow.Frame {
	return m.store.GetFrame()
}

func (m *Model) GetCurrentLayer() *glow.Layer {
	return m.store.GetCurrentLayer()
}

func (m *Model) SetCurrentLayer(i int) {
	m.store.SetCurrentLayer(i)
}

func (m *Model) LayerIndex() int {
	return m.store.LayerIndex
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

func (m *Model) WriteEffect() (err error) {
	m.store.UpdateHistory()

	for _, saveAction := range m.saveActions {
		saveAction()
	}

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
	m.store.HasPrevious.AddListener(listener)
}

func (m *Model) CanUndo() bool {
	return m.store.CanUndo(m.store.EffectName)
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

	fmt.Println("UndoEffect succeeded", frame.Interval, m.EffectName())
}

func (m *Model) onChangeFrame() {
	frame := m.GetFrame()
	summaries := make([]string, 0, len(frame.Layers))
	for i, layer := range frame.Layers {
		summaries = append(summaries, Summarize(&layer, i+1))
	}
	m.LayerSummaryList.Set(summaries)
}
