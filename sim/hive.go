package sim

import (
	"gamejam/types"
	"gamejam/util"
	"gamejam/vec2"
	"image"
	"math"
)

type Hive struct {
	*Building
	buildQueue *util.Queue[*QueuedItem]

	rallyPoint *image.Point
}

func NewAntHive() BuildingInterface {
	building := NewBuilding(0, 0, TileDimensions*2, TileDimensions*2, 0, types.BuildingTypeAntHive)
	h := &Hive{
		Building:   building,
		buildQueue: util.NewQueue[*QueuedItem](),
	}
	return h
}

func NewRoachHive() BuildingInterface {
	building := NewBuilding(0, 0, TileDimensions*2, TileDimensions*2, 0, types.BuildingTypeRoachHive)
	h := &Hive{
		Building:   building,
		buildQueue: util.NewQueue[*QueuedItem](),
	}
	return h
}

func (h *Hive) Update(sim *T) {
	if !h.buildQueue.IsEmpty() {
		queuedItem, err := h.buildQueue.Peek()
		if err != nil {
			return
		}
		if h.Stats.ProgressMax == 0 {
			h.Stats.ProgressMax = queuedItem.ConstructionTime
		}
		h.Stats.ProgressCurrent += 1
		if h.Stats.ProgressCurrent >= h.Stats.ProgressMax {
			queuedItem, err = h.buildQueue.Dequeue()
			if err != nil {
				return // todo handle?
			}
			newUnit := queuedItem.OnComplete(sim, h)
			if queuedItem.Type == types.QueuedItemTypeUnit && h.rallyPoint != nil {
				sim.issueGroupAction([]string{newUnit.ID.String()}, h.rallyPoint)
			}

			h.Stats.ProgressCurrent = 0
			h.Stats.ProgressMax = 0
		}
	}
}

func (h *Hive) DistanceTo(point vec2.T) uint {
	xDist := math.Abs(float64(h.Position.X - point.X))
	yDist := math.Abs(float64(h.Position.Y - point.Y))
	return uint(math.Sqrt(math.Pow(xDist, 2) + math.Pow(yDist, 2)))
}

func (h *Hive) AddItemToBuildQueue(item *QueuedItem) {
	h.buildQueue.Enqueue(item)
}

func (h *Hive) SetRallyPoint(rallypoint *image.Point) {
	h.rallyPoint = rallypoint
}
