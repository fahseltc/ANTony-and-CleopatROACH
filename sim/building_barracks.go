package sim

import "gamejam/types"

type BarracksBuilding struct {
	*Building
}

func NewBarracksBuilding(x, y int) BuildingInterface {
	building := NewBuilding(x, y, TileDimensions, TileDimensions, 1, types.BuildingTypeBridge, 0)

	bb := &BridgeBuilding{
		Building: building,
	}

	return bb
}
