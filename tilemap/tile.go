package tilemap

import "image"

var (
	TileTypePlain    = "plain"
	TileTypeResource = "resource"
)

type Tile struct {
	Type        string
	Coordinates *image.Point
	Rect        *image.Rectangle
	TileID      int
}
