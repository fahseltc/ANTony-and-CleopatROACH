package sim

import (
	"gamejam/tilemap"
)

type World struct {
	TileMap    *tilemap.Tilemap
	TileData   [][]*tilemap.Tile
	MapObjects []*tilemap.MapObject
}

// func (w *World) ScreenToTile(vec *vec2.T) {

// }
