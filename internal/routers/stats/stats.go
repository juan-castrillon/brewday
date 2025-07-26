package stats

import (
	"brewday/internal/summary"
	"brewday/internal/tools"
	"time"

	"github.com/labstack/echo/v4"
)

type StatsRouter struct {
	StatsStore StatsStore
}

func (r *StatsRouter) getStats() ([]StatEntry, error) {
	// if r.StatsStore != nil {
	// 	return r.StatsStore.GetAllStats()
	// }
	// return nil, errors.New("summary store not configured")
	rawStats := map[string]*summary.Statistics{
		"UmVjaXBlMQ==": {
			Evaporation:  16,
			Efficiency:   72,
			FinishedTime: time.Unix(150, 0),
		},
		"UmVjaXBlMg==": {
			Evaporation:  20,
			Efficiency:   60,
			FinishedTime: time.Unix(200000, 0),
		},
	}
	res := []StatEntry{}
	for rb64, s := range rawStats {
		name, err := tools.B64Decode(rb64)
		if err != nil {
			return nil, err
		}
		res = append(res, StatEntry{
			RecipeName:         name,
			Evaporation:        s.Evaporation,
			Efficiency:         s.Efficiency,
			FinishedTimeEpoch:  s.FinishedTime.Unix(),
			FinishedTimeString: s.FinishedTime.Format("2006-01-02"),
		})
	}
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
