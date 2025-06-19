package sim

import (
	"image"

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
}

func NewDefaultUnit() *Unit {
	return &Unit{
		ID: uuid.New(),
		Stats: &UnitStats{
			HPMax:     100,
			HPCur:     100,
			MoveSpeed: 10,
			Damage:    10,
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
		// if unit.NearestEnemy . InRange()
		// attack
		// else
		unit.MoveToDestination(world) // destination might be a unit?
	case AttackingAction:
	case HoldingPositionAction:
		// check range and attack if in range, otherwise stay in place
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
