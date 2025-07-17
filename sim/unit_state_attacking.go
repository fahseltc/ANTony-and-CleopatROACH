package sim

var (
	UnitAttackRange = 90 // px
)

type AttackingState struct{}

func (s *AttackingState) Enter(unit *Unit) {
	unit.Stats.AttackFramesCur = 0
}

func (s *AttackingState) Update(unit *Unit, sim *T) {
	if unit.NearestEnemy == nil || !unit.NearestEnemy.IsAlive() {
		unit.NearestEnemy = nil
		unit.ChangeState(&IdleState{})
		return
	}

	if unit.EdgeDistanceTo(unit.NearestEnemy.GetCenteredPosition()) <= uint(UnitAttackRange) {
		unit.Stats.AttackFramesCur += 1
		if unit.Stats.AttackFramesCur >= unit.Stats.AttackFrames {
			unit.NearestEnemy.Stats.HPCur -= unit.Stats.Damage
			unit.Stats.AttackFramesCur = 0
		}
	}

}

func (s *AttackingState) Exit(unit *Unit) {}
func (s *AttackingState) Name() string    { return "attacking" }

// 	if unit.NearestEnemy != nil && unit.NearestEnemy.IsAlive() && unit.TargetInRange(unit.NearestEnemy.GetCenteredPosition()) {
// 		if unit.IsInAttackZone(unit.NearestEnemy) {
// 			unit.Stats.AttackFramesCur++
// 			if unit.Stats.AttackFramesCur >= unit.Stats.AttackFrames {
// 				unit.NearestEnemy.Stats.HPCur -= unit.Stats.Damage
// 				unit.Stats.AttackFramesCur = 0
// 			}
// 		} else {
// 			toEnemy := unit.NearestEnemy.GetCenteredPosition().Sub(*unit.GetCenteredPosition())
// 			newPos := unit.NearestEnemy.GetCenteredPosition().Sub(toEnemy.Normalize().Scale(30))
// 			unit.Destinations.EnqueueFront(&newPos)
// 		}
// 	} else if unit.NearestEnemy == nil || !unit.NearestEnemy.IsAlive() {
// 		unit.findTarget(sim)
// 		if unit.NearestEnemy != nil && unit.NearestEnemy.IsAlive() {
// 			unit.Destinations.Enqueue(unit.NearestEnemy.GetCenteredPosition())
// 		}
// 	}
