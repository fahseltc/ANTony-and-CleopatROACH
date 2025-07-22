package sim

import (
	"gamejam/types"
	"math"
)

var (
	UnitHarvestDistance   = uint(110)
	UnitHarvestFrameCount = uint(30)
)

type HarvestingState struct {
	harvestTimer uint
}

func (s *HarvestingState) Enter(unit *Unit) {
	s.harvestTimer = 0
	unit.Destinations.Clear()
}
func (s *HarvestingState) Update(unit *Unit, sim *T) {
	dist := unit.EdgeDistanceTo(unit.LastResourcePos)
	if dist > UnitHarvestDistance {
		// Move back to the resource
		unit.Destinations.Clear()
		unit.Destinations.Enqueue(unit.LastResourcePos)
		unit.ChangeState(&MovingState{NextState: &HarvestingState{}})
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
