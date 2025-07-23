package stats

import "brewday/internal/summary"

// StatsStore represents a component that stores summaries
type StatsStore interface {
	// GetAllStats returns all the statistics mapped with the b64 representation of the recipe title
	GetAllStats() (map[string]*summary.Statistics, error)
}
