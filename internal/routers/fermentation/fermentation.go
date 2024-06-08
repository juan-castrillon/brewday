package fermentation

import (
	"brewday/internal/recipe"
	"brewday/internal/routers/common"
	"brewday/internal/tools"
	"brewday/internal/watcher"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const notificationNamePattern = "main_ferm_notification_"

type FermentationRouter struct {
	TLStore      TimelineStore
	SummaryStore SummaryStore
	Store        RecipeStore
	Notifier     Notifier
	watchersSet  map[string]bool // This keeps track if watches are set. In case of restart, it will go back to nil and force reconfig of watchers
}

// checkWatchers will check it watchers were set for a given recipe.
// If they were not, it will fetch the notification dates from the store and set them up again
// This method helps notifications be persistent in case of restarts.
// It should be called in handlers after the initial watcher setup (where a watcher set up is assumed)
func (r *FermentationRouter) checkWatchers(id string) error {
	reset := false
	if r.watchersSet == nil {
		reset = true
	} else {
		val, ok := r.watchersSet[id]
		if !ok {
			// In this case, the map exists but the recipe is not there
			// It could hae been that it got wiped, then restored by other recipe
			reset = true
		} else if !val {
			return errors.New("attempting to get watchers of recipe that does not expect them yet")
		}
	}
	if reset { // If the map is nil its been wiped, set up the watchers again
		re, err := r.Store.Retrieve(id)
		if err != nil {
			return err
		}
		dates, err := r.Store.RetrieveDates(id, notificationNamePattern)
		if err != nil {
			return err
		}
		var logMessage, notMessage string
		for _, date := range dates {
			if time.Until(*date) < 0 {
				logMessage = "Sending expired notification for recipe " + id
				notMessage = "Expired SG Measurement Notification. You should have measured on " + date.Format("2006-01-02")
			} else {
				logMessage = "notification"
				notMessage = "Measure SG"
			}
			watcher.NewWatcher(*date, func() error {
				log.Info().Str("id", id).Msg(logMessage)
				r.sendNotification(notMessage, "Main Fermentation "+re.Name, nil)
				return nil
			}).Start()
		}
		r.addWatchersSet(id)
	}
	return nil
}

func (r *FermentationRouter) addWatchersSet(id string) {
	if r.watchersSet == nil {
		r.watchersSet = make(map[string]bool)
	}
	r.watchersSet[id] = true
}

// sendNotification sends a notification if the notifier is available
func (r *FermentationRouter) sendNotification(message, title string, opts map[string]interface{}) error {
	if r.Notifier != nil {
		return r.Notifier.Send(message, title, opts)
	}
	return nil
}

// addTimelineEvent adds an event to the timeline
func (r *FermentationRouter) addTimelineEvent(id, message string) error {
	if r.TLStore != nil {
		return r.TLStore.AddEvent(id, message)
	}
	return nil
}

// addSummaryPreFermentation adds a pre fermentation summary
func (r *FermentationRouter) addSummaryPreFermentation(id string, volume, sg float32, notes string) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddPreFermentationVolume(id, volume, sg, notes)
	}
	return nil
}

// addSummaryEfficiency adds an efficiency summary
func (r *FermentationRouter) addSummaryEfficiency(id string, efficiencyPercentage float32) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddEfficiency(id, efficiencyPercentage)
	}
	return nil
}

// addSummaryYeastStart adds a yeast start summary
func (r *FermentationRouter) addSummaryYeastStart(id string, temperature, notes string) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddYeastStart(id, temperature, notes)
	}
	return nil
}

// addSummarySGMeasurement adds a SG measurement to the summary
func (r *FermentationRouter) addSummarySGMeasurement(id string, sg float32, date string, final bool, notes string) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddMainFermentationSGMeasurement(id, date, sg, final, notes)
	}
	return nil
}

// addSummaryAlcoholMainFermentation adds the alcohol after the main fermentation to the summary
func (r *FermentationRouter) addSummaryAlcoholMainFermentation(id string, alcohol float32) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddMainFermentationAlcohol(id, alcohol)
	}
	return nil
}

