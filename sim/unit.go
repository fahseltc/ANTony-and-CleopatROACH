package sim

import (
	"gamejam/util"
	"gamejam/vec2"
	"image"
	"math"
	"math/rand/v2"

	"github.com/google/uuid"
)

var ArrivalThreshold = 30
var MaxResourceCollectFrames = 30
var PlayerFaction = 0

var (
	TileSize     = 128.0
	HalfTileSize = 64.0
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
	Position    *vec2.T
	Type        UnitType
	Rect        *image.Rectangle
	MovingAngle float64

	Destinations          *util.Queue[*vec2.T]
	DestinationType       DestinationType
	Action                Action
	NearestEnemy          *Unit
	NearestHome           BuildingInterface
	LastResourcePos       *vec2.T
	CurrentAnim           string
	StuckFrames           int
	StuckSidestepAttempts int

	Faction uint
}

type UnitStats struct {
	HPMax     uint
	HPCur     uint
	MoveSpeed uint

	Damage          uint
	AttackRange     uint
	AttackFrames    uint
	AttackFramesCur uint

	MaxCarryCapactiy    uint
	ResourceCarried     uint
	ResourceTypeCarried ResourceType
	ResourceCollectTime uint

	VisionRange uint
}

func NewRoyalRoach() *Unit {
	u := NewDefaultAnt()
	u.Type = UnitTypeRoyalRoach
	size := 128 // match sprite
	u.Rect.Min = image.Point{0, 0}
	u.Rect.Max = image.Point{size, size}
	return u
}

func NewRoyalAnt() *Unit {
	u := NewDefaultAnt()
	u.Type = UnitTypeRoyalAnt
	size := 128 // match sprite
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

			Damage:       10,
			AttackRange:  100,
			AttackFrames: 30,

			MaxCarryCapactiy:    5,
			ResourceCarried:     0,
			ResourceTypeCarried: ResourceTypeNone,
			VisionRange:         4,
		},
		Position: &vec2.T{},
		Rect: &image.Rectangle{
			Min: image.Point{0, 0},
			Max: image.Point{int(HalfTileSize), int(HalfTileSize)},
		},
		Destinations: util.NewQueue[*vec2.T](),
		Action:       IdleAction,
		Faction:      uint(PlayerFaction),
	}
}

