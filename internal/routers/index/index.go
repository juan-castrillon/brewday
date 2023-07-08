package index

import "github.com/labstack/echo/v4"

// IndexRouter is the router for the index page
type IndexRouter struct {
}

// RegisterRoutes registers the routes of the index page
func (r *IndexRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	root.GET("", indexHandler)
}

// indexHandler is the handler for the index page
func indexHandler(c echo.Context) error {
	return c.Render(200, "index.html", map[string]interface{}{
		"message": "Hello World!",
	})
}
