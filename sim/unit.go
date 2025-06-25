package sim

import (
	"image"
	"math"

	"github.com/google/uuid"
)

var ArrivalThreshold = 15
var MaxResourceCollectFrames = 30
var PlayerFaction = 0

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

type DestinationType int

const (
	LocationDestination DestinationType = iota
	ResourceDestination
	EnemyDestination
)

type UnitType int

const (
	UnitTypeDefaultAnt UnitType = iota
	UnitTypeRoyalAnt
	UnitTypeDefaultRoach
	UnitTypeRoyalRoach
)

type Unit struct {
	ID          uuid.UUID
	Stats       *UnitStats
	Position    *image.Point
	Type        UnitType
	Rect        *image.Rectangle
	MovingAngle float64

	Destination     *image.Point
	DestinationType DestinationType
	Action          Action
	NearestEnemy    *Unit
	NearestHome     BuildingInterface
	LastResourcePos *image.Point
	CurrentAnim     string

	Faction uint
}

type UnitStats struct {
	HPMax     uint
	HPCur     uint
	MoveSpeed uint
	Damage    uint
	Range     uint

	MaxCarryCapactiy    uint
	ResourceCarried     uint
	ResourceTypeCarried string
	ResourceCollectTime uint
}

func NewRoyalRoach() *Unit {
	u := NewDefaultAnt()
	u.Type = UnitTypeRoyalRoach
	size := 192 // match sprite
	u.Rect.Min = image.Point{0, 0}
	u.Rect.Max = image.Point{size, size}
	return u
}

func NewRoyalAnt() *Unit {
	u := NewDefaultAnt()
	u.Type = UnitTypeRoyalAnt
	size := 192 // match sprite
	u.Rect.Min = image.Point{0, 0}
	u.Rect.Max = image.Point{size, size}
	return u
}

func NewDefaultRoach() *Unit {
	u := NewDefaultAnt()
	u.Type = UnitTypeDefaultRoach
	return u
}

func NewDefaultAnt() *Unit {
	return &Unit{
		ID:   uuid.New(),
		Type: UnitTypeDefaultAnt,
		Stats: &UnitStats{
			HPMax:     100,
			HPCur:     100,
			MoveSpeed: 10,
			Damage:    10,
			Range:     15,
			// acceleration / current speed?
			MaxCarryCapactiy:    5,
			ResourceCarried:     0,
			ResourceTypeCarried: "",
		},
		Position: &image.Point{0, 0},
		Rect: &image.Rectangle{
			Min: image.Point{0, 0},
			Max: image.Point{128, 128},
		},
		Destination: &image.Point{0, 0},
		Action:      IdleAction,
		Faction:     uint(PlayerFaction),
	}
}

