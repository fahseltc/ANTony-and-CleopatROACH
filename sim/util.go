package sim

import "gamejam/types"

func UtilBuildingTypeToBuilding(bt types.Building) BuildingInterface {
	return GetBuildingInstance(bt, uint(PlayerFaction))
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
