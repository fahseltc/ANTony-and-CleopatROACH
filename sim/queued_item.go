package sim

import (
	"gamejam/eventing"
	"gamejam/types"
)

type QueuedItem struct {
	ConstructionTime uint
	Type             types.QueuedItem
	Unit             *Unit
	TechID           TechID
}

func (qi *QueuedItem) OnComplete(s *T, constructingBuilding BuildingInterface) *Unit {
	switch qi.Type {
	case types.QueuedItemTypeTech:
		unlocked := s.GetPlayerState().TechTree.Unlock(qi.TechID, s.GetPlayerState())
		if !unlocked {
			s.EventBus.Publish(eventing.Event{
				Type: "NotificationEvent",
				Data: eventing.NotificationEvent{
					Message: s.GetPlayerState().TechTree.GetDescription(qi.TechID),
				},
			})
			return nil
		}
	case types.QueuedItemTypeUnit:
		newUnit := qi.Unit
		newUnit.SetPosition(constructingBuilding.GetNearbyPosition(s, newUnit.Rect.Dx()))
		s.AddUnit(newUnit)
		return newUnit
	}
	return nil
}
