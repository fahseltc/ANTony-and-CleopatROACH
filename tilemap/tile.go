package tilemap

import (
	"gamejam/types"
	"image"
)

type Tile struct {
	Type         types.Tile
	Coordinates  *image.Point
	Rect         *image.Rectangle
	TileID       int
	HasCollision bool
}
