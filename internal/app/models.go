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
