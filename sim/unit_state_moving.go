package sim

import (
	"gamejam/vec2"
	"image"
	"math"
	"math/rand/v2"
)

var UnitRepulsionWeight = 2.0
var UnitMaxStuckFrames = 120

type MovingState struct {
	NextState UnitStateInterface
}

func (s *MovingState) Enter(unit *Unit) {}
func (s *MovingState) Update(unit *Unit, sim *T) {
	unit.MoveToDestination(sim)
	if unit.Destinations.IsEmpty() {
		if s.NextState != nil {
			unit.ChangeState(s.NextState)
		} else {
			unit.ChangeState(&IdleState{})
		}
	}

}
func (s *MovingState) Exit(unit *Unit) {}
func (s *MovingState) GetName() string { return UnitStateMoving.ToString() }

func (unit *Unit) MoveToDestination(sim *T) {
	dest, err := unit.Destinations.Peek()
	if err != nil || dest == nil {
		return
	}
	speed := float64(unit.Stats.MoveSpeed)
	oldPos := unit.GetCenteredPosition()
	oldX := oldPos.X
	oldY := oldPos.Y

	dx := float64(dest.X - unit.GetCenteredPosition().X)
	dy := float64(dest.Y - unit.GetCenteredPosition().Y)
	// Repulsion avoidance
	repulsion := unit.computeRepulsion(sim)

	// Blend desired direction and repulsion
	toTarget := vec2.T{X: dx, Y: dy}
	if toTarget.Length() > 0 {
		toTarget = toTarget.Normalize()
	}

	arrived := unit.EdgeDistanceTo(dest) <= uint(ArrivalThreshold)
	// Reduce or ignore repulsion when very close.
	if arrived {
		if unit.IsWorker() {
			// Keep a gentle push for workers so a crowd converging on a resource
			// spreads into a ring instead of stacking on the same point.
			gentle := repulsion.Scale(0.5)
			repulsion = &gentle
		} else {
			repulsion = &vec2.T{} // disable repulsion
		}
	}

	moveVec := toTarget.Add(repulsion.Scale(UnitRepulsionWeight)).Normalize().Scale(speed)

	moveX := math.Copysign(math.Min(math.Abs(moveVec.X), speed), moveVec.X)
	moveY := math.Copysign(math.Min(math.Abs(moveVec.Y), speed), moveVec.Y)

	// Attempt X movement
	if moveX != 0 {
		newX := unit.Position.X + moveX
		newY := unit.Position.Y
		candidate := &image.Rectangle{
			Min: image.Point{X: int(newX), Y: int(newY)},
			Max: image.Point{X: int(newX) + unit.Rect.Dx(), Y: int(newY) + unit.Rect.Dy()},
		}
		if !unit.isColliding(candidate, sim) {
			unit.SetPosition(&vec2.T{X: newX, Y: unit.Position.Y})
		}
	}

	// Attempt Y movement
	if moveY != 0 {
		newY := unit.Position.Y + moveY
		newX := unit.Position.X
		candidate := &image.Rectangle{
			Min: image.Point{X: int(newX), Y: int(newY)},
			Max: image.Point{X: int(newX) + unit.Rect.Dx(), Y: int(newY) + unit.Rect.Dy()},
		}
		if !unit.isColliding(candidate, sim) {
			unit.SetPosition(&vec2.T{X: unit.Position.X, Y: newY})
		}
	}

	// Handle Rotation
	newCentered := unit.GetCenteredPosition()
	dxRot := float64(newCentered.X - oldX)
	dyRot := float64(newCentered.Y - oldY)
	if dxRot != 0 || dyRot != 0 { // update angle only if moved
		unit.StuckFrames = 0
		desiredAngle := math.Atan2(dyRot, dxRot) + math.Pi/2
		unit.RotateToward(desiredAngle, 1) // radians per frame
	}
	arrived = unit.EdgeDistanceTo(dest) <= uint(ArrivalThreshold)

	const stuckEpsilon = 1.5
	moved := math.Abs(dxRot) > stuckEpsilon || math.Abs(dyRot) > stuckEpsilon

	if !moved && !arrived {
		unit.StuckFrames++

		if unit.StuckFrames%30 == 0 {
			unit.NavigateAround(sim)
			//unit.TrySidestep(sim)
		}

		if unit.StuckFrames > UnitMaxStuckFrames {
			unit.StuckAttempts++
			if unit.StuckAttempts >= 3 {
				unit.StuckFrames = 0
				unit.StuckAttempts = 0

				// A worker heading TO a resource (carrying nothing) should never
				// silently abandon the job just because the approach was crowded.
				// If it's already close enough to the node, let it harvest right
				// here; otherwise re-plan a harvest approach to its own slot.
				//
				// A worker that IS carrying resources is on a delivery run, so we
				// leave its delivery destinations intact and just try to route
				// around the blockage below rather than redirecting it to harvest.
				harvestingWorker := unit.IsWorker() &&
					unit.LastResourcePos != nil &&
					unit.Stats.ResourcesCarried == 0
				if harvestingWorker {
					unit.Destinations.Clear()
					if unit.EdgeDistanceTo(unit.LastResourcePos) <= UnitHarvestDistance {
						unit.ChangeState(&HarvestingState{})
					} else {
						unit.Destinations.Enqueue(unit.HarvestApproachPos(unit.LastResourcePos))
						unit.ChangeState(&MovingState{NextState: &HarvestingState{}})
					}
					return
				}

				// Carrying workers: keep trying to reach home instead of dropping
				// to Idle (which would strand the resources). Non-workers give up.
				if unit.IsWorker() && unit.Stats.ResourcesCarried > 0 {
					unit.NavigateAround(sim)
					return
				}

				unit.ChangeState(&IdleState{})
				unit.Destinations.Clear()
				return
			}
			unit.NavigateAround(sim)
			unit.StuckFrames = 0
		}

	}
	if arrived && len(unit.Destinations.Items) >= 1 {
		unit.Destinations.Dequeue()
	}
}

