package app

import (
	"brewday/internal/recipe"
	"brewday/internal/routers/common"
	"brewday/internal/routers/cooling"
	"brewday/internal/routers/fermentation"
	"brewday/internal/routers/hopping"
	"brewday/internal/routers/import_recipe"
	"brewday/internal/routers/lautern"
	"brewday/internal/routers/mash"
	"brewday/internal/routers/recipes"
	secondaryferm "brewday/internal/routers/secondary_ferm"
	summary "brewday/internal/routers/summary"
	"context"
	"io/fs"
	"math"
	"net/url"

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
	TLStore  TimelineStore
	notifier Notifier
}

// AppComponents is the structure that contains the external components of the application
type AppComponents struct {
	Renderer     Renderer
	TL           TimelineStore
	Notifier     Notifier
	Store        RecipeStore
	SummaryStore SummaryStore
}

// NewApp creates a new App
func NewApp(staticFS fs.FS, components *AppComponents) (*App, error) {
	app := &App{
		staticFs: staticFS,
	}
	err := app.Initialize(components)
	if err != nil {
		return nil, err
	}
	return app, nil
}

// Initialize initializes the application
func (a *App) Initialize(components *AppComponents) error {
	a.server = echo.New()
	// Register global middlewares
	a.server.Use(middleware.Recover())
	// Initialize internal components
	store := components.Store
	a.renderer = components.Renderer
	a.TLStore = components.TL
	a.notifier = components.Notifier
	ss := components.SummaryStore
	timer := common.NewTimer(store, a.TLStore, a.notifier)
	// Register routers
	a.routers = []common.Router{
		&import_recipe.ImportRouter{
			Store:                store,
			SummaryRecorderStore: ss,
			TLStore:              a.TLStore,
		},
		&mash.MashRouter{
			Store:        store,
			TLStore:      a.TLStore,
			SummaryStore: ss,
			Timer:        timer,
		},
		&lautern.LauternRouter{
			Store:        store,
			TLStore:      a.TLStore,
			SummaryStore: ss,
			Timer:        timer,
		},
		&hopping.HoppingRouter{
			Store:        store,
			TLStore:      a.TLStore,
			SummaryStore: ss,
			Timer:        timer,
		},
		&cooling.CoolingRouter{
			Store:        store,
			TLStore:      a.TLStore,
			SummaryStore: ss,
			Timer:        timer,
		},
		&fermentation.FermentationRouter{
			TLStore:      a.TLStore,
			SummaryStore: ss,
			Store:        store,
			Notifier:     a.notifier,
		},
		&secondaryferm.SecondaryFermentationRouter{
			TLStore:      a.TLStore,
			SummaryStore: ss,
			Store:        store,
			Notifier:     a.notifier,
		},
		&summary.SummaryRouter{
			SummaryStore: ss,
			TLStore:      a.TLStore,
		},
		&recipes.RecipesRouter{
			Store:        store,
			TLStore:      a.TLStore,
			SummaryStore: ss,
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
	a.renderer.AddFunc("recipeStatus", func(r *recipe.Recipe) string {
		return r.GetStatusString()
	})
	a.renderer.AddFunc("urlEncode", func(s string) string {
		return url.QueryEscape(s)
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
	a.server.POST("/timeline/:recipe_id", a.postTimelineEvent).Name = "postTimelineEvent"
}

// Run starts the application
func (a *App) Run(address string) error {
	return a.server.Start(address)
}

// Stop stops the application
func (a *App) Stop(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
