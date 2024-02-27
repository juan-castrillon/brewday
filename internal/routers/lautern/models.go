package lautern

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

// SummaryRecorderStore represents a component that stores summary recorders
type SummaryRecorderStore interface {
	// AddLauternNotes adds lautern notes to the summary
	AddLaunternNotes(id, notes string) error
}

// ReqPostLautern represents the request body for the postLauternHandler
type ReqPostLautern struct {
	Notes string `json:"notes" form:"notes"`
}
