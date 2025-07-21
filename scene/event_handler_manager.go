package scene

import (
	"fmt"
	"gamejam/eventing"
	"gamejam/sim"
	"gamejam/types"
	"gamejam/ui"
	"image"
)

type EventHandlerManager struct {
	eventBus *eventing.EventBus
	scene    *PlayScene
}

func NewEventHandlerManager(bus *eventing.EventBus, scene *PlayScene) *EventHandlerManager {
	ehm := &EventHandlerManager{
		eventBus: bus,
		scene:    scene,
	}

	bus.Subscribe("NotificationEvent", ehm.HandleNotificationEvent)
	bus.Subscribe("NotEnoughResourcesEvent", ehm.HandleNotEnoughResourcesEvent)

	bus.Subscribe("MakeAntButtonClickedEvent", ehm.HandleMakeAntButtonClickedEvent)
	bus.Subscribe("BuildingButtonClickedEvent", ehm.HandleBuildingButtonClickedEvent)
	bus.Subscribe("BuildClickedEvent", ehm.HandleBuildClickedEvent)

	bus.Subscribe("ResearchButtonClickedEvent", ehm.HandleResearchButtonClickedEvent)
	bus.Subscribe("SetRallyPointEvent", ehm.HandleSetRallyPointEvent)

	return ehm
}

// Shows a UI notification based on an event
func (ehm *EventHandlerManager) HandleNotificationEvent(event eventing.Event) {
	ehm.scene.CurrentNotification = ui.NewNotification(&ehm.scene.fonts.Med, event.Data.(eventing.NotificationEvent).Message)
}

// Shows a UI notification about not having enough resources
func (ehm *EventHandlerManager) HandleNotEnoughResourcesEvent(event eventing.Event) {
	resName := event.Data.(eventing.NotEnoughResourcesEvent).ResourceName
	target := event.Data.(eventing.NotEnoughResourcesEvent).UnitBeingBuilt

	var str string
	if target == "Bridge" {
		str = fmt.Sprintf("Not enough %v to build %v\nOr builder is not close enough!", resName, target)
	} else {
		str = fmt.Sprintf("Not enough %v to build %v", resName, target)
	}

	ehm.scene.CurrentNotification = ui.NewNotification(&ehm.scene.fonts.Med, str)
}

// Handles when a UI button is clicked to create a unit. Adds Hive ID to the event and forwards to SIM for handling
func (ehm *EventHandlerManager) HandleMakeAntButtonClickedEvent(event eventing.Event) {
	if len(ehm.scene.selectedUnitIDs) == 1 {
		hiveID := ehm.scene.selectedUnitIDs[0]
		unitOrHiveString := ehm.scene.sim.DetermineUnitOrHiveById(hiveID)
		innerEvent := event.Data.(*eventing.MakeAntButtonClickedEvent)
		unitType := innerEvent.UnitType
		if unitOrHiveString == "hive" {
			ehm.eventBus.Publish(eventing.Event{
				Type: "ConstructUnitEvent",
				Data: eventing.ConstructUnitEvent{
					HiveID:   hiveID,
					UnitType: unitType,
				},
			})
		}
	}
}

// Handles when a UI button is clicked to create a building while a unit is selected. Enables construction ghost on mouse.
func (ehm *EventHandlerManager) HandleBuildingButtonClickedEvent(event eventing.Event) {
	if len(ehm.scene.selectedUnitIDs) >= 1 {
		innerEvent := event.Data.(eventing.BuildClickedEvent)
		unitID := ehm.scene.selectedUnitIDs[0]
		unitOrHiveString := ehm.scene.sim.DetermineUnitOrHiveById(unitID)
		if unitOrHiveString == "unit" {
			ehm.scene.constructionMouse.Enabled = true
			ehm.scene.constructionMouse.SetSprite(innerEvent.BuildingType)
			ehm.scene.drag.Enabled = false
		}
	}
}

// Handles when the mouse is clicked while in ghost construction mode. Actually removes resources and places the building.
func (ehm *EventHandlerManager) HandleBuildClickedEvent(event eventing.Event) {
	innerEvent := event.Data.(eventing.BuildClickedEvent)
	if len(ehm.scene.selectedUnitIDs) >= 1 {
		success := ehm.scene.sim.ConstructBuilding(innerEvent.TargetCoordinates, ehm.scene.selectedUnitIDs[0], innerEvent.BuildingType)
		if !success {
			ehm.eventBus.Publish(eventing.Event{
				Type: "NotEnoughResourcesEvent",
				Data: eventing.NotEnoughResourcesEvent{ // todo: add reason why, for example "unit not close enough" etc
					ResourceName:   "Wood",
					UnitBeingBuilt: "Bridge",
				},
			})
		}
	}
	ehm.scene.drag.Enabled = true
	ehm.scene.constructionMouse.Enabled = false
}

// Handles a UI button pressed while a hive is selected to begin research
func (ehm *EventHandlerManager) HandleResearchButtonClickedEvent(event eventing.Event) {
	innerEvent := event.Data.(eventing.ResearchButtonClickedEvent)
	hiveID := ehm.scene.selectedUnitIDs[0]
	unitOrHiveString := ehm.scene.sim.DetermineUnitOrHiveById(hiveID)
	if unitOrHiveString == "hive" { // double check they can afford it?
		hive, err := ehm.scene.sim.GetBuildingByID(hiveID)
		if err == nil {
			tech := ehm.scene.sim.GetPlayerState().TechTree.AvailableTech[sim.TechID(innerEvent.TechID)]
			qi := &sim.QueuedItem{
				Type:             types.QueuedItemTypeTech,
				TechID:           sim.TechID(innerEvent.TechID),
				ConstructionTime: tech.ResearchSeconds,
			}
			hive.AddItemToBuildQueue(qi)
		}
	} else {
		// error shouldnt happen? button isnt enabled unless a hive is selected
	}
}

// Handles a UI button press to set a hive's rally point
func (ehm *EventHandlerManager) HandleSetRallyPointEvent(event eventing.Event) {
	hiveID := ehm.scene.selectedUnitIDs[0]
	unitOrHiveString := ehm.scene.sim.DetermineUnitOrHiveById(hiveID)
	if unitOrHiveString == "hive" {
		hive, err := ehm.scene.sim.GetBuildingByID(hiveID)
		if err == nil {
			mouseWorldX, mouseWorldY := ehm.scene.Ui.Camera.MousePosToMapPos()
			pt := &image.Point{X: mouseWorldX, Y: mouseWorldY}
			ehm.scene.ActionIssuedLocation = pt
			hive.SetRallyPoint(pt)
		}
	}
}
