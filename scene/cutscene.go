package scene

import (
	"gamejam/sim"
	"gamejam/ui"
	"gamejam/vec2"
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

// NudgeCameraByTilesAction scrolls the camera by a RELATIVE offset measured in
// tiles from its current position (unlike PanCameraAction, which pans to an
// absolute tile position). Positive OffsetTilesX/Y nudge the view right/down.
//
// Unlike the old tile-nudge, this is zoom-aware: an offset of N tiles moves the
// view by N tiles as they appear on screen at the current zoom. It also
// accumulates sub-pixel movement so small speeds still progress instead of
// truncating to zero each frame.
//
// Speed is in on-screen pixels per second.
type NudgeCameraByTilesAction struct {
	OffsetTilesX, OffsetTilesY float64
	Speed                      float64

	resolved         bool
	targetX, targetY float64 // absolute viewport pixel target
	curX, curY       float64 // float-accumulated viewport position
}

func (a *NudgeCameraByTilesAction) Update(s *PlayScene, dt float64) bool {
	const tileSize = 128.0
	cam := s.Ui.Camera

	if !a.resolved {
		// A tile spans tileSize*zoom pixels on screen, so a relative tile offset
		// translates into that many viewport pixels. ViewPortX/Y are the draw
		// offset (negative of map position), and a positive tile offset should
		// move the view toward higher map coordinates, i.e. decrease ViewPort.
		a.curX = float64(cam.ViewPortX)
		a.curY = float64(cam.ViewPortY)
		a.targetX = a.curX - a.OffsetTilesX*tileSize*cam.ViewPortZoom
		a.targetY = a.curY - a.OffsetTilesY*tileSize*cam.ViewPortZoom
		a.resolved = true
	}

	dx := a.targetX - a.curX
	dy := a.targetY - a.curY
	dist := math.Hypot(dx, dy)

	if dist < 1 {
		cam.ViewPortX = int(a.targetX)
		cam.ViewPortY = int(a.targetY)
		cam.PanX(0) // clamp to map bounds
		cam.PanY(0)
		return true
	}

	step := a.Speed * dt
	if step >= dist {
		step = dist
	}
	a.curX += dx / dist * step
	a.curY += dy / dist * step

	prevX, prevY := cam.ViewPortX, cam.ViewPortY
	cam.ViewPortX = int(a.curX)
	cam.ViewPortY = int(a.curY)
	cam.PanX(0) // clamp to map bounds
	cam.PanY(0)

	// If the map-bounds clamp overrode our intended position, the camera is
	// against an edge. Resync the float accumulators to the clamped values (so
	// we don't keep pushing into the wall) and finish. Otherwise KEEP the float
	// accumulators as-is so sub-pixel movement builds up across frames and small
	// speeds still make progress instead of truncating to zero.
	if cam.ViewPortX != int(a.curX) || cam.ViewPortY != int(a.curY) {
		a.curX = float64(cam.ViewPortX)
		a.curY = float64(cam.ViewPortY)
		if cam.ViewPortX == prevX && cam.ViewPortY == prevY {
			return true // wedged against a bound, nothing more to do
		}
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

// RevealFogOfWarAction clears the fog of war over a rectangular region during a
// cutscene. TopLeft and BottomRight are TILE coordinates (inclusive), matching
// how other level-authored actions specify tiles.
//
// NOTE: this is a partial implementation for use in cutscene design. It reveals
// the region once and completes immediately. TODO(cutscene): consider an optional
// animated/progressive reveal and a "keep permanently visible" flag.
type RevealFogOfWarAction struct {
	TopLeft     *image.Point
	BottomRight *image.Point
}

func (a *RevealFogOfWarAction) Update(s *PlayScene, dt float64) bool {
	if a.TopLeft == nil || a.BottomRight == nil {
		return true
	}
	// sim.RevealFogOfWar works in tile coordinates via vec2.T; the fog grid is
	// tile-indexed, so pass tile coords straight through (no *128 conversion).
	topLeft := &vec2.T{X: float64(a.TopLeft.X), Y: float64(a.TopLeft.Y)}
	bottomRight := &vec2.T{X: float64(a.BottomRight.X), Y: float64(a.BottomRight.Y)}
	s.sim.RevealFogOfWar(topLeft, bottomRight)
	return true
}

type WaitAction struct {
	Duration float64
	Elapsed  float64
}

func (a *WaitAction) Update(s *PlayScene, dt float64) bool {
	a.Elapsed += dt
	return a.Elapsed >= a.Duration
}

// DisableInputAction locks out keyboard input (camera movement, unit hotkeys,
// button key activation) and drag selecting. Mouse clicks still work. It
// completes immediately; the lock stays in effect until an EnableInputAction
// runs. Note the cutscene loop re-enables drag when the whole cutscene ends.
type DisableInputAction struct{}

func (a *DisableInputAction) Update(s *PlayScene, dt float64) bool {
	s.inputDisabled = true
	return true
}

// EnableInputAction re-enables keyboard input and drag selecting.
type EnableInputAction struct{}

func (a *EnableInputAction) Update(s *PlayScene, dt float64) bool {
	s.inputDisabled = false
	return true
}

type IssueUnitCommandAction struct {
	unitID     string
	targetTile *image.Point
}

func (a *IssueUnitCommandAction) Update(s *PlayScene, dt float64) bool {
	a.targetTile.X = a.targetTile.X * 128
	a.targetTile.Y = a.targetTile.Y * 128
	s.sim.IssueAction([]string{a.unitID}, a.targetTile)
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
		a.spr.SetPosition(&vec2.T{
			X: float64(a.TargetPosition.X),
			Y: float64(a.TargetPosition.Y),
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

type DrawTemporarySpriteBetweenUnitsAction struct {
	spr             *ui.Sprite
	FirstUnit       *sim.Unit
	SecondUnit      *sim.Unit
	MaxDuration     int
	CurrentDuration int
}

func NewDrawTemporarySpriteBetweenUnitsAction(sprite *ui.Sprite, unit1, unit2 *sim.Unit, duration int) *DrawTemporarySpriteBetweenUnitsAction {
	return &DrawTemporarySpriteBetweenUnitsAction{
		spr:         sprite,
		FirstUnit:   unit1,
		SecondUnit:  unit2,
		MaxDuration: duration,
	}
}

func (a *DrawTemporarySpriteBetweenUnitsAction) Update(s *PlayScene, dt float64) bool {
	if a.CurrentDuration == 0 {
		// Calculate the midpoint between the two units
		pos1 := a.FirstUnit.GetCenteredPosition()
		pos2 := a.SecondUnit.GetCenteredPosition()
		midX := (pos1.X + pos2.X) / 2
		midY := (pos1.Y + pos2.Y) / 2

		a.spr.SetPosition(&vec2.T{
			X: midX + float64(a.spr.Rect.Dx())/2,
			Y: midY + float64(a.spr.Rect.Dy())/2,
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
