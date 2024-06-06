package cooling

import (
	"brewday/internal/recipe"
	"brewday/internal/routers/common"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type CoolingRouter struct {
	Store        RecipeStore
	TLStore      TimelineStore
	SummaryStore SummaryStore
	Timer        Timer
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
	cooling.GET("/timer/:recipe_id", r.getCoolingTimestamp).Name = "getCoolingTimestamp"
	cooling.POST("/timer/stop/:recipe_id", r.postCoolingStopTimer).Name = "postCoolingStopTimer"
	cooling.GET("/timer/duration/:recipe_id", r.getCoolingDuration).Name = "getCoolingDuration"
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
	started, stopped, err := r.Timer.GetBoolFlags(id, "cooling")
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "cooling.html", map[string]interface{}{
		"Title":            "Cooling",
		"RecipeID":         id,
		"StartClickedOnce": started,
		"Stopped":          stopped,
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

func (r *CoolingRouter) getCoolingTimestamp(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	started, _, err := r.Timer.GetBoolFlags(id, "cooling")
	if err != nil {
		return err
	}
	if !started {
		r.TLStore.AddEvent(id, "Started Cooling")
	}
	dur := 48 * time.Hour // Very long time to have a timer, not a countdown. Any normal wort will be cooled in two days :)
	return r.Timer.HandleStartTimer(c, id, dur, "cooling")
}

func (r *CoolingRouter) postCoolingStopTimer(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	tlEvent := "Stopped cooling"
	return r.Timer.HandleStopTimer(c, id, tlEvent, "", "", "cooling") //Notification will not be send as this will only stop manually
}

func (r *CoolingRouter) getCoolingDuration(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	return r.Timer.HandleRealDuration(c, id, "cooling")
}
