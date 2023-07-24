package app

import (
	"brewday/internal/recipe/mmum"
	"brewday/internal/routers/common"
	"brewday/internal/routers/import_recipe"
	"brewday/internal/routers/lautern"
	"brewday/internal/routers/mash"
	"brewday/internal/store/memory"
	"context"
	"io/fs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	// StaticFilesPath is the path to the static files
	StaticFilesPath = "/static"
)

// App is the application structure
// It encapsulates the web server, database, and other components
type App struct {
	server   *echo.Echo
	staticFs fs.FS
	routers  []common.Router
	renderer Renderer
	TL       Timeline
}

// NewApp creates a new App
func NewApp(staticFS fs.FS, renderer Renderer) (*App, error) {
	app := &App{
		staticFs: staticFS,
		renderer: renderer,
	}
	err := app.Initialize()
	if err != nil {
		return nil, err
	}
	return app, nil
}

// Initialize initializes the application
func (a *App) Initialize() error {
	a.server = echo.New()
	parser := mmum.MMUMParser{}
	store := memory.NewMemoryStore()
	a.routers = []common.Router{
		&import_recipe.ImportRouter{
			Parser: &parser,
			Store:  store,
		},
		&mash.MashRouter{
			Store: store,
			TL:    a.TL,
		},
		&lautern.LauternRouter{
			Store: store,
			TL:    a.TL,
		},
	}
	a.RegisterStaticFiles()
	err := a.RegisterTemplates()
	if err != nil {
		return err
	}
	a.server.Pre(middleware.RemoveTrailingSlash())
	a.RegisterRoutes()
	return nil
}

// RegisterStaticFiles registers the static files of the application
func (a *App) RegisterStaticFiles() {
	fs := echo.MustSubFS(a.staticFs, "web/static")
	a.server.StaticFS(StaticFilesPath, fs)
}

// RegisterTemplates registers the templates of the application
func (a *App) RegisterTemplates() error {
	a.renderer.AddFunc("static", func(path string) string {
		return StaticFilesPath + "/" + path
	})
	a.renderer.AddFunc("reverse", a.server.Reverse)

	fs := echo.MustSubFS(a.staticFs, "web/template")
	err := a.renderer.RegisterTemplates(fs)
	if err != nil {
		return err
	}
	a.server.Renderer = a.renderer
	return nil
}

// RegisterRoutes registers the routes of the application
func (a *App) RegisterRoutes() {
	group := a.server.Group("")
	for _, router := range a.routers {
		router.RegisterRoutes(a.server, group)
	}
	a.server.GET("/", func(c echo.Context) error {
		return c.Redirect(302, a.server.Reverse("getImport"))
	})
	a.server.POST("/timeline", a.postTimelineEvent).Name = "postTimelineEvent"
}

// Run starts the application
func (a *App) Run(address string) error {
	return a.server.Start(address)
}

// Stop stops the application
func (a *App) Stop(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

// addTimelineEvent adds an event to the timeline
func (a *App) addTimelineEvent(message string) {
	if a.TL != nil {
		a.TL.AddEvent(message)
	}
}

// postTimelineEvent is the handler for sent timeline events
func (a *App) postTimelineEvent(c echo.Context) error {
	var req ReqPostTimelineEvent
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	a.addTimelineEvent(req.Message)
	return c.NoContent(200)
}
