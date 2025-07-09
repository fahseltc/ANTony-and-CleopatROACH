package ui

import (
	"gamejam/sim"
	"gamejam/tilemap"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	MiniMapWidth     = 150
	MiniMapHeight    = 150
	MiniMapBottomPad = 10
	MiniMapLeftPad   = 10
)

type MiniMap struct {
	position image.Point
	rect     image.Rectangle
	image    *ebiten.Image
}

func NewMiniMap(tileMap *tilemap.Tilemap) *MiniMap {
	pos := image.Point{
		X: MiniMapLeftPad,
		Y: 600 - MiniMapHeight - MiniMapBottomPad,
	}

	return &MiniMap{
		position: pos,
		rect:     image.Rect(pos.X, pos.Y, MiniMapWidth, MiniMapHeight),
		image:    ebiten.NewImage(MiniMapWidth, MiniMapHeight),
	}
}
func (m *MiniMap) RenderFromTilemap(tileMap *tilemap.Tilemap, fog *sim.FogOfWar) {
	scaleX := float64(MiniMapWidth) / float64(tileMap.Width)
	scaleY := float64(MiniMapHeight) / float64(tileMap.Height)

	m.image.Clear()

	for y := 0; y < tileMap.Height; y++ {
		for x := 0; x < tileMap.Width; x++ {
			tile := tileMap.Tiles[x][y]

			var c color.Color

			// Determine base tile color
			if tile.HasCollision {
				c = color.RGBA{80, 80, 80, 255}
			} else {
				c = color.RGBA{120, 200, 120, 255}
			}

			// Apply fog effect if enabled
			if fog != nil && fog.Enabled {
				switch fog.Tiles[y][x] {
				case sim.FogUnexplored:
					c = color.RGBA{10, 10, 10, 255} // completely hidden
				case sim.FogMemory:
					// dim the color
					r, g, b, a := c.RGBA()
					c = color.RGBA{
						uint8(r >> 9), // ~25% brightness
						uint8(g >> 9),
						uint8(b >> 9),
						uint8(a >> 8),
					}
				case sim.FogVisible:
					// use normal color
				}
			}

			// Map tile to minimap pixels
			minX := int(float64(x) * scaleX)
			minY := int(float64(y) * scaleY)
			maxX := int(float64(x+1) * scaleX)
			maxY := int(float64(y+1) * scaleY)

			for py := minY; py < maxY; py++ {
				for px := minX; px < maxX; px++ {
					m.image.Set(px, py, c)
				}
			}
		}
	}
}

func (m *MiniMap) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(m.position.X), float64(m.position.Y))
	screen.DrawImage(m.image, op)
}
func (m *MiniMap) DrawUnits(units []*sim.Unit, tileMap *tilemap.Tilemap, fog *sim.FogOfWar) {
	scaleX := float64(MiniMapWidth) / float64(tileMap.Width*tileMap.TileSize)
	scaleY := float64(MiniMapHeight) / float64(tileMap.Height*tileMap.TileSize)

	unitRectSize := 4 // Size of the rectangle for each unit on the minimap

	for _, unit := range units {
		// Only show enemy units if they are visible in the fog of war
		if unit.Faction != uint(sim.PlayerFaction) && fog != nil && fog.Enabled {
			tileX := int(unit.Position.X) / tileMap.TileSize
			tileY := int(unit.Position.Y) / tileMap.TileSize
			if tileY < 0 || tileY >= len(fog.Tiles) || tileX < 0 || tileX >= len(fog.Tiles[0]) {
				continue
			}
			if fog.Tiles[tileY][tileX] != sim.FogVisible {
				continue
			}
		}

		x := int(unit.Position.X * scaleX)
		y := int(unit.Position.Y * scaleY)

		c := color.RGBA{255, 0, 0, 255} // enemy red
		if unit.Faction == uint(sim.PlayerFaction) {
			c = color.RGBA{0, 255, 0, 255} // friendly green
		}

		minX := x - unitRectSize/2
		minY := y - unitRectSize/2
		maxX := x + unitRectSize/2
		maxY := y + unitRectSize/2

		// Clamp to minimap bounds
		if minX < 0 {
			minX = 0
		}
		if minY < 0 {
			minY = 0
		}
		if maxX > MiniMapWidth {
			maxX = MiniMapWidth
		}
		if maxY > MiniMapHeight {
			maxY = MiniMapHeight
		}

		for py := minY; py < maxY; py++ {
			for px := minX; px < maxX; px++ {
				m.image.Set(px, py, c)
			}
		}
	}
}
func (m *MiniMap) DrawViewport(camera *Camera, tileMap *tilemap.Tilemap) {
	// Map size in pixels * zoom
	renderedMapWidth := float64(tileMap.Width*tileMap.TileSize) * camera.ViewPortZoom
	renderedMapHeight := float64(tileMap.Height*tileMap.TileSize) * camera.ViewPortZoom

	// Screen size in pixels
	screenW := 800.0
	screenH := 600.0

	// Camera world position (positive coords)
	cameraX := -float64(camera.ViewPortX)
	cameraY := -float64(camera.ViewPortY)

	// Scale factors from world pixels to minimap pixels
	scaleX := float64(MiniMapWidth) / renderedMapWidth
	scaleY := float64(MiniMapHeight) / renderedMapHeight

	// Map camera position to minimap position
	minimapX := int(cameraX * scaleX)
	minimapY := int(cameraY * scaleY)

	// Calculate size of viewport rectangle on minimap
	minimapW := int(screenW * scaleX)
	minimapH := int(screenH * scaleY)

	// Clamp minimapX/Y so viewport rectangle never leaves minimap bounds
	if minimapX < 0 {
		minimapX = 0
	}
	if minimapY < 0 {
		minimapY = 0
	}
	if minimapX+minimapW > MiniMapWidth {
		minimapX = MiniMapWidth - minimapW
	}
	if minimapY+minimapH > MiniMapHeight {
		minimapY = MiniMapHeight - minimapH
	}

	outlineColor := color.RGBA{255, 255, 255, 255}

	// Draw top and bottom lines
	for px := minimapX; px < minimapX+minimapW; px++ {
		if px >= 0 && px < MiniMapWidth {
			if minimapY >= 0 && minimapY < MiniMapHeight {
				m.image.Set(px, minimapY, outlineColor)
			}
			if minimapY+minimapH-1 >= 0 && minimapY+minimapH-1 < MiniMapHeight {
				m.image.Set(px, minimapY+minimapH-1, outlineColor)
			}
		}
	}

	// Draw left and right lines
	for py := minimapY; py < minimapY+minimapH; py++ {
		if py >= 0 && py < MiniMapHeight {
			if minimapX >= 0 && minimapX < MiniMapWidth {
				m.image.Set(minimapX, py, outlineColor)
			}
			if minimapX+minimapW-1 >= 0 && minimapX+minimapW-1 < MiniMapWidth {
				m.image.Set(minimapX+minimapW-1, py, outlineColor)
			}
		}
	}
}
