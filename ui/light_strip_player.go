package ui

import (
	"fmt"
	"glow-gui/glow"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ButtonItem struct {
	*widget.Button
}

func NewButtonItem(btn *widget.Button) (bi *ButtonItem) {
	btn.Importance = widget.LowImportance
	bi = &ButtonItem{
		Button: btn,
	}
	return
}

func (bi *ButtonItem) ToolbarObject() fyne.CanvasObject {
	return bi.Button
}

type LightStripPlayer struct {
	*widget.Toolbar

	strip       *LightStrip
	sourceFrame binding.Untyped

	playPauseButton *ButtonItem
	stepButton      *widget.ToolbarAction
	resetButton     *widget.ToolbarAction
	stopButton      *widget.ToolbarAction

	stopChan  chan int
	stepChan  chan int
	pauseChan chan int
	startChan chan int
	resetChan chan int
	frameChan chan *glow.Frame

	isPlaying bool
	isActive  bool
}

func NewLightStripPlayer(strip *LightStrip, sourceFrame binding.Untyped) *LightStripPlayer {
	sb := &LightStripPlayer{
		sourceFrame: sourceFrame,
		strip:       strip,
		stopChan:    make(chan int),
		stepChan:    make(chan int),
		pauseChan:   make(chan int),
		startChan:   make(chan int),
		resetChan:   make(chan int),
		frameChan:   make(chan *glow.Frame),
	}

	sb.playPauseButton = NewButtonItem(
		widget.NewButtonWithIcon("", theme.MediaPlayIcon(), sb.PlayPause))

	sb.stepButton = widget.NewToolbarAction(theme.MediaSkipNextIcon(), sb.Step)
	sb.resetButton = widget.NewToolbarAction(theme.MediaReplayIcon(), sb.Reset)
	sb.stopButton = widget.NewToolbarAction(theme.MediaStopIcon(), sb.Stop)

	sb.Toolbar = widget.NewToolbar(
		sb.playPauseButton,
		sb.stepButton,
		sb.resetButton,
		sb.stopButton,
	)

	sb.sourceFrame.AddListener(binding.NewDataListener(func() {
		sb.run()
		sb.frameChan <- sb.getFrame()
	}))
	return sb
}

func (sb *LightStripPlayer) getFrame() *glow.Frame {
	frame, _ := sb.sourceFrame.Get()
	return frame.(*glow.Frame)
}

func (sb *LightStripPlayer) OnExit() {
	sb.stopSpinner()
}

func (sb *LightStripPlayer) Stop() {
	sb.stopSpinner()
	sb.strip.TurnOff()
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

	spin := func() {
		var (
			isSpinning bool
			frame      *glow.Frame
			err        error
		)

		copyFrame := func(source *glow.Frame) {
			frame, err = glow.FrameDeepCopy(source)
			if err != nil {
				fmt.Println("startSpinner FrameCopy", err)
				panic("startSpinner")
			}
			frame.Setup(sb.strip.Length(), sb.strip.Rows(), sb.strip.Interval())
		}

		copyFrame(sb.getFrame())

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

			case <-sb.resetChan:
				copyFrame(sb.getFrame())
				frame.Spin(sb.strip)

			default:
				if isSpinning {
					frame.Spin(sb.strip)
				}
				time.Sleep(time.Duration(frame.Interval) * time.Millisecond)
			}
		}
	}
	go spin()
	sb.isActive = true
}
