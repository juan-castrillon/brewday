package cooling

import "brewday/internal/recipe"

// TimelineStore represents a component that stores timelines
type TimelineStore interface {
	// AddEvent adds an event to the timeline
	AddEvent(id, message string) error
}

// SummaryStore represents a component that stores summaries
type SummaryStore interface {
	// AddCooling adds a cooling to the summary and notes related to it
	AddCooling(id string, finalTemp, coolingTime float32, notes string) error
}

// RecipeStore represents a component that stores recipes
type RecipeStore interface {
	// UpdateStatus updates the status of a recipe in the store
	UpdateStatus(id string, status recipe.RecipeStatus, statusParams ...string) error
}

// ReqPostCooling represents the request to post a cooling
type ReqPostCooling struct {
	FinalTemp   float32 `form:"final_temp" json:"final_temp"`
	CoolingTime float32 `form:"cooling_time" json:"cooling_time"`
	Notes       string  `form:"notes" json:"notes"`
}
