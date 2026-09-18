package eventing

import (
	"gamejam/log"
	"gamejam/types"
	"image"
	"log/slog"
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

// type SceneCompletionEvent struct {
// 	RoyalAntID   string
// 	RoyalRoachID string
// }

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

type EventBus struct {
	subscribers map[string][]func(event Event)
	log         *slog.Logger
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]func(event Event)),
		log:         log.NewLogger().With("for", "EventBus"),
	}
}

func (eb *EventBus) Subscribe(eventType string, handler func(event Event)) {
	eb.log.Info("event subscribed", "eventType", eventType)
	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
}

func (eb *EventBus) GetSubscribers(eventType string) []func(event Event) {
	return eb.subscribers[eventType]
}

func (eb *EventBus) Unsubscribe(eventType string) {
	if _, exists := eb.subscribers[eventType]; exists {
		delete(eb.subscribers, eventType)
	} else {
		eb.log.Warn("no subscribers found for event type", "eventType", eventType)
	}
}

// Publish sends an event to all subscribers of a given event type
func (eb *EventBus) Publish(event Event) {
	handlers := eb.subscribers[event.Type]
	eb.log.Info("event published", "eventType", event.Type, "handlers", len(handlers))
	for _, handler := range handlers {
		eb.log.Debug("calling handler", "eventType", event.Type, "handlers", len(handlers))
		handler(event)
	}
}
