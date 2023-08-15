package basic

import (
	"sort"
	"sync"
	"time"
)

// BasicTimeline is a basic implementation of the Timeline interface.
// It stores the events in a slice and sorts them by time when needed.
type BasicTimeline struct {
	events []Event
	lock   sync.Mutex
}

// Event describes a single event in the timeline.
type Event struct {
	Timestamp time.Time
	Message   string
}

// NewBasicTimeline creates a new BasicTimeline.
func NewBasicTimeline() *BasicTimeline {
	return &BasicTimeline{
		events: make([]Event, 0),
	}
}

// AddEvent adds an event to the timeline
func (t *BasicTimeline) AddEvent(message string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	timestamp := time.Now()
	e := Event{
		Timestamp: timestamp,
		Message:   message,
	}
	t.events = append(t.events, e)
}

// GetTimeline returns the timeline as a slice of strings.
// It sorts the events by time before returning them.
func (t *BasicTimeline) GetTimeline() []string {
	t.lock.Lock()
	defer t.lock.Unlock()
	sort.Slice(t.events, func(i, j int) bool {
		return t.events[i].Timestamp.Before(t.events[j].Timestamp)
	})
	events := make([]string, len(t.events))
	for i, e := range t.events {
		events[i] = e.Timestamp.Format(time.RFC3339Nano) + " " + e.Message
	}
	return events
}