func (unit *Unit) RotateToward(targetAngle float64, maxDelta float64) {
	diff := targetAngle - unit.MovingAngle

	// Normalize angle to [-π, π]
	for diff > math.Pi {
		diff -= 2 * math.Pi
	}
	for diff < -math.Pi {
		diff += 2 * math.Pi
	}

	if math.Abs(diff) < maxDelta {
		unit.MovingAngle = targetAngle
	} else if diff > 0 {
		unit.MovingAngle += maxDelta
	} else {
		unit.MovingAngle -= maxDelta
	}
}

func (unit *Unit) isColliding(rect *image.Rectangle, sim *T) bool {
	futureUnitCenterX := float64(rect.Min.X+rect.Max.X) / 2
	futureUnitCenterY := float64(rect.Min.Y+rect.Max.Y) / 2
	futureUnitRadius := float64(rect.Dx()) / 2

	colliders := sim.GetAllCollidersOverlapping(rect)
	for _, collider := range colliders {
		if collider.OwnerID == unit.ID.String() {
			continue // skip self
		}

		collidingUnit, _ := sim.GetUnitByID(collider.OwnerID)

		// Skip unit-unit collision for workers that are not idle
		if unit.IsWorker() && unit.CurrentState.GetName() != UnitStateIdle.ToString() && collidingUnit != nil && collidingUnit.IsWorker() {
			continue
		}
		if collider.Radius > 0 && collider.Center != (image.Point{}) {
			dx := futureUnitCenterX - float64(collider.Center.X)
			dy := futureUnitCenterY - float64(collider.Center.Y)
			distance := math.Sqrt(dx*dx + dy*dy)
			if distance < futureUnitRadius+float64(collider.Radius) {
				return true
			}
		} else if collider.Rect.Overlaps(*rect) {
			return true
		}
	}
	for _, mo := range sim.world.MapObjects {
		closestX := math.Max(float64(mo.Rect.Min.X), math.Min(futureUnitCenterX, float64(mo.Rect.Max.X)))
		closestY := math.Max(float64(mo.Rect.Min.Y), math.Min(futureUnitCenterY, float64(mo.Rect.Max.Y)))
		dx := futureUnitCenterX - closestX
		dy := futureUnitCenterY - closestY
		if dx*dx+dy*dy <= futureUnitRadius*futureUnitRadius {
			return true
		}
	}
	return false
}

