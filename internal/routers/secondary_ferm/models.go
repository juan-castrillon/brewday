package secondaryferm

import (
	"brewday/internal/recipe"
	"brewday/internal/watcher"
)

// TimelineStore represents a component that stores timelines
type TimelineStore interface {
	// AddEvent adds an event to the timeline
	AddEvent(id, message string) error
}

// SummaryRecorderStore represents a component that stores summary recorders
type SummaryRecorderStore interface {
	// AddSummaryDryHop adds a summary of the dry hop
	AddSummaryDryHop(id string, name string, amount float32) error
	// AddSummaryPreBottle adds a summary of the pre bottling
	AddSummaryPreBottle(id string, volume float32) error
	// AddSummaryBottle adds a summary of the bottling
	AddSummaryBottle(id string, carbonation, alcohol, sugar, temp, vol float32, sugarType, notes string) error
	// AddSummarySecondary adds a summary of the secondary fermentation
	AddSummarySecondary(id string, days int, notes string) error
}

// RecipeStore represents a component that stores recipes
type RecipeStore interface {
	// Retrieve retrieves a recipe based on an identifier
	Retrieve(id string) (*recipe.Recipe, error)
	// UpdateStatus updates the status of a recipe in the store
	UpdateStatus(id string, status recipe.RecipeStatus, statusParams ...string) error
	// UpdateResults updates a certain result of a recipe
	UpdateResults(id string, resultType recipe.ResultType, value float32) error
}

// Notifier is the interface that helps decouple the notifier from the application
type Notifier interface {
	// Send sends a notification
	Send(message, title string, opts map[string]any) error
}

type DryHop struct {
	// id is a unique identifier for this dry hop within the recipe. This is useful is two of the same hop are used in different times for example
	id string
	// NotificationSet is an internal boolean flag used by the frontend to keep track of whether the notification has been set or not
	NotificationSet bool
	// In signalizes whether the hop has been added to the wort or not
	In bool
	// InDate is the date when the hop was added to the wort
	InDate   string
	Name     string
	Amount   float32
	Duration float32
}

// DryHopNotification is a map that relates a dry hop id to a watcher
type DryHopNotification = map[string]*watcher.Watcher

// DryHopMap is a map that relates a dry hop id to a dry hop
type DryHopMap = map[string]*DryHop

// SugarResult is the result of the sugar calculation
type SugarResult struct {
	// Water is the amount of water in liters
	Water float32
	// Amount is the amount of sugar in grams
	Amount float32
	// Alcohol is the estimated final alcohol content
	Alcohol float32
}

type SecondaryFermentationWatcher struct {
	watch *watcher.Watcher
}

type ReqPostDryHopStart struct {
	ID               string `json:"id" form:"id"`
	NotificationTime int    `json:"notification_time" form:"notification_time"`
	TimeUnit         string `json:"time_unit" form:"time_unit"`
}

type ReqPostDryHopConfirm struct {
	ID     string  `json:"id" form:"id"`
	Amount float32 `json:"real_amount" form:"real_amount"`
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
