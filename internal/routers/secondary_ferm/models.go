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

type DryHop struct {
	// id is a unique identifier for this dry hop within the recipe. This is useful is two of the same hop are used in different times for example
	id string
	// NotificationSet is an internal boolean flag used by the frontend to keep track of whether the notification has been set or not
	NotificationSet bool
	// In signalizes whether the hop has been added to the wort or not
	In       bool
	Name     string
	Amount   float32
	Duration float32
}

// DryHopNotification is a map that relates a dry hop id to a watcher
type DryHopNotification = map[string]*watcher.Watcher

// DryHopMap is a map that relates a dry hop id to a dry hop
type DryHopMap = map[string]*DryHop

type ReqPostDryHopStart struct {
	ID               string `json:"id" form:"id"`
	NotificationTime int    `json:"notification_time" form:"notification_time"`
	TimeUnit         string `json:"time_unit" form:"time_unit"`
}

type ReqPostDryHopConfirm struct {
	ID     string  `json:"id" form:"id"`
	Amount float32 `json:"real_amount" form:"real_amount"`
}
