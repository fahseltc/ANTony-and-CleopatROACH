package sim

import (
	"gamejam/util"
	"image"
	"math"
	"math/rand/v2"

	"github.com/google/uuid"
)

var TileDimensions = 128
var UnitConstructionTime = 120

type Hive struct {
	ID       uuid.UUID
	Position *image.Point
	Rect     *image.Rectangle
	Faction  uint

	buildQueue *util.Queue[*Unit]

	UnitContructing      bool
	constructionProgress uint
}

func NewHive(x, y int) *Hive {
	hive := &Hive{
		ID:       uuid.New(),
		Position: &image.Point{0, 0},
		Rect: &image.Rectangle{
			Min: image.Point{0, 0},
			Max: image.Point{TileDimensions * 2, TileDimensions * 2},
		},
		Faction:         1,
		UnitContructing: false,
		buildQueue:      util.NewQueue[*Unit](),
	}
	hive.SetTilePosition(x, y)
	return hive
}

func (h *Hive) SetTilePosition(x, y int) {
	h.Position = &image.Point{X: x * TileDimensions, Y: y * TileDimensions}
	h.Rect.Min = *h.Position
	h.Rect.Max = image.Point{X: h.Position.X + TileDimensions*2, Y: h.Position.Y + TileDimensions*2}
}

func (h *Hive) Update(sim *T) {
	if !h.buildQueue.IsEmpty() {
		h.UnitContructing = true
		h.constructionProgress += 1
		if h.constructionProgress >= uint(UnitConstructionTime) {
			u, err := h.buildQueue.Dequeue()
			if err != nil {
				return // todo handle?
			}
			// make sure position isnt colliding with anything and try again
			u.SetPosition(h.GetNearbyPosition(sim))
			sim.AddUnit(u)
			h.UnitContructing = false
			h.constructionProgress = 0
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
