package hopping

import "brewday/internal/recipe"

// RecipeStore represents a component that stores recipes
type RecipeStore interface {
	// Retrieve retrieves a recipe based on an identifier
	Retrieve(id string) (*recipe.Recipe, error)
}

// Timeline represents a timeline of events
type Timeline interface {
	// AddEvent adds an event to the timeline
	AddEvent(message string)
}

// SummaryRecorder represents a component that records a summary
type SummaryRecorder interface {
	// AddHopping adds a hopping to the summary and notes related to it
	AddHopping(name string, amount float32, alpha float32, notes string)
	// AddMeasuredVolume adds a measured volume to the summary
	AddMeasuredVolume(name string, amount float32, notes string)
}

// ReqPostStartHopping is the request for the start hopping route
type ReqPostStartHopping struct {
	InitialVolume float32 `json:"initial_volume" form:"initial_volume"`
	Notes         string  `json:"notes" form:"notes"`
}

// ReqPostEndHopping is the request for the end hopping route
type ReqPostEndHopping struct {
	FinalVolume float32 `json:"final_volume" form:"final_volume"`
	Notes       string  `json:"notes" form:"notes"`
}

// ReqPostHopping is the response for the hopping route
type ReqPostHopping struct {
	RealAmount   float32 `json:"real_amount" form:"real_amount"`
	RealDuration float32 `json:"real_duration" form:"real_duration"`
	RealAlpha    float32 `json:"real_alpha" form:"real_alpha"`
}