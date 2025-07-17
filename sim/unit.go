package sim

import (
	"gamejam/types"
	"gamejam/util"
	"gamejam/vec2"
	"image"
	"math"
	"math/rand/v2"

	"github.com/google/uuid"
)

var ArrivalThreshold = 80

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

	MaxCarryCapacity    uint
	ResourcesCarried    uint
	ResourceTypeCarried types.Resource

	VisionRange uint
}

func NewRoyalRoach() *Unit {
	u := NewDefaultAnt()
	u.Type = types.UnitTypeRoyalRoach
	size := 128 // match sprite
	u.Rect.Min = image.Point{0, 0}
	u.Rect.Max = image.Point{size, size}
	return u
}

func NewRoyalAnt() *Unit {
	u := NewDefaultAnt()
	u.Type = types.UnitTypeRoyalAnt
	size := 128 // match sprite
	u.Rect.Min = image.Point{0, 0}
	u.Rect.Max = image.Point{size, size}
	return u
}

func NewFighterAnt() *Unit {
	u := NewDefaultAnt()
	u.Type = types.UnitTypeFighterAnt
	size := 128 // match sprite
	u.Rect.Min = image.Point{0, 0}
	u.Rect.Max = image.Point{size, size}
	return u
}

func NewDefaultRoach() *Unit {
	u := NewDefaultAnt()
	u.Type = types.UnitTypeDefaultRoach
	return u
}

func NewDefaultAnt() *Unit {
	return &Unit{
		ID:           uuid.New(),
		Type:         types.UnitTypeDefaultAnt,
		CurrentState: nil,
		Stats: &UnitStats{
			HPMax:     100,
			HPCur:     100,
			MoveSpeed: 10,

			Damage:       10,
			AttackRange:  40,
			AttackFrames: 30,

			MaxCarryCapacity:    5,
			ResourcesCarried:    0,
			ResourceTypeCarried: types.ResourceTypeNone,
			VisionRange:         4,
		},
		Position: &vec2.T{},
		Rect: &image.Rectangle{
			Min: image.Point{0, 0},
			Max: image.Point{int(HalfTileSize), int(HalfTileSize)},
		},
		Destinations: util.NewQueue[*vec2.T](),
		Faction:      uint(PlayerFaction),
	}
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

	return val <= unit.Stats.AttackRange
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
				push = push.Add(deflect.Scale(1.5))
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

func (unit *Unit) IsAlive() bool {
	return unit.Stats.HPCur > 0
}

func (unit *Unit) IsWorker() bool {
	return unit.Type == types.UnitTypeDefaultAnt || unit.Type == types.UnitTypeDefaultRoach
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
