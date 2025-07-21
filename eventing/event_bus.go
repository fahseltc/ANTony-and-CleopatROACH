package eventing

import (
	"fmt"
	"gamejam/types"
	"image"
	"time"
)

// based on https://medium.com/@souravchoudhary0306/implementation-of-event-driven-architecture-in-go-golang-28d9a1c01f91

// Events
// StartSimulation
// StopSimulation
type Event struct {
	Type      string
	Timestamp time.Time
	Data      interface{}
}

type NotEnoughResourcesEvent struct {
	ResourceName   string
	UnitBeingBuilt string
}

type UnitNotUnlockedEvent struct {
	UnitName string
}

type SceneCompletionEvent struct {
	RoyalAntID   string
	RoyalRoachID string
}

type BuildClickedEvent struct {
	TargetCoordinates image.Point
	BuildingType      types.Building
}

type NotificationEvent struct {
	Message string
}

type ConstructUnitEvent struct {
	HiveID   string
	UnitType string
}

type MakeAntButtonClickedEvent struct {
	UnitType string
}

type ResearchButtonClickedEvent struct {
	TechID string
}

type ToggleRightSideHUDEvent struct {
	Show bool
}

type EventBus struct {
	subscribers map[string][]func(event Event)
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]func(event Event)),
	}
}

func (eb *EventBus) Subscribe(eventType string, handler func(event Event)) {
	fmt.Printf("EventBus: Event Subscribed To: %v\n", eventType)
	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
}

func (eb *EventBus) GetSubscribers(eventType string) []func(event Event) {
	return eb.subscribers[eventType]
}

func (eb *EventBus) Unsubscribe(eventType string) {
	if _, exists := eb.subscribers[eventType]; exists {
		delete(eb.subscribers, eventType)
	} else {
		fmt.Printf("EventBus: No subscribers found for event type: %s\n", eventType)
	}
}

// Publish sends an event to all subscribers of a given event type
func (eb *EventBus) Publish(event Event) {
	fmt.Printf("EventBus: Event Published: %v handlers: %v\n", event, len(eb.subscribers[event.Type]))
	handlers := eb.subscribers[event.Type]
	for _, handler := range handlers {
		fmt.Printf("EventBus: Calling handler for event type: %s, length: %v\n", event.Type, len(handlers))
		handler(event)
	}
}
