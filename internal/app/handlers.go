package app

import (
	"brewday/internal/routers/common"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// addTimelineEvent adds an event to the timeline
func (a *App) addTimelineEvent(id, message string) error {
	if a.TLStore != nil {
		return a.TLStore.AddEvent(id, message)
	}
	return nil
}

// postTimelineEvent is the handler for sent timeline events
func (a *App) postTimelineEvent(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	var req ReqPostTimelineEvent
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	err = a.addTimelineEvent(id, req.Message)
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add timeline event")
	}
	return c.NoContent(200)
}

// customErrorHandler is a custom error handler
func (a *App) customErrorHandler(err error, c echo.Context) {
	log.Error().Err(err).Msg(c.Request().RequestURI)
	notFound := strings.Contains(strings.ToLower(err.Error()), "not found")
	if err == common.ErrNoRecipeLoaded || err == common.ErrNoRecipeIDProvided || notFound {
		err2 := c.Render(404, "error_no_recipe_loaded.html", map[string]interface{}{
			"Title": "Error in recipe",
		})
		if err2 != nil {
			log.Error().Err(err2).Msg("error while rendering error page")
		}
	}
}
