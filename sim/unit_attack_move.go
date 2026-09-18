package sim

type AttackMoveState struct{}

func (s *AttackMoveState) Enter(unit *Unit) {}
func (s *AttackMoveState) Update(unit *Unit, sim *T) {
	idleUnitAggroRange := unit.Stats.VisionRange * uint(TileDimensions)

	// Try to find a target in range
	nearestEnemy := unit.findNearestEnemy(sim)

	if nearestEnemy != nil && unit.EdgeDistanceTo(nearestEnemy.GetCenteredPosition()) <= idleUnitAggroRange {
		unit.NearestEnemy = nearestEnemy
		unit.Destinations.EnqueueFront(unit.NearestEnemy.Position)
		if unit.TargetInAttackRange(nearestEnemy.GetCenteredPosition()) { // if its close enough to attack, attack it
			unit.ChangeState(&AttackingState{})
			return
		}
	}

	// Continue moving
	unit.MoveToDestination(sim)
	if unit.Destinations.IsEmpty() {
		unit.ChangeState(&IdleState{})
	}
}
func (s *AttackMoveState) Exit(unit *Unit) {}
func (s *AttackMoveState) GetName() string {
	return UnitStateAttackMove.ToString()
}
