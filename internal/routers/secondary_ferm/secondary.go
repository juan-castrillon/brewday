package secondaryferm

import (
	"brewday/internal/recipe"
	"brewday/internal/routers/common"
	"brewday/internal/tools"
	"brewday/internal/watcher"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type SecondaryFermentationRouter struct {
	TLStore         TimelineStore
	SummaryStore    SummaryStore
	Store           RecipeStore
	StatsStore      StatsStore
	Notifier        Notifier
	ingredientCache ingredientCache
	watchersSet     map[string]bool // This keeps track if watches are set. In case of restart, it will go back to nil and force reconfig of watchers
}

// RegisterRoutes adds routes to the web server
// It receives the root web server and a parent group
func (r *SecondaryFermentationRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	sf := parent.Group("/secondary_fermentation")
	sf.GET("/dry_hop/:recipe_id", r.getDryHopHandler).Name = "getDryHop"
	sf.POST("/dry_hop/:recipe_id", r.postDryHopInHandler).Name = "postDryHopIn"
	// sf.GET("/dry_hop/start/:recipe_id/load", r.getDryHopStartLoadHandler).Name = "getDryHopStartLoad"
	// sf.GET("/dry_hop/start/:recipe_id", r.getDryHopStartHandler).Name = "getDryHopStart"
	// sf.POST("/dry_hop/start/:recipe_id", r.postDryHopStartHandler).Name = "postDryHopStart"
	// sf.GET("/dry_hop/confirm/:recipe_id", r.getDryHopConfirmHandler).Name = "getDryHopConfirm"
	// sf.POST("/dry_hop/confirm/:recipe_id", r.postDryHopConfirmHandler).Name = "postDryHopConfirm"
	sf.GET("/pre_bottle/:recipe_id", r.getPreBottleHandler).Name = "getPreBottle"
	sf.POST("/pre_bottle/:recipe_id", r.postPreBottleHandler).Name = "postPreBottle"
	sf.GET("/bottle/:recipe_id", r.getBottleHandler).Name = "getBottle"
	sf.POST("/bottle/:recipe_id", r.postBottleHandler).Name = "postBottle"
	sf.GET("/start/:recipe_id", r.getSecondaryFermentationStartHandler).Name = "getSecondaryFermentationStart"
	sf.POST("/start/:recipe_id", r.postSecondaryFermentationStartHandler).Name = "postSecondaryFermentationStart"
	sf.GET("/end/:recipe_id", r.getSecondaryFermentationEndHandler).Name = "getSecondaryFermentationEnd"
	sf.POST("/end/:recipe_id", r.postSecondaryFermentationEndHandler).Name = "postSecondaryFermentationEnd"
	root.GET("/end/:recipe_id", r.getEndHandler).Name = "getEnd"
}

func (r *SecondaryFermentationRouter) getDryHopHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	err = r.Store.UpdateStatus(id, recipe.RecipeStatusFermenting, "dry_hop")
	if err != nil {
		return err
	}
	ings := r.ingredientCache.getIngredients(id, re)
	if len(ings) == 0 {
		log.Info().Str("id", id).Err(err).Msg("Recipe has no dry hops")
		return c.Redirect(http.StatusFound, c.Echo().Reverse("getPreBottle", id))
	}
	err = r.Store.AddBoolFlag(id, "has_dry_hops", true)
	if err != nil {
		return err
	}
	for i, ing := range ings {
		started, err := r.Store.RetrieveBoolFlag(id, "secondary_dry_hop_started_"+ing.SanitizedName)
		if err != nil {
			return err
		}
		ings[i].StartClickedOnce = started
		startedDates, err := r.Store.RetrieveDates(id, "secondary_dry_hop_"+ing.SanitizedName)
		if err != nil {
			return err
		}
		if len(startedDates) == 0 {
			ings[i].TimeElapsed = 0
		} else {
			startDate := *startedDates[0]
			since := time.Since(startDate).Hours()
			ings[i].TimeElapsed = float32(since)
		}
	}
	return c.Render(http.StatusOK, "secondary_dry_hop.html", map[string]interface{}{
		"Title":       "Dry Hopping",
		"Subtitle":    "Add hops and other ingredients into the fermented wort",
		"RecipeID":    id,
		"Ingredients": ings,
	})
}

