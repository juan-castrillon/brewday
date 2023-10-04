package lautern

import (
	"brewday/internal/recipe"
	"brewday/internal/routers/common"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type LauternRouter struct {
	TL           Timeline
	SummaryStore SummaryRecorderStore
	Store        RecipeStore
}

// RegisterRoutes adds routes to the web server
func (r *LauternRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	lautern := parent.Group("/lautern")
	lautern.GET("/:recipe_id", r.getLauternHandler).Name = "getLautern"
	lautern.POST("/:recipe_id", r.postLauternHandler).Name = "postLautern"
}

// addTimelineEvent adds an event to the timeline
func (r *LauternRouter) addTimelineEvent(message string) {
	if r.TL != nil {
		r.TL.AddEvent(message)
	}
}

// addSummaryLauternNotes adds lautern notes to the summary
func (r *LauternRouter) addSummaryLauternNotes(id, notes string) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddLaunternNotes(id, notes)
	}
	return nil
}

// getLauternHandler is the handler for the lautern page
func (r *LauternRouter) getLauternHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	r.addTimelineEvent("Started Läutern")
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	re.SetStatus(recipe.RecipeStatusLautering)
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
	r.addTimelineEvent("Finished Läutern")
	err = r.addSummaryLauternNotes(id, req.Notes)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to add lautern notes to summary")
	}
	return c.Redirect(302, c.Echo().Reverse("getStartHopping", id))
}
