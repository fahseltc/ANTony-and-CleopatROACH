package ui

import (
	"gamejam/util"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type ConstructionMouse struct {
	Enabled            bool
	constructingSprite *ebiten.Image
	placementRect      *image.Rectangle
}

func (cm *ConstructionMouse) Update() {
	if !cm.Enabled || cm.constructingSprite != nil {
		return
	}
}
func (cm *ConstructionMouse) Draw(screen *ebiten.Image, camera *Camera) {
	if !cm.Enabled || cm.constructingSprite != nil {
		return
	}

	mx, my := ebiten.CursorPosition()

	// Convert screen coordinates to map (world) coordinates
	mapX, mapY := camera.ScreenPosToMapPos(mx, my)

	// Snap to grid (assuming 128x128 grid size)
	gridSize := 128.0
	snappedX := float64(mapX/int(gridSize)) * gridSize
	snappedY := float64(mapY/int(gridSize)) * gridSize

	// Set placementRect in map (world) coordinates
	rect := image.Rect(
		int(snappedX),
		int(snappedY),
		int(snappedX)+int(gridSize),
		int(snappedY)+int(gridSize),
	)
	cm.placementRect = &rect

	// Convert snapped map position back to screen coordinates
	screenX, screenY := camera.MapPosToScreenPos(int(snappedX), int(snappedY))

	opts := &ebiten.DrawImageOptions{}
	// Set opacity to 50%
	opts.ColorM.Scale(1, 1, 1, 0.5)
	// Scale according to camera zoom
	opts.GeoM.Scale(camera.ViewPortZoom, camera.ViewPortZoom)
	// Draw at snapped position
	opts.GeoM.Translate(float64(screenX), float64(screenY))
	screen.DrawImage(cm.constructingSprite, opts)
}

func (cm *ConstructionMouse) SetSprite(sprite string) {
	cm.constructingSprite = util.LoadImage(sprite)
}
