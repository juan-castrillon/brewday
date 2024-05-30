package app

import (
	"brewday/internal/recipe"
	"brewday/internal/summary"
	"io"
	"io/fs"
	"time"

	"github.com/labstack/echo/v4"
)

// Renderer is the interface that helps decouple the renderer from the application
// It is used to render the templates and should implement the echo.Renderer interface
type Renderer interface {
	// Render renders a template document. It is the implementation of echo.Renderer
	Render(w io.Writer, name string, data any, c echo.Context) error
	// RegisterTemplates registers the templates based on a file system
	RegisterTemplates(fs fs.FS) error
	// AddFunc adds a function to the template
	AddFunc(name string, fn any)
}

// TimelineStore represents a component that stores timelines
type TimelineStore interface {
	// AddEvent adds an event to the timeline
	AddEvent(id, message string) error
	// GetTimeline returns a timeline of events
	GetTimeline(id string) ([]string, error)
	// AddTimeline adds a timeline to the store
	AddTimeline(recipeID string) error
	// DeleteTimeline deletes the timeline for the given recipe id
	DeleteTimeline(recipeID string) error
}

// Notifier is the interface that helps decouple the notifier from the application
type Notifier interface {
	// Send sends a notification
	Send(message, title string, opts map[string]any) error
}

// RecipeStore is the interface that helps decouple the recipe store from the application
// It represents a store that stores recipes
type RecipeStore interface {
	// Store stores a recipe and returns an identifier that can be used to retrieve it
	Store(recipe *recipe.Recipe) (string, error)
	// Retrieve retrieves a recipe based on an identifier
	Retrieve(id string) (*recipe.Recipe, error)
	// List lists all the recipes
	List() ([]*recipe.Recipe, error)
	// Delete deletes a recipe based on an identifier
	Delete(id string) error
	// UpdateStatus updates the status of a recipe in the store
	UpdateStatus(id string, status recipe.RecipeStatus, statusParams ...string) error
	// UpdateResult updates a certain result of a recipe
	UpdateResult(id string, resultType recipe.ResultType, value float32) error
	// RetrieveResult gets a certain result value from a recipe
	RetrieveResult(id string, resultType recipe.ResultType) (float32, error)
	// RetrieveResults gets the results from a certain recipe
	RetrieveResults(id string) (*recipe.RecipeResults, error)
	// AddMainFermSG adds a new specific gravity measurement to a given recipe
	AddMainFermSG(id string, m *recipe.SGMeasurement) error
	// RetrieveMainFermSGs returns all measured sgs for a recipe
	RetrieveMainFermSGs(id string) ([]*recipe.SGMeasurement, error)
	// AddDate allows to store a date with a certain purpose. It can be used to store notification dates, or timers
	AddDate(id string, date *time.Time, name string) error
	// RetrieveDates allows to retreive stored dates with its purpose (name).It can be used to store notification dates, or timers
	// It supports pattern in the name to retrieve multiple values
	RetrieveDates(id, namePattern string) ([]*time.Time, error)
	// AddSugarResult adds a new priming sugar result to a given recipe
	AddSugarResult(id string, r *recipe.PrimingSugarResult) error
	// RetrieveSugarResults returns all sugar results for a recipe
	RetrieveSugarResults(id string) ([]*recipe.PrimingSugarResult, error)
}

// SummaryStore is the interface that helps decouple the summary store from the application
// It represents a store that stores summaries
type SummaryStore interface {
	AddSummary(recipeID, title string) error
	DeleteSummary(recipeID string) error
	AddMashTemp(id string, temp float32, notes string) error
	AddRast(id string, temp float32, duration float32, notes string) error
	AddLauternNotes(id, notes string) error
	AddHopping(id string, name string, amount float32, alpha float32, duration float32, notes string) error
	AddVolumeBeforeBoil(id string, amount float32, notes string) error
	AddVolumeAfterBoil(id string, amount float32, notes string) error
	AddCooling(id string, finalTemp, coolingTime float32, notes string) error
	AddPreFermentationVolume(id string, volume float32, sg float32, notes string) error
	AddYeastStart(id string, temperature, notes string) error
	AddMainFermentationSGMeasurement(id string, date string, gravity float32, final bool, notes string) error
	AddMainFermentationAlcohol(id string, alcohol float32) error
	AddMainFermentationDryHop(id string, name string, amount, alpha, duration float32, notes string) error
	AddPreBottlingVolume(id string, volume float32) error
	AddBottling(id string, carbonation, alcohol, sugar, temp, vol float32, sugarType, notes string) error
	AddSummarySecondary(id string, days int, notes string) error
	AddEvaporation(id string, amount float32) error
	AddEfficiency(id string, efficiencyPercentage float32) error
	GetSummary(id string) (*summary.Summary, error)
}

// ReqPostTimelineEvent represents the request body for the postTimelineEvent
type ReqPostTimelineEvent struct {
	Message string `json:"message" form:"message"`
}

// ReqPostNotification represents the request body for the postNotification
type ReqPostNotification struct {
	Message string                 `json:"message" form:"message"`
	Title   string                 `json:"title" form:"title"`
	Options map[string]interface{} `json:"options" form:"options"`
}
