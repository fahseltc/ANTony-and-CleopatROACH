package sim

import (
	"image"
	"math"
	"math/rand/v2"

	"github.com/google/uuid"
)

var ArrivalThreshold = 25
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

	Destination           *image.Point
	DestinationType       DestinationType
	Action                Action
	NearestEnemy          *Unit
	NearestHome           BuildingInterface
	LastResourcePos       *image.Point
	CurrentAnim           string
	StuckFrames           int
	StuckSidestepAttempts int

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
			// Find the nearest hive and set it as the unit's home
			var nearest BuildingInterface
			minDist := uint(math.MaxUint32)
			for _, hive := range sim.GetAllBuildings() {
				if hive.GetFaction() == unit.Faction {
					dist := unit.DistanceTo(*hive.GetCenteredPosition())
					if nearest == nil || dist < minDist {
						nearest = hive
						minDist = dist
					}
				}
			}
			unit.NearestHome = nearest
			unit.LastResourcePos = unit.Destination
			unit.Destination = unit.NearestHome.GetClosestPosition(unit.Position.X, unit.Position.Y)
			unit.Action = DeliveringAction
		} else {
			// move to and collect resource
			unit.MoveToDestination(sim, false) // setting this to True causes jank behavior and its better as false?
			dist := unit.DistanceTo(*unit.Destination)
			if dist < 230 { // lots of tweaks needed here or fixes TODO
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
		if dist < 100 { // lots of tweaks needed here or fixes TODO
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
	oldPos := unit.GetCenteredPosition()
	oldX := oldPos.X
	oldY := oldPos.Y

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
			unit.SetPosition(&image.Point{X: newX, Y: unit.Position.Y})
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
			unit.SetPosition(&image.Point{X: unit.Position.X, Y: newY})
		}
	}

	// Handle Rotation
	newCentered := unit.GetCenteredPosition()
	dxRot := float64(newCentered.X - oldX)
	dyRot := float64(newCentered.Y - oldY)
	if dxRot != 0 || dyRot != 0 { // update angle only if moved
		unit.MovingAngle = math.Atan2(dyRot, dxRot) + math.Pi/2 // adjust for sprite orientation
	}
	arrived := math.Abs(dx) <= float64(ArrivalThreshold) && math.Abs(dy) <= float64(ArrivalThreshold)
	const stuckEpsilon = 1.5
	moved := math.Abs(dxRot) > stuckEpsilon || math.Abs(dyRot) > stuckEpsilon

	if !moved && !arrived && unit.Stats.ResourceCollectTime == 0 {
		unit.StuckFrames++

		if unit.StuckFrames%10 == 0 {

			unit.TrySidestep(sim)
			//unit.StuckSidestepAttempts++
		}

		if unit.StuckFrames > 200 { //|| unit.StuckSidestepAttempts > 3
			//Only sidestep if the destination itself isn't clearly blocked
			if unit.isDestinationBlocked(sim) {
				unit.Action = IdleAction
				unit.StuckFrames = 0
				return
			}
		}
	}
	// } else {
	// 	unit.StuckFrames = 0
	// 	unit.StuckSidestepAttempts = 0
	// }

	// Final snapping
	snapRect := &image.Rectangle{
		Min: *unit.Destination,
		Max: image.Point{
			X: unit.Destination.X + unit.Rect.Dx(),
			Y: unit.Destination.Y + unit.Rect.Dy(),
		},
	}
	if math.Abs(dx) <= float64(ArrivalThreshold) &&
		math.Abs(dy) <= float64(ArrivalThreshold) &&
		!unit.isColliding(snapRect, sim) {

		unit.SetPosition(unit.Destination)
		unit.Action = IdleAction
	}

	// if unit.EdgeDistanceTo(*unit.Destination) <= uint(ArrivalThreshold/2) {
	// 	unit.SetPosition(unit.Destination)
	// 	unit.Action = IdleAction
	// 	return
	// }
}

