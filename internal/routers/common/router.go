package common

import "github.com/labstack/echo/v4"

// Router represents a component that adds routes to the web server
// In the context of the application it encapsulates the different pages or functionalities
type Router interface {
	// RegisterRoutes adds routes to the web server
	// It receives the root web server and a parent group
	// The concrete implementations have the option to add routes to the root server or to a parent group
	// Middleware can be added to the parent group in the caller and its not a concern of the router
	RegisterRoutes(root *echo.Echo, parent *echo.Group)
}
