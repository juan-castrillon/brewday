package hopping

import "brewday/internal/recipe"

// RecipeStore represents a component that stores recipes
type RecipeStore interface {
	// Retrieve retrieves a recipe based on an identifier
	Retrieve(id string) (*recipe.Recipe, error)
	// UpdateStatus updates the status of a recipe in the store
	UpdateStatus(id string, status recipe.RecipeStatus, statusParams ...string) error
	// UpdateResult updates a certain result of a recipe
	UpdateResult(id string, resultType recipe.ResultType, value float32) error
	// RetrieveResult gets a certain result value from a recipe
	RetrieveResult(id string, resultType recipe.ResultType) (float32, error)
}

// TimelineStore represents a component that stores timelines
type TimelineStore interface {
	// AddEvent adds an event to the timeline
	AddEvent(id, message string) error
}

// SummaryStore represents a component that stores summaries
type SummaryStore interface {
	AddHopping(id string, name string, amount float32, alpha float32, duration float32, notes string) error
	AddVolumeBeforeBoil(id string, amount float32, notes string) error
	AddVolumeAfterBoil(id string, amount float32, notes string) error
	AddEvaporation(id string, amount float32) error
}

// ReqPostStartHopping is the request for the start hopping route
type ReqPostStartHopping struct {
	InitialVolume float32 `json:"initial_volume" form:"initial_volume"`
	Notes         string  `json:"notes" form:"notes"`
}

// ReqPostEndHopping is the request for the end hopping route
type ReqPostEndHopping struct {
	FinalVolume float32 `json:"final_volume" form:"final_volume"`
	Notes       string  `json:"notes" form:"notes"`
}

// ReqPostHopping is the response for the hopping route
type ReqPostHopping struct {
	RealAmount   float32 `json:"real_amount" form:"real_amount"`
	RealDuration float32 `json:"real_duration" form:"real_duration"`
	RealAlpha    float32 `json:"real_alpha" form:"real_alpha"`
	Notes        string  `json:"notes" form:"notes"`
}
