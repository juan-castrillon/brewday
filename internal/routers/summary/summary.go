package summary

import (
	"brewday/internal/recipe"
	"brewday/internal/routers/common"
	"brewday/internal/summary"
	"brewday/internal/summary/printer/tpl"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type SummaryRouter struct {
	SummaryStore  SummaryStore
	RecipeStore   RecipeStore
	TLStore       TimelineStore
	PrinterConfig map[string]SummaryPrinterConfiguration
}

// getSummary returns the summary
func (r *SummaryRouter) getSummary(id string) (*summary.Summary, error) {
	if r.SummaryStore != nil {
		return r.SummaryStore.GetSummary(id)
	}
	return nil, nil
}

// getTimeline returns the timeline
func (r *SummaryRouter) getTimeline(id string) ([]string, error) {
	if r.TLStore != nil {
		return r.TLStore.GetTimeline(id)
	}
	return []string{}, nil
}

// printSummary instanciates the correct summary printer based on the format and prints the summary as a string
func (r *SummaryRouter) printSummary(summ *summary.Summary, rec *recipe.Recipe, tl []string, tmplPath string) (string, error) {
	p := &tpl.TemplatePrinter{}
	return p.Print(summ, rec, tl, tmplPath)
}

// RegisterRoutes registers the routes for the summary router
func (r *SummaryRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	summary := parent.Group("/summary")
	root.GET("/end/:recipe_id", r.getEndHandler).Name = "getEnd"
	summary.GET("/:recipe_id", r.getSummaryHandler).Name = "getSummary"
}

// getEndHandler is the handler for the end page
func (r *SummaryRouter) getEndHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	return c.Render(http.StatusOK, "finished_day.html", map[string]interface{}{
		"Title":                  "End Fermentation",
		"RecipeID":               id,
		"SummaryTemplateConfigs": r.PrinterConfig,
		"Subtitle":               "Congratulations, you've finished the brew!",
	})
}

// getSummaryHandler handles the GET /summary/:recipe_id route
func (r *SummaryRouter) getSummaryHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	rec, err := r.RecipeStore.Retrieve(id)
	if err != nil {
		return err
	}
	tl, err := r.getTimeline(id)
	if err != nil {
		return err
	}
	tmplConfig := c.QueryParam("template_config")
	if tmplConfig == "" {
		return errors.New("Empty template configuration selected")
	}
	tmplConfig = strings.ToLower(tmplConfig)
	summ, err := r.getSummary(id)
	if err != nil {
		return err
	}
	tmplConfigData, ok := r.PrinterConfig[tmplConfig]
	if !ok {
		if tmplConfig == "default" {
			tmplConfigData = SummaryPrinterConfiguration{
				Path:       "", // So the default is taken in the printer
				FileFormat: "md",
			}
		} else {
			return errors.New("Template configuration not found")
		}
	}
	fileName := id + "." + tmplConfigData.FileFormat
	content, err := r.printSummary(summ, rec, tl, tmplConfigData.Path)
	if err != nil {
		return err
	}
	c.Response().Header().Set("Content-Type", "application/octet-stream")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	_, err = c.Response().Write([]byte(content))
	return err
}
