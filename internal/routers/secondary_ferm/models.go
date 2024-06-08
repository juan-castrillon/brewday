package secondaryferm

import (
	"brewday/internal/recipe"
	"time"
)

// TimelineStore represents a component that stores timelines
type TimelineStore interface {
	// AddEvent adds an event to the timeline
	AddEvent(id, message string) error
}

// SummaryStore represents a component that stores summaries
type SummaryStore interface {
	AddDryHopStart(id string, name string, amount, alpha float32, notes string) error
	AddDryHopEnd(id string, name string, durationHours float32) error
	AddPreBottlingVolume(id string, volume float32) error
	AddBottling(id string, carbonation, alcohol, sugar, temp, vol float32, sugarType, notes string) error
	AddSummarySecondary(id string, days int, notes string) error
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
	// AddSugarResult adds a new priming sugar result to a given recipe
	AddSugarResult(id string, r *recipe.PrimingSugarResult) error
	// RetrieveSugarResults returns all sugar results for a recipe
	RetrieveSugarResults(id string) ([]*recipe.PrimingSugarResult, error)
	// AddDate allows to store a date with a certain purpose. It can be used to store notification dates, or timers
	AddDate(id string, date *time.Time, name string) error
	// RetrieveDates allows to retreive stored dates with its purpose (name).It can be used to store notification dates, or timers
	// It supports pattern in the name to retrieve multiple values
	RetrieveDates(id, namePattern string) ([]*time.Time, error)
	// AddBoolFlag allows to store a given flag that can be true or false in the store with a unique name
	AddBoolFlag(id, name string, flag bool) error
	// RetrieveBoolFlag gets a bool flag from the store given its name
	RetrieveBoolFlag(id, name string) (bool, error)
}

// Notifier is the interface that helps decouple the notifier from the application
type Notifier interface {
	// Send sends a notification
	Send(message, title string, opts map[string]any) error
}

// SugarResult is the result of the sugar calculation
type SugarResult struct {
	// Water is the amount of water in liters
	Water float32
	// Amount is the amount of sugar in grams
	Amount float32
	// Alcohol is the estimated final alcohol content
	Alcohol float32
}

type ReqPostDryHopIn struct {
	IngredientName string  `json:"ingredient_name"`
	RealAmount     float32 `json:"real_amount,omitempty"`
	RealAlpha      float32 `json:"real_alpha,omitempty"`
}

type ReqPostPreBottle struct {
	Volume      float32 `json:"volume" form:"volume"`
	LostVolume  float32 `json:"lost" form:"lost"`
	SugarType   string  `json:"sugar_type" form:"sugar_type"`
	Temperature float32 `json:"temperature" form:"temperature"`
}

type ReqPostBottle struct {
	RealVolume  float32 `json:"real_volume" form:"real_volume"`
	SugarAmount float32 `json:"sugar_amount" form:"sugar_amount"`
	SugarType   string  `json:"sugar_type" form:"sugar_type"`
	Water       float32 `json:"water" form:"water"`
	Temperature float32 `json:"temperature" form:"temperature"`
	Time        float32 `json:"time" form:"time"`
	Notes       string  `json:"notes" form:"notes"`
}

type ReqPostSecondaryStart struct {
	NotificationTime int    `json:"notification_time" form:"notification_time"`
	TimeUnit         string `json:"time_unit" form:"time_unit"`
}

type ReqPostSecondaryEnd struct {
	Days  int    `json:"days" form:"days"`
	Notes string `json:"notes" form:"notes"`
}
