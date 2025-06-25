package scene

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type CutsceneAction interface {
	Update(s *PlayScene, dt float64) bool // returns true if finished
}

type PanCameraAction struct {
	TargetX, TargetY float64
	Speed            float64
}

func (a *PanCameraAction) Update(s *PlayScene, dt float64) bool {
	// Get the camera and screen details
	cam := s.Ui.Camera
	screenWidth := 800.0
	screenHeight := 600.0

	// Target camera position so that the target point appears centered
	targetX := -a.TargetX + screenWidth/2
	targetY := -a.TargetY + screenHeight/2

	dx := float64(targetX - float64(cam.ViewPortX))
	dy := float64(targetY - float64(cam.ViewPortY))

	dist := math.Hypot(dx, dy)
	if dist < 1 {
		cam.ViewPortX = int(targetX)
		cam.ViewPortY = int(targetY)
		return true
	}

	angle := math.Atan2(dy, dx)
	cam.ViewPortX += int(math.Cos(angle) * a.Speed * dt)
	cam.ViewPortY += int(math.Sin(angle) * a.Speed * dt)

	return false
}

type FadeCameraAction struct {
	Mode    string // "in" or "out"
	Speed   uint8
	Done    bool
	started bool
}

func (a *FadeCameraAction) Update(s *PlayScene, dt float64) bool {
	if a.Done {
		return true
	}
	switch a.Mode {
	case "in":
		if !a.started {
			s.Ui.Camera.FadeIn(a.Speed)
			a.started = true
		}
		a.Done = s.Ui.Camera.FadeAlpha == 0
	case "out":
		if !a.started {
			s.Ui.Camera.FadeOut(a.Speed)
			a.started = true
		}
		a.Done = s.Ui.Camera.FadeAlpha == 255
	}
	s.Ui.Camera.Update()
	return a.Done
}

type ShowTextAreaAction struct {
	Text string
	//TextArea ui.TextArea
	//Duration float64
	//Elapsed  float64
	Skipped bool
}

func (a *ShowTextAreaAction) Update(s *PlayScene, dt float64) bool {
	s.currentDialogText = a.Text // to be rendered later
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		a.Skipped = true
	}
	return a.Skipped
}

type WaitAction struct {
	Duration float64
	Elapsed  float64
}

func (a *WaitAction) Update(s *PlayScene, dt float64) bool {
	a.Elapsed += dt
	return a.Elapsed >= a.Duration
}
