package stats

import (
	"brewday/internal/summary"
	"errors"

	"github.com/labstack/echo/v4"
)

type StatsRouter struct {
	StatsStore StatsStore
}

func (r *StatsRouter) getStats() (map[string]*summary.Statistics, error) {
	if r.StatsStore != nil {
		return r.StatsStore.GetAllStats()
	}
	return nil, errors.New("summary store not configured")
}

// RegisterRoutes registers the routes for the stats router
func (r *StatsRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	stats := parent.Group("/stats")
	stats.GET("", r.getStatsHandler).Name = "getStats"
}

func (r *StatsRouter) getStatsHandler(c echo.Context) error {
	s, err := r.getStats()
	if err != nil {
		return err
	}
	return c.Render(200, "stats.html", map[string]any{
		"Title":    "Stats",
		"Subtitle": "Historical stats from saved summaries",
		"Stats":    s,
	})
}
