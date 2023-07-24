package lautern

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
	// AddLauternNotes adds lautern notes to the summary
	AddLaunternNotes(notes string)
}

// ReqPostLautern represents the request body for the postLauternHandler
type ReqPostLautern struct {
	Notes string `json:"notes" form:"notes"`
}