func (unit *Unit) Update(sim *T) {
	// Check for self death
	if unit == nil {
		return
	}

	if unit.Stats.HPCur <= 0 {
		sim.RemoveUnit(unit)
		// animation?
		return
	}

	switch unit.Action {
	case IdleAction:
		unit.Stats.ResourceCollectTime = 0
		return
	case MovingAction:
		unit.Stats.ResourceCollectTime = 0
		unit.MoveToDestination(sim)
	case AttackMovingAction:
		unit.Stats.ResourceCollectTime = 0
		if unit.NearestEnemy != nil && unit.NearestEnemy.IsAlive() && unit.TargetInRange(unit.NearestEnemy.GetCenteredPosition()) {
			unit.Stats.AttackFramesCur++
			if unit.Stats.AttackFramesCur >= unit.Stats.AttackFrames {
				unit.NearestEnemy.Stats.HPCur -= unit.Stats.Damage
				unit.Stats.AttackFramesCur = 0
				// Play SFX?
			}

		} else {
			unit.MoveToDestination(sim) // destination might be a unit?
			// check if there is an enemy unit in range and set it as NearestEnemy
			for _, enemy := range sim.GetAllEnemyUnits() {
				if enemy.GetCenteredPosition().Distance(*unit.GetCenteredPosition()) <= 300 {
					// TODO: make a list of all the nearby enemies and pick the closest?
					unit.NearestEnemy = enemy
					break
				}
			}
		}
	case HoldingPositionAction:
		if unit.NearestEnemy != nil && unit.TargetInRange(unit.Position) {
			unit.NearestEnemy.Stats.HPCur -= unit.Stats.Damage
			// pew pew animation
		}
	case CollectingAction:
		// if we are holding some resources, set home, then set DeliveringAction
		if unit.Stats.ResourceCarried > 0 {
			// Find the nearest hive and set it as the unit's home
			var nearest BuildingInterface
			minDist := uint(math.MaxUint32)
			for _, hive := range sim.GetAllBuildings() {
				if hive.GetFaction() == unit.Faction {
					dist := unit.DistanceTo(hive.GetCenteredPosition())
					if nearest == nil || dist < minDist {
						nearest = hive
						minDist = dist
					}
				}
			}
			unit.NearestHome = nearest
			unit.Stats.ResourceCollectTime = 0

			path := sim.FindClickedPath(unit.GetTileCoordinates(), unit.NearestHome.GetTilePosition())
			for _, p := range path {
				unit.Destinations.Enqueue(p.ToCenteredPixelCoordinates())
			}
			unit.Action = DeliveringAction
		} else {
			// move to and collect resource
			unit.MoveToDestination(sim)
			dist := unit.DistanceTo(unit.LastResourcePos)
			if dist < 150 { // lots of tweaks needed here or fixes TODO
				// TODO: play animation
				unit.Destinations.Clear()
				unit.Stats.ResourceCollectTime += 1
				if unit.Stats.ResourceCollectTime >= uint(MaxResourceCollectFrames) {
					unit.Stats.ResourceCollectTime = 0
					tile := sim.world.TileMap.GetTileByPosition(int(unit.LastResourcePos.X), int(unit.LastResourcePos.Y))
					if tile != nil && tile.Type != "none" {
						unit.Stats.ResourceCarried = 5
						var resType ResourceType
						switch tile.Type {
						case "plain":
							resType = ResourceTypeNone
						case "sucrose":
							resType = ResourceTypeSucrose
						case "wood":
							resType = ResourceTypeWood
						}
						unit.Stats.ResourceTypeCarried = resType
					}
				}
			}
		}
	case DeliveringAction: // return resource to home base
		unit.MoveToDestination(sim)
		dist := unit.EdgeDistanceTo(unit.NearestHome.GetCenteredPosition())
		if dist < 195 { // lots of tweaks needed here or fixes TODO
			sim.AddResource(uint(unit.Stats.ResourceCarried), unit.Stats.ResourceTypeCarried)
			unit.Stats.ResourceCarried = 0
			unit.Stats.ResourceTypeCarried = ResourceTypeNone
			unit.Destinations.Clear()

			resPos := &vec2.T{
				X: (unit.LastResourcePos.X - HalfTileSize) / TileSize,
				Y: (unit.LastResourcePos.Y - HalfTileSize) / TileSize,
			}
			path := sim.FindClickedPath(unit.GetTileCoordinates(), resPos)
			for _, p := range path {
				unit.Destinations.Enqueue(p.ToCenteredPixelCoordinates())
			}
			unit.Action = CollectingAction
		}

	}
}
func (unit *Unit) MoveToDestination(sim *T) {
	dest, err := unit.Destinations.Peek()
	if err != nil {
		return
	}
	speed := float64(unit.Stats.MoveSpeed)
	oldPos := unit.GetCenteredPosition()
	oldX := oldPos.X
	oldY := oldPos.Y

	dx := float64(dest.X - unit.GetCenteredPosition().X)
	dy := float64(dest.Y - unit.GetCenteredPosition().Y)
	// Repulsion avoidance
	repulsion := unit.computeRepulsion(sim)

	// Blend desired direction and repulsion
	toTarget := vec2.T{X: dx, Y: dy}
	if toTarget.Length() > 0 {
		toTarget = toTarget.Normalize()
	}
	moveVec := toTarget.Add(*repulsion).Normalize().Scale(speed)

	moveX := math.Copysign(math.Min(math.Abs(moveVec.X), speed), moveVec.X)
	moveY := math.Copysign(math.Min(math.Abs(moveVec.Y), speed), moveVec.Y)

	// Attempt X movement
	if moveX != 0 {
		newX := unit.Position.X + moveX
		newY := unit.Position.Y
		candidate := &image.Rectangle{
			Min: image.Point{X: int(newX), Y: int(newY)},
			Max: image.Point{X: int(newX) + unit.Rect.Dx(), Y: int(newY) + unit.Rect.Dy()},
		}
		if !unit.isColliding(candidate, sim) {
			unit.SetPosition(&vec2.T{X: newX, Y: unit.Position.Y})
		}
	}

	// Attempt Y movement
	if moveY != 0 {
		newY := unit.Position.Y + moveY
		newX := unit.Position.X
		candidate := &image.Rectangle{
			Min: image.Point{X: int(newX), Y: int(newY)},
			Max: image.Point{X: int(newX) + unit.Rect.Dx(), Y: int(newY) + unit.Rect.Dy()},
		}
		if !unit.isColliding(candidate, sim) {
			unit.SetPosition(&vec2.T{X: unit.Position.X, Y: newY})
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
			//unit.NavigateAround(sim)
		}

		if unit.StuckFrames > 2000 { //|| unit.StuckSidestepAttempts > 3
			//Only sidestep if the destination itself isn't clearly blocked
			if unit.isDestinationBlocked(sim) {
				unit.Action = IdleAction
				unit.StuckFrames = 0
				return
			}
		}
	}
	if arrived && len(unit.Destinations.Items) >= 1 {
		unit.Destinations.Dequeue()
	}
	if arrived && len(unit.Destinations.Items) == 0 {
		nearbyUnits := sim.GetAllNearbyFriendlyUnits(unit)
		for _, nearbyUnit := range nearbyUnits {
			nearbyUnit.SendMessage(sim, UnitMessageArrivedIdle)
		}
	}
}

