package stats

import "github.com/labstack/echo/v4"

type StatsRouter struct {
}

// RegisterRoutes registers the routes for the stats router
func (r *StatsRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	stats := parent.Group("/stats")
	stats.GET("", r.getStatsHandler).Name = "getStats"
}

func (r *StatsRouter) getStatsHandler(c echo.Context) error {
	return c.Render(200, "stats.html", map[string]any{
		"Title": "Stats",
		"Subtitle": "Historical stats from saved summaries",
	})
}
