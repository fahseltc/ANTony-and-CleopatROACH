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

func UtilUnitTypeToUnit(ut types.Unit) *Unit {
	switch ut {
	case types.UnitTypeDefaultAnt:
		return NewDefaultAnt()
	case types.UnitTypeRoyalAnt:
		return NewRoyalAnt()
	case types.UnitTypeFighterAnt:
		return NewFighterAnt()
	case types.UnitTypeDefaultRoach:
		return NewDefaultRoach()
	case types.UnitTypeRoyalRoach:
		return NewRoyalRoach()
	default:
		return NewDefaultAnt()
	}
}
