package import_recipe

import (
	"brewday/internal/recipe"
	"io"

	"github.com/labstack/echo/v4"
)

type ImportRouter struct {
	Parser RecipeParser
	Store  RecipeStore
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
			"Title":  "Import Recipe",
			"Recipe": nil,
		})
	}
	id, err := r.Store.Store(re)
	if err != nil {
		return err
	}
	return c.Render(200, "import.html", map[string]interface{}{
		"Title":    "Import Recipe",
		"Recipe":   re,
		"RecipeID": id,
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
	recipe, err := r.Parser.Parse(string(bytes))
	if err != nil {
		return err
	}
	c.Set("recipe", recipe)
	return r.getImportHandler(c)
}
