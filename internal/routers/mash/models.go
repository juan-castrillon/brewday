package mash

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
	// AddMashTemp adds a mash temperature to the summary and notes related to it
	AddMashTemp(temp float64, notes string)
	// AddRast adds a rast to the summary and notes related to it
	AddRast(temp float64, duration float64, notes string)
}

// ReqPostRasts represents the request body for the postRastsHandler
type ReqPostRasts struct {
	RealTemperature float64 `json:"real_temperature" form:"real_temp"`
	RealDuration    float64 `json:"real_duration" form:"real_duration"`
	Notes           string  `json:"notes" form:"notes"`
}

type ReqPostFirstRast struct {
	RealMashTemperature float64 `json:"real_mash_temperature" form:"real_mash_temp"`
	Notes               string  `json:"notes" form:"notes"`
}
