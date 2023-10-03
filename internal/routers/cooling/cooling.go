package cooling

import (
	"brewday/internal/recipe"
	"brewday/internal/routers/common"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CoolingRouter struct {
	Store   RecipeStore
	TL      Timeline
	Summary SummaryRecorder
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
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	re.SetStatus(recipe.RecipeStatusCooling)
	return c.Render(http.StatusOK, "cooling.html", map[string]interface{}{
		"Title":    "Cooling",
		"RecipeID": id,
	})
}

// postCoolingHandler handles the post request for the cooling page
func (r *CoolingRouter) postCoolingHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	var req ReqPostCooling
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	r.addSummaryCooling(req.FinalTemp, req.CoolingTime, req.Notes)
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getPreFermentation", id))
}
