package secondaryferm

import (
	"brewday/internal/recipe"
	"brewday/internal/routers/common"
	"brewday/internal/tools"
	"brewday/internal/watcher"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type SecondaryFermentationRouter struct {
	TLStore               TimelineStore
	SummaryStore          SummaryRecorderStore
	Store                 RecipeStore
	Notifier              Notifier
	hopWatchersLock       sync.Mutex
	secondaryWatchersLock sync.Mutex
	dryHopsLock           sync.Mutex
	sugarResultsLock      sync.Mutex
	// HopWatchers stores all dry hop watchers for a recipe id
	HopWatchers map[string]DryHopNotification
	// DryHops relates a list of dry hops with a recipe id
	DryHops map[string]DryHopMap
	// SugarResults stores the sugar results for a recipe id
	SugarResults map[string][]SugarResult
	// SecondaryWatchers stores the watchers for the secondary fermentation
	SecondaryWatchers map[string]SecondaryFermentationWatcher
}

// RegisterRoutes adds routes to the web server
// It receives the root web server and a parent group
func (r *SecondaryFermentationRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	sf := parent.Group("/secondary_fermentation")
	sf.GET("/dry_hop/start/:recipe_id/load", r.getDryHopStartLoadHandler).Name = "getDryHopStartLoad"
	sf.GET("/dry_hop/start/:recipe_id", r.getDryHopStartHandler).Name = "getDryHopStart"
	sf.POST("/dry_hop/start/:recipe_id", r.postDryHopStartHandler).Name = "postDryHopStart"
	sf.GET("/dry_hop/confirm/:recipe_id", r.getDryHopConfirmHandler).Name = "getDryHopConfirm"
	sf.POST("/dry_hop/confirm/:recipe_id", r.postDryHopConfirmHandler).Name = "postDryHopConfirm"
	sf.GET("/pre_bottle/:recipe_id", r.getPreBottleHandler).Name = "getPreBottle"
	sf.POST("/pre_bottle/:recipe_id", r.postPreBottleHandler).Name = "postPreBottle"
	sf.GET("/bottle/:recipe_id", r.getBottleHandler).Name = "getBottle"
	sf.POST("/bottle/:recipe_id", r.postBottleHandler).Name = "postBottle"
	sf.GET("/start/:recipe_id", r.getSecondaryFermentationStartHandler).Name = "getSecondaryFermentationStart"
	sf.POST("/start/:recipe_id", r.postSecondaryFermentationStartHandler).Name = "postSecondaryFermentationStart"
	sf.POST("/end/:recipe_id", r.postSecondaryFermentationEndHandler).Name = "postSecondaryFermentationEnd"
	root.GET("/end/:recipe_id", r.getEndHandler).Name = "getEnd"
}

// getDryHopStartLoadHandler is responsible for loading the dry hops in the internal structures and then
// redirecting to the dry hop start page. Its only called once per recipe
func (r *SecondaryFermentationRouter) getDryHopStartLoadHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	for i, h := range re.Hopping.Hops {
		if h.DryHop {
			hop := DryHop{
				id:       fmt.Sprintf("%s_%d", h.Name, i),
				Name:     h.Name,
				Amount:   h.Amount,
				Duration: h.Duration,
			}
			r.addDryHop(id, &hop)
		}
	}
	for i, a := range re.Fermentation.AdditionalIngredients {
		hop := DryHop{
			id:       fmt.Sprintf("%s_%d", a.Name, i),
			Name:     a.Name,
			Amount:   a.Amount,
			Duration: a.Duration,
		}
		r.addDryHop(id, &hop)
	}
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getDryHopStart", id))
}

// getDryHopStartHandler is the handler for the dry hop start page
// this page allows to set notifications for the different dry hops and their times
func (r *SecondaryFermentationRouter) getDryHopStartHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	err := r.Store.UpdateStatus(id, recipe.RecipeStatusFermenting, "dry_hop_start")
	if err != nil {
		return err
	}
	dr, err := r.getDryHops(id)
	if err != nil {
		log.Info().Str("id", id).Err(err).Msg("Recipe has no dry hops")
		return c.Redirect(http.StatusFound, c.Echo().Reverse("getPreBottle", id))
	}
	return c.Render(http.StatusOK, "secondary_dry_hop_start.html", map[string]interface{}{
		"Title":    "Dry Hop",
		"RecipeID": id,
		"Subtitle": "Set the notifications for the dry hops",
		"Hops":     dr,
	})
}

