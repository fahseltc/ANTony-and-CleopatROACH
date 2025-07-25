package sim

var BridgeBuildTime = 160

type InConstructionBuilding struct {
	*Building
	targetBuilding BuildingType
}

func NewInConstructionBuilding(x, y int, targetBuilding BuildingType) *InConstructionBuilding {
	building := NewBuilding(x, y, TileDimensions, TileDimensions, 1, BuildingTypeInConstruction, uint(BridgeBuildTime))

	icb := &InConstructionBuilding{
		Building:       building,
		targetBuilding: targetBuilding,
	}

	return icb
}

func (icb *InConstructionBuilding) Update(sim *T) {
	// check if there are ants around?
	icb.ProgressCurrent += 1
	if icb.ProgressCurrent <= icb.ProgressMax {
		return
	}
	// else create the new building
	icb.ProgressCurrent = 0
	sim.RemoveBuilding(icb)
	sim.world.TileMap.RemoveCollisionRect(icb.Rect)
	switch icb.targetBuilding {
	case BuildingTypeInConstruction: // shouldnt happen
	case BuildingTypeHive:
	case BuildingTypeBridge:
		bb := NewBridgeBuilding(icb.Position.X, icb.Position.Y)
		sim.AddBuilding(bb)

	}

}
