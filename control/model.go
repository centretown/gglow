package control

import (
	"fmt"
	"glow-gui/fileio"
	"glow-gui/glow"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

type Model struct {
	store            *fileio.Store
	LayerSummaryList binding.StringList
	saveActions      []func(*glow.Frame) `json:"-"`
	WindowHasContent bool
}

func NewModel(store *fileio.Store) *Model {
	m := &Model{
		store:            store,
		LayerSummaryList: binding.NewStringList(),
		saveActions:      make([]func(*glow.Frame), 0),
	}

	m.AddFrameListener(binding.NewDataListener(m.onChangeFrame))
	return m
}

func (m *Model) AddSaveAction(f func(*glow.Frame)) {
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
	frame := m.GetFrame()

	for _, saveAction := range m.saveActions {
		saveAction(frame)
	}

	err = m.store.WriteEffect(m.EffectName(), frame)
	current := *frame
	m.store.Frame.Set(&current)
	return
}

func (m *Model) ReadEffect(title string) (err error) {
	err = m.store.ReadEffect(title)
	if err != nil {
		fyne.LogError("ReadEffect", err)
		return
	}

	frame := m.GetFrame()
	fmt.Println("model frame", frame.Interval)

	return
}

func (m *Model) AddFrameListener(listener binding.DataListener) {
	m.store.Frame.AddListener(listener)
}

func (m *Model) AddLayerListener(listener binding.DataListener) {
	m.store.Layer.AddListener(listener)
}

func (m *Model) AddDirtyListener(listener binding.DataListener) {
	m.store.IsDirty.AddListener(listener)
}

func (m *Model) SetDirty() {
	if m.WindowHasContent {
		m.store.SetDirty(true)
	}
}

func (m *Model) IsDirty() bool {
	return m.store.GetDirty()
}

func (m *Model) UndoEffect() {
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
