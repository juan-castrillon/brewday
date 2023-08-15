package lautern

import (
	"errors"

	"github.com/labstack/echo/v4"
)

type LauternRouter struct {
	TL      Timeline
	Summary SummaryRecorder
	Store   RecipeStore
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
func (r *LauternRouter) addSummaryLauternNotes(notes string) {
	if r.Summary != nil {
		r.Summary.AddLaunternNotes(notes)
	}
}

// getLauternHandler is the handler for the lautern page
func (r *LauternRouter) getLauternHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return errors.New("no recipe id provided")
	}
	r.addTimelineEvent("Started Läutern")
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	return c.Render(200, "lautern.html", map[string]interface{}{
		"Title":       "Mash " + re.Name,
		"Subtitle":    "Läutern",
		"RecipeID":    id,
		"MashOutTemp": re.Mashing.MashOutTemperature,
		"Hops":        re.Hopping.Hops,
		"RestTime":    1,
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
		return errors.New("no recipe id provided")
	}
	r.addTimelineEvent("Finished Läutern")
	r.addSummaryLauternNotes(req.Notes)
	return c.Redirect(302, c.Echo().Reverse("getStartHopping", id))
}
