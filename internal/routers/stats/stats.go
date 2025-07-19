package stats

import "github.com/labstack/echo/v4"

type StatsRouter struct {
	SummaryStore SummaryStore
}

func (r *StatsRouter) getStats() ([]*RecipeStats, error) {
	// ids, err := r.SummaryStore.GetAllSummaries()
	// if err != nil {
	// 	return nil, err
	// }
	// res := []*RecipeStats{}
	// for _, id := range ids {
	// 	s, err := r.SummaryStore.GetSummary(id)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	res = append(res, &RecipeStats{
	// 		Evaporation: s.Statistics.Evaporation,
	// 		Efficiency:  s.Statistics.Efficiency,
	// 	})
	// }
	res := []*RecipeStats{{
		Evaporation: 30,
		Efficiency:  65,
	}, {
		Evaporation: 20,
		Efficiency:  72,
	}}
	return res, nil
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
