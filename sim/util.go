package sim

import "gamejam/types"

func UtilBuildingTypeToBuilding(bt types.Building) BuildingInterface {
	switch bt {
	case types.BuildingTypeBarracks:
		return NewBarracksBuilding(0, 0)
	case types.BuildingTypeBridge:
		return NewBridgeBuilding(0, 0)
	default:
		return NewHive()
	}

}
