package summary

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
)

type SummaryRouter struct {
	Summary SummaryRecorder
}

// getSummary returns the summary
func (r *SummaryRouter) getSummary() string {
	if r.Summary != nil {
		return r.Summary.GetSummary()
	}
	return ""
}

// getExtention returns the extention
func (r *SummaryRouter) getExtention() string {
	if r.Summary != nil {
		return r.Summary.GetExtention()
	}
	return ""
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
		return errors.New("no recipe id provided")
	}
	summ := r.getSummary()
	ext := r.getExtention()
	fileName := id + "." + ext
	c.Response().Header().Set("Content-Type", "application/octet-stream")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	_, err := c.Response().Write([]byte(summ))
	return err
}