func (unit *Unit) edgeDist(pos image.Point, goal image.Point) float64 {
	cx := pos.X + unit.Rect.Dx()/2
	cy := pos.Y + unit.Rect.Dy()/2
	dx := float64(goal.X - cx)
	dy := float64(goal.Y - cy)
	return math.Sqrt(dx*dx + dy*dy)
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

// func (unit *Unit) TrySidestep(sim *T) {
// 	speed := float64(unit.Stats.MoveSpeed)
// 	offsets := []image.Point{
// 		{X: 0, Y: -int(speed)}, // up
// 		{X: 0, Y: int(speed)},  // down
// 		{X: -int(speed), Y: 0}, // left
// 		{X: int(speed), Y: 0},  // right
// 	}

// 	for _, off := range offsets {
// 		newX := unit.Position.X + off.X
// 		newY := unit.Position.Y + off.Y
// 		candidate := &image.Rectangle{
// 			Min: image.Point{X: newX, Y: newY},
// 			Max: image.Point{X: newX + unit.Rect.Dx(), Y: newY + unit.Rect.Dy()},
// 		}
// 		if !unit.isColliding(candidate, sim) {
// 			unit.SetPosition(&image.Point{X: newX, Y: newY})
// 			break
// 		}
// 	}
// }

func (unit *Unit) TrySidestep(sim *T) {
	dest := unit.Destination
	bestOffset := image.Point{}
	shortestDist := unit.DistanceTo(*dest)

	// Try 8 directions (N, NE, E, SE, S, SW, W, NW)
	offsets := []image.Point{
		{X: -1, Y: 0}, {X: 1, Y: 0},
		{X: 0, Y: -1}, {X: 0, Y: 1},
		{X: -1, Y: -1}, {X: 1, Y: -1},
		{X: -1, Y: 1}, {X: 1, Y: 1},
	}

	// Shuffle offsets to avoid always biasing same direction
	rand.Shuffle(len(offsets), func(i, j int) {
		offsets[i], offsets[j] = offsets[j], offsets[i]
	})

	for _, off := range offsets {
		newX := unit.Position.X + off.X*int(unit.Stats.MoveSpeed)
		newY := unit.Position.Y + off.Y*int(unit.Stats.MoveSpeed)
		candidate := &image.Rectangle{
			Min: image.Point{X: newX, Y: newY},
			Max: image.Point{X: newX + unit.Rect.Dx(), Y: newY + unit.Rect.Dy()},
		}
		if !unit.isColliding(candidate, sim) {
			// Check if this move gets us closer to the destination
			newDist := unit.edgeDist(image.Point{X: newX, Y: newY}, *dest)
			if newDist < float64(shortestDist) {
				bestOffset = off
				shortestDist = uint(newDist)
			}
		}
	}

	// Apply best offset if found
	if bestOffset != (image.Point{}) {
		newX := unit.Position.X + bestOffset.X*int(unit.Stats.MoveSpeed)
		newY := unit.Position.Y + bestOffset.Y*int(unit.Stats.MoveSpeed)
		unit.SetPosition(&image.Point{X: newX, Y: newY})
	}
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

// func (unit *Unit) isDestinationReachable(sim *T) bool {
// 	destRect := &image.Rectangle{
// 		Min: image.Point{
// 			X: unit.Destination.X,
// 			Y: unit.Destination.Y,
// 		},
// 		Max: image.Point{
// 			X: unit.Destination.X + unit.Rect.Dx(),
// 			Y: unit.Destination.Y + unit.Rect.Dy(),
// 		},
// 	}
// 	return !unit.isColliding(destRect, sim)
// }

func (unit *Unit) isDestinationBlocked(sim *T) bool {
	destRect := &image.Rectangle{
		Min: *unit.Destination,
		Max: image.Point{
			X: unit.Destination.X + unit.Rect.Dx(),
			Y: unit.Destination.Y + unit.Rect.Dy(),
		},
	}

	// If the destination itself is colliding, it's likely invalid
	if unit.isColliding(destRect, sim) {
		return true
	}

	// Check 8 adjacent tiles for walls — if all are blocked, it's surrounded
	blockedSides := 0
	offsets := []image.Point{
		{X: -1, Y: 0}, {X: 1, Y: 0},
		{X: 0, Y: -1}, {X: 0, Y: 1},
		{X: -1, Y: -1}, {X: 1, Y: -1},
		{X: -1, Y: 1}, {X: 1, Y: 1},
	}
	for _, off := range offsets {
		pos := image.Point{
			X: unit.Destination.X + off.X*unit.Rect.Dx(),
			Y: unit.Destination.Y + off.Y*unit.Rect.Dy(),
		}
		rect := &image.Rectangle{
			Min: pos,
			Max: image.Point{X: pos.X + unit.Rect.Dx(), Y: pos.Y + unit.Rect.Dy()},
		}
		if unit.isColliding(rect, sim) {
			blockedSides++
		}
	}

	return blockedSides >= len(offsets) // surrounded
}
