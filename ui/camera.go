package ui

import (
	"gamejam/log"
	"log/slog"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

var MaxZoom = 1.0
var MinZoom = 0.3

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
		c.Zoom(.05)
	}
	if mouseWheelY < 0 {
		c.Zoom(-.05)
	}
}

func (c *Camera) Zoom(amount float64) {
	// c.ViewPortZoom += amount
	// if c.ViewPortZoom >= MaxZoom {
	// 	c.ViewPortZoom = MaxZoom
	// }
	// if c.ViewPortZoom <= MinZoom {
	// 	c.ViewPortZoom = MinZoom
	// }
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
	//c.log.Info("input", "x", x, "y", y)
	mapX := math.Abs(float64(c.ViewPortX)) + float64(x)
	mapY := math.Abs(float64(c.ViewPortY)) + float64(y)
	//c.log.Info("output", "mapX", mapX, "mapY", mapY)
	return int(mapX), int(mapY)
}

func (c *Camera) MapPosToScreenPos(x, y int) (int, int) {
	//c.log.Info("input", "mapX", x, "mapY", y)
	mapX := math.Abs(float64(c.ViewPortX)) - float64(x)
	mapY := math.Abs(float64(c.ViewPortY)) - float64(y)
	//c.log.Info("output", "mapX", mapX, "mapY", mapY)
	return int(mapX), int(mapY)
}
