package stores

import (
	"brewday/internal/timeline"
	"errors"
	"sync"
)

// TimelineMemoryStore represents a timeline store stored in memory.
// The recipe id is used as key
type TimelineMemoryStore struct {
	lock      sync.Mutex
	timelines map[string]timeline.Timeline
}

// NewTimelineMemoryStore creates a new TimelineStore
func NewTimelineMemoryStore() *TimelineMemoryStore {
	return &TimelineMemoryStore{
		timelines: make(map[string]timeline.Timeline),
	}
}

// getTimelineFromStore returns the timeline for the given recipe id
func (s *TimelineMemoryStore) getTimelineFromStore(recipeID string) (timeline.Timeline, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	tl, ok := s.timelines[recipeID]
	if !ok {
		return nil, errors.New("no timeline found for recipe id " + recipeID)
	}
	return tl, nil
}

// AddTimeline adds a timeline for the given recipe id
func (s *TimelineMemoryStore) AddTimeline(recipeID string, timelineType string) {
	timeline := timeline.TimelineFactory(timelineType)
	s.lock.Lock()
	defer s.lock.Unlock()
	s.timelines[recipeID] = timeline
}

// DeleteTimeline deletes the timeline for the given recipe id
func (s *TimelineMemoryStore) DeleteTimeline(recipeID string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.timelines, recipeID)
}

// AddEvent adds an event to the timeline for the given recipe id
func (s *TimelineMemoryStore) AddEvent(id string, message string) error {
	tl, err := s.getTimelineFromStore(id)
	if err != nil {
		return err
	}
	tl.AddEvent(message)
	return nil
}

// GetTimeline returns the timeline as a slice of strings for the given recipe id
func (s *TimelineMemoryStore) GetTimeline(id string) ([]string, error) {
	tl, err := s.getTimelineFromStore(id)
	if err != nil {
		return nil, err
	}
	return tl.GetTimeline(), nil
}
