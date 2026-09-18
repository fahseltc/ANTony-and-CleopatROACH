package types

type QueuedItem int

const (
	QueuedItemTypeUnit QueuedItem = iota
	QueuedItemTypeTech
)
