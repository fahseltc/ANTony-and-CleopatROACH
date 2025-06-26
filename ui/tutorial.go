package ui

import (
	"gamejam/sim"
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
	Image       *ebiten.Image
	TriggerFunc func(*sim.T) bool // Function to check if the step should be triggered
	Enabled     bool
	Completed   bool
}

func NewTutorialStep(image *ebiten.Image, rect *image.Rectangle, triggerFunc func(*sim.T) bool) *TutorialStep {
	return &TutorialStep{
		Image:       image,
		TriggerFunc: triggerFunc,
	}
}

func (ts *TutorialStep) CheckTrigger(sim *sim.T) {
	if ts.TriggerFunc != nil && ts.TriggerFunc(sim) {
		ts.Enabled = true
	}
	if ts.Enabled && !ts.Completed {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			ts.Completed = true
		}
	}
}

func (ts *TutorialStep) Draw(screen *ebiten.Image) {
	if ts.Enabled && !ts.Completed {
		screen.DrawImage(ts.Image, nil)
	}

}
