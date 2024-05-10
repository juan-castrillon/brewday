package recipes

import "brewday/internal/recipe"

// RecipeStore represents a component that stores recipes
type RecipeStore interface {
	// List lists all the recipes
	List() ([]*recipe.Recipe, error)
	// Retrieve retrieves a recipe based on an identifier
	Retrieve(id string) (*recipe.Recipe, error)
	// Delete deletes a recipe based on an identifier
	Delete(id string) error
}

// SummaryRecorderStore represents a component that stores summary recorders
// The recipe id is used as key
type SummaryRecorderStore interface {
	// DeleteSummaryRecorder deletes the summary recorder for the given recipe id
	DeleteSummaryRecorder(recipeID string) error
}

// TimelineStore represents a component that stores timelines
// The recipe id is used as key
type TimelineStore interface {
	// DeleteTimeline deletes the timeline for the given recipe id
	DeleteTimeline(recipeID string) error
}
