package import_recipe

import "brewday/internal/recipe"

// RecipeParser represents a component that parses recipes
type RecipeParser interface {
	// Parse parses a recipe from a string
	Parse(recipe string) (*recipe.Recipe, error)
}

// RecipeStore represents a component that stores recipes
type RecipeStore interface {
	// Store stores a recipe and returns an identifier that can be used to retrieve it
	Store(recipe *recipe.Recipe) (string, error)
	// UpdateStatus updates the status of a recipe in the store
	UpdateStatus(id string, status recipe.RecipeStatus, statusParams ...string) error
}

// SummaryRecorderStore represents a component that stores summary recorders
// The recipe id is used as key
type SummaryRecorderStore interface {
	AddSummaryRecorder(recipeID string, recorderType string)
}

// TimelineStore represents a component that stores timelines
// The recipe id is used as key
type TimelineStore interface {
	AddTimeline(recipeID string)
}
