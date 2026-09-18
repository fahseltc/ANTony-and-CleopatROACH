package sim

import (
	"gamejam/types"
	"gamejam/util"
	"gamejam/vec2"
	"image"
	"math"

	"github.com/google/uuid"
)

var ArrivalThreshold = 80
var UnitAttackRangeBuffer = 10

var PlayerFaction = 0

var (
	TileSize     = 128.0
	HalfTileSize = 64.0
)

type Unit struct {
	ID          uuid.UUID
	Stats       *UnitStats
	Position    *vec2.T
	Type        types.Unit
	Rect        *image.Rectangle
	MovingAngle float64

	Destinations    *util.Queue[*vec2.T]
	DestinationType types.Destination

	CurrentState UnitStateInterface

	NearestEnemy    *Unit
	NearestHome     BuildingInterface
	LastResourcePos *vec2.T
	CurrentAnim     string
	StuckFrames     int
	StuckAttempts   int

	Faction uint
}

type UnitStats struct {
	Name          string
	ToolTipString string
	HPMax         uint
	HPCur         uint
	MoveSpeed     uint
	SizePx        uint

	Damage          uint
	AttackRange     uint
	AttackFrames    uint
	AttackFramesCur uint

	MaxCarryCapacity    uint
	ResourcesCarried    uint
	ResourceTypeCarried types.Resource

	ConstructionTime uint
	ResourceCost     ResourceCost

	VisionRange uint
}

// RoyalUnitSize is the collision-rect size (px) of the royal ant/roach. It is
// deliberately larger than one tile (128px) so the royals are physically big:
// with radius = RoyalUnitSize/2 > 64px (half a tile), a royal centered on a
// single-tile-wide bridge still clips the unbridged water on either side, so it
// CANNOT cross a 1-wide bridge - it needs a wider crossing. Regular units
// (128px, radius 64) still fit a 1-wide bridge exactly.
const RoyalUnitSize = 176

func NewRoyalRoach() *Unit {
	u := GetUnitInstance(types.UnitTypeRoyalRoach, uint(PlayerFaction))
	u.Type = types.UnitTypeRoyalRoach
	u.Rect.Min = image.Point{0, 0}
	u.Rect.Max = image.Point{RoyalUnitSize, RoyalUnitSize}
	return u
}

func NewRoyalAnt() *Unit {
	u := GetUnitInstance(types.UnitTypeRoyalAnt, uint(PlayerFaction))
	u.Type = types.UnitTypeRoyalAnt
	u.Rect.Min = image.Point{0, 0}
	u.Rect.Max = image.Point{RoyalUnitSize, RoyalUnitSize}
	return u
}

func NewFighterAnt() *Unit {
	u := GetUnitInstance(types.UnitTypeFighterAnt, uint(PlayerFaction))
	return u
}

func NewDefaultRoach() *Unit {
	u := GetUnitInstance(types.UnitTypeDefaultRoach, uint(PlayerFaction))
	u.Type = types.UnitTypeDefaultRoach
	return u
}

func NewDefaultAnt() *Unit {
	return GetUnitInstance(types.UnitTypeDefaultAnt, uint(PlayerFaction))
}

func NewDefaultAntWithTilePosition(x int, y int) *Unit {
	u := GetUnitInstance(types.UnitTypeDefaultAnt, uint(PlayerFaction))
	u.SetTilePosition(x, y)
	return u
}

func (unit *Unit) findNearestEnemy(sim *T) *Unit {
	bestScore := math.Inf(-1)
	var bestTarget *Unit
	for _, enemy := range sim.GetAllEnemyUnitsByFaction(unit.Faction) {
		if enemy.ID == unit.ID {
			continue
		}
		score := unit.EvaluateEnemy(enemy)
		if score > bestScore {
			bestScore = score
			bestTarget = enemy
		}
	}
	if bestTarget != nil {
		return bestTarget
	}
	return nil
}

func (unit *Unit) EvaluateEnemy(enemy *Unit) float64 {
	if !enemy.IsAlive() {
		return math.Inf(-1)
	}

	score := 0.0
	distance := unit.GetCenteredPosition().Distance(*enemy.GetCenteredPosition())
	if distance > float64(unit.Stats.VisionRange*uint(TileSize)) {
		return math.Inf(-1)
	}

	// Prefer closer enemies
	score -= distance

	// Prefer low HP enemies
	score -= float64(enemy.Stats.HPCur) * 1.5

	// Prefer enemies attacking us
	if enemy.NearestEnemy != nil && enemy.NearestEnemy.ID == unit.ID {
		score += 300
	}

	// Prefer high damage enemies
	score += float64(enemy.Stats.Damage) * 2

	return score
}

