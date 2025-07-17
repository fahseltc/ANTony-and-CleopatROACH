package sim

import "gamejam/types"

type BarracksBuilding struct {
	*Building
}

func NewBarracksBuilding(x, y int) BuildingInterface {
	building := NewBuilding(x, y, TileDimensions, TileDimensions, 1, types.BuildingTypeBarracks, 0)

	bb := &BridgeBuilding{
		Building: building,
	}
	bb.Stats.Cost = ResourceCost{
		Wood:    100,
		Sucrose: 100,
	}

	return bb
}
