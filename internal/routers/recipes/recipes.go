package recipes

import (
	"brewday/internal/recipe"
	"brewday/internal/routers/common"
	"brewday/internal/tools"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RecipesRouter struct {
	Store        RecipeStore
	TLStore      TimelineStore
	SummaryStore SummaryRecorderStore
}

func (r *RecipesRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	recipes := parent.Group("/recipes")
	recipes.GET("", r.getRecipesHandler).Name = "getRecipes"
	recipes.GET("/continue/:recipe_id", r.getContinueHandler).Name = "getContinue"
	recipes.GET("/start/:recipe_id", r.getStartHandler).Name = "getRecipeStart"
	recipes.GET("/delete/:recipe_id", r.deleteRecipeHandler).Name = "deleteRecipe"
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

// getStartHandler is the handler for the starting a recipe
func (r *RecipesRouter) getStartHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	return c.Render(200, "recipe_start.html", map[string]interface{}{
		"Title":       "Starting Recipe",
		"Recipe":      re,
		"RecipeID":    id,
		"SquareColor": tools.EBCtoHex(re.ColorEBC),
	})
}

// getContinueHandler is the handler for the continue button on the recipes page
func (r *RecipesRouter) getContinueHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	redirectURL, err := r.statusRedirectURL(c, re, id)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, redirectURL)
}

// deleteRecipeHandler is the handler for the delete button on the recipes page
func (r *RecipesRouter) deleteRecipeHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	err := r.Store.Delete(id)
	if err != nil {
		return err
	}
	r.SummaryStore.DeleteSummaryRecorder(id)
	r.TLStore.DeleteTimeline(id)
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getRecipes"))
}

// statusRedirectURL returns the URL to redirect to based on the status of the recipe
func (r *RecipesRouter) statusRedirectURL(c echo.Context, re *recipe.Recipe, id string) (string, error) {
	status, params := re.GetStatus()
	switch status {
	case recipe.RecipeStatusCreated:
		return c.Echo().Reverse("getRecipeStart", id), nil
	case recipe.RecipeStatusMashing:
		param1, ok := params[0].(string)
		if !ok {
			return "", errors.New("invalid parameter for mashing status")
		}
		switch param1 {
		case "start":
			return c.Echo().Reverse("getMashStart", id), nil
		case "rast":
			rastNum, ok := params[1].(int)
			if !ok {
				return "", errors.New("invalid parameter for mashing status")
			}
			return c.Echo().Reverse("getRasts", id, rastNum), nil
		default:
			return "", errors.New("invalid parameter for mashing status")
		}
	case recipe.RecipeStatusLautering:
		return c.Echo().Reverse("getLautern", id), nil
	case recipe.RecipeStatusBoiling:
		param1, ok := params[0].(string)
		if !ok {
			return "", errors.New("invalid parameter for boiling status")
		}
		switch param1 {
		case "initialVol":
			return c.Echo().Reverse("getStartHopping", id), nil
		case "beforeBoil":
			return c.Echo().Reverse("getBoiling", id), nil
		case "lastBoil", "hop":
			ingrNum, ok := params[1].(int)
			if !ok {
				return "", errors.New("invalid parameter for boiling status")
			}
			return c.Echo().Reverse("getHopping", id, ingrNum), nil
		case "finalVol":
			return c.Echo().Reverse("getEndHopping", id), nil
		default:
			return "", errors.New("invalid parameter for boiling status")
		}
	case recipe.RecipeStatusCooling:
		return c.Echo().Reverse("getCooling", id), nil
	case recipe.RecipeStatusPreFermentation:
		param1, ok := params[0].(string)
		if !ok {
			return "", errors.New("invalid parameter for pre-fermentation status")
		}
		switch param1 {
		case "measure":
			return c.Echo().Reverse("getPreFermentation", id), nil
		case "water":
			volumeDiff, ok := params[1].(float64)
			if !ok {
				return "", errors.New("invalid parameter for pre-fermentation status")
			}
			sgDiff, ok := params[2].(float64)
			if !ok {
				return "", errors.New("invalid parameter for pre-fermentation status")
			}
			queryParams := fmt.Sprintf("?volumeDiff=%f&sgDiff=%f", volumeDiff, sgDiff)
			return c.Echo().Reverse("getPreFermentationWater", id) + queryParams, nil
		default:
			return "", errors.New("invalid parameter for pre-fermentation status")
		}
	case recipe.RecipeStatusFermenting:
		return c.Echo().Reverse("getFermentation", id), nil
	case recipe.RecipeStatusFinished:
		return c.Echo().Reverse("getEndFermentation", id), nil
	default:
		return "", errors.New("invalid recipe status")
	}
}
