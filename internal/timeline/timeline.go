package timeline

import (
	"brewday/internal/timeline/types/basic"
)

// Timeline represents a timeline of events.
type Timeline interface {
	// GetTimeline returns the timeline as a slice of strings.
	GetTimeline() []string
	// AddEvent adds an event to the timeline
	AddEvent(message string)
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
