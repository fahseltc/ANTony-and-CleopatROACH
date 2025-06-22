package sim

import (
	"image"

	"github.com/google/uuid"
)

// import "sync"

type Hive struct {
	ID             uuid.UUID
	Position       *image.Point
	Rect           *image.Rectangle
	ResourceAmount uint
}

var TileDimensions = 128

func NewHive() *Hive {
	return &Hive{
		ID:       uuid.New(),
		Position: &image.Point{0, 0},
		Rect: &image.Rectangle{
			Min: image.Point{0, 0},
			Max: image.Point{TileDimensions, TileDimensions},
		},
		ResourceAmount: 0,
	}
}

func (h *Hive) SetTilePosition(x, y int) {
	h.Position = &image.Point{X: x * TileDimensions, Y: y * TileDimensions}
	h.Rect.Min = *h.Position
	h.Rect.Max = image.Point{X: h.Position.X + TileDimensions, Y: h.Position.Y + TileDimensions}
}

// func (h *Hive) AddResource(amount uint) {
// 	h.ResourceAmount += amount
// }

// func (h *Hive) SpendResource(amount uint) bool {
// 	if h.ResourceAmount-amount < 0 {
// 		return false
// 	}
// 	h.ResourceAmount -= amount
// 	return true
// }
