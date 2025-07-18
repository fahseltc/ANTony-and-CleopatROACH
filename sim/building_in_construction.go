package sim

import (
	"gamejam/eventing"
	"gamejam/types"
)

var BridgeBuildTime = 160

type InConstructionBuilding struct {
	*Building
	targetBuilding types.Building
}

func NewInConstructionBuilding(x, y int, targetBuilding types.Building) *InConstructionBuilding {
	building := NewBuilding(x, y, TileDimensions, TileDimensions, uint(PlayerFaction), types.BuildingTypeInConstruction)

	icb := &InConstructionBuilding{
		Building:       building,
		targetBuilding: targetBuilding,
	}
	b := UtilBuildingTypeToBuilding(targetBuilding)
	icb.Stats.ProgressMax = b.GetStats().ProgressMax

	return icb
}

func (icb *InConstructionBuilding) Update(sim *T) {
	// check if there are ants around?
	icb.Stats.ProgressCurrent += 1
	if icb.Stats.ProgressCurrent <= icb.Stats.ProgressMax {
		return
	}
	// else create the new building
	icb.Stats.ProgressCurrent = 0
	sim.RemoveBuilding(icb)
	sim.world.TileMap.RemoveCollisionRect(icb.Rect)

	switch icb.targetBuilding {
	case types.BuildingTypeInConstruction: // shouldnt ever happen
	case types.BuildingTypeHive:
	case types.BuildingTypeBarracks:
		bb := NewBarracksBuilding(int(icb.Position.X), int(icb.Position.Y))
		sim.AddBuilding(bb)
		state := sim.GetPlayerState()
		state.TechTree.Unlock(TechBuildFighterUnit, sim.GetPlayerState())
		sim.world.TileMap.AddCollisionRect(bb.GetRect())
		sim.EventBus.Publish(eventing.Event{
			Type: "NotificationEvent",
			Data: eventing.NotificationEvent{
				Message: "Fighter units unlocked!",
			},
		})
	case types.BuildingTypeBridge:
		bb := NewBridgeBuilding(int(icb.Position.X), int(icb.Position.Y))
		sim.AddBuilding(bb)
		//sim.world.TileMap.AddCollisionRect()
	}

}
