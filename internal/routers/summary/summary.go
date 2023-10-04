package summary

import (
	"brewday/internal/routers/common"
	"fmt"

	"github.com/labstack/echo/v4"
)

type SummaryRouter struct {
	SummaryStore SummaryRecorderStore
	TL           Timeline
}

// getSummary returns the summary
func (r *SummaryRouter) getSummary(id string) (string, error) {
	if r.SummaryStore != nil {
		return r.SummaryStore.GetSummary(id)
	}
	return "", nil
}

// getExtention returns the extension
func (r *SummaryRouter) getExtension(id string) (string, error) {
	if r.SummaryStore != nil {
		return r.SummaryStore.GetExtension(id)
	}
	return "", nil
}

// addTimeline adds a timeline
func (r *SummaryRouter) addTimeline(id string, tl []string) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddTimeline(id, tl)
	}
	return nil
}

// closeSummary closes the summary
func (r *SummaryRouter) closeSummary(id string) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.Close(id)
	}
	return nil
}

// getTimeline returns the timeline
func (r *SummaryRouter) getTimeline() []string {
	if r.TL != nil {
		return r.TL.GetTimeline()
	}
	return []string{}
}

// RegisterRoutes registers the routes for the summary router
func (r *SummaryRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	summary := parent.Group("/summary")
	summary.GET("/:recipe_id", r.getSummaryHandler).Name = "getSummary"
}

// getSummaryHandler handles the GET /summary/:recipe_id route
func (r *SummaryRouter) getSummaryHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	tl := r.getTimeline()
	if len(tl) > 0 {
		err := r.addTimeline(id, tl)
		if err != nil {
			return err
		}
	}
	err := r.closeSummary(id)
	if err != nil {
		return err
	}
	summ, err := r.getSummary(id)
	if err != nil {
		return err
	}
	ext, err := r.getExtension(id)
	if err != nil {
		return err
	}
	fileName := id + "." + ext
	c.Response().Header().Set("Content-Type", "application/octet-stream")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	_, err = c.Response().Write([]byte(summ))
	return err
}
