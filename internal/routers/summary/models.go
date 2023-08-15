package summary

// SummaryRecorder represents a component that records a summary
type SummaryRecorder interface {
	AddTimeline(timeline []string)
	GetSummary() string
	GetExtention() string
	Close()
}

// Timeline represents a timeline of events
type Timeline interface {
	// GetTimeline returns a timeline of events
	GetTimeline() []string
}
