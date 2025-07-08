package sim

import (
	"gamejam/vec2"
	"image"
	"math"

	"github.com/google/uuid"
)

var TileDimensions = 128

type Building struct {
	ID       uuid.UUID
	Type     BuildingType
	Position *vec2.T
	Rect     *image.Rectangle
	Faction  uint

	ProgressMax     uint
	ProgressCurrent uint
}

type BuildingType int

const (
	BuildingTypeInConstruction BuildingType = iota
	BuildingTypeHive
	BuildingTypeRoachHive
	BuildingTypeBridge
)

type BuildingInterface interface {
	//SetPosition(x, y, width, height int)
	SetTilePosition(x, y int)
	GetTilePosition() *vec2.T
	GetAdjacentCoordinates() []*vec2.T
	GetID() uuid.UUID
	GetType() BuildingType
	GetPosition() *vec2.T
	GetCenteredPosition() *vec2.T
	GetClosestPosition(*vec2.T) *vec2.T
	GetRect() *image.Rectangle
	GetFaction() uint

	Update(sim *T) // if buildings have an Update behavior
	DistanceTo(point vec2.T) uint
	GetProgress() float64

	// needs to be implemented by inheriting building
	AddUnitToBuildQueue()
}

func NewBuilding(x, y, width, height int, faction uint, bt BuildingType, progressMax uint) *Building {
	rect := &image.Rectangle{
		Min: image.Point{X: x, Y: y},
		Max: image.Point{X: x + width, Y: y + height},
	}
	return &Building{
		ID:          uuid.New(),
		Type:        bt,
		Position:    &vec2.T{X: float64(x), Y: float64(y)},
		Rect:        rect,
		Faction:     faction,
		ProgressMax: progressMax,
	}
}

func (b *Building) SetTilePosition(x, y int) {
	b.Position = &vec2.T{X: float64(x * TileDimensions), Y: float64(y * TileDimensions)}
	b.Rect.Min = image.Point{X: x * TileDimensions, Y: y * TileDimensions}
	b.Rect.Max = image.Point{X: int(b.Position.X) + TileDimensions*2, Y: int(b.Position.Y) + TileDimensions*2}
}

func (b *Building) GetTilePosition() *vec2.T {
	tileX := float64(b.Position.X) / float64(TileDimensions)
	tileY := float64(b.Position.Y) / float64(TileDimensions)
	return &vec2.T{
		X: tileX,
		Y: tileY,
	}
}
func (b *Building) GetAdjacentCoordinates() []*vec2.T {
	if b.Rect == nil {
		return nil
	}
	// Get tile bounds
	minTileX := b.Rect.Min.X / TileDimensions
	minTileY := b.Rect.Min.Y / TileDimensions
	maxTileX := (b.Rect.Max.X - 1) / TileDimensions
	maxTileY := (b.Rect.Max.Y - 1) / TileDimensions

	adjacent := []*vec2.T{}

	// Top edge
	for x := minTileX; x <= maxTileX; x++ {
		adjacent = append(adjacent, &vec2.T{X: float64(x), Y: float64(minTileY - 1)})
	}
	// Bottom edge
	for x := minTileX; x <= maxTileX; x++ {
		adjacent = append(adjacent, &vec2.T{X: float64(x), Y: float64(maxTileY + 1)})
	}
	// Left edge
	for y := minTileY; y <= maxTileY; y++ {
		adjacent = append(adjacent, &vec2.T{X: float64(minTileX - 1), Y: float64(y)})
	}
	// Right edge
	for y := minTileY; y <= maxTileY; y++ {
		adjacent = append(adjacent, &vec2.T{X: float64(maxTileX + 1), Y: float64(y)})
	}

	return adjacent
}

func (b *Building) GetID() uuid.UUID      { return b.ID }
func (b *Building) GetType() BuildingType { return b.Type }
func (b *Building) GetPosition() *vec2.T  { return b.Position }
func (b *Building) GetCenteredPosition() *vec2.T {
	if b.Rect == nil {
		return b.Position
	}
	centerX := b.Rect.Min.X + (b.Rect.Dx() / 2)
	centerY := b.Rect.Min.Y + (b.Rect.Dy() / 2)
	return &vec2.T{X: float64(centerX), Y: float64(centerY)}
}

func (b *Building) GetClosestPosition(pos *vec2.T) *vec2.T {
	clampedX := pos.X
	clampedY := pos.X

	if clampedX < float64(b.Rect.Min.X) {
		clampedX = float64(b.Rect.Min.X)
	} else if clampedX > float64(b.Rect.Max.X) {
		clampedX = float64(b.Rect.Max.X)
	}

	if clampedY < float64(b.Rect.Min.Y) {
		clampedY = float64(b.Rect.Min.Y)
	} else if clampedY > float64(b.Rect.Max.Y) {
		clampedY = float64(b.Rect.Max.Y)
	}

	return &vec2.T{X: float64(clampedX), Y: float64(clampedY)}
}

func (b *Building) GetRect() *image.Rectangle { return b.Rect }
func (b *Building) GetFaction() uint          { return b.Faction }
func (b *Building) DistanceTo(point vec2.T) uint {
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
