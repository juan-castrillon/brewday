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

// SummaryRecorderStore represents a component that stores summary recorders
type SummaryRecorderStore interface {
	// AddSummaryPreFermentation adds a summary of the pre fermentation
	AddSummaryPreFermentation(id string, volume float32, sg float32, notes string) error
	// AddEfficiency adds the efficiency (sudhausausbeute) to the summary
	AddEfficiency(id string, efficiencyPercentage float32) error
	// AddYeastStart adds the yeast start to the summary
	AddYeastStart(id string, temperature, notes string) error
	// AddSGMeasurement adds a SG measurement to the summary
	AddSGMeasurement(id string, date string, gravity float32, final bool, notes string) error
	// AddAlcoholMainFermentation adds the alcohol after the main fermentation to the summary
	AddAlcoholMainFermentation(id string, alcohol float32) error
}

// RecipeStore represents a component that stores recipes
type RecipeStore interface {
	// Retrieve retrieves a recipe based on an identifier
	Retrieve(id string) (*recipe.Recipe, error)
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
