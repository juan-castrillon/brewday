package summary

// SummaryRecorderStore represents a component that stores summary recorders
type SummaryRecorderStore interface {
	// AddTimeline adds a timeline to the summary
	AddTimeline(id string, timeline []string) error
	// GetSummary returns the summary
	GetSummary(id string) (string, error)
	// GetExtension returns the extension of the summary
	GetExtension(id string) (string, error)
	// Close closes the summary recorder
	Close(id string) error
}

// Timeline represents a timeline of events
type Timeline interface {
	// GetTimeline returns a timeline of events
	GetTimeline() []string
}
