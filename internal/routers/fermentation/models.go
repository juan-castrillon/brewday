package fermentation

import (
	"brewday/internal/recipe"
	"brewday/internal/watcher"
)

type SGMeasurement struct {
	Date    string
	Gravity float32
}

type FermentationStatus struct {
	MinDaysPassed bool
	InitialWatch  *watcher.Watcher
}

// TimelineStore represents a component that stores timelines
type TimelineStore interface {
	// AddEvent adds an event to the timeline
	AddEvent(id, message string) error
}

// SummaryStore represents a component that stores summaries
type SummaryStore interface {
	AddPreFermentationVolume(id string, volume float32, sg float32, notes string) error
	AddYeastStart(id string, temperature, notes string) error
	AddMainFermentationSGMeasurement(id string, date string, gravity float32, final bool, notes string) error
	AddMainFermentationAlcohol(id string, alcohol float32) error
	AddMainFermentationDryHop(id string, name string, amount, alpha, duration float32, notes string) error
	AddEfficiency(id string, efficiencyPercentage float32) error
}

// RecipeStore represents a component that stores recipes
type RecipeStore interface {
	// Retrieve retrieves a recipe based on an identifier
	Retrieve(id string) (*recipe.Recipe, error)
	// UpdateStatus updates the status of a recipe in the store
	UpdateStatus(id string, status recipe.RecipeStatus, statusParams ...string) error
	// UpdateResult updates a certain result of a recipe
	UpdateResult(id string, resultType recipe.ResultType, value float32) error
	// RetrieveResults gets the results from a certain recipe
	RetrieveResults(id string) (*recipe.RecipeResults, error)
	// AddMainFermSG adds a new specific gravity measurement to a given recipe
	AddMainFermSG(id string, m *recipe.SGMeasurement) error
	// RetrieveMainFermSGs returns all measured sgs for a recipe
	RetrieveMainFermSGs(id string) ([]*recipe.SGMeasurement, error)
}

// Notifier is the interface that helps decouple the notifier from the application
type Notifier interface {
	// Send sends a notification
	Send(message, title string, opts map[string]any) error
}

// ReqPostPreFermentation represents the request for the post pre fermentation page
type ReqPostPreFermentation struct {
	Volume float32 `json:"volume" form:"volume"`
	SG     float32 `json:"sg" form:"sg"`
	Notes  string  `json:"notes" form:"notes"`
}

// WaterOption represent an option for adding water
type WaterOption struct {
	ToAdd        float32 `json:"to_add"`
	FinalVolume  float32 `json:"final_volume"`
	FinalSG      float32 `json:"final_sg"`
	FinalSGPlato float32 `json:"final_sg_plato"`
}

// ReqPostPreFermentationWater represents the request for the post pre fermentation water page
type ReqPostPreFermentationWater struct {
	FinalVolume float32 `json:"final_volume" form:"final_volume"`
	FinalSG     float32 `json:"final_sg" form:"final_sg"`
	Notes       string  `json:"notes" form:"notes"`
}

// ReqPostFermentationYeast represents the request for the post yeast fermentation page
type ReqPostFermentationYeast struct {
	Temperature string `json:"temperature" form:"temperature"` // string because it can be a range
	Notes       string `json:"notes" form:"notes"`
}

// ReqPostFermentationStart represents the request for the post fermentation start page
type ReqPostFermentationStart struct {
	NotificationDays       int    `json:"notification_days" form:"notification_days"`
	NotificationDaysBefore int    `json:"notification_days_before" form:"notification_days_before"`
	TimeUnit               string `json:"time_unit" form:"time_unit"`
}

// ReqPostMainFermentation represents the request for the post main fermentation page
type ReqPostMainFermentation struct {
	SG    float32 `json:"sg" form:"sg"`
	Final bool    `json:"final" form:"final"`
	Notes string  `json:"notes" form:"notes"`
}
