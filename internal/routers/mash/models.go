package mash

import "brewday/internal/recipe"

// RecipeStore represents a component that stores recipes
type RecipeStore interface {
	// Retrieve retrieves a recipe based on an identifier
	Retrieve(id string) (*recipe.Recipe, error)
	// UpdateStatus updates the status of a recipe in the store
	UpdateStatus(id string, status recipe.RecipeStatus, statusParams ...string) error
}

// TimelineStore represents a component that stores timelines
type TimelineStore interface {
	// AddEvent adds an event to the timeline
	AddEvent(id, message string) error
}

// SummaryStore represents a component that stores summaries
// The recipe id is used as key
type SummaryStore interface {
	AddMashTemp(id string, temp float32, notes string) error
	AddRast(id string, temp float32, duration float32, notes string) error
}

// ReqPostRasts represents the request body for the postRastsHandler
type ReqPostRasts struct {
	RealTemperature float32 `json:"real_temperature" form:"real_temp"`
	RealDuration    float32 `json:"real_duration" form:"real_duration"`
	Notes           string  `json:"notes" form:"notes"`
}

type ReqPostFirstRast struct {
	RealMashTemperature float32 `json:"real_mash_temperature" form:"real_mash_temp"`
	Notes               string  `json:"notes" form:"notes"`
}
