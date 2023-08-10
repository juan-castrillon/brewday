package fermentation

import "brewday/internal/recipe"

// Timeline represents a timeline of events
type Timeline interface {
	// AddEvent adds an event to the timeline
	AddEvent(message string)
}

// SummaryRecorder represents a component that records a summary
type SummaryRecorder interface {
	AddSummaryPreFermentation(volume float32, sg float32, notes string)
	AddEfficiency(efficiencyPercentage float32)
	AddYeastStart(temperature, notes string)
	// Close closes the summary recorder
	Close()
}

// RecipeStore represents a component that stores recipes
type RecipeStore interface {
	// Retrieve retrieves a recipe based on an identifier
	Retrieve(id string) (*recipe.Recipe, error)
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

// ReqPostFermentation represents the request for the post fermentation page
type ReqPostFermentation struct {
	Temperature string `json:"temperature" form:"temperature"` // string because it can be a range
	Notes       string `json:"notes" form:"notes"`
}