func (unit *Unit) Update(sim *T) {
	switch unit.Action {
	case IdleAction:
		return
	case MovingAction:
		unit.MoveToDestination(sim, false)
	case AttackMovingAction:
		if unit.NearestEnemy != nil && unit.TargetInRange(*unit.Position) {
			unit.NearestEnemy.Stats.HPCur -= unit.Stats.Damage
			// pew pew animation
		} else {
			unit.MoveToDestination(sim, false) // destination might be a unit?
		}
	case HoldingPositionAction:
		if unit.NearestEnemy != nil && unit.TargetInRange(*unit.Position) {
			unit.NearestEnemy.Stats.HPCur -= unit.Stats.Damage
			// pew pew animation
		}
	case CollectingAction:
		// if we are holding some resources, set home, then set deliveringAction
		if unit.Stats.ResourceCarried > 0 { // better logic so it doesnt always bring back minimal resource amount
			if unit.NearestHome == nil {
				for _, hive := range sim.GetAllBuildings() {
					if hive.GetFaction() == unit.Faction {
						unit.NearestHome = hive
						break
					}
				}
			}
			unit.LastResourcePos = unit.Destination
			unit.Destination = unit.NearestHome.GetClosestPosition(unit.Position.X, unit.Position.Y)
			unit.Action = DeliveringAction
		} else {
			// move to and collect resource
			unit.MoveToDestination(sim, false) // setting this to True causes jank behavior and its better as false?
			dist := unit.DistanceTo(*unit.Destination)
			if dist < 200 { // lots of tweaks needed here or fixes TODO
				// TODO: play animation and wait some time to harvest?
				unit.Stats.ResourceCollectTime += 1
				if unit.Stats.ResourceCollectTime >= uint(MaxResourceCollectFrames) {
					unit.Stats.ResourceCollectTime = 0
					tile := sim.world.TileMap.GetTileByPosition(unit.Destination.X, unit.Destination.Y)
					if tile != nil && tile.Type != "none" {
						unit.Stats.ResourceCarried = 5
						unit.Stats.ResourceTypeCarried = tile.Type
					}
				}
			}
		}
	case DeliveringAction:
		// return resource to home base
		// set home if unset
		unit.MoveToDestination(sim, false) // setting this to True causes jank behavior and its better as false?
		dist := unit.EdgeDistanceTo(*unit.Destination)
		if dist < 140 { // lots of tweaks needed here or fixes TODO
			if unit.Stats.ResourceTypeCarried == "wood" {
				sim.AddWood(unit.Stats.ResourceCarried)
				unit.Stats.ResourceCarried = 0
				unit.Stats.ResourceTypeCarried = ""
			} else if unit.Stats.ResourceTypeCarried == "sucrose" {
				sim.AddSucrose(unit.Stats.ResourceCarried)
				unit.Stats.ResourceCarried = 0
				unit.Stats.ResourceTypeCarried = ""
			}
			unit.Destination = unit.LastResourcePos
			unit.Action = CollectingAction
		}

	}
}
func (unit *Unit) MoveToDestination(sim *T, harvesting bool) {
	speed := float64(unit.Stats.MoveSpeed)
	oldX := unit.Position.X
	oldY := unit.Position.Y

	dx := float64(unit.Destination.X - unit.Position.X)
	dy := float64(unit.Destination.Y - unit.Position.Y)

	// Movement request
	moveX := math.Copysign(math.Min(math.Abs(dx), speed), dx) // move by at most `speed` towards target X
	moveY := math.Copysign(math.Min(math.Abs(dy), speed), dy) // move by at most `speed` towards target Y

	// Attempt X movement
	if moveX != 0 {
		newX := unit.Position.X + int(moveX)
		newY := unit.Position.Y
		candidate := &image.Rectangle{
			Min: image.Point{X: newX, Y: newY},
			Max: image.Point{X: newX + unit.Rect.Dx(), Y: newY + unit.Rect.Dy()},
		}

		if !unit.isColliding(candidate, sim) {
			unit.SetPosition(&image.Point{X: unit.Position.X + int(moveX), Y: unit.Position.Y})
		}
	}

	// Attempt Y movement
	if moveY != 0 {
		newY := unit.Position.Y + int(moveY)
		newX := unit.Position.X
		candidate := &image.Rectangle{
			Min: image.Point{X: newX, Y: newY},
			Max: image.Point{X: newX + unit.Rect.Dx(), Y: newY + unit.Rect.Dy()},
		}
		if !unit.isColliding(candidate, sim) {
			unit.SetPosition(&image.Point{X: unit.Position.X, Y: unit.Position.Y + int(moveY)})
		}
	}

	dxRot := float64(unit.Position.X - oldX)
	dyRot := float64(unit.Position.Y - oldY)
	if dxRot != 0 || dyRot != 0 { // update angle only if moved
		unit.MovingAngle = math.Atan2(dyRot, dxRot) + math.Pi/2 // adjust for sprite orientation
	}

	// Final snapping
	if math.Abs(dx) <= float64(ArrivalThreshold) && math.Abs(dy) <= float64(ArrivalThreshold) {
		unit.SetPosition(&image.Point{X: unit.Destination.X, Y: unit.Destination.Y})
		unit.Action = IdleAction
	}
}

func (unit *Unit) isColliding(rect *image.Rectangle, sim *T) bool {
	colliders := sim.GetAllCollidersOverlapping(rect)
	for _, collider := range colliders {
		if collider.OwnerID == unit.ID.String() {
			continue // skip self
		}
		if collider.Rect.Overlaps(*rect) {
			return true
		}
	}
	for _, mo := range sim.world.MapObjects {
		if mo.Rect.Overlaps(*rect) {
			return true
		}
	}
	return false
}

func (unit *Unit) SetNearestEnemy(target *Unit) {
	unit.NearestEnemy = target
}

func (unit *Unit) DistanceTo(point image.Point) uint {
	selfCentered := unit.GetCenteredPosition()
	xDist := math.Abs(float64(selfCentered.X - point.X))
	yDist := math.Abs(float64(selfCentered.Y - point.Y))
	return uint(math.Sqrt(math.Pow(xDist, 2) + math.Pow(yDist, 2)))
}

func (unit *Unit) EdgeDistanceTo(point image.Point) uint {
	// Calculate the shortest distance from any edge of unit.Rect to the given point.
	rect := unit.Rect
	px, py := point.X, point.Y

	// Clamp point to the rectangle to find the closest point on the edge
	clampedX := math.Max(float64(rect.Min.X), math.Min(float64(px), float64(rect.Max.X)))
	clampedY := math.Max(float64(rect.Min.Y), math.Min(float64(py), float64(rect.Max.Y)))

	dx := float64(px) - clampedX
	dy := float64(py) - clampedY

	return uint(math.Sqrt(dx*dx + dy*dy))
}

func (unit *Unit) TargetInRange(point image.Point) bool {
	return unit.DistanceTo(point) <= unit.Stats.Range
}

func (unit *Unit) SetPosition(pos *image.Point) {
	sizeX := unit.Rect.Dx()
	sizeY := unit.Rect.Dy()
	unit.Position = pos
	unit.Rect.Min = *pos
	unit.Rect.Max = image.Point{
		X: pos.X + sizeX,
		Y: pos.Y + sizeY,
	}
}

func (unit *Unit) SetTilePosition(x, y int) {
	unit.SetPosition(&image.Point{X: x * 128, Y: y * 128})
}

func (unit *Unit) GetCenteredPosition() *image.Point {
	return &image.Point{
		X: unit.Position.X + unit.Rect.Dx()/2,
		Y: unit.Position.Y + unit.Rect.Dy()/2,
	}
}
