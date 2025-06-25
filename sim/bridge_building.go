package sim

type BridgeBuilding struct {
	*Building
}

func NewBridgeBuilding(x, y int) BuildingInterface {
	building := NewBuilding(x, y, TileDimensions, TileDimensions, 1, BuildingTypeBridge, 0)

	bb := &BridgeBuilding{
		Building: building,
	}

	return bb
}
