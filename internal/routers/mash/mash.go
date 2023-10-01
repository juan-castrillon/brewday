package mash

import (
	"brewday/internal/recipe"
	"brewday/internal/routers/common"
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"
)

type MashRouter struct {
	Store   RecipeStore
	TL      Timeline
	Summary SummaryRecorder
	recipe  *recipe.Recipe
}

func (r *MashRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	mash := parent.Group("/mash")
	mash.GET("/start/:recipe_id", r.getMashStartHandler).Name = "getMashStart"
	mash.POST("/rasts/:recipe_id/:rast_num", r.postRastsHandler).Name = "postRasts"
}

// addTimelineEvent adds an event to the timeline
func (r *MashRouter) addTimelineEvent(message string) {
	if r.TL != nil {
		r.TL.AddEvent(message)
	}
}

// addSummaryMashTemp adds a mash temperature to the summary and notes related to it
func (r *MashRouter) addSummaryMashTemp(temp float64, notes string) {
	if r.Summary != nil {
		r.Summary.AddMashTemp(temp, notes)
	}
}

// addSummaryRast adds a rast to the summary and notes related to it
func (r *MashRouter) addSummaryRast(temp float64, duration float64, notes string) {
	if r.Summary != nil {
		r.Summary.AddRast(temp, duration, notes)
	}
}

// getMashStartHandler is the handler for the mash start page
func (r *MashRouter) getMashStartHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	r.recipe = re
	r.addTimelineEvent("Started mashing")
	return c.Render(200, "mash_start.html", map[string]interface{}{
		"Title":        "Mash " + re.Name,
		"MainWater":    re.Mashing.MainWaterVolume,
		"MashTemp":     re.Mashing.MashTemperature,
		"NextRastTemp": re.Mashing.Rasts[0].Temperature,
		"RecipeID":     id,
	})
}

// postRastsHandler is the handler for the mash rasts page
func (r *MashRouter) postRastsHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	if r.recipe == nil {
		return common.ErrNoRecipeLoaded
	}
	rastNumStr := c.Param("rast_num")
	if rastNumStr == "" {
		return errors.New("no rast number provided")
	}
	rastNum, err := strconv.Atoi(rastNumStr)
	if err != nil {
		return err
	}
	var nextRastNum int
	switch rastNum {
	case 0:
		r.addTimelineEvent("Finished Einmaischen")
		var req ReqPostFirstRast
		err := c.Bind(&req)
		if err != nil {
			return err
		}
		r.addSummaryMashTemp(req.RealMashTemperature, req.Notes)
		nextRastNum = 1
	case len(r.recipe.Mashing.Rasts):
		r.addTimelineEvent("Finished mashing")
		return c.Redirect(302, c.Echo().Reverse("getLautern", id))
	default:
		var req ReqPostRasts
		err := c.Bind(&req)
		if err != nil {
			return err
		}
		r.addSummaryRast(req.RealTemperature, req.RealDuration, req.Notes)
		nextRastNum = rastNum + 1
	}
	missing := r.recipe.Mashing.Rasts[rastNum+1:]
	missingDuration := float32(0.0)
	if len(missing) > 0 {
		for _, rast := range missing {
			missingDuration += rast.Duration
		}
	}
	return c.Render(200, "mash_rasts.html", map[string]interface{}{
		"Title":                "Mash " + r.recipe.Name,
		"Rast":                 r.recipe.Mashing.Rasts[rastNum],
		"RastNumber":           rastNum,
		"NextRast":             nextRastNum,
		"MissingRasts":         missing,
		"MissingRastsDuration": missingDuration,
		"Nachguss":             r.recipe.Mashing.Nachguss,
		"RecipeID":             id,
	})
}
