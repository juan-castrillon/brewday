package app

import (
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

// ReqPostTimelineEvent represents the request body for the postTimelineEvent
type ReqPostTimelineEvent struct {
	Message string `json:"message" form:"message"`
}
