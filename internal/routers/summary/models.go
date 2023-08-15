package summary

// SummaryRecorder represents a component that records a summary
type SummaryRecorder interface {
	GetSummary() string
	GetExtention() string
}
