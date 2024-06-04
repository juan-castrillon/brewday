package mash

import (
	"brewday/internal/recipe"
	"brewday/internal/routers/common"
	"brewday/internal/tools"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type MashRouter struct {
	Store        RecipeStore
	TLStore      TimelineStore
	SummaryStore SummaryStore
	Timer        Timer
}

func (r *MashRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	mash := parent.Group("/mash")
	mash.GET("/start/:recipe_id", r.getMashStartHandler).Name = "getMashStart"
	mash.GET("/rasts/:recipe_id/:rast_num", r.getRastsHandler).Name = "getRasts"
	mash.POST("/rasts/:recipe_id/:rast_num", r.postRastsHandler).Name = "postRasts"
	mash.GET("/rasts/timer/:recipe_id/:rast_num", r.getRastTimestamp).Name = "getMashRastTimestamp"
	mash.POST("/rasts/timer/stop/:recipe_id/:rast_num", r.postRastStopTimer).Name = "postMashRastStopTimer"
	mash.GET("/rasts/timer/duration/:recipe_id/:rast_num", r.getRastRealDuration).Name = "getMashRastDuration"
}

// addTimelineEvent adds an event to the timeline
func (r *MashRouter) addTimelineEvent(id, message string) error {
	if r.TLStore != nil {
		return r.TLStore.AddEvent(id, message)
	}
	return nil
}

// addSummaryMashTemp adds a mash temperature to the summary and notes related to it
func (r *MashRouter) addSummaryMashTemp(id string, temp float32, notes string) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddMashTemp(id, temp, notes)
	}
	return nil
}

// addSummaryRast adds a rast to the summary and notes related to it
func (r *MashRouter) addSummaryRast(id string, temp float32, duration float32, notes string) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddRast(id, temp, duration, notes)
	}
	return nil
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
	err = r.Store.UpdateStatus(id, recipe.RecipeStatusMashing, "start")
	if err != nil {
		return err
	}
	err = r.addTimelineEvent(id, "Started mashing")
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add timeline event")
	}
	return c.Render(200, "mash_start.html", map[string]interface{}{
		"Title":        "Mash " + re.Name,
		"MainWater":    re.Mashing.MainWaterVolume,
		"MashTemp":     re.Mashing.MashTemperature,
		"NextRastTemp": re.Mashing.Rasts[0].Temperature,
		"RecipeID":     id,
	})
}

// getRastsHandler is the handler for the mash rasts page
func (r *MashRouter) getRastsHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	rastNumStr := c.Param("rast_num")
	if rastNumStr == "" {
		return errors.New("no rast number provided")
	}
	rastNum, err := strconv.Atoi(rastNumStr)
	if err != nil {
		return err
	}
	missing := re.Mashing.Rasts[rastNum+1:]
	missingDuration := float32(0.0)
	if len(missing) > 0 {
		for _, rast := range missing {
			missingDuration += rast.Duration
		}
	}
	nextRastNum := rastNum + 1
	err = r.Store.UpdateStatus(id, recipe.RecipeStatusMashing, "rast", tools.AnyToString(rastNum))
	if err != nil {
		return err
	}
	started, stopped, err := r.Timer.GetBoolFlags(id, "mashing_rast", rastNumStr)
	if err != nil {
		return err
	}
	return c.Render(200, "mash_rasts.html", map[string]interface{}{
		"Title":                "Mash " + re.Name,
		"Rast":                 re.Mashing.Rasts[rastNum],
		"RastNumber":           rastNum,
		"NextRast":             nextRastNum,
		"MissingRasts":         missing,
		"MissingRastsDuration": missingDuration,
		"Nachguss":             re.Mashing.Nachguss,
		"RecipeID":             id,
		"StartClickedOnce":     started,
		"Stopped":              stopped,
	})
}

// postRastsHandler is the handler for the mash rasts page
func (r *MashRouter) postRastsHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	rastNumStr := c.Param("rast_num")
	if rastNumStr == "" {
		return errors.New("no rast number provided")
	}
	rastNum, err := strconv.Atoi(rastNumStr)
	if err != nil {
		return err
	}
	switch rastNum {
	case 0:
		err = r.addTimelineEvent(id, "Finished Einmaischen")
		if err != nil {
			log.Error().Str("id", id).Err(err).Msg("could not add timeline event")
		}
		var req ReqPostFirstRast
		err := c.Bind(&req)
		if err != nil {
			return err
		}
		err = r.addSummaryMashTemp(id, req.RealMashTemperature, req.Notes)
		if err != nil {
			log.Error().Str("id", id).Err(err).Msg("could not add mash temp to summary")
		}
	case len(re.Mashing.Rasts):
		err = r.addTimelineEvent(id, "Finished mashing")
		if err != nil {
			log.Error().Str("id", id).Err(err).Msg("could not add timeline event")
		}
		return c.Redirect(302, c.Echo().Reverse("getLautern", id))
	default:
		var req ReqPostRasts
		err := c.Bind(&req)
		if err != nil {
			return err
		}
		err = r.addSummaryRast(id, req.RealTemperature, req.RealDuration, req.Notes)
		if err != nil {
			log.Error().Str("id", id).Err(err).Msg("could not add rast to summary")
		}
	}
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getRasts", id, rastNum))
}

// getRastTimestamp returns the final timestamp for the frontend timers for a rast
func (r *MashRouter) getRastTimestamp(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	rastNumStr := c.Param("rast_num")
	if rastNumStr == "" {
		return errors.New("no rast number provided")
	}
	rastNum, err := strconv.Atoi(rastNumStr)
	if err != nil {
		return err
	}
	rast := re.Mashing.Rasts[rastNum]
	duration := time.Duration(rast.Duration * float32(time.Minute))
	return r.Timer.HandleStartTimer(c, id, duration, "mashing_rast", rastNumStr)
}

// postRastStopTimer handles the post request when a timer stops
func (r *MashRouter) postRastStopTimer(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	rastNumStr := c.Param("rast_num")
	if rastNumStr == "" {
		return errors.New("no rast number provided")
	}
	return r.Timer.HandleStopTimer(c, id, "mashing_rast", rastNumStr)
}

// getRastRealDuration handles the get request to send the real duration of a rast
func (r *MashRouter) getRastRealDuration(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	rastNumStr := c.Param("rast_num")
	if rastNumStr == "" {
		return errors.New("no rast number provided")
	}
	return r.Timer.HandleRealDuration(c, id, "mashing_rast", rastNumStr)
}
