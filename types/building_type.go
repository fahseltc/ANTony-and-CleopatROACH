package types

type Building int

const (
	BuildingTypeNone Building = iota
	BuildingTypeInConstruction
	BuildingTypeHive
	BuildingTypeBarracks
	BuildingTypeRoachHive
	BuildingTypeBridge
)
