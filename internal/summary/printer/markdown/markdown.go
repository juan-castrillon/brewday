package markdown

import (
	"brewday/internal/recipe"
	"brewday/internal/summary"
	"bytes"
	_ "embed"
	"strings"
	"text/template"
	"time"
)

//go:embed md.tmpl
var tmpl string

type RecipeSummary struct {
	RealData       *summary.Summary
	OriginalRecipe *recipe.Recipe
}

type MarkdownPrinter struct {
}

func (m *MarkdownPrinter) Print(s *summary.Summary, r *recipe.Recipe, timeline []string) (string, error) {
	t, err := template.New("md.tmpl").Funcs(template.FuncMap{
		"SplitString": func(st, sep string) []string {
			return strings.Split(st, sep)
		},
	}).Parse(tmpl)
	if err != nil {
		return "", err
	}
	s.GenerationDate = time.Now().Format("2006-01-02 15:04:05")
	s.Timeline = timeline
	var buf bytes.Buffer
	rs := &RecipeSummary{
		RealData:       s,
		OriginalRecipe: r,
	}
	err = t.Execute(&buf, rs)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
