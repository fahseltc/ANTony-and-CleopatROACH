package sim

type UnitStateInterface interface {
	Enter(unit *Unit)
	Update(unit *Unit, sim *T)
	Exit(unit *Unit)
	GetName() string
}

type UnitState uint

const (
	UnitStateAttacking UnitState = iota
	UnitStateIdle
	UnitStateMoving
	UnitStateAttackMove
	UnitStateHarvesting
	UnitStateDelivering
	UnitStateConstructing
)

func (us UnitState) ToString() string {
	switch us {
	case UnitStateAttacking:
		return "attacking"
	case UnitStateIdle:
		return "idle"
	case UnitStateMoving:
		return "moving"
	case UnitStateAttackMove:
		return "attackmove"
	case UnitStateHarvesting:
		return "harvesting"
	case UnitStateDelivering:
		return "delivering"
	case UnitStateConstructing:
		return "constructing"
	default:
		return "idle"
	}
}