func (r *SecondaryFermentationRouter) postDryHopInHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	var req ReqPostDryHopIn
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	now := time.Now()
	err = r.Store.AddDate(id, &now, "secondary_dry_hop_"+req.IngredientName)
	if err != nil {
		return err
	}
	err = r.Store.AddBoolFlag(id, "secondary_dry_hop_started_"+req.IngredientName, true)
	if err != nil {
		return err
	}
	err = r.TLStore.AddEvent(id, "Added dry ingredient "+req.IngredientName)
	if err != nil {
		return err
	}
	err = r.SummaryStore.AddDryHopStart(id, req.IngredientName, req.RealAmount, req.RealAlpha, "") //TODO: Support notes here?
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

// getPreBottleHandler is the handler for the pre bottle page
// This page will ask for the volume and the type of priming sugar
func (r *SecondaryFermentationRouter) getPreBottleHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	err := r.Store.UpdateStatus(id, recipe.RecipeStatusFermenting, "pre_bottle")
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "secondary_pre_bottle.html", map[string]interface{}{
		"Title":    "Bottle",
		"Subtitle": "Set the volume and the type of priming sugar",
		"RecipeID": id,
	})
}

// postPreBottleHandler is the handler for the pre bottle page
func (r *SecondaryFermentationRouter) postPreBottleHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	var req ReqPostPreBottle
	err = c.Bind(&req)
	if err != nil {
		return err
	}
	res, err := r.Store.RetrieveResults(id)
	if err != nil {
		return err
	}
	vol := req.Volume - req.LostVolume
	err = r.calculateSugar(
		id, vol,
		re.Fermentation.Carbonation, req.Temperature,
		res.Alcohol,
		req.SugarType,
	)
	if err != nil {
		return err
	}
	redirect := "getBottle"
	queryParams := fmt.Sprintf("?type=%s", req.SugarType)
	err = r.addSummaryPreBottle(id, req.Volume)
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add summary pre bottle")
	}
	return c.Redirect(http.StatusFound, c.Echo().Reverse(redirect, id)+queryParams)
}

// getBottleHandler is the handler for the bottle page
func (r *SecondaryFermentationRouter) getBottleHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	t := c.QueryParam("type")
	err := r.Store.UpdateStatus(id, recipe.RecipeStatusFermenting, "bottle", t)
	if err != nil {
		return err
	}
	sugarResults, err := r.Store.RetrieveSugarResults(id)
	if err != nil {
		return err
	}
	s := sugarResults[0]
	return c.Render(http.StatusOK, "secondary_bottle.html", map[string]interface{}{
		"Title":        "Bottle",
		"Subtitle":     "Create the sugar solution",
		"RecipeID":     id,
		"SugarResults": sugarResults,
		"Sugar":        s.Amount,
		"SugarType":    t,
	})
}

// postBottleHandler is the handler for the bottle page
func (r *SecondaryFermentationRouter) postBottleHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	var req ReqPostBottle
	err = c.Bind(&req)
	if err != nil {
		return err
	}
	res, err := r.Store.RetrieveResults(id)
	if err != nil {
		return err
	}
	volumeBeforeSugar := req.RealVolume - req.Water
	_, realAlcohol := tools.SugarForCarbonation(
		volumeBeforeSugar, re.Fermentation.Carbonation, req.Temperature,
		res.Alcohol, req.Water, req.SugarType,
	)
	realCO2 := tools.CarbonationForSugar(req.RealVolume, req.SugarAmount, req.Temperature, req.SugarType)
	err = r.addSummaryBottle(id, realCO2, realAlcohol, req.SugarAmount, req.Temperature, req.RealVolume, req.SugarType, req.Notes)
	// TODO: Add the bottling time now that i have it, add also to the stats
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add summary bottle")
	}
	err = r.Store.UpdateResult(id, recipe.ResultAlcohol, realAlcohol)
	if err != nil {
		return err
	}
	r.addTimelineEvent(id, "Bottled")
	// Now, we calculate the duration of the dry hops based on the time now minus the time it took to bottle
	// we assume when bottling starts, dry hopping ends
	hasDryHops, err := r.Store.RetrieveBoolFlag(id, "has_dry_hops")
	if err != nil {
		return err
	}
	if hasDryHops {
		re, err := r.Store.Retrieve(id)
		if err != nil {
			return err
		}
		ings := r.ingredientCache.getIngredients(id, re)
		for _, ing := range ings {
			var since float32
			startedDates, err := r.Store.RetrieveDates(id, "secondary_dry_hop_"+ing.SanitizedName)
			if err != nil {
				return err
			}
			if len(startedDates) == 0 {
				since = 0
			} else {
				startDate := *startedDates[0]
				since = float32(time.Since(startDate).Minutes())
			}
			realDuration := since - req.Time
			err = r.SummaryStore.AddDryHopEnd(id, ing.SanitizedName, realDuration/60)
			if err != nil {
				return err
			}
		}
	}
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getSecondaryFermentationStart", id))
}

