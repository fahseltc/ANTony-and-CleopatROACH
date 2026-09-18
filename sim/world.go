package sim

import (
	"gamejam/tilemap"
	"gamejam/vec2"
	"math"
)

type World struct {
	TileMap    *tilemap.Tilemap
	TileData   [][]*tilemap.Tile
	MapObjects []*tilemap.MapObject
	FogOfWar   *FogOfWar
}

// RevealFogOfWar clears the fog over the rectangular region bounded (inclusively)
// by topLeft and bottomRight. Coordinates are TILE coordinates (not world pixels),
// matching the fog grid's indexing.
//
// NOTE: this is a partial implementation intended for use from cutscene design.
// Because FogOfWar.Update runs every tick and demotes FogVisible -> FogMemory
// before re-deriving live visibility from units/buildings, we reveal to FogMemory
// so the region stays "explored" rather than being wiped on the next frame.
// TODO(cutscene): decide whether some cutscenes need a truly persistent
// "always visible" reveal, which would require a dedicated flag on FogOfWar.
func (w *World) RevealFogOfWar(topLeft *vec2.T, bottomRight *vec2.T) {
	if w.FogOfWar == nil || topLeft == nil || bottomRight == nil {
		return
	}
	fog := w.FogOfWar

	// Normalize the rectangle so we can iterate regardless of corner order.
	minX := int(math.Min(topLeft.X, bottomRight.X))
	minY := int(math.Min(topLeft.Y, bottomRight.Y))
	maxX := int(math.Max(topLeft.X, bottomRight.X))
	maxY := int(math.Max(topLeft.Y, bottomRight.Y))

	// Clamp to the grid bounds.
	if minX < 0 {
		minX = 0
	}
	if minY < 0 {
		minY = 0
	}
	if maxX > fog.Width-1 {
		maxX = fog.Width - 1
	}
	if maxY > fog.Height-1 {
		maxY = fog.Height - 1
	}

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if fog.Tiles[y][x] == FogUnexplored {
				fog.Tiles[y][x] = FogMemory
			}
		}
	}
}
