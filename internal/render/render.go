package render

import (
	"html/template"
	"io"
	"io/fs"

	"github.com/labstack/echo/v4"
)

// TemplateRenderer is the renderer for html/template templates
type TemplateRenderer struct {
	templates *template.Template
	funcMap   map[string]any
}

// NewTemplateRenderer creates a new TemplateRenderer
func NewTemplateRenderer() *TemplateRenderer {
	return &TemplateRenderer{
		funcMap: make(map[string]any),
	}
}

// Render renders a template document. It is the implementation of echo.Renderer
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// RegisterTemplates registers the templates based on a file system
func (t *TemplateRenderer) RegisterTemplates(fs fs.FS) error {
	tmp, err := template.New("").Funcs(t.funcMap).ParseFS(fs, "*.html")
	//tmp, err := t.templates.ParseFS(fs, "*.html")
	if err != nil {
		return err
	}
	t.templates = tmp
	return nil
}

// AddFunc adds a function to the template
// All functions must be added before the templates are parsed (registered)
func (t *TemplateRenderer) AddFunc(name string, fn any) {
	t.funcMap[name] = fn
}
