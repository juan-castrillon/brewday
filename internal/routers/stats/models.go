package stats

import "brewday/internal/summary"

// StatsStore represents a component that stores summaries
type StatsStore interface {
	// GetAllStats returns all the statistics mapped with the b64 representation of the recipe title
	GetAllStats() (map[string]*summary.Statistics, error)
	// AddStatsExternal adds statistics from recipes outside the app
	AddStatsExternal(recipeName string, stats *summary.Statistics) error
}

type StatEntry struct {
	RecipeName         string
	Evaporation        *float32
	Efficiency         *float32
	FinishedTimeString string
	FinishedTimeEpoch  int64
}

// ReqPostAddStat represents the request for adding a external stat
type ReqPostAddStat struct {
	RecipeName         string  `json:"name" form:"name"`
	Evaporation        float32 `json:"evaporation" form:"evaporation"`
	Efficiency         float32 `json:"efficiency" form:"efficiency"`
	FinishedTimeString string  `json:"finished" form:"finished"`
}