func (unit *Unit) Update(sim *T) {
	// Check for self death
	if unit == nil || unit.Stats.HPCur <= 0 {
		sim.RemoveUnit(unit) // this accepts nil unit and just returns
		return
	}

	if unit.CurrentState == nil {
		unit.ChangeState(&IdleState{}) // default state
	}

	unit.CurrentState.Update(unit, sim)
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

func (unit *Unit) TargetInAttackRange(point *vec2.T) bool {
	val := unit.EdgeDistanceTo(point)
	if val == 0 {
		return false
	}

	return val <= unit.Stats.AttackRange+uint(UnitAttackRangeBuffer)
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

func (unit *Unit) GetCenteredPosition() *vec2.T {
	return &vec2.T{
		X: unit.Position.X + float64(unit.Rect.Dx())/2,
		Y: unit.Position.Y + float64(unit.Rect.Dy())/2,
	}
}

func (unit *Unit) IsAlive() bool {
	return unit.Stats.HPCur > 0
}

func (unit *Unit) IsWorker() bool {
	return unit.Type == types.UnitTypeDefaultAnt || unit.Type == types.UnitTypeDefaultRoach
}

// HarvestApproachPos returns the pixel-center of a walkable tile adjacent to the
// resource tile at resourceCenter, i.e. an actual spot on the map the unit can
// stand to harvest. Workers are distributed across the available adjacent tiles
// deterministically by unit ID, so a group sent to one resource fans out onto
// different neighbouring tiles instead of all targeting the same point.
//
// Unlike a purely geometric ring offset, this only ever returns tiles that are
// on the map, not impassable, not themselves resource tiles, and — importantly —
// actually reachable from the unit's current position via pathfinding. This
// prevents a worker from being assigned a tile that is walkable in isolation but
// boxed in behind the resource, which made it wiggle in place instead of
// harvesting. If no adjacent tile is reachable it falls back to the closest
// walkable neighbour, and finally to the resource center.
func (unit *Unit) HarvestApproachPos(sim *T, resourceCenter *vec2.T) *vec2.T {
	if resourceCenter == nil {
		return resourceCenter
	}

	// Resource tile coordinates from the center pixel.
	rx := int(resourceCenter.X / TileSize)
	ry := int(resourceCenter.Y / TileSize)

	// Candidate stand tiles: the 8 neighbours around the resource tile.
	dirs := []struct{ dx, dy int }{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1}, // cardinals first (preferred)
		{-1, -1}, {1, -1}, {-1, 1}, {1, 1}, // diagonals
	}

	// Collect walkable, non-resource neighbour tiles, split by whether the unit
	// can actually path to them. A tile can be walkable in isolation but boxed
	// in behind the resource/other obstacles so the unit can never reach it —
	// picking such a tile is what makes a worker wiggle behind the resource
	// instead of harvesting. We only assign reachable tiles.
	unitTile := unit.GetTileCoordinates()
	var reachable []*vec2.T
	var walkableButUnreachable []*vec2.T
	for _, d := range dirs {
		nx, ny := rx+d.dx, ry+d.dy
		tile := sim.world.TileMap.GetTileByCoordinates(nx, ny)
		if tile == nil || tile.HasCollision {
			continue
		}
		// Don't stand on another resource tile — units can't occupy those and
		// it's the source of the "sent inside another resource" bug.
		if tile.Type == types.TileTypeWood || tile.Type == types.TileTypeSucrose {
			continue
		}
		center := &vec2.T{
			X: float64(nx)*TileSize + HalfTileSize,
			Y: float64(ny)*TileSize + HalfTileSize,
		}
		// Reachability: does a path exist from the unit's tile to this tile?
		// (If the unit is already standing on the candidate tile, it's trivially
		// reachable.)
		if (int(unitTile.X) == nx && int(unitTile.Y) == ny) ||
			len(sim.world.TileMap.FindPath(unitTile, &vec2.T{X: float64(nx), Y: float64(ny)})) > 0 {
			reachable = append(reachable, center)
		} else {
			walkableButUnreachable = append(walkableButUnreachable, center)
		}
	}

	// Prefer reachable tiles; distribute deterministically per-unit so a group
	// sent to one resource fans out but each unit's target is stable.
	if len(reachable) > 0 {
		var seed uint32
		for _, b := range unit.ID {
			seed = seed*31 + uint32(b)
		}
		return reachable[seed%uint32(len(reachable))]
	}

	// Nothing reachable. Aim at the closest walkable-but-unreachable neighbour if
	// any (pathing will get the unit as close as it can); otherwise the resource
	// center. Either way this avoids committing to a far unreachable slot that
	// leaves the unit stuck on the wrong side of the resource.
	if len(walkableButUnreachable) > 0 {
		closest := walkableButUnreachable[0]
		minDist := unit.GetCenteredPosition().Distance(*closest)
		for _, c := range walkableButUnreachable[1:] {
			if d := unit.GetCenteredPosition().Distance(*c); d < minDist {
				minDist = d
				closest = c
			}
		}
		return closest
	}
	return resourceCenter
}

func (unit *Unit) ChangeState(newState UnitStateInterface) {
	if unit.CurrentState != nil {
		unit.CurrentState.Exit(unit)
	}
	unit.CurrentState = newState
	if unit.CurrentState != nil {
		unit.CurrentState.Enter(unit)
	}
}

func (unit *Unit) isDestinationBlocked(sim *T) bool {
	dest, err := unit.Destinations.Peek()
	if err != nil || dest == nil {
		return false
	}

	destRect := &image.Rectangle{
		Min: image.Point{X: int(dest.X), Y: int(dest.Y)},
		Max: image.Point{X: int(dest.X) + unit.Rect.Dx(), Y: int(dest.Y) + unit.Rect.Dy()},
	}

	return unit.isColliding(destRect, sim)
}
