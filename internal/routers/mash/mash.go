package mash

import (
	"brewday/internal/recipe"

	"github.com/labstack/echo/v4"
)

type MashRouter struct {
	Store  RecipeStore
	recipe *recipe.Recipe
}

func (r *MashRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	mash := parent.Group("/mash")
	mash.GET("/start/:recipe_id", r.getMashStartHandler).Name = "getMashStart"
}

// getMashStartHandler is the handler for the mash start page
func (r *MashRouter) getMashStartHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	r.recipe = re
	return c.Render(200, "mash_start.html", map[string]interface{}{
		"Title":  "Mash " + re.Name,
		"Recipe": r.recipe,
	})
}