// getSecondaryFermentationStartHandler is the handler for the secondary fermentation start page
func (r *SecondaryFermentationRouter) getSecondaryFermentationStartHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	err := r.Store.UpdateStatus(id, recipe.RecipeStatusFermenting, "start_secondary")
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "secondary_start.html", map[string]interface{}{
		"Title":    "Start Secondary Fermentation",
		"RecipeID": id,
		"Subtitle": "First, let the bottles at warm temperature",
		"MinDays":  5, //TODO: make configurable
	})
}

// postSecondaryFermentationStartHandler is the handler for the secondary fermentation start page
func (r *SecondaryFermentationRouter) postSecondaryFermentationStartHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	var req ReqPostSecondaryStart
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	now := time.Now()
	var notificationDate time.Time
	switch req.TimeUnit {
	case "days":
		notificationDate = now.AddDate(0, 0, req.NotificationTime)
	case "seconds": // This is mainly for testing
		notificationDate = now.Add(time.Duration(req.NotificationTime) * time.Second)
	default:
		return fmt.Errorf("unknown time unit %s", req.TimeUnit)
	}
	w := watcher.NewWatcher(notificationDate, func() error {
		log.Info().Msgf("secondary fermentation notification triggered")
		return r.sendNotification("Time to put bottles in the fridge", "Secondary Fermentation", nil)
	})
	w.Start()
	err = r.Store.AddDate(id, &notificationDate, "secondary_ferm_notification")
	if err != nil {
		return err
	}
	r.addTimelineEvent(id, "Secondary Fermentation Started")
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getSecondaryFermentationEnd", id))
}

// getSecondaryFermentationEndHandler handles serving the end page for the secondary fermentation
func (r *SecondaryFermentationRouter) getSecondaryFermentationEndHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	err := r.checkWatchers(id) // TODO: maybe move this and fermentation to its own function/handler so it can be called on app start
	if err != nil {
		return err
	}
	notDates, err := r.Store.RetrieveDates(id, "secondary_ferm_notification")
	if err != nil {
		return err
	}
	missing := time.Until(*notDates[0])
	if missing > 0 {
		// Not finished yet
		err := r.Store.UpdateStatus(id, recipe.RecipeStatusFermenting, "wait_secondary")
		if err != nil {
			return err
		}
		return c.Render(http.StatusOK, "secondary_wait.html", map[string]interface{}{
			"Title":    "Wait for Secondary Fermentation",
			"RecipeID": id,
			"Subtitle": "Let the bottles at warm temperature",
			"Missing":  missing.String(),
		})
	} else {
		err := r.Store.UpdateStatus(id, recipe.RecipeStatusFermenting, "end_secondary")
		if err != nil {
			return err
		}
		return c.Render(http.StatusOK, "secondary_notes.html", map[string]interface{}{
			"Title":    "End Secondary Fermentation",
			"RecipeID": id,
			"Subtitle": "Enter your notes",
		})
	}
}

// postSecondaryFermentationEndHandler is the handler for the secondary fermentation end page
func (r *SecondaryFermentationRouter) postSecondaryFermentationEndHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	var req ReqPostSecondaryEnd
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	r.addTimelineEvent(id, "Secondary Fermentation Ended")
	err = r.addSummarySecondaryFermentation(id, req.Days, req.Notes)
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add summary secondary fermentation")
	}
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getEnd", id))
}

// getEndHandler is the handler for the end page
func (r *SecondaryFermentationRouter) getEndHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	err := r.addTimelineEvent(id, "Finished Day")
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add timeline event")
	}
	err = r.Store.UpdateStatus(id, recipe.RecipeStatusFinished)
	if err != nil {
		return err
	}
	err = r.addSummaryFinishedTime(id, time.Now())
	if err != nil {
		return err
	}
	err = r.addStats(id)
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "finished_day.html", map[string]interface{}{
		"Title":    "End Fermentation",
		"RecipeID": id,
		"Subtitle": "Congratulations, you've finished the brew!",
	})
}
