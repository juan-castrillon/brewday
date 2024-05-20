package summary

import "brewday/internal/summary"

// SummaryStore represents a component that stores summaries
type SummaryStore interface {
	GetSummary(id string) (*summary.Summary, error)
}

// TimelineStore represents a component that stores timelines
type TimelineStore interface {
	// GetTimeline returns a timeline of events
	GetTimeline(id string) ([]string, error)
}

// SummaryPrinter represents a component that outputs a summary as a certain document (string)
type SummaryPrinter interface {
	Print(s *summary.Summary, timeline []string) (string, error)
}
