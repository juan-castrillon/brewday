package app

import (
	"brewday/internal/recipe/mmum"
	"brewday/internal/render"
	"brewday/internal/routers/common"
	"brewday/internal/routers/cooling"
	"brewday/internal/routers/fermentation"
	"brewday/internal/routers/hopping"
	"brewday/internal/routers/import_recipe"
	"brewday/internal/routers/lautern"
	"brewday/internal/routers/mash"
	summary "brewday/internal/routers/summary"
	"brewday/internal/store/memory"
	"brewday/internal/summary_recorder/markdown"
	"brewday/internal/timeline/basic"
	"context"
	"io/fs"
	"math"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
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
func NewApp(staticFS fs.FS) (*App, error) {
	app := &App{
		staticFs: staticFS,
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
	// Register global middlewares
	a.server.Use(middleware.Recover())
	// Initialize internal components
	parser := mmum.MMUMParser{}
	store := memory.NewMemoryStore()
	summ := markdown.NewMarkdownSummaryRecorder()
	r := render.NewTemplateRenderer()
	tl := basic.NewBasicTimeline()
	a.renderer = r
	a.TL = tl
	// Register routers
	a.routers = []common.Router{
		&import_recipe.ImportRouter{
			Parser: &parser,
			Store:  store,
		},
		&mash.MashRouter{
			Store:   store,
			TL:      a.TL,
			Summary: summ,
		},
		&lautern.LauternRouter{
			Store:   store,
			TL:      a.TL,
			Summary: summ,
		},
		&hopping.HoppingRouter{
			Store:   store,
			TL:      a.TL,
			Summary: summ,
		},
		&cooling.CoolingRouter{
			TL:      a.TL,
			Summary: summ,
		},
		&fermentation.FermentationRouter{
			TL:      a.TL,
			Summary: summ,
			Store:   store,
		},
		&summary.SummaryRouter{
			Summary: summ,
			TL:      a.TL,
		},
	}
	a.RegisterStaticFiles()
	err := a.RegisterTemplates()
	if err != nil {
		return err
	}
	a.server.Pre(middleware.RemoveTrailingSlash())
	a.server.HTTPErrorHandler = a.customErrorHandler
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
	a.renderer.AddFunc("truncateFloat", func(f float32, decimals int) float64 {
		f64 := float64(f)
		return math.Round(f64*(math.Pow10(decimals))) / math.Pow10(decimals)
	})

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

// customErrorHandler is a custom error handler
func (a *App) customErrorHandler(err error, c echo.Context) {
	log.Error().Err(err).Msg(c.Request().RequestURI)
	notFound := strings.Contains(strings.ToLower(err.Error()), "not found")
	if err == common.ErrNoRecipeLoaded || err == common.ErrNoRecipeIDProvided || notFound {
		err2 := c.Render(404, "error_no_recipe_loaded.html", map[string]interface{}{
			"Title": "Error in recipe",
		})
		if err2 != nil {
			log.Error().Err(err2).Msg("error while rendering error page")
		}
	}
}
