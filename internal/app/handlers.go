package app

import (
	"brewday/internal/routers/common"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// addTimelineEvent adds an event to the timeline
func (a *App) addTimelineEvent(message string) {
	if a.TL != nil {
		a.TL.AddEvent(message)
	}
}

// postTimelineEvent is the handler for sent timeline events
func (a *App) postTimelineEvent(c echo.Context) error {
	var req ReqPostTimelineEvent
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	a.addTimelineEvent(req.Message)
	return c.NoContent(200)
}

// sendNotification sends a notification if the notifier is available
func (a *App) sendNotification(message, title string, opts map[string]interface{}) error {
	if a.notifier != nil {
		return a.notifier.Send(message, title, opts)
	}
	return nil
}

// postNotification is the handler for sending notifications
func (a *App) postNotification(c echo.Context) error {
	var req ReqPostNotification
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	err = a.sendNotification(req.Message, req.Title, req.Options)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
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
