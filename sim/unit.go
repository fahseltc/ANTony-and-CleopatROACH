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
	Destination  *image.Point
	Action       Action
	NearestEnemy *Unit
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
		Position:    &image.Point{0, 0},
		Destination: &image.Point{0, 0},
		Action:      IdleAction,
	}
}

func (unit *Unit) Update(world *World) {
	switch unit.Action {
	case IdleAction:
		return
	case MovingAction:
		unit.MoveToDestination(world)
	case AttackMovingAction:
		if unit.NearestEnemy != nil && unit.TargetInRange(*unit.Position) {
			unit.NearestEnemy.Stats.HPCur -= unit.Stats.Damage
			// pew pew animation
		} else {
			unit.MoveToDestination(world) // destination might be a unit?
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

func (unit *Unit) MoveToDestination(world *World) {
	// TODO pathfinding
	if unit.Position.X < unit.Destination.X {
		unit.Position.X += int(unit.Stats.MoveSpeed)
	} else if unit.Position.X > unit.Destination.X {
		unit.Position.X -= int(unit.Stats.MoveSpeed)
	}
	if unit.Position.Y < unit.Destination.Y {
		unit.Position.Y += int(unit.Stats.MoveSpeed)
	} else if unit.Position.Y > unit.Destination.Y {
		unit.Position.Y -= int(unit.Stats.MoveSpeed)
	}
}

func (unit *Unit) SetNearestEnemy(target *Unit) {
	unit.NearestEnemy = target
}

func (unit *Unit) DistanceTo(point image.Point) uint {
	xDist := math.Abs(float64(unit.Position.X - point.X))
	yDist := math.Abs(float64(unit.Position.Y - point.Y))
	return uint(math.Pow(xDist, 2) + math.Pow(yDist, 2))
}

func (unit *Unit) TargetInRange(point image.Point) bool {
	return unit.DistanceTo(point) <= unit.Stats.Range
}
