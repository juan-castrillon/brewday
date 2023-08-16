package import_recipe

import (
	"brewday/internal/recipe"
	"brewday/internal/recipe/braureka_json"
	"brewday/internal/recipe/mmum"
	"brewday/internal/tools"
	"io"

	"github.com/labstack/echo/v4"
)

var parsers = map[string]RecipeParser{
	"braureka_json": &braureka_json.BraurekaJSONParser{},
	"mmum":          &mmum.MMUMParser{},
}

type ImportRouter struct {
	Store RecipeStore
}

func (r *ImportRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	parent.GET("/import", r.getImportHandler).Name = "getImport"
	parent.POST("/import", r.postImportHandler).Name = "postMashImport"
}

// getImportHandler is the handler for the import page
func (r *ImportRouter) getImportHandler(c echo.Context) error {
	re, ok := c.Get("recipe").(*recipe.Recipe)
	if !ok || re == nil {
		return c.Render(200, "import.html", map[string]interface{}{
			"Title":       "Import Recipe",
			"Recipe":      nil,
			"SquareColor": "#000000",
		})
	}
	id, err := r.Store.Store(re)
	if err != nil {
		return err
	}
	return c.Render(200, "import.html", map[string]interface{}{
		"Title":       "Import Recipe",
		"Recipe":      re,
		"RecipeID":    id,
		"SquareColor": tools.EBCtoHex(re.ColorEBC),
	})
}

// postImportHandler is the handler for the import form
func (r *ImportRouter) postImportHandler(c echo.Context) error {
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
		return r.getImportHandler(c)
	}
	recipe, err := parser.Parse(string(bytes))
	if err != nil {
		return err
	}
	c.Set("recipe", recipe)
	return r.getImportHandler(c)
}
