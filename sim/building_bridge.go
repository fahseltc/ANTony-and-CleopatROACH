package sim

type BridgeBuilding struct {
	*Building
}

// func NewBridgeBuilding(x, y int) BuildingInterface {
// 	building := NewBuilding(x, y, TileDimensions, TileDimensions, uint(PlayerFaction), types.BuildingTypeBridge)

// 	bb := &BridgeBuilding{
// 		Building: building,
// 	}
// 	bb.Stats.ResourceCost = ResourceCost{
// 		Wood: 50,
// 	}
// 	bb.Stats.ProgressMax = 90

// 	return bb
// }
