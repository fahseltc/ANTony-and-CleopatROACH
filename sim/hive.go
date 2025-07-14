package sim

import (
	"gamejam/util"
	"gamejam/vec2"
	"image"
	"math"
	"sort"
)

var UnitConstructionTime = 120

type Hive struct {
	*Building
	buildQueue      *util.Queue[*Unit]
	UnitContructing bool
}

func NewHive() BuildingInterface {
	building := NewBuilding(0, 0, TileDimensions*2, TileDimensions*2, 0, BuildingTypeHive, uint(UnitConstructionTime))
	h := &Hive{
		Building:        building,
		UnitContructing: false,
		buildQueue:      util.NewQueue[*Unit](),
	}
	return h
}

func NewRoachHive() BuildingInterface {
	building := NewBuilding(0, 0, TileDimensions*2, TileDimensions*2, 0, BuildingTypeRoachHive, uint(UnitConstructionTime))
	h := &Hive{
		Building:        building,
		UnitContructing: false,
		buildQueue:      util.NewQueue[*Unit](),
	}
	return h
}

func (h *Hive) Update(sim *T) {
	if !h.buildQueue.IsEmpty() {
		h.UnitContructing = true
		h.Stats.ProgressCurrent += 1
		if h.Stats.ProgressCurrent >= uint(UnitConstructionTime) {
			u, err := h.buildQueue.Dequeue()
			if err != nil {
				return // todo handle?
			}
			// make sure position isnt colliding with anything and try again
			u.SetPosition(h.GetNearbyPosition(sim, 128)) // unit size static for now but could change later
			sim.AddUnit(u)
			h.UnitContructing = false
			h.Stats.ProgressCurrent = 0
		}
	}
}

func (h *Hive) DistanceTo(point vec2.T) uint {
	xDist := math.Abs(float64(h.Position.X - point.X))
	yDist := math.Abs(float64(h.Position.Y - point.Y))
	return uint(math.Sqrt(math.Pow(xDist, 2) + math.Pow(yDist, 2)))
}

func (h *Hive) AddUnitToBuildQueue() {
	var unit *Unit
	switch h.Type {
	case BuildingTypeHive:
		unit = NewDefaultAnt()
	case BuildingTypeRoachHive:
		unit = NewDefaultRoach()
	}
	h.buildQueue.Enqueue(unit)
}

func (h *Hive) GetNearbyPosition(sim *T, unitSize int) *vec2.T {
	const maxRadius = 3
	const tileSize = 128
	center := h.GetCenteredPosition()

	type spawnCandidate struct {
		point   image.Point
		rect    *image.Rectangle
		density int
	}
	// TODO this sucks make it better

	var candidates []spawnCandidate

	for dx := -maxRadius; dx <= maxRadius; dx++ {
		for dy := -maxRadius; dy <= maxRadius; dy++ {
			x := int(center.X) + dx*tileSize
			y := int(center.Y) + dy*tileSize

			rect := &image.Rectangle{
				Min: image.Point{X: x - unitSize/2, Y: y - unitSize/2},
				Max: image.Point{X: x + unitSize/2, Y: y + unitSize/2},
			}

			if rect.Overlaps(*h.Rect) {
				continue
			}

			// Score this tile by number of units overlapping or nearby
			density := 0
			for _, unit := range sim.GetAllUnits() {
				if unit == nil || unit.ID.String() == h.ID.String() {
					continue
				}
				// Check proximity (not just overlap)
				unitCenter := unit.GetCenteredPosition()
				dist := math.Hypot(float64(int(unitCenter.X)-x), float64(int(unitCenter.Y)-y))
				if dist < float64(tileSize*2) { // count units within 2-tile radius
					density++
				}
			}

			candidates = append(candidates, spawnCandidate{
				point:   image.Point{X: x - unitSize/2, Y: y - unitSize/2},
				rect:    rect,
				density: density,
			})
		}
	}

	// Sort by least density first
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].density < candidates[j].density
	})

	// Try the best candidates in order
	for _, c := range candidates {
		colliders := sim.GetAllCollidersOverlapping(c.rect)
		collision := false
		for _, collider := range colliders {
			if collider.OwnerID == h.ID.String() {
				continue
			}
			if c.rect.Overlaps(*collider.Rect) {
				collision = true
				break
			}
		}

		if !collision {
			return &vec2.T{
				X: float64(c.point.X),
				Y: float64(c.point.Y),
			}
		}
	}

	// Fallback position
	return &vec2.T{
		X: float64(center.X + tileSize),
		Y: float64(center.Y),
	}
}
