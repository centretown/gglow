package effectio

import (
	"fmt"
	"gglow/glow"
	"gglow/iohandler"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

const PathSeparator = "/"

type EffectIo struct {
	iohandler.IoHandler

	Accessor *iohandler.Accessor
	frame    *glow.Frame
	layer    *glow.Layer

	selection   string
	folderName  string
	effectName  string
	layerIndex  int
	summaryList []string

	folderWatch binding.Int
	frameWatch  binding.Int
	layerWatch  binding.Int
	hasChanged  binding.Bool
	data        binding.BoolTree

	isActive    bool
	saveActions []func(*glow.Frame)
}

func NewEffect(io iohandler.IoHandler, preferences fyne.Preferences, accessor *iohandler.Accessor) *EffectIo {

	eff := &EffectIo{
		IoHandler:   io,
		folderWatch: binding.NewInt(),
		frameWatch:  binding.NewInt(),
		layerWatch:  binding.NewInt(),

		hasChanged:  binding.NewBool(),
		saveActions: make([]func(*glow.Frame), 0),
		summaryList: make([]string, 0),
		Accessor:    accessor,
	}

	eff.frame = glow.NewFrame()
	eff.layer = eff.frame.Layers[0]

	folder := accessor.Folder
	effect := accessor.Effect

	if folder != "" {
		eff.folderName = folder
	}
	if len(effect) > 0 {
		eff.LoadEffect(effect)
	}
	eff.data = eff.BuildTreeData()
	return eff
}

func (eff *EffectIo) IsFolder(title string) bool {
	return iohandler.IsFolder(title)
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

func (eff *EffectIo) ListEffects(folder string) []string {
	return eff.data.ChildIDs(folder)
}

func (eff *EffectIo) ListFolders() []string {
	return eff.data.ChildIDs(binding.DataTreeRootID)
}

func (eff *EffectIo) Select(selection string) {
	// eff.selection = selection
	// if iohandler.IsFolder(selection) {
	// 	eff.alertFolder()
	// 	return
	// }

	split := strings.Split(selection, PathSeparator)
	if len(split) < 1 {
		return
	}
	eff.folderName = split[0]
	if len(split) > 1 {
		eff.LoadEffect(split[1])
		return
	}
	eff.alertFolder()
}

func (eff *EffectIo) LoadEffect(title string) error {
	frame, err := eff.IoHandler.ReadEffect(eff.folderName, title)
	if err != nil {
		return err
	}

	eff.effectName = title
	eff.setFrame(frame, 0)
	eff.SetUnchanged()
	return nil
}

// func (eff *EffectIo) loadFolder(folder string) []string {
// 	ls, err := eff.IoHandler.SetCurrentFolder(folder)
// 	if err != nil {
// 		fyne.LogError("loadfolder", err)
// 		return ls
// 	}
// 	eff.alertFolder()
// 	return ls
// }

func (eff *EffectIo) GetFrame() *glow.Frame {
	return eff.frame
}

func (eff *EffectIo) setFrame(frame *glow.Frame, layerIndex int) {
	eff.frame = frame

	summaries := make([]string, 0, len(frame.Layers))
	for i, layer := range frame.Layers {
		summaries = append(summaries, SummarizeLayer(layer, i+1))
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
	eff.layer = eff.frame.Layers[index]
}

func (eff *EffectIo) SetCurrentLayer(i int) {
	eff.setLayer(i)
	eff.alertLayer()
}

func (eff *EffectIo) GetCurrentLayer() *glow.Layer {
	return eff.layer
}

func (eff *EffectIo) InsertLayer() {
	eff.insertLayer(eff.layerIndex)
}

func (eff *EffectIo) insertLayer(pos int) {
	count := len(eff.frame.Layers)
	if pos < 0 {
		pos = 0
	} else if pos > count {
		pos = count
	}

	layer := glow.NewLayer()
	layers := make([]*glow.Layer, count+1)

	j := 0
	for i := range layers {
		if i == pos {
			layers[i] = layer
		} else {
			layers[i] = eff.frame.Layers[j]
			j++
		}
	}
	eff.frame.Layers = layers
	eff.setLayer(pos)
	eff.alertLayer()
	eff.SetChanged()
}

func (eff *EffectIo) AddLayer() {
	eff.insertLayer(len(eff.frame.Layers))
}

func (eff *EffectIo) RemoveLayer() {
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

func (eff *EffectIo) HasChanged() bool {
	b, _ := eff.hasChanged.Get()
	return b
}

func (eff *EffectIo) FolderExists(folder string) (exists bool) {
	m, _, _ := eff.data.Get()
	_, exists = m[folder]
	return
}

func (eff *EffectIo) EffectExists(effect string) bool {
	ls := eff.data.ChildIDs(eff.folderName)
	for _, e := range ls {
		if e == effect {
			return true
		}
	}
	return false
}

func (eff *EffectIo) ValidateNewFolderName(title string) error {
	if eff.FolderExists(title) {
		return fmt.Errorf("%s already exists", title)
	}
	return ValidateFolderName(title)
}

func (eff *EffectIo) CreateNewFolder(folder string) error {
	if eff.FolderExists(folder) {
		return fmt.Errorf("%s already exists", folder)
	}
	return eff.CreateFolder(folder)
}

func (eff *EffectIo) AddFolder(title string) (err error) {
	err = eff.CreateNewFolder(title)
	if err != nil {
		return
	}
	eff.data.Append(binding.DataTreeRootID, title, false)
	eff.alertFolder()
	return
}

func (eff *EffectIo) ListCurrent() []string {
	if iohandler.IsFolder(eff.selection) {
		return eff.ListFolders()
	}
	size := len(eff.data.ChildIDs(eff.folderName)) + 1
	s := make([]string, 0, size)
	s = append(s, iohandler.AsFolder())
	s = append(s, eff.data.ChildIDs(eff.folderName)...)
	return s
}

func (eff *EffectIo) ValidateNewEffectName(title string) error {
	if eff.EffectExists(title) {
		return fmt.Errorf("%s already exists", title)
	}
	return ValidateEffectName(title)
}

func (eff *EffectIo) CreateNewEffect(title string, frame *glow.Frame) error {
	if eff.EffectExists(title) {
		return fmt.Errorf("%s already exists", title)
	}
	return eff.CreateEffect(eff.folderName, title, frame)
}

func (eff *EffectIo) AddEffect(title string, frame *glow.Frame) (err error) {
	err = eff.CreateNewEffect(title, frame)
	if err != nil {
		fyne.LogError(title, err)
	}
	eff.effectName = title
	eff.frame = frame
	eff.data.Append(eff.folderName, title, false)
	eff.alertFolder()
	eff.alertFrame()
	return
}

func (eff *EffectIo) SaveEffect() (err error) {
	for _, saveAction := range eff.saveActions {
		saveAction(eff.frame)
	}

	title := eff.EffectName()
	err = eff.IoHandler.UpdateEffect(eff.folderName, title, eff.frame)
	if err != nil {
		fyne.LogError("SaveEffect", err)
		return
	}

	eff.setFrame(eff.frame, eff.layerIndex)
	eff.SetUnchanged()
	return
}

func (eff *EffectIo) OnSave(f func(*glow.Frame)) {
	eff.saveActions = append(eff.saveActions, f)
}

func (eff *EffectIo) EffectName() string {
	return eff.effectName
}

func (eff *EffectIo) FolderName() string {
	return eff.folderName
}

func (eff *EffectIo) TreeData() binding.BoolTree {
	return eff.data
}

func (eff *EffectIo) BuildTreeData() binding.BoolTree {
	var data binding.BoolTree = binding.NewBoolTree()
	folders, _ := eff.ListKeys(iohandler.AsFolder())
	for _, folder := range folders {
		data.Append(binding.DataTreeRootID, folder.Folder, false)
	}

	for _, folder := range folders {
		ls, _ := eff.ListKeys(folder.Folder)
		for _, l := range ls {
			if !iohandler.IsFolder(l.Folder) {
				val := l.Effect + PathSeparator + l.Folder
				data.Append(folder.Folder, val, false)
			}
		}
	}
	return data
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

func (eff *EffectIo) alertFolder() {
	Alert(eff.folderWatch)
}

func (eff *EffectIo) alertFrame() {
	Alert(eff.frameWatch)
}

func (eff *EffectIo) alertLayer() {
	Alert(eff.layerWatch)
}
