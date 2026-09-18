package scene

import (
	"gamejam/ui"
	"gamejam/util"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TutorialAlignment controls where on the screen a tutorial modal is drawn.
type TutorialAlignment int

const (
	TutorialCenter TutorialAlignment = iota
	TutorialTopLeft
	TutorialTopRight
	TutorialBottomLeft
	TutorialBottomRight
	TutorialTopCenter    // centered along the top edge
	TutorialBottomCenter // centered along the bottom edge
	TutorialLeftCenter   // centered along the left edge
	TutorialRightCenter  // centered along the right edge
)

// TutorialSize controls how large the tutorial modal is drawn. The image keeps
// its native aspect ratio and is fit inside a box of the size below.
type TutorialSize int

const (
	TutorialSizeTiny TutorialSize = iota
	TutorialSizeSmall
	TutorialSizeMedium
	TutorialSizeFullscreen
)

// TutorialMargin is the gap (px) between a corner-aligned modal and the screen
// edges. Center and fullscreen modals ignore it.
var TutorialMargin = 20

// tutorialSizeBox returns the max width/height (px) a modal of the given size
// may occupy. The image is scaled to fit inside this box preserving aspect.
func tutorialSizeBox(size TutorialSize) (w, h int) {
	switch size {
	case TutorialSizeTiny:
		return ui.GameResolutionW / 3, ui.GameResolutionH / 3
	case TutorialSizeSmall:
		return ui.GameResolutionW / 2, ui.GameResolutionH / 2
	case TutorialSizeMedium:
		return (ui.GameResolutionW * 3) / 4, (ui.GameResolutionH * 3) / 4
	case TutorialSizeFullscreen:
		return ui.GameResolutionW, ui.GameResolutionH
	default:
		return ui.GameResolutionW / 2, ui.GameResolutionH / 2
	}
}

// fitPreservingAspect returns the largest (w, h) that fits inside (boxW, boxH)
// while keeping the srcW:srcH aspect ratio.
func fitPreservingAspect(srcW, srcH, boxW, boxH int) (int, int) {
	if srcW <= 0 || srcH <= 0 {
		return boxW, boxH
	}
	scale := float64(boxW) / float64(srcW)
	if hScale := float64(boxH) / float64(srcH); hScale < scale {
		scale = hScale
	}
	return int(float64(srcW) * scale), int(float64(srcH) * scale)
}

// alignRect positions a w×h rect on screen according to the given alignment.
func alignRect(w, h int, alignment TutorialAlignment) *image.Rectangle {
	m := TutorialMargin
	var minX, minY int
	switch alignment {
	case TutorialCenter:
		minX = (ui.GameResolutionW - w) / 2
		minY = (ui.GameResolutionH - h) / 2
	case TutorialTopLeft:
		minX, minY = m, m
	case TutorialTopRight:
		minX, minY = ui.GameResolutionW-w-m, m
	case TutorialBottomLeft:
		minX, minY = m, ui.GameResolutionH-h-m
	case TutorialBottomRight:
		minX, minY = ui.GameResolutionW-w-m, ui.GameResolutionH-h-m
	case TutorialTopCenter:
		minX, minY = (ui.GameResolutionW-w)/2, m
	case TutorialBottomCenter:
		minX, minY = (ui.GameResolutionW-w)/2, ui.GameResolutionH-h-m
	case TutorialLeftCenter:
		minX, minY = m, (ui.GameResolutionH-h)/2
	case TutorialRightCenter:
		minX, minY = ui.GameResolutionW-w-m, (ui.GameResolutionH-h)/2
	default:
		minX = (ui.GameResolutionW - w) / 2
		minY = (ui.GameResolutionH - h) / 2
	}
	rect := image.Rect(minX, minY, minX+w, minY+h)
	return &rect
}

type Tutorial interface {
	CheckTrigger(s *PlayScene)
	Draw(screen *ebiten.Image)
	IsComplete() bool
	// WantsDismissOnPauseOpen reports whether this step should be dismissed when
	// the pause menu is opened (used only by the pause tutorial, so opening pause
	// with Escape teaches the key and clears the prompt at the same time).
	WantsDismissOnPauseOpen() bool
	// Dismiss marks the step complete (if it has been triggered/enabled).
	Dismiss()
}

type TutorialStep struct {
	Image        *ebiten.Image
	Rect         *image.Rectangle
	TriggerFunc  func(*PlayScene) bool // Function to check if the step should be triggered
	CompleteFunc func(*PlayScene) bool // Function to check if the step is completed
	Enabled      bool
	Completed    bool
	// DismissOnPauseOpen, when true, lets opening the pause menu complete this
	// step. Only the pause tutorial sets this; all other steps are unaffected.
	DismissOnPauseOpen bool
}

// NewTutorialStep builds a tutorial step. Instead of hand-positioning a pixel
// rectangle, callers pick a screen alignment (center or a corner) and a size;
// the image is scaled to fit that size box while preserving its aspect ratio
// and placed at the chosen alignment.
func NewTutorialStep(imagePath string, alignment TutorialAlignment, size TutorialSize, triggerFunc func(*PlayScene) bool, completeFunc func(*PlayScene) bool) *TutorialStep {
	img := util.LoadImage(imagePath)
	boxW, boxH := tutorialSizeBox(size)
	w, h := fitPreservingAspect(img.Bounds().Dx(), img.Bounds().Dy(), boxW, boxH)
	rect := alignRect(w, h, alignment)

	tutorial := &TutorialStep{
		Image: util.ScaleImage(img, float32(w), float32(h)),
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

// NewTutorialStepNoClick builds a tutorial step that never advances on a
// left-click. Use it for steps whose action IS a left-click on the game world
// (e.g. "left-click to start construction"): the placement click would
// otherwise also dismiss the modal, so these steps must be dismissed by a game
// condition instead. A completeFunc is therefore REQUIRED; passing nil would
// leave the step with no way to complete, so we guard against that.
func NewTutorialStepNoClick(imagePath string, alignment TutorialAlignment, size TutorialSize, triggerFunc func(*PlayScene) bool, completeFunc func(*PlayScene) bool) *TutorialStep {
	if completeFunc == nil {
		panic("NewTutorialStepNoClick requires a non-nil completeFunc; a no-click step needs a game condition to complete")
	}
	return NewTutorialStep(imagePath, alignment, size, triggerFunc, completeFunc)
}

// dismissOnAnyKey is a completeFunc that advances when the player presses any
// keyboard key (but NOT a mouse button). Use it for read-and-continue info
// panels shown while a mouse action is pending (e.g. with the build cursor
// active), so the pending left-click can't also skip the panel.
func dismissOnAnyKey(*PlayScene) bool {
	return len(inpututil.AppendJustPressedKeys(nil)) > 0
}

// NewPauseTutorialStep is like NewTutorialStep but additionally opts the step
// into being dismissed when the player opens the pause menu. Use it for the
// tutorial that teaches the pause key so pressing Escape both pauses the game
// and clears the prompt. It still dismisses on left-click via the default
// completeFunc when completeFunc is nil.
func NewPauseTutorialStep(imagePath string, alignment TutorialAlignment, size TutorialSize, triggerFunc func(*PlayScene) bool, completeFunc func(*PlayScene) bool) *TutorialStep {
	ts := NewTutorialStep(imagePath, alignment, size, triggerFunc, completeFunc)
	ts.DismissOnPauseOpen = true
	return ts
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

func (ts *TutorialStep) WantsDismissOnPauseOpen() bool {
	return ts.DismissOnPauseOpen
}

// Dismiss marks the step complete, but only once it has actually been triggered
// (Enabled). This prevents dismissing a step the player hasn't seen yet.
func (ts *TutorialStep) Dismiss() {
	if ts.Enabled {
		ts.Completed = true
	}
}
