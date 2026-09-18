package sim

import (
	"gamejam/types"
	"math"
)

var (
	// UnitHarvestDistance is the edge-to-resource-center distance within which a
	// worker may harvest. A worker stands on a tile ADJACENT to the resource, so
	// when correctly placed its nearest rect edge is only ~half a tile (~64px)
	// from the resource center; a diagonal neighbour and the movement arrival
	// slack (ArrivalThreshold) push that a little further. 140px comfortably
	// covers a worker standing on any adjacent tile without letting one harvest
	// from ~two tiles away, which is what made ants look far off yet still
	// gather. Keep this a bit above one tile (128px) so a worker settled just
	// past its slot doesn't ping-pong between Moving and Harvesting.
	UnitHarvestDistance   = uint(140)
	UnitHarvestFrameCount = uint(30)
)

// MaxHarvestApproachAttempts caps how many times a worker will re-walk toward
// the resource before giving up and harvesting from wherever it ended up. This
// lets a genuinely boxed-in worker (can't physically get adjacent) still gather
// instead of bouncing forever, without letting an ordinary worker harvest from
// far away after a single approach.
var MaxHarvestApproachAttempts = 3

type HarvestingState struct {
	harvestTimer uint
	// approachAttempts counts how many times we've walked toward the resource
	// this session. Unlike the old single "approached" commit, we keep trying to
	// get within range across a few attempts, so a worker that stopped short
	// (movement arrival slack, crowding) closes the gap instead of harvesting
	// from a tile or two away. Only after MaxHarvestApproachAttempts do we let
	// it harvest out of range (it's genuinely blocked).
	approachAttempts int
}

func (s *HarvestingState) Enter(unit *Unit) {
	s.harvestTimer = 0
	unit.Destinations.Clear()
}
func (s *HarvestingState) Update(unit *Unit, sim *T) {
	// Distance is measured to the resource center itself so any worker standing
	// within harvest range of the node can gather, regardless of which adjacent
	// tile it settled on. This lets many workers harvest the same resource
	// concurrently instead of fighting for one exact pixel.
	dist := unit.EdgeDistanceTo(unit.LastResourcePos)
	// Walk toward the resource while we're out of harvest range, retrying up to
	// MaxHarvestApproachAttempts times. Requiring in-range to harvest (rather
	// than committing after a single approach) stops workers gathering from a
	// tile or two away when they stopped short of their approach slot; the retry
	// cap prevents an endless Moving<->Harvesting bounce for a boxed-in worker.
	if dist > UnitHarvestDistance && s.approachAttempts < MaxHarvestApproachAttempts {
		s.approachAttempts++
		unit.Destinations.Clear()
		unit.Destinations.Enqueue(unit.HarvestApproachPos(sim, unit.LastResourcePos))
		// Re-enter THIS same state (not a fresh one) after moving, so the
		// attempt counter persists and we don't reset progress each approach.
		unit.ChangeState(&MovingState{NextState: s})
		return
	}
	s.harvestTimer += 1
	var finalResourceCollectionTime uint
	if sim.GetPlayerState().TechTree.UnlockedTech[TechFasterGathering] {
		finalResourceCollectionTime = uint(float64(UnitHarvestFrameCount) * 0.8) // 20% faster
	} else {
		finalResourceCollectionTime = uint(UnitHarvestFrameCount)
	}

	if s.harvestTimer >= finalResourceCollectionTime {
		s.harvestTimer = 0
		unit.Stats.ResourceTypeCarried = s.determineResourceHarvested(unit, sim)
		unit.Stats.ResourcesCarried = 5

		// Set destination to nearest home
		nearestHive := s.determineNearestHive(unit, sim)
		if nearestHive != nil {
			unit.NearestHome = nearestHive
			path := sim.FindClickedPath(unit.GetTileCoordinates(), nearestHive.GetTilePosition())
			for _, p := range path {
				unit.Destinations.Enqueue(p.ToCenteredPixelCoordinatesDouble())
			}

			// 		//sim.FindNearestSurroundingWalkableTiles(unit.Position, unit.NearestHome.GetTilePosition())

			// 		path := sim.FindClickedPath(unit.GetTileCoordinates(), unit.NearestHome.GetTilePosition())
			// 		for _, p := range path {
			// 			unit.Destinations.Enqueue(p.ToCenteredPixelCoordinatesDouble())
			// 		}

			unit.ChangeState(&DeliveringState{})
			return
		} else {
			unit.ChangeState(&IdleState{})
		}
	}

}
func (s *HarvestingState) Exit(unit *Unit) {}
func (s *HarvestingState) GetName() string { return UnitStateHarvesting.ToString() }