// postDryHopStartHandler is the handler for the dry hop start page
func (r *SecondaryFermentationRouter) postDryHopStartHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	var req ReqPostDryHopStart
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	dr, err := r.getDryHops(id)
	if err != nil {
		return err
	}
	dh, ok := dr[req.ID]
	if !ok {
		return fmt.Errorf("dry hop with id %s not found", req.ID)
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
		log.Info().Msgf("dry hop notification triggered for hop %s", dh.Name)
		return r.sendNotification(fmt.Sprintf("Dry hop %s", dh.Name), "Dry hop", nil)
	})
	err = r.addDryHopNotification(id, req.ID, w)
	if err != nil {
		return err
	}
	w.Start()
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getDryHopStart", id))
}

// getDryHopConfirmHandler is the handler for the dry hop confirm page
// this page asks for confirmation of the dry hop and shows remaining time of the notification
func (r *SecondaryFermentationRouter) getDryHopConfirmHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	err := r.Store.UpdateStatus(id, recipe.RecipeStatusFermenting, "dry_hop_confirm")
	if err != nil {
		return err
	}
	not, err := r.getDryHopNotifications(id)
	if err != nil {
		return err
	}
	dr, err := r.getDryHops(id)
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "secondary_dry_hop_confirm.html", map[string]interface{}{
		"Title":         "Dry Hop",
		"Subtitle":      "Confirm the dry hop",
		"RecipeID":      id,
		"Notifications": not,
		"Hops":          dr,
	})
}

// postDryHopConfirmHandler is the handler for the dry hop confirm page
func (r *SecondaryFermentationRouter) postDryHopConfirmHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	var req ReqPostDryHopConfirm
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	dr, err := r.getDryHops(id)
	if err != nil {
		return err
	}
	dh, ok := dr[req.ID]
	if !ok {
		return fmt.Errorf("dry hop with id %s not found", req.ID)
	}
	not, err := r.getDryHopNotifications(id)
	if err != nil {
		return err
	}
	w, ok := not[req.ID]
	if ok {
		w.Stop()
	}
	dh.In = true
	dh.InDate = time.Now().Format("2006-01-02 15:04")
	err = r.addTimelineEvent(id, fmt.Sprintf("Added dry hop %s", dh.Name))
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add timeline event")
	}
	err = r.addSummaryDryHop(id, dh.Name, req.Amount)
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add summary dry hop")
	}
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getDryHopConfirm", id))
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
	res := re.GetResults()
	vol := req.Volume - req.LostVolume
	r.calculateSugar(
		id, vol,
		re.Fermentation.Carbonation, req.Temperature,
		res.Alcohol,
		req.SugarType,
	)
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
	sugarResults, err := r.getSugarResults(id)
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
	res := re.GetResults()
	_, realAlcohol := tools.SugarForCarbonation(
		req.RealVolume, re.Fermentation.Carbonation, req.Temperature,
		res.Alcohol, req.Water, req.SugarType,
	)
	realCO2 := tools.CarbonationForSugar(req.RealVolume, req.SugarAmount, req.Temperature, req.SugarType)
	err = r.addSummaryBottle(id, realCO2, realAlcohol, req.SugarAmount, req.Temperature, req.RealVolume, req.SugarType, req.Notes)
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add summary bottle")
	}
	re.SetAlcohol(realAlcohol)
	r.addTimelineEvent(id, "Bottled")
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getSecondaryFermentationStart", id))
}

// getSecondaryFermentationStartHandler is the handler for the secondary fermentation start page
func (r *SecondaryFermentationRouter) getSecondaryFermentationStartHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	watch := r.getSecondaryWatcher(id)
	missingTime := ""
	isDone := false
	isSet := false
	if watch != nil {
		isSet = true
		if watch.IsDone() {
			isDone = true
		} else {
			missingTime = watch.MissingTime().String()
		}
	}
	err := r.Store.UpdateStatus(id, recipe.RecipeStatusFermenting, "start_secondary")
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "secondary_start.html", map[string]interface{}{
		"Title":    "Start Secondary Fermentation",
		"RecipeID": id,
		"Subtitle": "First, let the bottles at warm temperature",
		"MinDays":  5,
		"Missing":  missingTime,
		"IsDone":   isDone,
		"IsSet":    isSet,
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
		return r.sendNotification("Secondary Fermentation", "Time to put bottles in the fridge", nil)
	})
	w.Start()
	err = r.addSecondaryWatcher(id, w)
	if err != nil {
		return err
	}
	r.addTimelineEvent(id, "Secondary Fermentation Started")
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getSecondaryFermentationStart", id))
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
	return c.Render(http.StatusOK, "finished_day.html", map[string]interface{}{
		"Title":    "End Fermentation",
		"RecipeID": id,
		"Subtitle": "Congratulations, you've finished the brew!",
	})
}
