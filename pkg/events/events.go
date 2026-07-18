package events

import "sync"

type EventType string

const (
	ServiceStarted       EventType = "service_started"
	ServiceStopped       EventType = "service_stopped"
	ServiceStatusChanged EventType = "service_status_changed"
	ServiceError         EventType = "service_error"
	PHPVersionChanged    EventType = "php_version_changed"
	PackageInstalled     EventType = "package_installed"
	PackageRemoved       EventType = "package_removed"
)

type Event struct {
	Type    EventType
	Payload any
}

type Handler func(Event)

type Bus struct {
	mu       sync.RWMutex
	handlers map[EventType][]Handler
}

func NewBus() *Bus {
	return &Bus{
		handlers: make(map[EventType][]Handler),
	}
}

func (b *Bus) Subscribe(eventType EventType, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *Bus) Publish(event Event) {
	b.mu.RLock()
	handlers := b.handlers[event.Type]
	b.mu.RUnlock()

	for _, h := range handlers {
		h(event)
	}
}

func (b *Bus) Unsubscribe(eventType EventType, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	handlers := b.handlers[eventType]
	for i, h := range handlers {
		if &h == &handler {
			b.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}
