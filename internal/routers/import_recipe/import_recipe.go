package import_recipe

import (
	"brewday/internal/recipe"
	"brewday/internal/recipe/braureka_json"
	"brewday/internal/recipe/brewday"
	"brewday/internal/recipe/mmum"
	"brewday/internal/routers/common"
	"brewday/internal/tools"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

var parsers = map[string]RecipeParser{
	"braureka_json": &braureka_json.BraurekaJSONParser{},
	"mmum":          &mmum.MMUMParser{},
	"brewday":       &brewday.BrewdayParser{},
}

// errorMsgs maps parsing errors to messages to display
var errorMsgs = map[string]string{
	"max_vol_surpassed": "The recipe batch size is bigger than the brewing system max volume",
}

type ImportRouter struct {
	Store                RecipeStore
	SummaryRecorderStore SummaryStore
	TLStore              TimelineStore
	TempCache            map[string]*recipe.Recipe
	InfoProvider         InfoProvider
}

// storeRecipe stores a recipe in the temporary cache
func (r *ImportRouter) storeRecipe(re *recipe.Recipe) string {
	if r.TempCache == nil {
		r.TempCache = make(map[string]*recipe.Recipe)
	}
	id := idFromRecipe(re.Name)
	r.TempCache[id] = re
	return id
}

// getRecipe retrieves a recipe from the temporary cache
func (r *ImportRouter) getRecipe(id string) *recipe.Recipe {
	if r.TempCache == nil {
		return nil
	}
	return r.TempCache[id]
}

func (r *ImportRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	imp := parent.Group("/import")
	imp.GET("", r.getImportHandler).Name = "getImport"
	imp.POST("/preview", r.postImportPreviewHandler).Name = "postImportPreview"
	imp.GET("/:recipe_id/:next_action", r.getImportNextHandler).Name = "getImportNext"
	imp.POST("/set_bs/:recipe_id/:next_action", r.postImportSetBSHandler).Name = "postSetBS"
}

// getImportHandler is the handler for the import page
func (r *ImportRouter) getImportHandler(c echo.Context) error {
	id := c.QueryParam("recipe")
	errRaw := c.QueryParam("error")
	re := r.getRecipe(id)
	if re == nil {
		return c.Render(200, "import.html", map[string]interface{}{
			"Title":       "Import Recipe",
			"Recipe":      nil,
			"SquareColor": "#000000",
		})
	}
	if errRaw != "" {
		err, ok := errorMsgs[errRaw]
		if !ok {
			err = "Unknown error"
		}
		errRaw = err
	}
	return c.Render(200, "import.html", map[string]interface{}{
		"Title":       "Import Recipe",
		"Recipe":      re,
		"RecipeID":    id,
		"Error":       errRaw,
		"SquareColor": tools.EBCtoHex(re.ColorEBC),
	})
}

// postImportPreviewHandler is the handler for the import form preview
func (r *ImportRouter) postImportPreviewHandler(c echo.Context) error {
	file, err := c.FormFile("recipe_file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	bytes, err := io.ReadAll(src)
	if err != nil {
		return err
	}
	parserType := c.FormValue("parser_type")
	parser, ok := parsers[parserType]
	if !ok {
		return errors.New("invalid parser type")
	}
	recipe, err := parser.Parse(string(bytes))
	if err != nil {
		return err
	}
	id := r.storeRecipe(recipe)
	redirect := "getImport"
	idEncoded := url.QueryEscape(id)
	queryParams := "?recipe=" + idEncoded
	return c.Redirect(http.StatusFound, c.Echo().Reverse(redirect)+queryParams)
}

// idFromRecipe returns the identifier of a recipe based on its name
func idFromRecipe(name string) string {
	return name
}

// getImportNextHandler is the handler for importing and starting a recipe (or continuing)
func (r *ImportRouter) getImportNextHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	decodedID, err := url.QueryUnescape(id)
	if err != nil {
		return err
	}
	nextAction := c.Param("next_action")
	if nextAction == "" {
		return errors.New("no next action provided")
	}
	re := r.getRecipe(decodedID)
	if re == nil {
		return errors.New("no recipe found")
	}
	if r.InfoProvider.HasSystems() {
		if re.BrewingSystem == "" {
			available := r.InfoProvider.GetSystemNames()
			return c.Render(200, "import_add_bs.html", map[string]interface{}{
				"Title":          "Add brewing system",
				"Subtitle":       "No brewing system detected, add one?",
				"RecipeID":       decodedID,
				"NextAction":     nextAction,
				"BrewingSystems": available,
			})
		} else if re.BrewingSystem != "undefined" {
			maxVol, err := r.InfoProvider.GetMaxVol(re.BrewingSystem)
			if err != nil {
				return err
			}
			if re.BatchSize > maxVol {
				redirect := "getImport"
				params := url.Values{}
				params.Add("recipe", decodedID)
				params.Add("error", "max_vol_surpassed")
				return c.Redirect(http.StatusFound, c.Echo().Reverse(redirect)+"?"+params.Encode())
			}
		}
	}
	id, err = r.Store.Store(re)
	if err != nil {
		return err
	}
	err = r.Store.UpdateStatus(id, recipe.RecipeStatusCreated)
	if err != nil {
		return err
	}
	// Once stored, we can delete it from the cache
	delete(r.TempCache, decodedID)
	err = r.SummaryRecorderStore.AddSummary(id, re.Name)
	if err != nil {
		return err
	}
	err = r.SummaryRecorderStore.AddBrewingSystem(id, re.BrewingSystem)
	if err != nil {
		return err
	}
	err = r.TLStore.AddTimeline(id)
	if err != nil {
		return err
	}
	switch nextAction {
	case "start":
		return c.Redirect(http.StatusFound, c.Echo().Reverse("getRecipeStart", id))
	case "continue":
		return c.Redirect(http.StatusFound, c.Echo().Reverse("getImport"))
	default:
		return errors.New("invalid next action")
	}
}

// postImportSetBSHandler is the handler for adding a brewing system to a recipe
func (r *ImportRouter) postImportSetBSHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	decodedID, err := url.QueryUnescape(id)
	if err != nil {
		return err
	}
	nextAction := c.Param("next_action")
	if nextAction == "" {
		return errors.New("no next action provided")
	}
	re := r.getRecipe(decodedID)
	if re == nil {
		return errors.New("no recipe found")
	}
	bs := c.FormValue("bs")
	re.BrewingSystem = bs
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getImportNext", id, nextAction))
}
