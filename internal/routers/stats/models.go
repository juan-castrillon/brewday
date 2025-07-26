package stats

import "brewday/internal/summary"

// StatsStore represents a component that stores summaries
type StatsStore interface {
	// GetAllStats returns all the statistics mapped with the b64 representation of the recipe title
	GetAllStats() (map[string]*summary.Statistics, error)
}

type StatEntry struct {
	RecipeName         string
	Evaporation        float32
	Efficiency         float32
	FinishedTimeString string
	FinishedTimeEpoch  int64
}
