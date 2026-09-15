package sim

import (
	"gamejam/types"
	"gamejam/vec2"
)

var (
	UnitResourceDeliverRange = uint(195)
)

type DeliveringState struct{}

func (s *DeliveringState) Enter(unit *Unit) {
	// Assume destinations already set by previous state
}

func (s *DeliveringState) Update(unit *Unit, sim *T) {
	unit.MoveToDestination(sim)

	homeDist := unit.EdgeDistanceTo(unit.NearestHome.GetCenteredPosition())

	if homeDist <= UnitResourceDeliverRange {
		sim.AddResource(uint(unit.Stats.ResourcesCarried), unit.Stats.ResourceTypeCarried)
		unit.Stats.ResourcesCarried = 0
		unit.Stats.ResourceTypeCarried = types.ResourceTypeNone

		// Go back to the resource location
		resPos := &vec2.T{
			X: (unit.LastResourcePos.X - HalfTileSize) / TileSize,
			Y: (unit.LastResourcePos.Y - HalfTileSize) / TileSize,
		}
		path := sim.FindClickedPath(unit.GetTileCoordinates(), resPos)
		for _, p := range path {
			unit.Destinations.Enqueue(p.ToCenteredPixelCoordinates())
		}
		// Replace the final waypoint (the shared resource tile center) with this
		// worker's own approach slot so returning workers fan out around the node
		// instead of converging on one pixel.
		if n := len(unit.Destinations.Items); n > 0 {
			unit.Destinations.Items = unit.Destinations.Items[:n-1]
		}
		unit.Destinations.Enqueue(unit.HarvestApproachPos(unit.LastResourcePos))
		unit.ChangeState(&HarvestingState{})
		return
	}

	if homeDist > UnitResourceDeliverRange && unit.Destinations.IsEmpty() {
		// retry moving to nearest
		unit.Destinations.Clear()
		walkableTileCoords := sim.FindNearestSurroundingWalkableTiles(unit.GetTileCoordinates(), unit.NearestHome.GetTilePosition())
		unit.Destinations.Enqueue(walkableTileCoords.ToCenteredPixelCoordinates())
		unit.ChangeState(&MovingState{NextState: &DeliveringState{}})
		return
	}

}

func (s *DeliveringState) Exit(unit *Unit) {}
func (s *DeliveringState) GetName() string { return UnitStateDelivering.ToString() }
