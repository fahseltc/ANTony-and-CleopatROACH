package scene

import (
	"gamejam/ui"
	"image"
	"math"
)

type CutsceneAction interface {
	Update(s *PlayScene, dt float64) bool // returns true if finished
}
type ZoomCameraAction struct {
	TargetZoom float64
	Speed      float64
	FocusX     float64 // world X coordinate to zoom towards
	FocusY     float64 // world Y coordinate to zoom towards
}

func (a *ZoomCameraAction) Update(s *PlayScene, dt float64) bool {
	currentZoom := s.Ui.Camera.ViewPortZoom
	// Calculate how much to adjust the zoom towards the target
	diff := a.TargetZoom - currentZoom
	step := a.Speed * dt
	if math.Abs(diff) < step {
		// Finalize zoom
		step = diff
	}
	// Convert world focus point to screen coordinates
	screenX, screenY := s.Ui.Camera.MapPosToScreenPos(int(a.FocusX), int(a.FocusY))
	// Perform the zoom
	s.Ui.Camera.Zoom(step, screenX, screenY)
	// Done when zoom is very close to target
	return math.Abs(s.Ui.Camera.ViewPortZoom-a.TargetZoom) < 0.001
}

type PanCameraAction struct {
	TargetX, TargetY float64
	Speed            float64
}

func (a *PanCameraAction) Update(s *PlayScene, dt float64) bool {
	// Get camera and screen details
	cam := s.Ui.Camera
	screenWidth := 800.0
	screenHeight := 600.0
	tileSize := 128.0

	// Target camera position (centered) based on tile coordinates
	targetMapX := a.TargetX * tileSize
	targetMapY := a.TargetY * tileSize
	targetX := -int(targetMapX - screenWidth/2)
	targetY := -int(targetMapY - screenHeight/2)

	// Current camera position
	dx := float64(targetX - cam.ViewPortX)
	dy := float64(targetY - cam.ViewPortY)

	dist := math.Hypot(dx, dy)

	// Arrival threshold
	if dist < 1 {
		cam.ViewPortX = targetX
		cam.ViewPortY = targetY
		// Final clamp
		cam.PanX(0)
		cam.PanY(0)
		return true
	}

	// Move towards target
	angle := math.Atan2(dy, dx)
	prevX, prevY := cam.ViewPortX, cam.ViewPortY
	cam.ViewPortX += int(math.Cos(angle) * a.Speed * dt)
	cam.ViewPortY += int(math.Sin(angle) * a.Speed * dt)

	// After move, clamp the camera
	cam.PanX(0)
	cam.PanY(0)

	// Check if camera is stuck (can't move further due to bounds)
	if cam.ViewPortX == prevX && cam.ViewPortY == prevY {
		return true
	}

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

type ShowPortraitTextAreaAction struct {
	portraitTextArea *ui.PortraitTextArea
}

func (a *ShowPortraitTextAreaAction) Update(s *PlayScene, dt float64) bool {
	s.currentDialog = a.portraitTextArea
	if a.portraitTextArea.Ta.Dismissed {
		s.currentDialog = nil
		return true
	}
	return a.portraitTextArea.Ta.Dismissed
}

type WaitAction struct {
	Duration float64
	Elapsed  float64
}

func (a *WaitAction) Update(s *PlayScene, dt float64) bool {
	a.Elapsed += dt
	return a.Elapsed >= a.Duration
}

type IssueUnitCommandAction struct {
	unitID     string
	targetTile *image.Point
}

func (a *IssueUnitCommandAction) Update(s *PlayScene, dt float64) bool {
	a.targetTile.X = a.targetTile.X * 128
	a.targetTile.Y = a.targetTile.Y * 128
	s.sim.IssueAction(a.unitID, a.targetTile)
	return true
}

type DrawTemporarySpriteAction struct {
	spr             *ui.Sprite
	TargetPosition  *image.Point
	MaxDuration     int
	CurrentDuration int
}

func NewDrawTemporarySpriteAction(sprite *ui.Sprite, targetTile *image.Point, duration int) *DrawTemporarySpriteAction {
	return &DrawTemporarySpriteAction{
		spr:            sprite,
		TargetPosition: targetTile,
		MaxDuration:    duration,
	}
}

func (a *DrawTemporarySpriteAction) Update(s *PlayScene, dt float64) bool {
	if a.CurrentDuration == 0 {
		a.spr.SetPosition(&image.Point{
			X: a.TargetPosition.X,
			Y: a.TargetPosition.Y,
		})
		s.Sprites[a.spr.Id.String()] = a.spr
	}
	a.CurrentDuration++
	if a.CurrentDuration >= a.MaxDuration {
		delete(s.Sprites, a.spr.Id.String())
		a.CurrentDuration = 0
		return true
	}

	return false
}