func (unit *Unit) isColliding(rect *image.Rectangle, sim *T) bool {
	futureUnitCenterX := float64(rect.Min.X+rect.Max.X) / 2
	futureUnitCenterY := float64(rect.Min.Y+rect.Max.Y) / 2
	futureUnitRadius := float64(rect.Dx()) / 2

	colliders := sim.GetAllCollidersOverlapping(rect)
	for _, collider := range colliders {
		if collider.OwnerID == unit.ID.String() {
			continue // skip self
		}
		if collider.Radius > 0 && collider.Center != (image.Point{}) {
			dx := futureUnitCenterX - float64(collider.Center.X)
			dy := futureUnitCenterY - float64(collider.Center.Y)
			distance := math.Sqrt(dx*dx + dy*dy)
			if distance < futureUnitRadius+float64(collider.Radius) {
				return true
			}
		} else if collider.Rect.Overlaps(*rect) {
			return true
		}
	}
	for _, mo := range sim.world.MapObjects {
		closestX := math.Max(float64(mo.Rect.Min.X), math.Min(futureUnitCenterX, float64(mo.Rect.Max.X)))
		closestY := math.Max(float64(mo.Rect.Min.Y), math.Min(futureUnitCenterY, float64(mo.Rect.Max.Y)))
		dx := futureUnitCenterX - closestX
		dy := futureUnitCenterY - closestY
		if dx*dx+dy*dy <= futureUnitRadius*futureUnitRadius {
			return true
		}
	}
	return false
}

func (unit *Unit) TrySidestep(sim *T) {
	bestOffset := vec2.T{}

	// Try 8 directions (N, NE, E, SE, S, SW, W, NW)
	offsets := []vec2.T{
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
		newX := unit.Position.X + off.X*float64(unit.Stats.MoveSpeed)
		newY := unit.Position.Y + off.Y*float64(unit.Stats.MoveSpeed)
		candidate := &image.Rectangle{
			Min: image.Point{X: int(newX), Y: int(newY)},
			Max: image.Point{X: int(newX) + unit.Rect.Dx(), Y: int(newY) + unit.Rect.Dy()},
		}
		if !unit.isColliding(candidate, sim) {
			bestOffset = off
		}
	}

	// Apply best offset if found
	if bestOffset != (vec2.T{}) {
		newX := unit.Position.X + bestOffset.X*float64(unit.Stats.MoveSpeed)
		newY := unit.Position.Y + bestOffset.Y*float64(unit.Stats.MoveSpeed)
		unit.SetPosition(&vec2.T{X: newX, Y: newY})
	}
}

func (unit *Unit) SetNearestEnemy(target *Unit) {
	unit.NearestEnemy = target
}

