package ui

import (
	"gamejam/log"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	MaxZoom       = 1.0
	MinZoom       = 0.3
	ZoomIncrement = 0.05
	CameraSpeed   = 10.0
)

type Camera struct {
	log *slog.Logger

	X, Y          float64
	Width, Height int
	Zoom          float64
	MoveSpeed     float64

	// Effective Width/Height after applying Zoom.
	// Cached unless Zoom changes to save cycles
	effWidth, effHeight float64

	mapWidth, mapHeight int
}

func NewCamera(screenWidth, screenHeight, mapWidth, mapHeight int) *Camera {
	return &Camera{
		log:       log.NewLogger().With("for", "camera"),
		X:         0,
		Y:         0,
		effWidth:  float64(screenWidth) / 1.0,
		effHeight: float64(screenHeight) / 1.0,
		mapWidth:  mapWidth,
		mapHeight: mapHeight,
		Zoom:      1.0,
	}
}

func (c *Camera) Update() {
	var dx, dy float64
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		dy = -c.MoveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		dx = -c.MoveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		dy = c.MoveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		dx = c.MoveSpeed
	}

	c.Move(dx, dy)

	_, mouseWheelY := ebiten.Wheel()
	if mouseWheelY > 0 {
		c.SetZoom(ZoomIncrement)
	}
	if mouseWheelY < 0 {
		c.SetZoom(-ZoomIncrement)
	}
}

func (c *Camera) Move(dx, dy float64) {
	c.X += dx
	c.Y += dy

	// Constrain camera to world boundaries
	if c.X < 0 {
		c.X = 0
	}
	if c.Y < 0 {
		c.Y = 0
	}
	maxWidth := float64(c.mapWidth) - c.effWidth
	maxHeight := float64(c.mapHeight) - c.effHeight
	if c.X > maxWidth {
		c.X = maxWidth
	}
	if c.Y > maxHeight {
		c.Y = maxHeight
	}
}

func (c *Camera) SetZoom(amount float64) {
	c.Zoom += amount
	if c.Zoom >= MaxZoom {
		c.Zoom = MaxZoom
	}
	if c.Zoom <= MinZoom {
		c.Zoom = MinZoom
	}

	// Re-calculate effective viewport size
	c.effWidth = float64(c.Width) / c.Zoom
	c.effHeight = float64(c.Height) / c.Zoom

	// Re-calculate move speed
	c.MoveSpeed = CameraSpeed * c.Zoom
}

// GetViewportBounds returns the world bounds of the current viewport
func (c *Camera) GetViewportBounds() (x, y, width, height float64) {
	return c.X, c.Y, c.effWidth, c.effHeight
}

func (c *Camera) MapPosToScreenPos(mapX, mapY float64) (screenX, screenY float64) {
	return (mapX - c.X) * c.Zoom, (mapY - c.Y) * c.Zoom
}

func (c *Camera) ScreenPosToMapPos(x, y int) (int, int) {
	//c.log.Info("input", "x", x, "y", y)
	mapX := (float64(x) + float64(c.X)) * c.Zoom
	mapY := (float64(y) + float64(c.Y)) * c.Zoom
	//c.log.Info("output", "mapX", mapX, "mapY", mapY)
	return int(mapX), int(mapY)
}

// IsVisible checks if a sprite at given position is visible in camera viewport
func (c *Camera) IsVisible(x, y, width, height float64) bool {
	return x+width > c.X && x < c.X+c.effWidth &&
		y+height > c.Y && y < c.Y+c.effHeight
}
