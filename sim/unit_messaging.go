package sim

type UnitMessage int

const (
	UnitMessageNone UnitMessage = iota
	UnitMessageArrivedIdle
	UnitMessageArrivedAttack
)

func (unit *Unit) SendMessage(sim *T, msg UnitMessage) {
	switch msg {
	case UnitMessageNone:
	case UnitMessageArrivedIdle:
		unit.Destinations.Clear()
		unit.Action = IdleAction
		nearbyUnits := sim.GetAllNearbyFriendlyUnits(unit)
		for _, nearbyUnit := range nearbyUnits {
			if nearbyUnit.Action != IdleAction {
				nearbyUnit.SendMessage(sim, UnitMessageArrivedIdle)
			}
		}
	case UnitMessageArrivedAttack:

	}
}
