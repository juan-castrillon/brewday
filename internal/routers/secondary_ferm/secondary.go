package secondaryferm

import (
	"brewday/internal/recipe"
	"brewday/internal/routers/common"
	"brewday/internal/watcher"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type SecondaryFermentationRouter struct {
	TLStore         TimelineStore
	SummaryStore    SummaryRecorderStore
	Store           RecipeStore
	Notifier        Notifier
	hopWatchersLock sync.Mutex
	dryHopsLock     sync.Mutex
	// HopWatchers stores all dry hop watchers for a recipe id
	HopWatchers map[string]DryHopNotification
	// DryHops relates a list of dry hops with a recipe id
	DryHops map[string]DryHopMap
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
	sf.GET("/bottle/:recipe_id", r.getBottleHandler).Name = "getBottle"
	sf.POST("/bottle/:recipe_id", r.postBottleHandler).Name = "postBottle"
	sf.GET("/start/:recipe_id", r.getSecondaryFermentationStartHandler).Name = "getSecondaryFermentationStart"
	sf.POST("/start/:recipe_id", r.postSecondaryFermentationStartHandler).Name = "postSecondaryFermentationStart"
	sf.GET("/fridge/:recipe_id", r.getFridgeHandler).Name = "getFridge"
	sf.POST("/fridge/:recipe_id", r.postFridgeHandler).Name = "postFridge"
	root.GET("/end/:recipe_id", r.getEndHandler).Name = "getEnd"
	//http://localhost:8080/secondary_fermentation/dry_hop/start/47696e67657220576974/load
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
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	re.SetStatus(recipe.RecipeStatusFermenting, "dry_hop_start")
	dr, err := r.getDryHops(id)
	if err != nil {
		return err
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
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	re.SetStatus(recipe.RecipeStatusFermenting, "dry_hop_confirm")
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

// getBottleHandler is the handler for the bottle page
func (r *SecondaryFermentationRouter) getBottleHandler(c echo.Context) error {
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getEnd", c.Param("recipe_id")))
}

// postBottleHandler is the handler for the bottle page
func (r *SecondaryFermentationRouter) postBottleHandler(c echo.Context) error {
	return nil
}

// getSecondaryFermentationStartHandler is the handler for the secondary fermentation start page
func (r *SecondaryFermentationRouter) getSecondaryFermentationStartHandler(c echo.Context) error {
	return nil
}

// postSecondaryFermentationStartHandler is the handler for the secondary fermentation start page
func (r *SecondaryFermentationRouter) postSecondaryFermentationStartHandler(c echo.Context) error {
	return nil
}

// getFridgeHandler is the handler for the fridge page
func (r *SecondaryFermentationRouter) getFridgeHandler(c echo.Context) error {
	return nil
}

// postFridgeHandler is the handler for the fridge page
func (r *SecondaryFermentationRouter) postFridgeHandler(c echo.Context) error {
	return nil
}

// getEndHandler is the handler for the end page
func (r *SecondaryFermentationRouter) getEndHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	err = r.addTimelineEvent(id, "Finished Day")
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add timeline event")
	}
	re.SetStatus(recipe.RecipeStatusFinished)
	return c.Render(http.StatusOK, "finished_day.html", map[string]interface{}{
		"Title":    "End Fermentation",
		"RecipeID": id,
		"Subtitle": "Congratulations, you've finished the brew day!",
	})
}
