package summary

import (
	"brewday/internal/routers/common"
	"brewday/internal/summary"
	"brewday/internal/summary/printer/markdown"
	"errors"
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
)

type SummaryRouter struct {
	SummaryStore SummaryStore
	TLStore      TimelineStore
}

// getSummary returns the summary
func (r *SummaryRouter) getSummary(id string) (*summary.Summary, error) {
	if r.SummaryStore != nil {
		return r.SummaryStore.GetSummary(id)
	}
	return nil, nil
}

// getExtention returns the extension
func (r *SummaryRouter) getExtension(format string) string {
	switch format {
	case "markdown":
		return "md"
	default:
		return "md"
	}
}

// getTimeline returns the timeline
func (r *SummaryRouter) getTimeline(id string) ([]string, error) {
	if r.TLStore != nil {
		return r.TLStore.GetTimeline(id)
	}
	return []string{}, nil
}

// printSummary instanciates the correct summary printer based on the format and prints the summary as a string
func (r *SummaryRouter) printSummary(format string, summ *summary.Summary, tl []string) (string, error) {
	var p SummaryPrinter
	switch format {
	case "markdown":
		p = &markdown.MarkdownPrinter{}
	default:
		return "", errors.New("could not find suitable printer for format " + format)
	}
	return p.Print(summ, tl)
}

// RegisterRoutes registers the routes for the summary router
func (r *SummaryRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	summary := parent.Group("/summary")
	summary.GET("/:recipe_id", r.getSummaryHandler).Name = "getSummary"
}

// getSummaryHandler handles the GET /summary/:recipe_id route
func (r *SummaryRouter) getSummaryHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	tl, err := r.getTimeline(id)
	if err != nil {
		return err
	}
	format := c.QueryParam("format")
	if format == "" {
		format = "markdown"
	}
	format = strings.ToLower(format)
	summ, err := r.getSummary(id)
	if err != nil {
		return err
	}
	ext := r.getExtension(format)
	fileName := id + "." + ext
	content, err := r.printSummary(format, summ, tl)
	if err != nil {
		return err
	}
	c.Response().Header().Set("Content-Type", "application/octet-stream")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	_, err = c.Response().Write([]byte(content))
	return err
}
