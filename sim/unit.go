package sim

import (
	"image"
	"math"

	"github.com/google/uuid"
)

type Action int

const (
	IdleAction Action = iota
	MovingAction
	AttackMovingAction
	AttackingAction
	HoldingPositionAction
	CollectingAction
	DeliveringAction
)

type Unit struct {
	ID           uuid.UUID
	Stats        *UnitStats
	Position     *image.Point
	Rect         *image.Rectangle
	Destination  *image.Point
	Action       Action
	NearestEnemy *Unit

	Faction uint
}

type UnitStats struct {
	HPMax     uint
	HPCur     uint
	MoveSpeed uint
	Damage    uint
	Range     uint
}

func NewDefaultUnit() *Unit {
	return &Unit{
		ID: uuid.New(),
		Stats: &UnitStats{
			HPMax:     100,
			HPCur:     100,
			MoveSpeed: 10,
			Damage:    10,
			Range:     15,
			// acceleration / current speed?
		},
		Position: &image.Point{0, 0},
		Rect: &image.Rectangle{
			Min: image.Point{0, 0},
			Max: image.Point{128, 128},
		},
		Destination: &image.Point{0, 0},
		Action:      IdleAction,
		Faction:     0,
	}
}

func (unit *Unit) Update(sim *T) {
	switch unit.Action {
	case IdleAction:
		return
	case MovingAction:
		unit.MoveToDestination(sim)
	case AttackMovingAction:
		if unit.NearestEnemy != nil && unit.TargetInRange(*unit.Position) {
			unit.NearestEnemy.Stats.HPCur -= unit.Stats.Damage
			// pew pew animation
		} else {
			unit.MoveToDestination(sim) // destination might be a unit?
		}
	case HoldingPositionAction:
		if unit.NearestEnemy != nil && unit.TargetInRange(*unit.Position) {
			unit.NearestEnemy.Stats.HPCur -= unit.Stats.Damage
			// pew pew animation
		}
	case CollectingAction:
		// go get resource and pick it up
	case DeliveringAction:
		// return resource to home base
	}
}

func (unit *Unit) MoveToDestination(sim *T) {
	// TODO pathfinding
	//oldX := unit.Position.X
	if unit.Position.X != unit.Destination.X {
		newX := 0
		if unit.Position.X < unit.Destination.X {
			newX = unit.Position.X + int(unit.Stats.MoveSpeed)
		} else if unit.Position.X > unit.Destination.X {
			newX = unit.Position.X - int(unit.Stats.MoveSpeed)
		}
		newRect := &image.Rectangle{
			Min: image.Point{
				X: newX,
				Y: unit.Position.Y,
			},
			Max: image.Point{
				X: newX + unit.Rect.Dx(),
				Y: unit.Position.Y + unit.Rect.Dy(),
			},
		}
		// check X collision
		collidesX := false
		for _, worldUnit := range sim.GetAllNearbyUnits(unit.Position.X, unit.Position.Y) {
			if worldUnit.Rect.Overlaps(*newRect) {
				collidesX = true
				break
			}
		}
		if !collidesX {
			// check X collision with world
			for _, worldCollision := range sim.world.CollisionRects {
				if worldCollision.Overlaps(*newRect) {
					collidesX = true
					break
				}
			}
		}
		if collidesX {
			// dont move in X then
		} else {
			unit.SetPosition(&image.Point{X: newX, Y: unit.Position.Y})
			// if the distance to desitination is smaller than movespeed, set Pos to Dest to prevent flicker/wobbling
			if math.Abs(float64(unit.Position.X-unit.Destination.X)) <= float64(unit.Stats.MoveSpeed) {
				unit.SetPosition(&image.Point{X: unit.Destination.X, Y: unit.Position.Y})
			}

		}
	}

	// move in Y
	if unit.Position.Y != unit.Destination.Y {
		newY := 0
		if unit.Position.Y < unit.Destination.Y {
			newY = unit.Position.Y + int(unit.Stats.MoveSpeed)
		} else if unit.Position.Y > unit.Destination.Y {
			newY = unit.Position.Y - int(unit.Stats.MoveSpeed)
		}
		newRect := &image.Rectangle{
			Min: image.Point{
				X: unit.Position.X,
				Y: newY,
			},
			Max: image.Point{
				X: unit.Position.X + unit.Rect.Dx(),
				Y: newY + unit.Rect.Dy(),
			},
		}
		// check Y collision with units
		collidesY := false
		for _, worldUnit := range sim.GetAllNearbyUnits(unit.Position.X, unit.Position.Y) {
			if worldUnit.Rect.Overlaps(*newRect) {
				collidesY = true
				break
			}
		}
		if !collidesY {
			// check Y collision with world
			for _, worldCollision := range sim.world.CollisionRects {
				if worldCollision.Overlaps(*newRect) {
					collidesY = true
					break
				}
			}
		}
		if collidesY {
			// dont move in Y then
		} else {
			unit.SetPosition(&image.Point{X: unit.Position.X, Y: newY})
			if math.Abs(float64(unit.Position.Y-unit.Destination.Y)) <= float64(unit.Stats.MoveSpeed) {
				unit.SetPosition(&image.Point{X: unit.Position.X, Y: unit.Destination.Y})
			}
		}
	}

	if unit.Position == unit.Destination { // If we've arrived, go Idle
		unit.Action = IdleAction
	}
}

func (unit *Unit) SetNearestEnemy(target *Unit) {
	unit.NearestEnemy = target
}

func (unit *Unit) DistanceTo(point image.Point) uint {
	xDist := math.Abs(float64(unit.Position.X - point.X))
	yDist := math.Abs(float64(unit.Position.Y - point.Y))
	return uint(math.Sqrt(math.Pow(xDist, 2) + math.Pow(yDist, 2)))
}

func (unit *Unit) TargetInRange(point image.Point) bool {
	return unit.DistanceTo(point) <= unit.Stats.Range
}

func (unit *Unit) SetPosition(pos *image.Point) {
	unit.Position = pos
	unit.Rect = &image.Rectangle{
		Min: *pos,
		Max: image.Point{
			X: pos.X + unit.Rect.Dx(),
			Y: pos.Y + unit.Rect.Dy(),
		},
	}
}
