package sim

type Action int

const (
	Idle Action = iota
	Moving
	Collecting
	Delivering
	Building
	Attacking
)

type Unit struct {
	ID           uint32
	Health       int
	X, Y         float64
	MoveSpeed    float64
	Action       Action
	DestX, DestY float64
}