// registerRoutes registers the routes for the fermentation router
func (r *FermentationRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	fermentation := parent.Group("/fermentation")
	fermentation.GET("/pre/:recipe_id", r.getPreFermentationHandler).Name = "getPreFermentation"
	fermentation.POST("/pre/:recipe_id", r.postPreFermentationHandler).Name = "postPreFermentation"
	fermentation.GET("/pre/water/:recipe_id", r.getPreFermentationWaterHandler).Name = "getPreFermentationWater"
	fermentation.POST("/pre/water/:recipe_id", r.postPreFermentationWaterHandler).Name = "postPreFermentationWater"
	fermentation.GET("/yeast/:recipe_id", r.getFermentationYeastHandler).Name = "getFermentationYeast"
	fermentation.POST("/yeast/:recipe_id", r.postFermentationYeastHandler).Name = "postFermentationYeast"
	fermentation.GET("/start/:recipe_id", r.getMainFermentationStartHandler).Name = "getMainFermentationStart"
	fermentation.POST("/start/:recipe_id", r.postMainFermentationStartHandler).Name = "postMainFermentationStart"
	fermentation.GET("/main/:recipe_id", r.getMainFermentationHandler).Name = "getMainFermentation"
	fermentation.POST("/main/:recipe_id", r.postMainFermentationHandler).Name = "postMainFermentation"
}

// getPreFermentationHandler returns the handler for the pre fermentation page
func (r *FermentationRouter) getPreFermentationHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	err := r.addTimelineEvent(id, "Started Pre Fermentation")
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add timeline event")
	}
	err = r.Store.UpdateStatus(id, recipe.RecipeStatusPreFermentation, "measure")
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "fermentation_pre.html", map[string]interface{}{
		"Title":    "Pre Fermentation",
		"RecipeID": id,
	})
}

// postPreFermentationHandler handles the post request for the pre fermentation page
func (r *FermentationRouter) postPreFermentationHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	var req ReqPostPreFermentation
	err = c.Bind(&req)
	if err != nil {
		return err
	}
	err = r.addSummaryPreFermentation(id, req.Volume, req.SG, req.Notes)
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add pre fermentation summary")
	}
	// TODO: this is int he wrong place, it must be in the "hot" wort
	results, err := r.Store.RetrieveResults(id)
	if err != nil {
		return err
	}
	hotWortVol := results.HotWortVolume
	eff := tools.CalculateEfficiencySG(req.SG, hotWortVol, re.Mashing.GetTotalMaltWeight())
	err = r.addSummaryEfficiency(id, eff)
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add efficiency to summary")
	}
	// TODO: is this right? I am asking for the yeast lose somewhere else
	volumeDiff := req.Volume - (re.BatchSize + 1) // +1 for the 1l of yeast
	sgDiff := re.InitialSG - req.SG
	redirect := "getPreFermentationWater"
	queryParams := fmt.Sprintf("?volumeDiff=%f&sgDiff=%f", volumeDiff, sgDiff)
	return c.Redirect(http.StatusFound, c.Echo().Reverse(redirect, id)+queryParams)
}

