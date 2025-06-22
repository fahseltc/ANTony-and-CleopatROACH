package ui

import (
	"gamejam/log"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
)

var MaxZoom = 1.0
var MinZoom = 0.3
var ZoomIncrement = 0.05

type Camera struct {
	log          *slog.Logger
	ViewPortX    int
	ViewPortY    int
	ViewPortZoom float64

	// mapWidth  int
	// mapHeight int
}

func NewCamera() *Camera {
	return &Camera{
		log:          log.NewLogger().With("for", "camera"),
		ViewPortX:    0,
		ViewPortY:    0,
		ViewPortZoom: 1,
	}
}

func (c *Camera) Update() {
	scrollSpeed := 10
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		c.PanY(scrollSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		c.PanX(scrollSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		c.PanY(-scrollSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		c.PanX(-scrollSpeed)
	}
	_, mouseWheelY := ebiten.Wheel()
	if mouseWheelY > 0 {
		c.Zoom(ZoomIncrement)
	}
	if mouseWheelY < 0 {
		c.Zoom(-ZoomIncrement)
	}
}

func (c *Camera) Zoom(amount float64) {
	c.ViewPortZoom += amount
	if c.ViewPortZoom >= MaxZoom {
		c.ViewPortZoom = MaxZoom
	}
	if c.ViewPortZoom <= MinZoom {
		c.ViewPortZoom = MinZoom
	}
}
func (c *Camera) PanX(amount int) {
	c.ViewPortX += amount
	if c.ViewPortX >= 0 {
		c.ViewPortX = 0
	}
}
func (c *Camera) PanY(amount int) {
	c.ViewPortY += amount
	if c.ViewPortY >= 0 {
		c.ViewPortY = 0
	}
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