func (unit *Unit) DistanceTo(point *vec2.T) uint {
	selfCentered := unit.GetCenteredPosition()
	xDist := math.Abs(float64(selfCentered.X - point.X))
	yDist := math.Abs(float64(selfCentered.Y - point.Y))
	return uint(math.Sqrt(math.Pow(xDist, 2) + math.Pow(yDist, 2)))
}

func (unit *Unit) EdgeDistanceTo(point *vec2.T) uint {
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

func (unit *Unit) TargetInRange(point *vec2.T) bool {
	return unit.DistanceTo(point) <= unit.Stats.AttackRange
}

func (unit *Unit) SetPosition(pos *vec2.T) {
	sizeX := unit.Rect.Dx()
	sizeY := unit.Rect.Dy()
	unit.Position = pos
	unit.Rect.Min = image.Point{
		X: int(pos.X),
		Y: int(pos.Y),
	}
	unit.Rect.Max = image.Point{
		X: int(pos.X) + sizeX,
		Y: int(pos.Y) + sizeY,
	}
}

func (unit *Unit) SetTilePosition(x, y int) {
	unit.SetPosition(&vec2.T{X: float64(x * int(TileSize)), Y: float64(y * int(TileSize))})
}

func (unit *Unit) GetTileCoordinates() *vec2.T {
	return &vec2.T{
		X: math.Round(unit.Position.X / TileSize),
		Y: math.Round(unit.Position.Y / TileSize),
	}
}
func (unit *Unit) computeRepulsion(sim *T) *vec2.T {
	repulsion := vec2.T{}
	myCenter := unit.GetCenteredPosition()

	for _, other := range sim.GetAllUnits() {
		if other.ID == unit.ID {
			continue
		}
		otherCenter := other.GetCenteredPosition()
		dir := myCenter.Sub(*otherCenter)
		dist := dir.Length()

		if dist < 160 && dist > 0.1 {
			// Normalize direction and scale force
			push := dir.Normalize().Scale((160 - dist) / 160)

			// Apply sideways deflection if mostly head-on
			if math.Abs(push.X) < 0.2 && math.Abs(push.Y) > 0.5 {
				// Deflect left or right
				deflect := vec2.T{X: 1, Y: 0}
				if rand.IntN(2) == 0 {
					deflect.X = -1
				}
				push = push.Add(deflect.Scale(0.5))
			}

			repulsion = repulsion.Add(push)
		}
	}

	return &repulsion
}

func (unit *Unit) GetCenteredPosition() *vec2.T {
	return &vec2.T{
		X: unit.Position.X + float64(unit.Rect.Dx())/2,
		Y: unit.Position.Y + float64(unit.Rect.Dy())/2,
	}
}

func (unit *Unit) isDestinationBlocked(sim *T) bool {
	dest, _ := unit.Destinations.Peek()
	destRect := &image.Rectangle{
		Min: dest.ToPoint(),
		Max: image.Point{
			X: int(dest.X) + unit.Rect.Dx(),
			Y: int(dest.Y) + unit.Rect.Dy(),
		},
	}

	// If the destination itself is colliding, it's likely invalid
	if unit.isColliding(destRect, sim) {
		return true
	}

	// Check 8 adjacent tiles for walls — if all are blocked, it's surrounded
	blockedSides := 0
	offsets := []vec2.T{
		{X: -1, Y: 0}, {X: 1, Y: 0},
		{X: 0, Y: -1}, {X: 0, Y: 1},
		{X: -1, Y: -1}, {X: 1, Y: -1},
		{X: -1, Y: 1}, {X: 1, Y: 1},
	}
	for _, off := range offsets {
		pos := vec2.T{
			X: dest.X + off.X*float64(unit.Rect.Dx()),
			Y: dest.Y + off.Y*float64(unit.Rect.Dy()),
		}
		rect := &image.Rectangle{
			Min: pos.ToPoint(),
			Max: image.Point{X: int(pos.X) + unit.Rect.Dx(), Y: int(pos.Y) + unit.Rect.Dy()},
		}
		if unit.isColliding(rect, sim) {
			blockedSides++
		}
	}

	return blockedSides >= len(offsets) // surrounded
}

func (unit *Unit) IsAlive() bool {
	return unit.Stats.HPCur > 0
}
