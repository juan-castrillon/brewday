package cooling

// Timeline represents a timeline of events
type Timeline interface {
	// AddEvent adds an event to the timeline
	AddEvent(message string)
}

// SummaryRecorder represents a component that records a summary
type SummaryRecorder interface {
	// AddCooling adds a cooling to the summary and notes related to it
	AddCooling(finalTemp, coolingTime float32, notes string)
}

// ReqPostCooling represents the request to post a cooling
type ReqPostCooling struct {
	FinalTemp   float32 `form:"final_temp" json:"final_temp"`
	CoolingTime float32 `form:"cooling_time" json:"cooling_time"`
	Notes       string  `form:"notes" json:"notes"`
}
