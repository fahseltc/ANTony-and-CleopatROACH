package sim

type Action int

const (
	Idle Action = iota
	Moving
	AttackMoving
	Attacking
	HoldingPosition
	Collecting
	Delivering
	//Building
)

type Unit struct {
	ID           uint32
	Health       int
	X, Y         int
	MoveSpeed    float64
	Action       Action
	DestX, DestY int
}