// internal helper methods for harvesting
func (s *HarvestingState) determineResourceHarvested(unit *Unit, sim *T) types.Resource {
	tile := sim.world.TileMap.GetTileByPosition(
		int(unit.LastResourcePos.X),
		int(unit.LastResourcePos.Y),
	)
	return tile.Type.ToResourceType()
}

func (s *HarvestingState) determineNearestHive(unit *Unit, sim *T) BuildingInterface {
	var nearest BuildingInterface
	minDist := uint(math.MaxUint32)
	for _, building := range sim.GetAllBuildings() {
		if building.GetFaction() == unit.Faction && (building.GetType() == types.BuildingTypeAntHive || building.GetType() == types.BuildingTypeRoachHive) {
			dist := unit.DistanceTo(building.GetCenteredPosition())
			if nearest == nil || dist < minDist {
				nearest = building
				minDist = dist
			}
		}
	}
	return nearest
}

// case CollectingAction:
// 	// if we are holding some resources, set home, then set DeliveringAction
// 	if unit.Stats.ResourcesCarried >= unit.Stats.MaxCarryCapacity {
// 		// Find the nearest hive and set it as the unit's TileTypePlain
// 		var nearest BuildingInterface
// 		minDist := uint(math.MaxUint32)
// 		for _, hive := range sim.GetAllBuildings() {
// 			if hive.GetFaction() == unit.Faction {
// 				dist := unit.DistanceTo(hive.GetCenteredPosition())
// 				if nearest == nil || dist < minDist {
// 					nearest = hive
// 					minDist = dist
// 				}
// 			}
// 		}
// 		unit.NearestHome = nearest
// 		unit.Stats.ResourceCollectTime = 0

// 		//sim.FindNearestSurroundingWalkableTiles(unit.Position, unit.NearestHome.GetTilePosition())

// 		path := sim.FindClickedPath(unit.GetTileCoordinates(), unit.NearestHome.GetTilePosition())
// 		for _, p := range path {
// 			unit.Destinations.Enqueue(p.ToCenteredPixelCoordinatesDouble())
// 		}
// 		unit.Action = DeliveringAction
// 	} else {
// 		// move to and collect resource
// 		unit.MoveToDestination(sim)
// 		dist := unit.DistanceTo(unit.LastResourcePos)
// 		if unit.Destinations.IsEmpty() && int(dist) > ResourceCollectionDistance {
// 			// move a little closer
// 			unit.Destinations.Enqueue(unit.LastResourcePos)
// 		}

// 		if int(dist) < ResourceCollectionDistance { // lots of tweaks needed here or fixes TODO
// 			// TODO: play animation
// 			unit.Destinations.Clear()
// 			unit.Stats.ResourceCollectTime += 1
// 			var finalResourceCollectionTime uint
// 			if sim.GetPlayerState().TechTree.UnlockedTech[TechFasterGathering] {
// 				finalResourceCollectionTime = uint(float64(MaxResourceCollectFrames) * 0.8) // 20% faster
// 			} else {
// 				finalResourceCollectionTime = uint(MaxResourceCollectFrames)
// 			}
// 			if unit.Stats.ResourceCollectTime >= finalResourceCollectionTime {
// 				unit.Stats.ResourceCollectTime = 0
// 				tile := sim.world.TileMap.GetTileByPosition(int(unit.LastResourcePos.X), int(unit.LastResourcePos.Y))
// 				if tile != nil && tile.Type != types.TileTypePlain {
// 					unit.Stats.ResourcesCarried = unit.Stats.MaxCarryCapacity
// 					unit.Stats.ResourceTypeCarried = tile.Type.ToResourceType()
// 				}
// 			}
// 		}
// 	}