func (unit *Unit) TrySidestep(sim *T) {
	bestOffset := vec2.T{}

	// Try 8 directions (N, NE, E, SE, S, SW, W, NW)
	offsets := []vec2.T{
		{X: -1, Y: 0}, {X: 1, Y: 0},
		{X: 0, Y: -1}, {X: 0, Y: 1},
		{X: -1, Y: -1}, {X: 1, Y: -1},
		{X: -1, Y: 1}, {X: 1, Y: 1},
	}

	// Shuffle offsets to avoid always biasing same direction
	rand.Shuffle(len(offsets), func(i, j int) {
		offsets[i], offsets[j] = offsets[j], offsets[i]
	})

	for _, off := range offsets {
		newX := unit.Position.X + off.X*float64(unit.Stats.MoveSpeed)
		newY := unit.Position.Y + off.Y*float64(unit.Stats.MoveSpeed)
		candidate := &image.Rectangle{
			Min: image.Point{X: int(newX), Y: int(newY)},
			Max: image.Point{X: int(newX) + unit.Rect.Dx(), Y: int(newY) + unit.Rect.Dy()},
		}
		if !unit.isColliding(candidate, sim) {
			bestOffset = off
		}
	}

	// Apply best offset if found
	if bestOffset != (vec2.T{}) {
		newX := unit.Position.X + bestOffset.X*float64(unit.Stats.MoveSpeed)
		newY := unit.Position.Y + bestOffset.Y*float64(unit.Stats.MoveSpeed)
		unit.SetPosition(&vec2.T{X: newX, Y: newY})
	}
}

func (unit *Unit) NavigateAround(sim *T) {
	// First, try to recompute a real path from the unit's current tile to the
	// tile of its FINAL destination. The grid pathfinder respects collision
	// MapObjects, so this routes the unit around the obstacle it's wedged
	// against instead of blindly stepping backward into it again.
	if finalDest, err := unit.Destinations.PeekBack(); err == nil && finalDest != nil {
		start := unit.GetTileCoordinates()
		endTile := &vec2.T{
			X: math.Floor(finalDest.X / TileSize),
			Y: math.Floor(finalDest.Y / TileSize),
		}
		// Don't bother re-pathing to the tile we're already on.
		if int(start.X) != int(endTile.X) || int(start.Y) != int(endTile.Y) {
			path := sim.FindClickedPath(start, endTile)
			if len(path) > 0 {
				// Preserve the exact final approach waypoint (e.g. a harvest
				// slot or precise drop-off pixel) which the caller enqueued and
				// which may sit off tile-center.
				unit.Destinations.Clear()
				for _, p := range path {
					unit.Destinations.Enqueue(&vec2.T{
						X: p.X*TileSize + HalfTileSize,
						Y: p.Y*TileSize + HalfTileSize,
					})
				}
				unit.Destinations.Enqueue(finalDest)
				return
			}
		}
	}

	// Fallback (no viable re-path found): step one tile backward from current
	// facing to try to break a local deadlock, then let normal movement retry.
	angle := unit.MovingAngle - math.Pi/2 // undo +π/2 used earlier

	backwards := vec2.T{
		X: -math.Cos(angle),
		Y: -math.Sin(angle),
	}.Normalize()

	tileSize := 128.0
	backDest := unit.GetCenteredPosition().Add(backwards.Scale(tileSize))

	// Clamp to map boundaries if needed
	backDest.X = math.Max(0, math.Min(backDest.X, float64(sim.world.TileMap.Width*sim.world.TileMap.TileSize-unit.Rect.Dx())))
	backDest.Y = math.Max(0, math.Min(backDest.Y, float64(sim.world.TileMap.Height*sim.world.TileMap.TileSize-unit.Rect.Dy())))

	unit.Destinations.EnqueueFront(&vec2.T{
		X: backDest.X,
		Y: backDest.Y,
	})
}
func (unit *Unit) computeRepulsion(sim *T) *vec2.T {
	repulsion := vec2.T{}
	myCenter := unit.GetCenteredPosition()
	maxPush := 3.0 // cap max influence of any one unit

	for _, other := range sim.GetAllUnits() {
		if other.ID == unit.ID {
			continue
		}

		otherCenter := other.GetCenteredPosition()
		dir := myCenter.Sub(*otherCenter)
		dist := dir.Length()

		if dist < 160 && dist > 0.1 {
			strength := (160 - dist) / 160
			falloff := math.Pow(strength, 2) // aggressive falloff
			push := dir.Normalize().Scale(math.Min(falloff*5, maxPush))

			// Optional: Apply sideways deflection to break deadlocks
			// But reduce strength at very low distances to avoid chaos
			if falloff > 0.3 && math.Abs(push.X) < 0.2 && math.Abs(push.Y) > 0.5 {
				deflect := vec2.T{X: 1, Y: 0}
				if rand.IntN(2) == 0 {
					deflect.X = -1
				}
				push = push.Add(deflect.Scale(0.5)) // reduced from 1.5 to 0.5
			}

			repulsion = repulsion.Add(push)
		}
	}

	// Final clamp on total repulsion to avoid extreme jittering
	if repulsion.Length() > maxPush {
		repulsion = repulsion.Normalize().Scale(maxPush)
	}

	return &repulsion
}