// getPreFermentationWaterHandler returns the handler for the pre fermentation water page
func (r *FermentationRouter) getPreFermentationWaterHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	err = r.addTimelineEvent(id, "Started Pre Fermentation Water")
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add timeline event")
	}
	volumeDiffRaw := c.QueryParam("volumeDiff")
	sgDiffRaw := c.QueryParam("sgDiff")
	options := []WaterOption{}
	sgDiff, err := strconv.ParseFloat(sgDiffRaw, 32)
	if err != nil {
		return err
	}
	volumeDiff, err := strconv.ParseFloat(volumeDiffRaw, 32)
	if err != nil {
		return err
	}
	currentSG := re.InitialSG - float32(sgDiff)
	currentVol := re.BatchSize + float32(volumeDiff) + 1
	if sgDiff < 0.0 {
		toAdd, finalVol := tools.WaterForGravity(currentSG, re.InitialSG, currentVol)
		options = append(options, WaterOption{
			ToAdd:        toAdd,
			FinalVolume:  finalVol,
			FinalSG:      re.InitialSG,
			FinalSGPlato: tools.SGToPlato(re.InitialSG),
		})
		if volumeDiff < 0.0 {
			targetVol := re.BatchSize + 1
			toAdd, finalSG := tools.WaterForVolume(currentVol, targetVol, currentSG)
			options = append(options, WaterOption{
				ToAdd:        toAdd,
				FinalVolume:  targetVol,
				FinalSG:      finalSG,
				FinalSGPlato: tools.SGToPlato(finalSG),
			})
		}
	}
	err = r.Store.UpdateStatus(id, recipe.RecipeStatusPreFermentation, "water", tools.AnyToString(volumeDiff), tools.AnyToString(sgDiff))
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "fermentation_pre_water.html", map[string]interface{}{
		"Title":         "Pre Fermentation Water",
		"RecipeID":      id,
		"RecipeVolume":  re.BatchSize + 1,
		"RecipeSG":      re.InitialSG,
		"CurrentSG":     currentSG,
		"CurrentVolume": currentVol,
		"Options":       options,
	})
}

// postPreFermentationWaterHandler handles the post request for the pre fermentation water page
func (r *FermentationRouter) postPreFermentationWaterHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	var req ReqPostPreFermentationWater
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	err = r.addTimelineEvent(id, "Finished Adding Water")
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add timeline event")
	}
	err = r.addSummaryPreFermentation(id, req.FinalVolume, req.FinalSG, req.Notes)
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add pre fermentation summary")
	}
	err = r.Store.UpdateResult(id, recipe.ResultOriginalGravity, req.FinalSG)
	if err != nil {
		return err
	}
	err = r.Store.UpdateResult(id, recipe.ResultMainFermentationVolume, req.FinalVolume)
	if err != nil {
		return err
	}
	err = r.addTimelineEvent(id, "Finished Pre Fermentation")
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add timeline event")
	}
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getFermentationYeast", id))
}

// getFermentationYeastHandler returns the handler for the start fermentation (yeast) page
func (r *FermentationRouter) getFermentationYeastHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	err = r.addTimelineEvent(id, "Started Fermentation")
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add timeline event")
	}
	err = r.Store.UpdateStatus(id, recipe.RecipeStatusFermenting, "yeast")
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "fermentation_yeast.html", map[string]interface{}{
		"Title":       "Fermentation",
		"Subtitle":    "Start Fermentation",
		"RecipeID":    id,
		"Yeast":       re.Fermentation.Yeast,
		"Temperature": re.Fermentation.Temperature,
	})
}

// postFermentationYeastHandler handles the post request for the start fermentation (yeast) page
func (r *FermentationRouter) postFermentationYeastHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	var req ReqPostFermentationYeast
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	err = r.addTimelineEvent(id, "Inserted Yeast")
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add timeline event")
	}
	err = r.addSummaryYeastStart(id, req.Temperature, req.Notes)
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add yeast start to summary")
	}
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getMainFermentationStart", id))
}

// getMainFermentationHandler returns the handler for the main fermentation page
func (r *FermentationRouter) getMainFermentationStartHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	err := r.Store.UpdateStatus(id, recipe.RecipeStatusFermenting, "start")
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "fermentation_start.html", map[string]interface{}{
		"Title":              "Fermentation",
		"Subtitle":           "Set notification",
		"RecipeID":           id,
		"RecommendedMinDays": 8,
		"RecommendedDays":    10,
	})
}

