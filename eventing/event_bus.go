package eventing

import (
	"fmt"
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

type ConstructUnitEvent struct {
	HiveID string
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
	fmt.Printf("EventBus: Event Subscribed To: %v", eventType)
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
