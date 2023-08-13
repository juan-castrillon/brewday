package cooling

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CoolingRouter struct {
	TL      Timeline
	Summary SummaryRecorder
}

// addTimelineEvent adds an event to the timeline
func (r *CoolingRouter) addTimelineEvent(message string) {
	if r.TL != nil {
		r.TL.AddEvent(message)
	}
}

// addSummaryCooling adds a cooling to the summary and notes related to it
func (r *CoolingRouter) addSummaryCooling(finalTemp, coolingTime float32, notes string) {
	if r.Summary != nil {
		r.Summary.AddCooling(finalTemp, coolingTime, notes)
	}
}

// RegisterRoutes registers the routes for the cooling router
func (r *CoolingRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	cooling := parent.Group("/cooling")
	cooling.GET("/:recipe_id", r.getCoolingHandler).Name = "getCooling"
	cooling.POST("/:recipe_id", r.postCoolingHandler).Name = "postCooling"
}

// getCoolingHandler returns the handler for the cooling page
func (r *CoolingRouter) getCoolingHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return errors.New("no recipe id provided")
	}
	r.addTimelineEvent("Started Cooling")
	return c.Render(http.StatusOK, "cooling.html", map[string]interface{}{
		"Title":    "Cooling",
		"RecipeID": id,
	})
}

// postCoolingHandler handles the post request for the cooling page
func (r *CoolingRouter) postCoolingHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return errors.New("no recipe id provided")
	}
	var req ReqPostCooling
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	r.addSummaryCooling(req.FinalTemp, req.CoolingTime, req.Notes)
	r.addTimelineEvent("Finished Cooling")
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getPreFermentation", id))
}
