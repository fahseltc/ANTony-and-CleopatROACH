package sim

import (
	"gamejam/util"
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
		h.ProgressCurrent += 1
		if h.ProgressCurrent >= uint(UnitConstructionTime) {
			u, err := h.buildQueue.Dequeue()
			if err != nil {
				return // todo handle?
			}
			// make sure position isnt colliding with anything and try again
			u.SetPosition(h.GetNearbyPosition(sim, 128)) // unit size static for now but could change later
			sim.AddUnit(u)
			h.UnitContructing = false
			h.ProgressCurrent = 0
		}
	}
}

func (h *Hive) DistanceTo(point image.Point) uint {
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
func (h *Hive) GetNearbyPosition(sim *T, unitSize int) *image.Point {
	const maxRadius = 3
	const tileSize = 128
	center := h.GetCenteredPosition()

	type spawnCandidate struct {
		point   image.Point
		rect    *image.Rectangle
		density int
	}

	var candidates []spawnCandidate

	for dx := -maxRadius; dx <= maxRadius; dx++ {
		for dy := -maxRadius; dy <= maxRadius; dy++ {
			x := center.X + dx*tileSize
			y := center.Y + dy*tileSize

			rect := &image.Rectangle{
				Min: image.Point{X: x - unitSize/2, Y: y - unitSize/2},
				Max: image.Point{X: x + unitSize/2, Y: y + unitSize/2},
			}

			if rect.Overlaps(*h.Rect) {
				continue
			}

			// Score this tile by number of units overlapping or nearby
			density := 0
			for _, unit := range append(sim.playerUnits, sim.enemyUnits...) {
				if unit == nil || unit.ID.String() == h.ID.String() {
					continue
				}
				// Check proximity (not just overlap)
				unitCenter := unit.GetCenteredPosition()
				dist := math.Hypot(float64(unitCenter.X-x), float64(unitCenter.Y-y))
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
			return &c.point
		}
	}

	// Fallback position
	return &image.Point{X: center.X + tileSize, Y: center.Y}
}
