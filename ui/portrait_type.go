package ui

type PortraitType int

const (
	PortraitTypeRoyalAnt PortraitType = iota
	PortraitTypeRoyalRoach
)

var portraitTypeToString = map[PortraitType]string{
	PortraitTypeRoyalAnt:   "portraits/ant-royalty.png",
	PortraitTypeRoyalRoach: "portraits/cockroach-royalty.png",
}

func (pt PortraitType) String() string {
	if s, ok := portraitTypeToString[pt]; ok {
		return s
	}
	return "unknown"
}
