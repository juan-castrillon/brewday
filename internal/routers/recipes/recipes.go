package recipes

import (
	"brewday/internal/recipe"

	"github.com/labstack/echo/v4"
)

type RecipesRouter struct {
	Store RecipeStore
}

func (r *RecipesRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	recipes := parent.Group("/recipes")
	recipes.GET("", r.getRecipesHandler).Name = "getRecipes"
}

// getRecipeList returns the list of recipes
func (r *RecipesRouter) getRecipeList() ([]*recipe.Recipe, error) {
	if r.Store != nil {
		return r.Store.List()
	}
	return nil, nil
}

// getRecipesHandler is the handler for the recipes page
func (r *RecipesRouter) getRecipesHandler(c echo.Context) error {
	recipes, err := r.getRecipeList()
	if err != nil {
		return err
	}
	return c.Render(200, "recipes.html", map[string]interface{}{
		"Title":    "Recipes",
		"Subtitle": "Loaded Recipes",
		"Recipes":  recipes,
	})
}
