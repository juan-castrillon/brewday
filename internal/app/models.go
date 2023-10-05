package app

import (
	"brewday/internal/recipe"
	"io"
	"io/fs"

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
	AddTimeline(recipeID string, timelineType string)
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
}

// SummaryRecorderStore is the interface that helps decouple the summary recorder store from the application
// It represents a store that stores summary recorders
type SummaryRecorderStore interface {
	// AddSummaryRecorder adds a summary recorder to the store
	AddSummaryRecorder(recipeID string, recorderType string)
	// AddMashTemp adds a mash temperature to the summary and notes related to it
	AddMashTemp(id string, temp float64, notes string) error
	// AddRast adds a rast to the summary and notes related to it
	AddRast(id string, temp float64, duration float64, notes string) error
	// AddLauternNotes adds lautern notes to the summary
	AddLaunternNotes(id string, notes string) error
	// AddHopping adds a hopping to the summary and notes related to it
	AddHopping(id string, name string, amount float32, alpha float32, duration float32, notes string) error
	// AddMeasuredVolume adds a measured volume to the summary
	AddMeasuredVolume(id string, name string, amount float32, notes string) error
	// AddEvaporation adds an evaporation to the summary
	AddEvaporation(id string, amount float32) error
	// AddCooling adds a cooling to the summary and notes related to it
	AddCooling(id string, finalTemp, coolingTime float32, notes string) error
	// AddSummaryPreFermentation adds a summary of the pre fermentation
	AddSummaryPreFermentation(id string, volume float32, sg float32, notes string) error
	// AddEfficiency adds the efficiency (sudhausausbeute) to the summary
	AddEfficiency(id string, efficiencyPercentage float32) error
	// AddYeastStart adds the yeast start to the summary
	AddYeastStart(id string, temperature, notes string) error
	// AddTimeline adds a timeline to the summary
	AddTimeline(id string, timeline []string) error
	// GetSummary returns the summary
	GetSummary(id string) (string, error)
	// GetExtension returns the extension of the summary
	GetExtension(id string) (string, error)
	// Close closes the summary recorder
	Close(id string) error
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
