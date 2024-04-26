package memory

import (
	"errors"
	"sort"
	"sync"
	"time"
)

// BasicTimeline is a basic implementation of the Timeline interface.
// It stores the events in a slice and sorts them by time when needed.
type BasicTimeline struct {
	events []Event
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

// TimelineMemoryStore represents a timeline store stored in memory.
// The recipe id is used as key
type TimelineMemoryStore struct {
	lock      sync.Mutex
	timelines map[string]*BasicTimeline
}

// NewTimelineMemoryStore creates a new TimelineStore
func NewTimelineMemoryStore() *TimelineMemoryStore {
	return &TimelineMemoryStore{
		timelines: make(map[string]*BasicTimeline),
	}
}

// getTimelineFromStore returns the timeline for the given recipe id
func (s *TimelineMemoryStore) getTimelineFromStore(recipeID string) (*BasicTimeline, error) {
	tl, ok := s.timelines[recipeID]
	if !ok {
		return nil, errors.New("no timeline found for recipe id " + recipeID)
	}
	return tl, nil
}

// AddTimeline adds a timeline for the given recipe id
func (s *TimelineMemoryStore) AddTimeline(recipeID string) error {
	timeline := NewBasicTimeline()
	s.lock.Lock()
	defer s.lock.Unlock()
	s.timelines[recipeID] = timeline
	return s.AddEvent(recipeID, "Initialized Recipe")
}

// DeleteTimeline deletes the timeline for the given recipe id
func (s *TimelineMemoryStore) DeleteTimeline(recipeID string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.timelines, recipeID)
}

// AddEvent adds an event to the timeline for the given recipe id
func (s *TimelineMemoryStore) AddEvent(id string, message string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	tl, err := s.getTimelineFromStore(id)
	if err != nil {
		return err
	}
	timestamp := time.Now()
	e := Event{
		Timestamp: timestamp,
		Message:   message,
	}
	tl.events = append(tl.events, e)
	return nil
}

// GetTimeline returns the timeline as a slice of strings for the given recipe id
func (s *TimelineMemoryStore) GetTimeline(id string) ([]string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	tl, err := s.getTimelineFromStore(id)
	if err != nil {
		return nil, err
	}
	sort.Slice(tl.events, func(i, j int) bool {
		return tl.events[i].Timestamp.Before(tl.events[j].Timestamp)
	})
	events := make([]string, len(tl.events))
	for i, e := range tl.events {
		events[i] = e.Timestamp.Format(time.RFC3339Nano) + " " + e.Message
	}
	return events, nil
}
