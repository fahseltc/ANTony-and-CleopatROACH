package ui

import (
	"gamejam/log"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
)

var MapScrollSpeed = 10

var MaxZoom = 1.0
var MinZoom = 0.3
var ZoomIncrement = 0.05

type Camera struct {
	log          *slog.Logger
	ViewPortX    int
	ViewPortY    int
	ViewPortZoom float64

	mapWidth  int
	mapHeight int
}

func NewCamera(TileWidthCount, TileHeightCount int) *Camera {
	return &Camera{
		log:          log.NewLogger().With("for", "camera"),
		ViewPortX:    0,
		ViewPortY:    0,
		ViewPortZoom: 1,

		mapWidth:  TileWidthCount * TileDimensions,
		mapHeight: TileHeightCount * TileDimensions,
	}
}

func (c *Camera) Update() {
	mx, my := ebiten.CursorPosition()
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		c.PanY(MapScrollSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		c.PanX(MapScrollSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		c.PanY(-MapScrollSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		c.PanX(-MapScrollSpeed)
	}
	_, mouseWheelY := ebiten.Wheel()
	if mouseWheelY > 0 {
		c.Zoom(ZoomIncrement, mx, my)
	}
	if mouseWheelY < 0 {
		c.Zoom(-ZoomIncrement, mx, my)
	}
}
func (c *Camera) Zoom(amount float64, mouseX, mouseY int) {
	oldZoom := c.ViewPortZoom
	newZoom := c.ViewPortZoom + amount
	if newZoom > MaxZoom {
		newZoom = MaxZoom
	}
	if newZoom < MinZoom {
		newZoom = MinZoom
	}
	// Mouse position in world coordinates before zoom
	mapX := (float64(mouseX) - float64(c.ViewPortX)) / oldZoom
	mapY := (float64(mouseY) - float64(c.ViewPortY)) / oldZoom

	// Apply new zoom
	c.ViewPortZoom = newZoom

	// Adjust camera so that the mouse stays over the same map point
	c.ViewPortX = mouseX - int(mapX*c.ViewPortZoom)
	c.ViewPortY = mouseY - int(mapY*c.ViewPortZoom)

	// Now clamp
	c.PanX(0) // Will enforce constraints
	c.PanY(0) // Will enforce constraints
}
func (c *Camera) PanX(amount int) {
	c.ViewPortX += amount

	// Final rendered map width
	renderedMapWidth := float64(c.mapWidth) * c.ViewPortZoom

	// Leftmost position
	if c.ViewPortX > 0 {
		c.ViewPortX = 0
	}
	// Rightmost position
	minX := 800 - int(renderedMapWidth) // e.g., if map is smaller than screen, this can be > 0
	if c.ViewPortX < minX {
		c.ViewPortX = minX
	}
}
func (c *Camera) PanY(amount int) {
	c.ViewPortY += amount

	// Final rendered map height
	renderedMapHeight := float64(c.mapHeight) * c.ViewPortZoom

	// Topmost position
	if c.ViewPortY > 0 {
		c.ViewPortY = 0
	}
	// Bottommost position
	minY := 600 - int(renderedMapHeight)
	if c.ViewPortY < minY {
		c.ViewPortY = minY
	}
}

// VisibleMapPixels returns the width and height in map pixels currently visible in the viewport.
func (c *Camera) VisibleMapPixels() (int, int) {
	width := int(float64(800) / c.ViewPortZoom)
	height := int(float64(600) / c.ViewPortZoom)
	return width, height
}

func (c *Camera) ScreenPosToMapPos(x, y int) (int, int) {
	mapX := (float64(x) - float64(c.ViewPortX)) / c.ViewPortZoom
	mapY := (float64(y) - float64(c.ViewPortY)) / c.ViewPortZoom
	return int(mapX), int(mapY)
}

func (c *Camera) MapPosToScreenPos(x, y int) (int, int) {
	screenX := float64(x)*c.ViewPortZoom + float64(c.ViewPortX)
	screenY := float64(y)*c.ViewPortZoom + float64(c.ViewPortY)
	return int(screenX), int(screenY)
}

func (c *Camera) SetPosition(x, y int) {
	c.ViewPortX = -x
	c.ViewPortY = -y
}
func (c *Camera) SetZoom(amount float64) {
	c.ViewPortZoom = amount
}
