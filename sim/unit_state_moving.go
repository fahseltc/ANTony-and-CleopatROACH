package sim

import (
	"gamejam/vec2"
	"image"
	"math"
	"math/rand/v2"
)

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
func (s *MovingState) Name() string    { return "moving" }

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
	repulsionWeight := 2.0
	moveVec := toTarget.Add(repulsion.Scale(repulsionWeight)).Normalize().Scale(speed)

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
	arrived := unit.EdgeDistanceTo(dest) <= uint(ArrivalThreshold)

	const stuckEpsilon = 1.5
	moved := math.Abs(dxRot) > stuckEpsilon || math.Abs(dyRot) > stuckEpsilon

	if !moved && !arrived {
		unit.StuckFrames++

		if unit.StuckFrames%30 == 0 {
			unit.NavigateAround(sim)
			//unit.TrySidestep(sim)
		}

		// if unit.StuckFrames > 2000 { //|| unit.StuckSidestepAttempts > 3
		// 	//Only sidestep if the destination itself isn't clearly blocked
		// 	if unit.isDestinationBlocked(sim) {
		// 		unit.Action = IdleAction
		// 		unit.StuckFrames = 0
		// 		return
		// 	}
		// }
	}
	if arrived && len(unit.Destinations.Items) >= 1 {
		unit.Destinations.Dequeue()
	}
	// if arrived && len(unit.Destinations.Items) == 0 &&
	// 	(unit.CurrentState.Name() != "collecting" ||
	// 		unit.CurrentState.Name() != "delivering") &&
	// 	unit.Stats.ResourcesCarried != 0 {
	// 	nearbyUnits := sim.GetAllNearbyFriendlyUnits(unit)
	// 	for _, nearbyUnit := range nearbyUnits {
	// 		nearbyUnit.SendMessage(sim, UnitMessageArrivedIdle)
	// 	}
	// }
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
		if unit.IsWorker() && unit.CurrentState.Name() != "idle" && collidingUnit != nil && collidingUnit.IsWorker() {
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
	// Get direction unit was moving last
	angle := unit.MovingAngle - math.Pi/2 // undo +π/2 used earlier

	// Calculate "backwards" vector (opposite of current direction)
	backwards := vec2.T{
		X: -math.Cos(angle),
		Y: -math.Sin(angle),
	}.Normalize()

	// Move one tile back — assuming tiles are 128x128
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
