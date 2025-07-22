package sim

import (
	"gamejam/vec2"
	"math"
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

	enemyPos := unit.NearestEnemy.GetCenteredPosition()
	selfPos := unit.GetCenteredPosition()

	// --- Rotate toward target ---
	dx := enemyPos.X - selfPos.X
	dy := enemyPos.Y - selfPos.Y
	desiredAngle := math.Atan2(dy, dx) + math.Pi/2
	unit.RotateToward(desiredAngle, 0.1) // max 0.1 rad/frame

	// --- Only attack if in range and facing ---
	inRange := unit.TargetInAttackRange(enemyPos)
	inCone := unit.IsTargetInFrontalCone(unit.NearestEnemy, math.Pi/3)

	if inRange && inCone {
		unit.Stats.AttackFramesCur += 1
		if unit.Stats.AttackFramesCur >= unit.Stats.AttackFrames {
			if unit.Stats.Damage >= unit.NearestEnemy.Stats.HPCur { // avoid UNIT underflow
				unit.NearestEnemy.Stats.HPCur = 0
			} else {
				unit.NearestEnemy.Stats.HPCur -= unit.Stats.Damage
			}

			unit.Stats.AttackFramesCur = 0
		}
	} else {
		// Optional: Reset attack timer if out of arc or range
		unit.Stats.AttackFramesCur = 0
	}
}

func (unit *Unit) IsTargetInFrontalCone(target *Unit, coneAngle float64) bool {
	toTarget := target.GetCenteredPosition().Sub(*unit.GetCenteredPosition()).Normalize()
	unitFacing := vec2.T{
		X: math.Cos(unit.MovingAngle - math.Pi/2),
		Y: math.Sin(unit.MovingAngle - math.Pi/2),
	}
	dot := toTarget.Dot(unitFacing)
	angle := math.Acos(dot)
	return angle <= coneAngle/2
}

func (s *AttackingState) Exit(unit *Unit) {}
func (s *AttackingState) GetName() string { return UnitStateAttacking.ToString() }

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
