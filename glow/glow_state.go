package glow

type GlowState struct {
	Frame      *Frame
	LayerIndex int
}

func NewGlowState(frame *Frame, layerIndex int) *GlowState {
	return &GlowState{Frame: frame, LayerIndex: layerIndex}
}

func (gs *GlowState) Copy(source *GlowState) (err error) {
	gs.LayerIndex = source.LayerIndex
	gs.Frame, err = FrameDeepCopy(source.Frame)
	return
}

func (gs *GlowState) Layer() *Layer {
	return &gs.Frame.Layers[gs.LayerIndex]
}

func (gs *GlowState) ValidIndex() bool {
	return gs.LayerIndex < len(gs.Frame.Layers)
}
