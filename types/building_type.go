package types

type Building int

const (
	BuildingTypeNone Building = iota
	BuildingTypeInConstruction
	BuildingTypeAntHive
	BuildingTypeRoachHive
	BuildingTypeBarracks

	BuildingTypeBridge
)

func (bt Building) ToString() string {
	switch bt {
	case BuildingTypeNone:
		return "none"
	case BuildingTypeInConstruction:
		return "in_construction"
	case BuildingTypeAntHive:
		return "ant_hive"
	case BuildingTypeRoachHive:
		return "roach_hive"
	case BuildingTypeBarracks:
		return "barracks"
	case BuildingTypeBridge:
		return "bridge"

	default:
		return "none"
	}
}

func UtilBuildingTypeFromString(buildingString string) Building {
	switch buildingString {
	case "none":
		return BuildingTypeNone
	case "in_construction":
		return BuildingTypeInConstruction
	case "ant_hive":
		return BuildingTypeAntHive
	case "roach_hive":
		return BuildingTypeRoachHive
	case "barracks":
		return BuildingTypeBarracks
	case "bridge":
		return BuildingTypeBridge

	default:
		return BuildingTypeNone
	}
}
