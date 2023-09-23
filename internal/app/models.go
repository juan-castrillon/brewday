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

// Timeline represents a timeline of events
type Timeline interface {
	// AddEvent adds an event to the timeline
	AddEvent(message string)
	// GetTimeline returns a timeline of events
	GetTimeline() []string
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
}

// SummaryRecorder is the interface that helps decouple the summary recorder from the application
// It represents a recorder that records summaries
type SummaryRecorder interface {
	// AddMashTemp adds a mash temperature to the summary and notes related to it
	AddMashTemp(temp float64, notes string)
	// AddRast adds a rast to the summary and notes related to it
	AddRast(temp float64, duration float64, notes string)
	// AddLauternNotes adds lautern notes to the summary
	AddLaunternNotes(notes string)
	// AddHopping adds a hopping to the summary and notes related to it
	AddHopping(name string, amount float32, alpha float32, duration float32, notes string)
	// AddMeasuredVolume adds a measured volume to the summary
	AddMeasuredVolume(name string, amount float32, notes string)
	// AddEvaporation adds an evaporation to the summary
	AddEvaporation(amount float32)
	// AddCooling adds a cooling to the summary and notes related to it
	AddCooling(finalTemp, coolingTime float32, notes string)
	// AddSummaryPreFermentation adds a summary of the pre fermentation
	AddSummaryPreFermentation(volume float32, sg float32, notes string)
	// AddEfficiency adds the efficiency (sudhausausbeute) to the summary
	AddEfficiency(efficiencyPercentage float32)
	// AddYeastStart adds the yeast start to the summary
	AddYeastStart(temperature, notes string)
	// AddTimeline adds a timeline to the summary
	AddTimeline(timeline []string)
	// GetSummary returns the summary
	GetSummary() string
	// GetExtention returns the extension of the summary
	GetExtention() string
	// Close closes the summary recorder
	Close()
}

// ReqPostTimelineEvent represents the request body for the postTimelineEvent
type ReqPostTimelineEvent struct {
	Message string `json:"message" form:"message"`
}
