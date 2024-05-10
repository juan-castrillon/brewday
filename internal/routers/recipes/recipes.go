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
	err = r.SummaryStore.DeleteSummaryRecorder(id)
	if err != nil {
		return err
	}
	err = r.TLStore.DeleteTimeline(id)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getRecipes"))
}

// statusRedirectURL returns the URL to redirect to based on the status of the recipe
func (r *RecipesRouter) statusRedirectURL(c echo.Context, re *recipe.Recipe, id string) (string, error) {
	status, params := re.GetStatus()
	switch status {
	case recipe.RecipeStatusCreated:
		return c.Echo().Reverse("getRecipeStart", id), nil
	case recipe.RecipeStatusMashing:
		switch params[0] {
		case "start":
			return c.Echo().Reverse("getMashStart", id), nil
		case "rast":
			return c.Echo().Reverse("getRasts", id, params[1]), nil
		default:
			return "", errors.New("invalid parameter for mashing status")
		}
	case recipe.RecipeStatusLautering:
		return c.Echo().Reverse("getLautern", id), nil
	case recipe.RecipeStatusBoiling:
		switch params[0] {
		case "initialVol":
			return c.Echo().Reverse("getStartHopping", id), nil
		case "beforeBoil":
			return c.Echo().Reverse("getBoiling", id), nil
		case "lastBoil", "hop":
			return c.Echo().Reverse("getHopping", id, params[1]), nil
		case "finalVol":
			return c.Echo().Reverse("getEndHopping", id), nil
		default:
			return "", errors.New("invalid parameter for boiling status")
		}
	case recipe.RecipeStatusCooling:
		return c.Echo().Reverse("getCooling", id), nil
	case recipe.RecipeStatusPreFermentation:
		switch params[0] {
		case "measure":
			return c.Echo().Reverse("getPreFermentation", id), nil
		case "water":
			volumeDiff := params[1]
			sgDiff := params[2]
			queryParams := fmt.Sprintf("?volumeDiff=%s&sgDiff=%s", volumeDiff, sgDiff)
			return c.Echo().Reverse("getPreFermentationWater", id) + queryParams, nil
		default:
			return "", errors.New("invalid parameter for pre-fermentation status")
		}
	case recipe.RecipeStatusFermenting:
		switch params[0] {
		case "yeast":
			return c.Echo().Reverse("getFermentationYeast", id), nil
		case "start":
			return c.Echo().Reverse("getMainFermentationStart", id), nil
		case "wait", "main":
			return c.Echo().Reverse("getMainFermentation", id), nil
		case "dry_hop_start":
			return c.Echo().Reverse("getDryHopStart", id), nil
		case "dry_hop_confirm":
			return c.Echo().Reverse("getDryHopConfirm", id), nil
		case "pre_bottle":
			return c.Echo().Reverse("getPreBottle", id), nil
		case "bottle":
			queryParams := fmt.Sprintf("?type=%s", params[1])
			return c.Echo().Reverse("getBottle", id) + queryParams, nil
		case "start_secondary":
			return c.Echo().Reverse("getSecondaryFermentationStart", id), nil
		default:
			return "", errors.New("invalid parameter for fermentation status")
		}
	case recipe.RecipeStatusFinished:
		return c.Echo().Reverse("getEnd", id), nil
	default:
		return "", errors.New("invalid recipe status")
	}
}
