package types

type Tile uint

const (
	TileTypePlain Tile = iota
	TileTypeSucrose
	TileTypeWood
)

func (t Tile) ToResourceType() Resource {
	switch t {
	case TileTypePlain:
		return ResourceTypeNone
	case TileTypeSucrose:
		return ResourceTypeSucrose
	case TileTypeWood:
		return ResourceTypeWood
	default:
		return ResourceTypeNone
	}
}
