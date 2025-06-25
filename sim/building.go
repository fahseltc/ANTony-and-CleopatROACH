package sim

import (
	"image"
	"math"

	"github.com/google/uuid"
)

var TileDimensions = 128

type Building struct {
	ID       uuid.UUID
	Type     BuildingType
	Position *image.Point
	Rect     *image.Rectangle
	Faction  uint

	ProgressMax     uint
	ProgressCurrent uint
}

type BuildingType int

const (
	BuildingTypeInConstruction BuildingType = iota
	BuildingTypeHive
	BuildingTypeBridge
)

type BuildingInterface interface {
	SetPosition(x, y, width, height int)
	SetTilePosition(x, y int)
	GetID() uuid.UUID
	GetType() BuildingType
	GetPosition() *image.Point
	GetCenteredPosition() *image.Point
	GetClosestPosition(x, y int) *image.Point
	GetRect() *image.Rectangle
	GetFaction() uint

	Update(sim *T) // if buildings have an Update behavior
	DistanceTo(point image.Point) uint
	GetProgress() float64

	// needs to be implemented by inheriting building
	AddUnitToBuildQueue()
}

func NewBuilding(x, y, width, height int, faction uint, bt BuildingType, progressMax uint) *Building {
	pos := &image.Point{X: x, Y: y}
	rect := &image.Rectangle{
		Min: *pos,
		Max: image.Point{X: x + width, Y: y + height},
	}
	return &Building{
		ID:          uuid.New(),
		Type:        bt,
		Position:    pos,
		Rect:        rect,
		Faction:     faction,
		ProgressMax: progressMax,
	}
}

func (b *Building) SetPosition(x, y, width, height int) {
	b.Position = &image.Point{X: x, Y: y}
	b.Rect.Min = *b.Position
	b.Rect.Max = image.Point{X: x + width, Y: y + height}
}
func (b *Building) SetTilePosition(x, y int) {
	b.Position = &image.Point{X: x * TileDimensions, Y: y * TileDimensions}
	b.Rect.Min = *b.Position
	b.Rect.Max = image.Point{X: b.Position.X + TileDimensions*2, Y: b.Position.Y + TileDimensions*2}
}

func (b *Building) GetID() uuid.UUID          { return b.ID }
func (b *Building) GetType() BuildingType     { return b.Type }
func (b *Building) GetPosition() *image.Point { return b.Position }
func (b *Building) GetCenteredPosition() *image.Point {
	if b.Rect == nil {
		return b.Position
	}
	centerX := b.Rect.Min.X + (b.Rect.Dx() / 2)
	centerY := b.Rect.Min.Y + (b.Rect.Dy() / 2)
	return &image.Point{X: centerX, Y: centerY}
}

func (b *Building) GetClosestPosition(x, y int) *image.Point {
	clampedX := x
	clampedY := y

	if clampedX < b.Rect.Min.X {
		clampedX = b.Rect.Min.X
	} else if clampedX > b.Rect.Max.X {
		clampedX = b.Rect.Max.X
	}

	if clampedY < b.Rect.Min.Y {
		clampedY = b.Rect.Min.Y
	} else if clampedY > b.Rect.Max.Y {
		clampedY = b.Rect.Max.Y
	}

	return &image.Point{X: clampedX, Y: clampedY}
}

func (b *Building) GetRect() *image.Rectangle { return b.Rect }
func (b *Building) GetFaction() uint          { return b.Faction }
func (b *Building) DistanceTo(point image.Point) uint {
	xDist := math.Abs(float64(b.Position.X - point.X))
	yDist := math.Abs(float64(b.Position.Y - point.Y))
	return uint(math.Sqrt(xDist*xDist + yDist*yDist))
}
func (b *Building) GetProgress() float64 {
	if b.ProgressMax == 0 {
		return 0
	}
	return float64(b.ProgressCurrent) / float64(b.ProgressMax)
}
func (b *Building) Update(_ *T)          {} // Default no-op
func (b *Building) AddUnitToBuildQueue() {}
