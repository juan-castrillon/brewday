package fermentation

import (
	"brewday/internal/recipe"
	"brewday/internal/routers/common"
	"brewday/internal/tools"
	"brewday/internal/watcher"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type FermentationRouter struct {
	TLStore       TimelineStore
	SummaryStore  SummaryRecorderStore
	Store         RecipeStore
	Notifier      Notifier
	statusMapLock sync.Mutex
	sgMapLock     sync.Mutex
	// MainFermentationStatus is a map of recipe id to whether the main fermentation minimum days have passed
	MainFermentationStatus map[string]*FermentationStatus
	// SGMeasurements is a map of recipe id to the SG measurements
	SGMeasurements map[string][]SGMeasurement
}

// sendNotification sends a notification if the notifier is available
func (r *FermentationRouter) sendNotification(message, title string, opts map[string]interface{}) error {
	if r.Notifier != nil {
		return r.Notifier.Send(message, title, opts)
	}
	return nil
}

// addMainFermentationStatus adds a main fermentation status to a recipe id
func (r *FermentationRouter) addMainFermentationStatus(id string, status bool) {
	r.statusMapLock.Lock()
	defer r.statusMapLock.Unlock()
	if r.MainFermentationStatus == nil {
		r.MainFermentationStatus = make(map[string]*FermentationStatus)
	}
	_, ok := r.MainFermentationStatus[id]
	if !ok {
		r.MainFermentationStatus[id] = &FermentationStatus{}
	}
	r.MainFermentationStatus[id].MinDaysPassed = status
}

// addMainFermentationWatcher adds a watcher to a recipe id
func (r *FermentationRouter) addMainFermentationWatcher(id string, w *watcher.Watcher) {
	r.statusMapLock.Lock()
	defer r.statusMapLock.Unlock()
	if r.MainFermentationStatus == nil {
		r.MainFermentationStatus = make(map[string]*FermentationStatus)
	}
	_, ok := r.MainFermentationStatus[id]
	if !ok {
		r.MainFermentationStatus[id] = &FermentationStatus{}
	}
	r.MainFermentationStatus[id].InitialWatch = w
}

// getMainFermentationStatus returns the main fermentation status for a recipe id
func (r *FermentationRouter) getMainFermentationStatus(id string) (*FermentationStatus, error) {
	r.statusMapLock.Lock()
	defer r.statusMapLock.Unlock()
	status, ok := r.MainFermentationStatus[id]
	if !ok {
		return nil, fmt.Errorf("no status for id %s", id)
	}
	return status, nil
}

// addSGMeasurement adds an SG measurement to a recipe id
func (r *FermentationRouter) addSGMeasurement(id string, sg SGMeasurement) {
	r.sgMapLock.Lock()
	defer r.sgMapLock.Unlock()
	if r.SGMeasurements == nil {
		r.SGMeasurements = make(map[string][]SGMeasurement)
	}
	_, ok := r.SGMeasurements[id]
	if !ok {
		r.SGMeasurements[id] = []SGMeasurement{}
	}
	r.SGMeasurements[id] = append(r.SGMeasurements[id], sg)
}

// getSGMeasurements returns the SG measurements for a recipe id
func (r *FermentationRouter) getSGMeasurements(id string) ([]SGMeasurement, error) {
	r.sgMapLock.Lock()
	defer r.sgMapLock.Unlock()
	sgs, ok := r.SGMeasurements[id]
	if !ok {
		return []SGMeasurement{}, nil
	}
	return sgs, nil
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
		return r.SummaryStore.AddSummaryPreFermentation(id, volume, sg, notes)
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
func (r *FermentationRouter) addSummarySGMeasurement(id string, m SGMeasurement, final bool, notes string) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddSGMeasurement(id, m.Date, m.Gravity, final, notes)
	}
	return nil
}

// addSummaryAlcoholMainFermentation adds the alcohol after the main fermentation to the summary
func (r *FermentationRouter) addSummaryAlcoholMainFermentation(id string, alcohol float32) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddAlcoholMainFermentation(id, alcohol)
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
	err = r.Store.UpdateResults(id, recipe.ResultOriginalGravity, req.FinalSG)
	if err != nil {
		return err
	}
	err = r.Store.UpdateResults(id, recipe.ResultMainFermentationVolume, req.FinalVolume)
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
	r.addMainFermentationStatus(id, false)
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
		var w *watcher.Watcher
		if i == 0 {
			w = watcher.NewWatcher(notificationDate, func() error {
				log.Info().Str("id", id).Msg("first notification")
				r.sendNotification("Measure SG for the first time", "Main Fermentation "+re.Name, nil)
				r.addMainFermentationStatus(id, true)
				return nil
			})
			r.addMainFermentationWatcher(id, w)
		} else {
			w = watcher.NewWatcher(notificationDate, func() error {
				log.Info().Str("id", id).Msg("notification")
				r.sendNotification("Measure SG", "Main Fermentation "+re.Name, nil)
				return nil
			})
		}
		w.Start()
	}
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getMainFermentation", id))
}

// getMainFermentationHandler returns the handler for the main fermentation page
func (r *FermentationRouter) getMainFermentationHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	status, err := r.getMainFermentationStatus(id)
	if err != nil {
		return err
	}
	if !status.MinDaysPassed {
		err = r.Store.UpdateStatus(id, recipe.RecipeStatusFermenting, "wait")
		if err != nil {
			return err
		}
		wt := status.InitialWatch.MissingTime()
		return c.Render(http.StatusOK, "fermentation_wait.html", map[string]interface{}{
			"Title":       "Fermentation",
			"Subtitle":    "Main Fermentation",
			"RecipeID":    id,
			"MissingTime": wt.String(),
		})
	} else {
		// This should ask for the SGs and once user clicks on its stable for me lead to
		// sugar calculation
		err = r.Store.UpdateStatus(id, recipe.RecipeStatusFermenting, "main")
		if err != nil {
			return err
		}
		measurements, err := r.getSGMeasurements(id)
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
	m := SGMeasurement{
		Date:    time.Now().Format("2006-01-02"),
		Gravity: req.SG,
	}
	r.addSGMeasurement(id, m)
	err = r.addSummarySGMeasurement(id, m, req.Final, req.Notes)
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add sg measurement to summary")
	}
	if req.Final {
		err = r.Store.UpdateResults(id, recipe.ResultFinalGravity, req.SG)
		if err != nil {
			return err
		}
		results, err := r.Store.RetrieveResults(id)
		if err != nil {
			return err
		}
		og := results.OriginalGravity
		alc := tools.CalculateAlcohol(og, req.SG)
		err = r.Store.UpdateResults(id, recipe.ResultAlcohol, alc)
		if err != nil {
			return err
		}
		err = r.addSummaryAlcoholMainFermentation(id, alc)
		if err != nil {
			log.Error().Str("id", id).Err(err).Msg("could not add alcohol to summary")
		}
		return c.Redirect(http.StatusFound, c.Echo().Reverse("getDryHopStartLoad", id))
	}
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getMainFermentation", id))
}
