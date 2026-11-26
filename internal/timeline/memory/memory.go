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
	Tag       string
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
	s.timelines[recipeID] = timeline
	s.lock.Unlock()
	return s.AddEvent(recipeID, "Initialized Recipe")
}

// DeleteTimeline deletes the timeline for the given recipe id
func (s *TimelineMemoryStore) DeleteTimeline(recipeID string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.timelines[recipeID]
	if !ok {
		return errors.New("timeline not found")
	}
	delete(s.timelines, recipeID)
	return nil
}

// AddEvent adds an event to the timeline for the given recipe id
// The event can be tagged for retrieval using a tag argument (jus the first one will be used)
func (s *TimelineMemoryStore) AddEvent(id string, message string, tag ...string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	tl, err := s.getTimelineFromStore(id)
	if err != nil {
		return err
	}
	t := ""
	if len(tag) > 0 {
		t = tag[0]
	}
	timestamp := time.Now()
	e := Event{
		Timestamp: timestamp,
		Message:   message,
		Tag:       t,
	}
	tl.events = append(tl.events, e)
	return nil
}

// GetTimeline returns the timeline as a slice of strings for the given recipe id
// Events are represented in the format <Timestamp>@<Event>
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
		events[i] = e.Timestamp.Format(time.RFC3339Nano) + "@" + e.Message
	}
	return events, nil
}

// GetTimestampByTag retrieves a timestamp from from the recipe (id) timeline using a tag
// This means only tagged events can be retrieved
func (s *TimelineMemoryStore) GetTimestampByTag(id, tag string) (time.Time, error) {
	// TODO: Implement, test also what happens if not found
	return time.Now(), nil
}
