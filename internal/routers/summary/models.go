package summary

import (
	"brewday/internal/recipe"
	"brewday/internal/summary"
)

// SummaryStore represents a component that stores summaries
type SummaryStore interface {
	GetSummary(id string) (*summary.Summary, error)
}

// RecipeStore represents a component that stores recipes
type RecipeStore interface {
	// Retrieve retrieves a recipe based on an identifier
	Retrieve(id string) (*recipe.Recipe, error)
}

// TimelineStore represents a component that stores timelines
type TimelineStore interface {
	// GetTimeline returns a timeline of events
	GetTimeline(id string) ([]string, error)
}

// SummaryPrinter represents a component that outputs a summary as a certain document (string)
type SummaryPrinter interface {
	Print(s *summary.Summary, r *recipe.Recipe, timeline []string) (string, error)
}
