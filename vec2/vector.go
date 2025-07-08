package vec2

import (
	"image"
	"math"
)

var (
	TileSize     = 128.0
	HalfTileSize = 64.0
)

type T struct {
	X, Y float64
}

func (a T) Add(b T) T            { return T{a.X + b.X, a.Y + b.Y} }
func (a T) Sub(b T) T            { return T{a.X - b.X, a.Y - b.Y} }
func (a T) Scale(s float64) T    { return T{a.X * s, a.Y * s} }
func (a T) Distance(b T) float64 { return math.Hypot(a.X-b.X, a.Y-b.Y) }
func (a T) Length() float64      { return math.Hypot(a.X, a.Y) }
func (a T) Normalize() T {
	len := math.Hypot(a.X, a.Y)
	if len == 0 {
		return T{0, 0}
	}
	return T{a.X / len, a.Y / len}
}
func (a T) ToPoint() image.Point {
	return image.Point{X: int(a.X), Y: int(a.Y)}
}
func (a T) RoundToGrid() *T {
	return &T{
		X: math.Floor(a.X/128) * 128,
		Y: math.Floor(a.Y/128) * 128,
	}
}

func (a T) ToCenteredPixelCoordinates() *T {
	return &T{
		X: a.X*TileSize + HalfTileSize,
		Y: a.Y*TileSize + HalfTileSize,
	}
}
