package import_recipe

import (
	"brewday/internal/recipe"
	"brewday/internal/recipe/braureka_json"
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
}

type ImportRouter struct {
	Store                RecipeStore
	SummaryRecorderStore SummaryRecorderStore
	TLStore              TimelineStore
	TempCache            map[string]*recipe.Recipe
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
}

// getImportHandler is the handler for the import page
func (r *ImportRouter) getImportHandler(c echo.Context) error {
	id := c.QueryParam("recipe")
	re := r.getRecipe(id)
	if re == nil {
		return c.Render(200, "import.html", map[string]interface{}{
			"Title":       "Import Recipe",
			"Recipe":      nil,
			"SquareColor": "#000000",
		})
	}
	return c.Render(200, "import.html", map[string]interface{}{
		"Title":       "Import Recipe",
		"Recipe":      re,
		"RecipeID":    id,
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
	// TODO: make this configurable probably via the UI
	r.SummaryRecorderStore.AddSummaryRecorder(id, "markdown")
	r.TLStore.AddTimeline(id, "basic")
	switch nextAction {
	case "start":
		return c.Redirect(http.StatusFound, c.Echo().Reverse("getRecipeStart", id))
	case "continue":
		return c.Redirect(http.StatusFound, c.Echo().Reverse("getImport"))
	default:
		return errors.New("invalid next action")
	}
}
