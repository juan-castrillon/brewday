package secondaryferm

import "github.com/labstack/echo/v4"

type SecondaryFermentationRouter struct {
	TLStore      TimelineStore
	SummaryStore SummaryRecorderStore
	Store        RecipeStore
	Notifier     Notifier
}

// RegisterRoutes adds routes to the web server
// It receives the root web server and a parent group
func (r *SecondaryFermentationRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {}
