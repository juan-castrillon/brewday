package cooling

import (
	"brewday/internal/recipe"
	"time"

	"github.com/labstack/echo/v4"
)

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

// ReqPostCooling represents the request to post a cooling
type ReqPostCooling struct {
	FinalTemp   float32 `form:"final_temp" json:"final_temp"`
	CoolingTime float32 `form:"cooling_time" json:"cooling_time"`
	Notes       string  `form:"notes" json:"notes"`
}
