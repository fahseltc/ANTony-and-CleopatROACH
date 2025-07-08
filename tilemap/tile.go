package tilemap

import "image"

var (
	TileTypePlain   = "plain"
	TileTypeSucrose = "sucrose"
	TileTypeWood    = "wood"
)

type Tile struct {
	Type         string
	Coordinates  *image.Point
	Rect         *image.Rectangle
	TileID       int
	HasCollision bool
}
