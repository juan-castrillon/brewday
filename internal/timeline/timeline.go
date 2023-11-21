package timeline

// Timeline represents a timeline of events.
type Timeline interface {
	// GetTimeline returns the timeline as a slice of strings.
	GetTimeline() []string
	// AddEvent adds an event to the timeline
	AddEvent(message string)
}
