package cooling

import (
	"brewday/internal/recipe"
	"brewday/internal/routers/common"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type CoolingRouter struct {
	Store        RecipeStore
	TLStore      TimelineStore
	SummaryStore SummaryStore
}

// addSummaryCooling adds a cooling to the summary and notes related to it
func (r *CoolingRouter) addSummaryCooling(id string, finalTemp, coolingTime float32, notes string) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddCooling(id, finalTemp, coolingTime, notes)
	}
	return nil
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
	err := r.Store.UpdateStatus(id, recipe.RecipeStatusCooling)
	if err != nil {
		return err
	}
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
	err = r.addSummaryCooling(id, req.FinalTemp, req.CoolingTime, req.Notes)
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add cooling to summary")
	}
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getPreFermentation", id))
}
