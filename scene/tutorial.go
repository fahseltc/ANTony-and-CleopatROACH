package scene

import (
	"gamejam/util"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Tutorial interface {
	CheckTrigger(s *PlayScene)
	Draw(screen *ebiten.Image)
	IsComplete() bool
}

type TutorialStep struct {
	Image        *ebiten.Image
	Rect         *image.Rectangle
	TriggerFunc  func(*PlayScene) bool // Function to check if the step should be triggered
	CompleteFunc func(*PlayScene) bool // Function to check if the step is completed
	Enabled      bool
	Completed    bool
}

func NewTutorialStep(imagePath string, rect *image.Rectangle, triggerFunc func(*PlayScene) bool, completeFunc func(*PlayScene) bool) *TutorialStep {
	tutorial := &TutorialStep{
		Image: util.ScaleImage(util.LoadImage(imagePath), float32(rect.Bounds().Dx()), float32(rect.Bounds().Dy())),
		Rect:  rect,
	}

	if triggerFunc == nil {
		tutorial.TriggerFunc = func(*PlayScene) bool { return true }
	} else {
		tutorial.TriggerFunc = triggerFunc
	}
	if completeFunc == nil {
		tutorial.CompleteFunc = func(*PlayScene) bool { return inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) }
	} else {
		tutorial.CompleteFunc = completeFunc
	}

	return tutorial
}

func (ts *TutorialStep) CheckTrigger(s *PlayScene) {
	if ts.TriggerFunc != nil && ts.TriggerFunc(s) {
		ts.Enabled = true
	}
	if ts.Enabled && !ts.Completed {
		ts.Completed = ts.CompleteFunc(s)
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
