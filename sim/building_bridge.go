package sim

import "gamejam/types"

type BridgeBuilding struct {
	*Building
}

func NewBridgeBuilding(x, y int) BuildingInterface {
	building := NewBuilding(x, y, TileDimensions, TileDimensions, uint(PlayerFaction), types.BuildingTypeBridge)

	bb := &BridgeBuilding{
		Building: building,
	}
	bb.Stats.Cost = ResourceCost{
		Wood: 50,
	}
	bb.Stats.ProgressMax = 90

	return bb
}
