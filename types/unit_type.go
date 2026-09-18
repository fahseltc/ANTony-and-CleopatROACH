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

// DisplayName returns a player-facing name for the unit, suitable for showing
// in UI notifications (unlike ToString, which returns the internal enum name).
func (ut Unit) DisplayName() string {
	switch ut {
	case UnitTypeDefaultAnt:
		return "Worker Ant"
	case UnitTypeRoyalAnt:
		return "Royal Ant"
	case UnitTypeFighterAnt:
		return "Fighter Ant"
	case UnitTypeDefaultRoach:
		return "Worker Roach"
	case UnitTypeRoyalRoach:
		return "Royal Roach"
	default:
		return "Worker Ant"
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
