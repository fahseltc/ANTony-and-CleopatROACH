package types

type Building int

const (
	BuildingTypeInConstruction Building = iota
	BuildingTypeHive
	BuildingTypeRoachHive
	BuildingTypeBridge
)
