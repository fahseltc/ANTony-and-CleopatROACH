package sim

import (
	"gamejam/types"
	"gamejam/util"
	"gamejam/vec2"
	"math"
)

type Hive struct {
	*Building
	buildQueue *util.Queue[*QueuedItem]
}

func NewHive() BuildingInterface {
	building := NewBuilding(0, 0, TileDimensions*2, TileDimensions*2, 0, types.BuildingTypeHive)
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
			queuedItem.OnComplete(sim, h) // returns bool - we could check it if neeed

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
	// 	switch h.Type {
	// case types.BuildingTypeHive:
	// 	h.buildQueue.Enqueue(item)
	// case types.BuildingTypeRoachHive:
	// 	switch unitType {
	// 	case types.UnitTypeDefaultAnt:
	// 		unit = NewDefaultRoach()
	// 	// case UnitTypeFighterAnt:
	// 	// 	unit = NewFighterAnt()
	// 	default:
	// 		unit = NewDefaultRoach()
	// 	}
	// }
	// h.buildQueue.Enqueue(unit)
}
