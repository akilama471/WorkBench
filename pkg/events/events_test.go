package events

import (
	"sync"
	"testing"
)

func TestEventBusSubscribeAndPublish(t *testing.T) {
	bus := NewBus()

	var received Event
	var wg sync.WaitGroup
	wg.Add(1)

	bus.Subscribe(ServiceStarted, func(e Event) {
		received = e
		wg.Done()
	})

	bus.Publish(Event{Type: ServiceStarted, Payload: "apache"})

	wg.Wait()

	if received.Type != ServiceStarted {
		t.Errorf("received event type = %v, want %v", received.Type, ServiceStarted)
	}
	if received.Payload != "apache" {
		t.Errorf("received payload = %v, want \"apache\"", received.Payload)
	}
}

func TestEventBusNoCrossSubscription(t *testing.T) {
	bus := NewBus()

	serviceReceived := false
	runtimeReceived := false

	bus.Subscribe(ServiceStarted, func(e Event) {
		serviceReceived = true
	})

	bus.Subscribe(PHPVersionChanged, func(e Event) {
		runtimeReceived = true
	})

	bus.Publish(Event{Type: ServiceStarted, Payload: "apache"})

	if !serviceReceived {
		t.Error("service handler not called")
	}
	if runtimeReceived {
		t.Error("runtime handler should not be called for service event")
	}
}

func TestEventBusMultipleSubscribers(t *testing.T) {
	bus := NewBus()

	callCount := 0
	var mu sync.Mutex

	handler1 := func(e Event) {
		mu.Lock()
		callCount++
		mu.Unlock()
	}
	handler2 := func(e Event) {
		mu.Lock()
		callCount++
		mu.Unlock()
	}

	bus.Subscribe(ServiceStarted, handler1)
	bus.Subscribe(ServiceStarted, handler2)

	bus.Publish(Event{Type: ServiceStarted, Payload: "apache"})

	mu.Lock()
	defer mu.Unlock()
	if callCount != 2 {
		t.Errorf("callCount = %d, want 2", callCount)
	}
}
