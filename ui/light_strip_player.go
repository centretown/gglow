package ui

import (
	"fmt"
	"glow-gui/glow"
	"time"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type LightStripPlayer struct {
	*widget.Toolbar

	strip  *LightStrip
	source *glow.Frame

	sep         *widget.ToolbarSeparator
	playButton  *widget.ToolbarAction
	pauseButton *widget.ToolbarAction
	stepButton  *widget.ToolbarAction
	resetButton *widget.ToolbarAction
	stopButton  *widget.ToolbarAction

	stopChan  chan int
	stepChan  chan int
	pauseChan chan int
	startChan chan int

	isSpinning bool
	isActive   bool
}

func NewLightStripPlayer(strip *LightStrip) *LightStripPlayer {
	sb := &LightStripPlayer{
		source:    &glow.Frame{},
		strip:     strip,
		stopChan:  make(chan int),
		stepChan:  make(chan int),
		pauseChan: make(chan int),
		startChan: make(chan int),
	}

	var items []widget.ToolbarItem

	sb.sep = widget.NewToolbarSeparator()
	sb.playButton = widget.NewToolbarAction(theme.MediaPlayIcon(), sb.Play)
	sb.pauseButton = widget.NewToolbarAction(theme.MediaPauseIcon(), sb.Pause)
	sb.stepButton = widget.NewToolbarAction(theme.MediaSkipNextIcon(), sb.Step)
	sb.resetButton = widget.NewToolbarAction(theme.MediaReplayIcon(), sb.Reset)
	sb.stopButton = widget.NewToolbarAction(theme.MediaStopIcon(), sb.Stop)

	items = append(items,
		sb.playButton,
		sb.stepButton,
		sb.pauseButton,
		sb.resetButton,
		sb.stopButton,
	)

	sb.Toolbar = widget.NewToolbar(items...)
	return sb
}

func (sb *LightStripPlayer) Pause() {
	if sb.isActive {
		sb.pauseChan <- 0
	}
}

func (sb *LightStripPlayer) Step() {
	sb.isSpinning = false
	if !sb.isActive {
		sb.startSpinner()
	}
	sb.stepChan <- 0
}

func (sb *LightStripPlayer) Play() {
	sb.isSpinning = true
	if !sb.isActive {
		sb.startSpinner()
	}
	sb.startChan <- 0
}

func (sb *LightStripPlayer) Stop() {
	sb.stopSpinner()
	sb.strip.TurnOff()
}

func (sb *LightStripPlayer) Reset() {
	isSpinning := sb.isSpinning
	isActive := sb.isActive

	sb.stopSpinner()

	if isActive {
		if isSpinning {
			sb.Play()
		} else {
			sb.Step()
		}
	}
}

func (sb *LightStripPlayer) stopSpinner() {
	if sb.isActive {
		sb.stopChan <- 0
	}
}

func (sb *LightStripPlayer) startSpinner() {
	sb.stopSpinner()

	var (
		isSpinning = sb.isSpinning
		frame      *glow.Frame
		err        error
	)
	frame, err = glow.FrameCopy(sb.source)
	if err != nil {
		fmt.Println("startSpinner FrameCopy", err)
		panic("startSpinner")
	}

	frame.Setup(sb.strip.Length(), sb.strip.Rows(), sb.strip.Interval())

	spin := func() {
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

func (sb *LightStripPlayer) OnExit() {
	sb.stopSpinner()
}

func (sb *LightStripPlayer) SetFrame(frame *glow.Frame) {
	sb.source = frame
	sb.source.Setup(sb.strip.Length(), sb.source.Rows, sb.strip.Interval())
	sb.Reset()
}
