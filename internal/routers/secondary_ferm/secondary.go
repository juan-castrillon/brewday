package secondaryferm

import (
	"brewday/internal/recipe"
	"brewday/internal/routers/common"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type SecondaryFermentationRouter struct {
	TLStore      TimelineStore
	SummaryStore SummaryRecorderStore
	Store        RecipeStore
	Notifier     Notifier
}

// RegisterRoutes adds routes to the web server
// It receives the root web server and a parent group
func (r *SecondaryFermentationRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	sf := parent.Group("/secondary_fermentation")
	sf.GET("/dry_hop/start/:recipe_id", r.getDryHopStartHandler).Name = "getDryHopStart"
	sf.POST("/dry_hop/start/:recipe_id", r.postDryHopStartHandler).Name = "postDryHopStart"
	sf.GET("/dry_hop/confirm/:recipe_id", r.getDryHopConfirmHandler).Name = "getDryHopConfirm"
	sf.POST("/dry_hop/confirm/:recipe_id", r.postDryHopConfirmHandler).Name = "postDryHopConfirm"
	sf.GET("/dry_hop/:recipe_id", r.getDryHopHandler).Name = "getDryHop"
	sf.POST("/dry_hop/:recipe_id", r.postDryHopHandler).Name = "postDryHop"
	sf.GET("/bottle/:recipe_id", r.getBottleHandler).Name = "getBottle"
	sf.POST("/bottle/:recipe_id", r.postBottleHandler).Name = "postBottle"
	sf.GET("/start/:recipe_id", r.getSecondaryFermentationStartHandler).Name = "getSecondaryFermentationStart"
	sf.POST("/start/:recipe_id", r.postSecondaryFermentationStartHandler).Name = "postSecondaryFermentationStart"
	sf.GET("/fridge/:recipe_id", r.getFridgeHandler).Name = "getFridge"
	sf.POST("/fridge/:recipe_id", r.postFridgeHandler).Name = "postFridge"
	root.GET("/end/:recipe_id", r.getEndHandler).Name = "getEnd"
}

// getDryHopStartHandler is the handler for the dry hop start page
// this page allows to set notifications for the different dry hops and their times
func (r *SecondaryFermentationRouter) getDryHopStartHandler(c echo.Context) error {
	return nil
}

// postDryHopStartHandler is the handler for the dry hop start page
func (r *SecondaryFermentationRouter) postDryHopStartHandler(c echo.Context) error {
	return nil
}

// getDryHopConfirmHandler is the handler for the dry hop confirm page
// this page asks for confirmation of the dry hop and shows remaining time of the notification
func (r *SecondaryFermentationRouter) getDryHopConfirmHandler(c echo.Context) error {
	return nil
}

// postDryHopConfirmHandler is the handler for the dry hop confirm page
func (r *SecondaryFermentationRouter) postDryHopConfirmHandler(c echo.Context) error {
	return nil
}

func (r *SecondaryFermentationRouter) getDryHopHandler(c echo.Context) error {
	return nil
}

// postDryHopHandler is the handler for the dry hop page
func (r *SecondaryFermentationRouter) postDryHopHandler(c echo.Context) error {
	return nil
}

// getBottleHandler is the handler for the bottle page
func (r *SecondaryFermentationRouter) getBottleHandler(c echo.Context) error {
	return nil
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
	var hops []recipe.Hops
	for _, h := range re.Hopping.Hops {
		if h.DryHop {
			hops = append(hops, h)
		}
	}
	err = r.addTimelineEvent(id, "Finished Day")
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add timeline event")
	}
	re.SetStatus(recipe.RecipeStatusFinished)
	return c.Render(http.StatusOK, "finished_day.html", map[string]interface{}{
		"Title":     "End Fermentation",
		"RecipeID":  id,
		"Subtitle":  "Congratulations, you've finished the brew day!",
		"Hops":      hops,
		"Additions": re.Fermentation.AdditionalIngredients,
	})
}
