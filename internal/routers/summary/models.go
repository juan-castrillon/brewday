package summary

// SummaryRecorder represents a component that records a summary
type SummaryRecorder interface {
	// AddTimeline adds a timeline to the summary
	AddTimeline(timeline []string)
	// GetSummary returns the summary
	GetSummary() string
	// GetExtention returns the extension of the summary
	GetExtention() string
	// Close closes the summary recorder
	Close()
}

// Timeline represents a timeline of events
type Timeline interface {
	// GetTimeline returns a timeline of events
	GetTimeline() []string
}