// postMainFermentationHandler handles the post request for the main fermentation page
func (r *FermentationRouter) postMainFermentationStartHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	var req ReqPostFermentationStart
	err = c.Bind(&req)
	if err != nil {
		return err
	}
	now := time.Now()
	timeDiff := req.NotificationDays - req.NotificationDaysBefore
	var notificationDate time.Time
	for i := 0; i <= timeDiff; i++ {
		switch req.TimeUnit {
		case "days":
			notificationDate = now.AddDate(0, 0, req.NotificationDaysBefore+i)
		case "seconds": // This is mainly for testing
			notificationDate = now.Add(time.Duration(req.NotificationDaysBefore+i) * time.Second)
		default:
			return fmt.Errorf("unknown time unit %s", req.TimeUnit)
		}
		err = r.Store.AddDate(id, &notificationDate, fmt.Sprintf(notificationNamePattern+"%d", i))
		if err != nil {
			return err
		}
		var logMessage, notMessage string
		if i == 0 {
			logMessage = "first notification"
			notMessage = "Measure SG for the first time"
		} else {
			logMessage = "notification"
			notMessage = "Measure SG"
		}
		watcher.NewWatcher(notificationDate, func() error {
			log.Info().Str("id", id).Msg(logMessage)
			r.sendNotification(notMessage, "Main Fermentation "+re.Name, nil)
			return nil
		}).Start()
		r.addWatchersSet(id)
	}
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getMainFermentation", id))
}

// getMainFermentationHandler returns the handler for the main fermentation page
func (r *FermentationRouter) getMainFermentationHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	err := r.checkWatchers(id)
	if err != nil {
		return err
	}
	minDate, err := r.Store.RetrieveDates(id, notificationNamePattern+"0")
	if err != nil {
		return err
	}
	missing := time.Until(*minDate[0])
	if missing > 0 {
		err = r.Store.UpdateStatus(id, recipe.RecipeStatusFermenting, "wait")
		if err != nil {
			return err
		}
		return c.Render(http.StatusOK, "fermentation_wait.html", map[string]interface{}{
			"Title":       "Fermentation",
			"Subtitle":    "Main Fermentation",
			"RecipeID":    id,
			"MissingTime": missing.String(),
		})
	} else {
		// This should ask for the SGs and once user clicks on its stable for me lead to
		// sugar calculation
		err = r.Store.UpdateStatus(id, recipe.RecipeStatusFermenting, "main")
		if err != nil {
			return err
		}
		measurements, err := r.Store.RetrieveMainFermSGs(id)
		if err != nil {
			return err
		}
		return c.Render(http.StatusOK, "fermentation_main.html", map[string]interface{}{
			"Title":            "Fermentation",
			"Subtitle":         "Main Fermentation",
			"RecipeID":         id,
			"PastMeasurements": measurements,
		})
	}
}

// postMainFermentationHandler handles the post request for the main fermentation page
func (r *FermentationRouter) postMainFermentationHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	var req ReqPostMainFermentation
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	err = r.addTimelineEvent(id, "Added SG Measurement")
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add timeline event")
	}
	m := recipe.SGMeasurement{
		Date:  time.Now().Format("2006-01-02"),
		Value: req.SG,
	}
	err = r.Store.AddMainFermSG(id, &m)
	if err != nil {
		return err
	}
	err = r.addSummarySGMeasurement(id, m.Value, m.Date, req.Final, req.Notes)
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add sg measurement to summary")
	}
	if req.Final {
		err = r.Store.UpdateResult(id, recipe.ResultFinalGravity, req.SG)
		if err != nil {
			return err
		}
		results, err := r.Store.RetrieveResults(id)
		if err != nil {
			return err
		}
		og := results.OriginalGravity
		alc := tools.CalculateAlcohol(og, req.SG)
		err = r.Store.UpdateResult(id, recipe.ResultAlcohol, alc)
		if err != nil {
			return err
		}
		err = r.addSummaryAlcoholMainFermentation(id, alc)
		if err != nil {
			log.Error().Str("id", id).Err(err).Msg("could not add alcohol to summary")
		}
		return c.Redirect(http.StatusFound, c.Echo().Reverse("getDryHop", id))
	}
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getMainFermentation", id))
}
