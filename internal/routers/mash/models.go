package mash

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

// TimelineStore represents a component that stores timelines
type TimelineStore interface {
	// AddEvent adds an event to the timeline
	AddEvent(id, message string) error
}

// SummaryStore represents a component that stores summaries
// The recipe id is used as key
type SummaryStore interface {
	AddMashTemp(id string, temp float32, notes string) error
	AddRast(id string, temp float32, duration float32, notes string) error
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

// ReqPostRasts represents the request body for the postRastsHandler
type ReqPostRasts struct {
	RealTemperature float32 `json:"real_temperature" form:"real_temp"`
	RealDuration    float32 `json:"real_duration" form:"real_duration"`
	Notes           string  `json:"notes" form:"notes"`
}

type ReqPostFirstRast struct {
	RealMashTemperature float32 `json:"real_mash_temperature" form:"real_mash_temp"`
	Notes               string  `json:"notes" form:"notes"`
}
