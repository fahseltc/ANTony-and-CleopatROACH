package sim

import (
	"image"
	"math"

	"github.com/google/uuid"
)

// import "sync"

type Hive struct {
	ID       uuid.UUID
	Position *image.Point
	Rect     *image.Rectangle
	Faction  uint

	Resource uint
}

var TileDimensions = 128

func NewHive(x, y int) *Hive {
	hive := &Hive{
		ID:       uuid.New(),
		Position: &image.Point{0, 0},
		Rect: &image.Rectangle{
			Min: image.Point{0, 0},
			Max: image.Point{TileDimensions * 2, TileDimensions * 2},
		},
		Faction:  1,
		Resource: 0,
	}
	hive.SetTilePosition(x, y)
	return hive
}

func (h *Hive) SetTilePosition(x, y int) {
	h.Position = &image.Point{X: x * TileDimensions, Y: y * TileDimensions}
	h.Rect.Min = *h.Position
	h.Rect.Max = image.Point{X: h.Position.X + TileDimensions*2, Y: h.Position.Y + TileDimensions*2}
}

func (h *Hive) Update(sim *T) {}

func (h *Hive) DistanceTo(point image.Point) uint {
	xDist := math.Abs(float64(h.Position.X - point.X))
	yDist := math.Abs(float64(h.Position.Y - point.Y))
	return uint(math.Sqrt(math.Pow(xDist, 2) + math.Pow(yDist, 2)))
}
