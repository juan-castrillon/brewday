package lautern

import (
	"brewday/internal/recipe"
	"brewday/internal/routers/common"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type LauternRouter struct {
	TLStore      TimelineStore
	SummaryStore SummaryStore
	Store        RecipeStore
}

// RegisterRoutes adds routes to the web server
func (r *LauternRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	lautern := parent.Group("/lautern")
	lautern.GET("/:recipe_id", r.getLauternHandler).Name = "getLautern"
	lautern.POST("/:recipe_id", r.postLauternHandler).Name = "postLautern"
}

// addTimelineEvent adds an event to the timeline
func (r *LauternRouter) addTimelineEvent(id, message string) error {
	if r.TLStore != nil {
		return r.TLStore.AddEvent(id, message)
	}
	return nil
}

// addSummaryLauternNotes adds lautern notes to the summary
func (r *LauternRouter) addSummaryLauternNotes(id, notes string) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddLauternNotes(id, notes)
	}
	return nil
}

// getLauternHandler is the handler for the lautern page
func (r *LauternRouter) getLauternHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	err := r.addTimelineEvent(id, "Started Läutern")
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add timeline event")
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	err = r.Store.UpdateStatus(id, recipe.RecipeStatusLautering)
	if err != nil {
		return err
	}
	return c.Render(200, "lautern.html", map[string]interface{}{
		"Title":       "Mash " + re.Name,
		"Subtitle":    "Läutern",
		"RecipeID":    id,
		"MashOutTemp": re.Mashing.MashOutTemperature,
		"Hops":        re.Hopping.Hops,
		"RestTime":    15,
	})
}

// postLauternHandler is the handler for the lautern page
func (r *LauternRouter) postLauternHandler(c echo.Context) error {
	var req ReqPostLautern
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	err = r.addTimelineEvent(id, "Finished Läutern")
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add timeline event")
	}
	err = r.addSummaryLauternNotes(id, req.Notes)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to add lautern notes to summary")
	}
	return c.Redirect(302, c.Echo().Reverse("getStartHopping", id))
}
