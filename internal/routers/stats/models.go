package stats

import "brewday/internal/summary"

// SummaryStore represents a component that stores summaries
type SummaryStore interface {
	GetAllSummaries() ([]string, error)
	GetSummary(id string) (*summary.Summary, error)
}

// RecipeStats represent important statistics for a recipe
type RecipeStats struct {
	Evaporation float32
	Efficiency  float32
}
