package sim

import (
	"gamejam/util"
	"image"
	"math"
	"math/rand/v2"
)

var UnitConstructionTime = 120

type Hive struct {
	*Building
	buildQueue      *util.Queue[*Unit]
	unitContructing bool
}

func NewHive() BuildingInterface {
	building := NewBuilding(0, 0, TileDimensions*2, TileDimensions*2, 0, BuildingTypeHive, uint(UnitConstructionTime))

	h := &Hive{
		Building:        building,
		unitContructing: false,
		buildQueue:      util.NewQueue[*Unit](),
	}
	return h
}

func (h *Hive) Update(sim *T) {
	if !h.buildQueue.IsEmpty() {
		h.unitContructing = true
		h.ProgressCurrent += 1
		if h.ProgressCurrent >= uint(UnitConstructionTime) {
			u, err := h.buildQueue.Dequeue()
			if err != nil {
				return // todo handle?
			}
			// make sure position isnt colliding with anything and try again
			u.SetPosition(h.GetNearbyPosition(sim))
			sim.AddUnit(u)
			h.unitContructing = false
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
	h.buildQueue.Enqueue(NewDefaultAnt())
}
func (h *Hive) GetNearbyPosition(sim *T) *image.Point {
	directions := []image.Point{
		{X: 300, Y: 0},
		{X: -150, Y: 0},
		{X: 0, Y: 300},
		{X: 0, Y: -150},
	}

	// Optional shuffle directions to spread spawn points randomly
	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})

	for _, dir := range directions {
		candidate := &image.Rectangle{
			Min: image.Point{X: h.Rect.Min.X + dir.X, Y: h.Rect.Min.Y + dir.Y},
			Max: image.Point{X: h.Rect.Max.X + dir.X, Y: h.Rect.Max.Y + dir.Y},
		}

		// skip if overlaps hive itself
		if candidate.Overlaps(*h.Rect) {
			continue
		}

		colliders := sim.GetAllCollidersOverlapping(candidate)
		collision := false
		for _, collider := range colliders {
			if collider.OwnerID == h.ID.String() {
				continue
			}
			if candidate.Overlaps(*collider.Rect) {
				collision = true
				break
			}
		}

		if !collision {
			return &candidate.Min
		}
	}

	// fallback spot, just offset further away
	return &image.Point{X: h.Rect.Min.X + 300, Y: h.Rect.Min.Y}
}
