package types

type Destination int

const (
	DestinationTypeTile Destination = iota
	DestinationTypeResource
	DestinationTypeEnemy
	// building?
)
