package lautern

import (
	"brewday/internal/recipe"
	"time"

	"github.com/labstack/echo/v4"
)

// RecipeStore represents a component that stores recipes
type RecipeStore interface {
	// Retrieve retrieves a recipe based on an identifier
	Retrieve(id string) (*recipe.Recipe, error)
	// UpdateStatus updates the status of a recipe in the store
	UpdateStatus(id string, status recipe.RecipeStatus, statusParams ...string) error
}

type Timer interface {
	// GetBoolFlags returns whether the timer has started and has been stopped. Only the first suffix is used
	GetBoolFlags(id string, prefix string, suffix ...string) (bool, bool, error)
	//HandleStartTimer will respond with the correct json for the timer template to work. Only the first suffix is used
	HandleStartTimer(c echo.Context, id string, duration time.Duration, prefix string, suffix ...string) error
	//HandleStopTimer will mark the timer as stopped. Only the first suffix is used
	HandleStopTimer(c echo.Context, id string, timelineEvent string, notificationMessage string, notificationTitle string, prefix string, suffix ...string) error
	//HandleRealDuration will return the real duration to the timer template. Only the first suffix is used
	HandleRealDuration(c echo.Context, id string, prefix string, suffix ...string) error
}

// TimelineStore represents a component that stores timelines
type TimelineStore interface {
	// AddEvent adds an event to the timeline
	AddEvent(id, message string) error
}

// SummaryStore represents a component that stores summaries
type SummaryStore interface {
	// AddLauternNotes adds lautern notes to the summary
	AddLauternNotes(id, notes string) error
}

// ReqPostLautern represents the request body for the postLauternHandler
type ReqPostLautern struct {
	Notes        string  `json:"notes" form:"notes"`
	RealDuration float32 `json:"real_duration" form:"real_duration"`
}
