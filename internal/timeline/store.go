package timeline

import (
	"brewday/internal/timeline/basic"
	"errors"
	"sync"
)

// TimelineStore represents a timeline store.
// The recipe id is used as key
type TimelineStore struct {
	lock      sync.Mutex
	timelines map[string]Timeline
}

// NewTimelineStore creates a new TimelineStore
func NewTimelineStore() *TimelineStore {
	return &TimelineStore{
		timelines: make(map[string]Timeline),
	}
}

// TimelineFactory is a factory for Timeline implementations
func TimelineFactory(timelineType string) Timeline {
	switch timelineType {
	case "simple":
		return basic.NewBasicTimeline()
	default:
		return basic.NewBasicTimeline()
	}
}

// getTimelineFromStore returns the timeline for the given recipe id
func (s *TimelineStore) getTimelineFromStore(recipeID string) (Timeline, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	tl, ok := s.timelines[recipeID]
	if !ok {
		return nil, errors.New("no timeline found for recipe id " + recipeID)
	}
	return tl, nil
}

// AddTimeline adds a timeline for the given recipe id
func (s *TimelineStore) AddTimeline(recipeID string, timelineType string) {
	timeline := TimelineFactory(timelineType)
	s.lock.Lock()
	defer s.lock.Unlock()
	s.timelines[recipeID] = timeline
}

// DeleteTimeline deletes the timeline for the given recipe id
func (s *TimelineStore) DeleteTimeline(recipeID string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.timelines, recipeID)
}

// AddEvent adds an event to the timeline for the given recipe id
func (s *TimelineStore) AddEvent(id string, message string) error {
	tl, err := s.getTimelineFromStore(id)
	if err != nil {
		return err
	}
	tl.AddEvent(message)
	return nil
}

// GetTimeline returns the timeline as a slice of strings for the given recipe id
func (s *TimelineStore) GetTimeline(id string) ([]string, error) {
	tl, err := s.getTimelineFromStore(id)
	if err != nil {
		return nil, err
	}
	return tl.GetTimeline(), nil
}
