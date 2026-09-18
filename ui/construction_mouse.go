package ui

import (
	"gamejam/eventing"
	"gamejam/sim"
	"gamejam/tilemap"
	"gamejam/types"
	"gamejam/util"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var GridSize = 128.0

type ConstructionMouse struct {
	Enabled bool

	currentSprite  *ebiten.Image
	bridgeSprite   *ebiten.Image // limited to 1x1 buildings
	barracksSprite *ebiten.Image

	placementRect *image.Rectangle

	placingBuildingType types.Building

	placementValid bool
}

func NewConstructionMouse() *ConstructionMouse {
	cm := &ConstructionMouse{
		Enabled:        false,
		bridgeSprite:   util.LoadImage("buildings/bridge.png"),
		barracksSprite: util.LoadImage("buildings/barracks.png"),
	}
	return cm
}

func (cm *ConstructionMouse) Update(tm *tilemap.Tilemap, sim *sim.T, camera *Camera) {
	if !cm.Enabled || cm.currentSprite == nil {
		return
	}

	// for _, mo := range tm.MapObjects {
	// 	if mo.IsBuildable && cm.placementRect != nil {
	// 		if cm.placementRect.Overlaps(*mo.Rect) {
	// 			cm.placementValid = true
	// 			break
	// 		} else {
	// 			cm.placementValid = false // GROSS - needs fixing later
	// 		}
	// 	}
	// }

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) && cm.placingBuildingType != types.BuildingTypeNone {
		mx, my := ebiten.CursorPosition()
		mapX, mapY := camera.ScreenPosToMapPos(mx, my)
		tileX := mapX / 128
		tileY := mapY / 128
		sim.EventBus.Publish(eventing.Event{
			Type: "BuildClickedEvent",
			Data: eventing.BuildClickedEvent{
				TargetCoordinates: image.Point{X: tileX, Y: tileY},
				BuildingType:      cm.placingBuildingType,
			},
		})
		// for _, mo := range tm.MapObjects {
		// 	if !mo.IsBuildable {
		// 		continue
		// 	}
		// 	rect := mo.Rect
		// 	if rect.Min == cm.placementRect.Min && rect.Max == cm.placementRect.Max {
		// 		matchingRect := rect
		// 		fmt.Printf("Found matching buildable rect: %v\n", matchingRect)
		// 		sim.EventBus.Publish(eventing.Event{
		// 			Type: "BuildClickedEvent",
		// 			Data: eventing.BuildClickedEvent{
		// 				TargetRect:   matchingRect,
		// 				BuildingType: cm.placingBuildingType,
		// 			},
		// 		})
		// 	}
		// }
		cm.Enabled = false
	}
}

func (cm *ConstructionMouse) Draw(screen *ebiten.Image, camera *Camera) {
	if !cm.Enabled || cm.currentSprite == nil {
		return
	}

	mx, my := ebiten.CursorPosition()

	// Convert screen coordinates to map (world) coordinates
	mapX, mapY := camera.ScreenPosToMapPos(mx, my)

	// Snap to grid (assuming 128x128 grid size)
	snappedX := float64(mapX/int(GridSize)) * GridSize
	snappedY := float64(mapY/int(GridSize)) * GridSize

	// Set placementRect in map (world) coordinates
	rect := image.Rect(
		int(snappedX),
		int(snappedY),
		int(snappedX)+int(GridSize),
		int(snappedY)+int(GridSize),
	)
	cm.placementRect = &rect

	// Convert snapped map position back to screen coordinates
	screenX, screenY := camera.MapPosToScreenPos(int(snappedX), int(snappedY))

	opts := &ebiten.DrawImageOptions{}
	// Set opacity to 50%
	// if cm.placementValid {
	// 	opts.ColorM.Scale(0, 1, 0, 0.5) // green, 50% opacity
	// } else {
	// 	opts.ColorM.Scale(1, 0, 0, 0.5) // red, 50% opacity
	// }
	// Scale according to camera zoom
	opts.GeoM.Scale(camera.ViewPortZoom, camera.ViewPortZoom)
	// Draw at snapped position
	opts.GeoM.Translate(float64(screenX), float64(screenY))
	screen.DrawImage(cm.currentSprite, opts)
}

func (cm *ConstructionMouse) SetSprite(bt types.Building) {
	switch bt {
	case types.BuildingTypeBarracks:
		cm.currentSprite = cm.barracksSprite
		cm.placingBuildingType = types.BuildingTypeBarracks
	case types.BuildingTypeBridge:
		cm.currentSprite = cm.bridgeSprite
		cm.placingBuildingType = types.BuildingTypeBridge
	default:
		cm.currentSprite = nil
		cm.placingBuildingType = types.BuildingTypeNone
	}
}
