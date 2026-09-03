package stats

import (
	"brewday/internal/summary"
	"errors"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/labstack/echo/v4"
)

type StatsRouter struct {
	StatsStore   StatsStore
	InfoProvider InfoProvider
}

func (r *StatsRouter) getStats(bs string) ([]StatEntry, error) {
	if r.StatsStore == nil {
		return nil, errors.New("summary store not configured")
	}
	rawStats, err := r.StatsStore.GetAllStats()
	if err != nil {
		return nil, err
	}
	res := []StatEntry{}
	for name, s := range rawStats {
		if s.BrewingSystem == bs {
			res = append(res, StatEntry{
				RecipeName:         name,
				Evaporation:        nullIf0(s.Evaporation),
				Efficiency:         nullIf0(s.Efficiency),
				FinishedTimeEpoch:  s.FinishedTime.Unix(),
				FinishedTimeString: s.FinishedTime.Format("2006-01-02"),
			})
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].FinishedTimeEpoch < res[j].FinishedTimeEpoch
	})
	return res, nil

}

func (r *StatsRouter) addStats(req *ReqPostAddStat) error {
	if r.StatsStore == nil {
		return errors.New("summary store not configured")
	}
	finished, err := time.Parse("2006-01-02", req.FinishedTimeString)
	if err != nil {
		return err
	}
	s := &summary.Statistics{
		Evaporation:  req.Evaporation,
		Efficiency:   req.Efficiency,
		FinishedTime: finished,
	}
	return r.StatsStore.AddStatsExternal(req.RecipeName, s)
}

// RegisterRoutes registers the routes for the stats router
func (r *StatsRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	stats := parent.Group("/stats")
	stats.GET("", r.getStatsHandler).Name = "getStats"
	stats.POST("/add", r.postAddExtStatHandler).Name = "postAddExtStat"
	stats.POST("/delete", r.deleteStatsHandler).Name = "deleteStats"
	stats.POST("/bs", r.setBsHandler).Name = "postStatsSetBS"
}

func (r *StatsRouter) getStatsHandler(c echo.Context) error {
	var s []StatEntry
	bs := c.QueryParam("bs")
	if bs != "" {
		stats, err := r.getStats(bs)
		if err != nil {
			return err
		}
		s = stats
	}
	return c.Render(200, "stats.html", map[string]any{
		"Title":          "Stats",
		"Subtitle":       "Historical stats from saved summaries",
		"Stats":          s,
		"BrewingSystems": r.InfoProvider.GetSystemNames(),
		"BS":             bs,
	})
}

func (r *StatsRouter) postAddExtStatHandler(c echo.Context) error {
	var req ReqPostAddStat
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	err = r.addStats(&req)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getStats"))
}

func (r *StatsRouter) deleteStatsHandler(c echo.Context) error {
	if r.StatsStore == nil {
		return errors.New("summary store not configured")
	}
	var req ReqPostDeleteStat
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	err = r.StatsStore.DeleteStats(req.RecipeTitle)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getStats"))
}

func (r *StatsRouter) setBsHandler(c echo.Context) error {
	var req ReqPostSetBs
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	redirect := "getStats"
	params := url.Values{}
	params.Add("bs", req.BrewingSystem)
	return c.Redirect(http.StatusFound, c.Echo().Reverse(redirect)+"?"+params.Encode())
}
