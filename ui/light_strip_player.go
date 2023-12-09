package ui

import (
	"glow-gui/effects"
	"glow-gui/glow"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type LightStripPlayer struct {
	*widget.Toolbar

	effect effects.Effect
	// sourceFrame binding.Untyped
	sourceStrip binding.Untyped
	strip       *LightStrip

	playPauseButton *ButtonItem
	stepButton      *widget.ToolbarAction
	resetButton     *widget.ToolbarAction
	stopButton      *widget.ToolbarAction
	layoutButton    *widget.ToolbarAction

	stopChan     chan int
	stepChan     chan int
	pauseChan    chan int
	startChan    chan int
	resetChan    chan int
	stripChan    chan int
	intervalChan chan int
	frameChan    chan *glow.Frame

	isPlaying bool
	isActive  bool
}

func NewLightStripPlayer(sourceStrip binding.Untyped, effect effects.Effect,
	lightStripLayout *LightStripLayout) *LightStripPlayer {

	sb := &LightStripPlayer{
		// sourceFrame:  sourceFrame,
		effect:       effect,
		sourceStrip:  sourceStrip,
		stopChan:     make(chan int),
		stepChan:     make(chan int),
		pauseChan:    make(chan int),
		startChan:    make(chan int),
		resetChan:    make(chan int),
		stripChan:    make(chan int),
		intervalChan: make(chan int),
		frameChan:    make(chan *glow.Frame),
	}

	sb.playPauseButton = NewButtonItem(
		widget.NewButtonWithIcon("", theme.MediaPlayIcon(), sb.PlayPause))

	sb.stepButton = widget.NewToolbarAction(theme.MediaSkipNextIcon(), sb.Step)
	sb.resetButton = widget.NewToolbarAction(theme.MediaReplayIcon(), sb.Reset)
	sb.stopButton = widget.NewToolbarAction(theme.MediaStopIcon(), sb.Stop)

	sb.layoutButton = widget.NewToolbarAction(theme.SettingsIcon(), func() {
		lightStripLayout.CustomDialog.Show()
	})

	sb.Toolbar = widget.NewToolbar(
		sb.playPauseButton,
		sb.stepButton,
		sb.resetButton,
		sb.stopButton,
		sb.layoutButton,
	)

	sb.strip = sb.getStrip()
	sb.effect.AddFrameListener(binding.NewDataListener(sb.frameListener))
	return sb
}

// func (sb *LightStripPlayer) getFrame() *glow.Frame {
// 	frame, _ := sb.sourceFrame.Get()
// 	return frame.(*glow.Frame)
// }

func (sb *LightStripPlayer) getStrip() *LightStrip {
	strip, _ := sb.sourceStrip.Get()
	return strip.(*LightStrip)
}

func (sb *LightStripPlayer) ResetStrip() {
	sb.run()
	sb.stripChan <- 0
}

func (sb *LightStripPlayer) frameListener() {
	sb.run()
	sb.frameChan <- sb.effect.GetFrame()
}

func (sb *LightStripPlayer) OnExit() {
	sb.stopSpinner()
}

func (sb *LightStripPlayer) Stop() {
	sb.pause()
	sb.stopSpinner()
	strip := sb.getStrip()
	strip.TurnOff()
}

func (sb *LightStripPlayer) Step() {
	if sb.isPlaying {
		sb.pause()
	} else {
		sb.run()
	}
	sb.stepChan <- 0
}

func (sb *LightStripPlayer) PlayPause() {
	if sb.isPlaying {
		sb.pause()
	} else {
		sb.play()
	}
}

func (sb *LightStripPlayer) Reset() {
	sb.run()
	sb.resetChan <- 0
}

func (sb *LightStripPlayer) pause() {
	sb.playPauseButton.SetIcon(theme.MediaPlayIcon())
	sb.isPlaying = false
	sb.run()
	sb.pauseChan <- 0
}

func (sb *LightStripPlayer) play() {
	sb.playPauseButton.SetIcon(theme.MediaPauseIcon())
	sb.isPlaying = true
	sb.run()
	sb.startChan <- 0
}

func (sb *LightStripPlayer) run() {
	if !sb.isActive {
		sb.startSpinner()
	}
}

func (sb *LightStripPlayer) stopSpinner() {
	if sb.isActive {
		sb.stopChan <- 0
	}
}

func (sb *LightStripPlayer) startSpinner() {
	sb.stopSpinner()
	go sb.spin()
	sb.isActive = true
}

func (sb *LightStripPlayer) spin() {
	var (
		isSpinning bool
		frame      *glow.Frame
		err        error
	)

	copyFrame := func(source *glow.Frame) {
		frame, err = glow.FrameDeepCopy(source)
		if err != nil {
			reason := "startSpinner FrameDeepCopy"
			fyne.LogError(reason, err)
			reason += " " + err.Error()
			panic(reason)
		}
		frame.Setup(sb.strip.Length(), sb.strip.Rows())
		if frame.Interval == 0 {
			frame.Interval = glow.DefaultInterval
		}
	}

	copyFrame(sb.effect.GetFrame())

	for {
		select {
		case <-sb.stopChan:
			sb.isActive = false
			return

		case <-sb.startChan:
			isSpinning = true

		case <-sb.pauseChan:
			isSpinning = false

		case <-sb.stepChan:
			isSpinning = false
			frame.Spin(sb.strip)

		case f := <-sb.frameChan:
			copyFrame(f)
			frame.Spin(sb.strip)

		case <-sb.stripChan:
			sb.strip = sb.getStrip()
			copyFrame(sb.effect.GetFrame())
			frame.Spin(sb.strip)

		case <-sb.resetChan:
			copyFrame(sb.effect.GetFrame())
			frame.Spin(sb.strip)

		default:
			if isSpinning {
				frame.Spin(sb.strip)
			}
			time.Sleep(time.Duration(frame.Interval) * time.Millisecond)
		}
	}

}
