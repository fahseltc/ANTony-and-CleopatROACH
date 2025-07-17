package sim

import "gamejam/types"

type BarracksBuilding struct {
	*Building
}

func NewBarracksBuilding(x, y int) BuildingInterface {
	building := NewBuilding(x, y, TileDimensions, TileDimensions, uint(PlayerFaction), types.BuildingTypeBarracks)

	bb := &BarracksBuilding{
		Building: building,
	}
	bb.Stats.Cost = ResourceCost{
		Wood:    100,
		Sucrose: 100,
	}

	return bb
}
