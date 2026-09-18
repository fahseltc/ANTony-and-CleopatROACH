package types

type Resource uint

const (
	ResourceTypeNone Resource = iota
	ResourceTypeSucrose
	ResourceTypeWood
)

func (rt Resource) ToString() string {
	switch rt {
	case ResourceTypeNone:
		return ""
	case ResourceTypeSucrose:
		return "sucrose"
	case ResourceTypeWood:
		return "wood"
	default:
		return ""
	}
}
