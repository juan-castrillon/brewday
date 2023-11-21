package mash

import "brewday/internal/recipe"

// RecipeStore represents a component that stores recipes
type RecipeStore interface {
	// Retrieve retrieves a recipe based on an identifier
	Retrieve(id string) (*recipe.Recipe, error)
}

// TimelineStore represents a component that stores timelines
type TimelineStore interface {
	// AddEvent adds an event to the timeline
	AddEvent(id, message string) error
}

// SummaryRecorderStore represents a component that stores summary recorders
// The recipe id is used as key
type SummaryRecorderStore interface {
	// AddMashTemp adds a mash temperature to the summary and notes related to it
	AddMashTemp(id string, temp float64, notes string) error
	// AddRast adds a rast to the summary and notes related to it
	AddRast(id string, temp float64, duration float64, notes string) error
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
