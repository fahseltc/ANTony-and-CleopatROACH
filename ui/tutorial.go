package ui

import (
	"gamejam/sim"
	"gamejam/util"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Tutorial interface {
	CheckTrigger(sim *sim.T)
	Draw(screen *ebiten.Image)
	IsComplete() bool
}

type TutorialStep struct {
	Image        *ebiten.Image
	Rect         *image.Rectangle
	TriggerFunc  func(*sim.T) bool // Function to check if the step should be triggered
	CompleteFunc func(*sim.T) bool // Function to check if the step is completed
	Enabled      bool
	Completed    bool
}

func NewTutorialStep(imagePath string, rect *image.Rectangle, triggerFunc func(*sim.T) bool, completeFunc func(*sim.T) bool) *TutorialStep {
	tutorial := &TutorialStep{
		Image: util.ScaleImage(util.LoadImage(imagePath), float32(rect.Bounds().Dx()), float32(rect.Bounds().Dy())),
		Rect:  rect,
	}

	if triggerFunc == nil {
		tutorial.TriggerFunc = func(*sim.T) bool { return true }
	} else {
		tutorial.TriggerFunc = triggerFunc
	}
	if completeFunc == nil {
		tutorial.CompleteFunc = func(*sim.T) bool { return inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) }
	} else {
		tutorial.CompleteFunc = completeFunc
	}

	return tutorial
}

func (ts *TutorialStep) CheckTrigger(sim *sim.T) {
	if ts.TriggerFunc != nil && ts.TriggerFunc(sim) {
		ts.Enabled = true
	}
	if ts.Enabled && !ts.Completed {
		ts.Completed = ts.CompleteFunc(sim)
	}
}

func (ts *TutorialStep) Draw(screen *ebiten.Image) {
	if ts.Enabled && !ts.Completed {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(ts.Rect.Min.X), float64(ts.Rect.Min.Y))
		screen.DrawImage(ts.Image, opts)
	}
}
func (ts *TutorialStep) IsComplete() bool {
	return ts.Completed
}
