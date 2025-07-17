package types

type Unit int

const (
	UnitTypeDefaultAnt Unit = iota
	UnitTypeRoyalAnt
	UnitTypeFighterAnt
	UnitTypeDefaultRoach
	UnitTypeRoyalRoach
)

func (ut Unit) ToString() string {
	switch ut {
	case UnitTypeDefaultAnt:
		return "UnitTypeDefaultAnt"
	case UnitTypeRoyalAnt:
		return "UnitTypeRoyalAnt"
	case UnitTypeFighterAnt:
		return "UnitTypeFighterAnt"
	case UnitTypeDefaultRoach:
		return "UnitTypeDefaultRoach"
	case UnitTypeRoyalRoach:
		return "UnitTypeRoyalRoach"
	default:
		return "UnitTypeDefaultAnt"
	}
}

func UtilUnitTypeFromString(unitString string) Unit {
	switch unitString {
	case "worker":
		return UnitTypeDefaultAnt
	// case "UnitTypeRoyalAnt":
	// 	return UnitTypeRoyalAnt
	case "fighter":
		return UnitTypeFighterAnt
	case "UnitTypeDefaultRoach":
		return UnitTypeDefaultRoach
	// case "UnitTypeRoyalRoach":
	// 	return UnitTypeRoyalRoach
	default:
		return UnitTypeDefaultAnt
	}
}
