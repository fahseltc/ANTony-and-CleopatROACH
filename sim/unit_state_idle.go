package sim

type IdleState struct{}

func (s *IdleState) Enter(unit *Unit) {}

func (s *IdleState) Update(unit *Unit, sim *T) {
	idleUnitAggroRange := unit.Stats.VisionRange * uint(TileDimensions)
	nearestEnemy := unit.findNearestEnemy(sim)
	if nearestEnemy == nil {
		return
	}

	if unit.EdgeDistanceTo(nearestEnemy.GetCenteredPosition()) <= idleUnitAggroRange {
		unit.NearestEnemy = nearestEnemy
		if unit.TargetInAttackRange(nearestEnemy.GetCenteredPosition()) { // if its close enough to attack, attack it
			unit.ChangeState(&AttackingState{})
			return
		} else { // otherwise just move towards it
			unit.Destinations.Enqueue(unit.NearestEnemy.GetCenteredPosition()) // TODO: this might need pathfinding?
			unit.ChangeState(&MovingState{})
			return
		}
	}
}

func (s *IdleState) Exit(unit *Unit) {}
func (s *IdleState) Name() string    { return "idle" }
