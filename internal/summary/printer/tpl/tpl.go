package tpl

import (
	"brewday/internal/recipe"
	"brewday/internal/summary"
	"bytes"
	_ "embed"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

//go:embed md_default.tmpl
var tmpl string

type RecipeSummary struct {
	RealData       *summary.Summary
	OriginalRecipe *recipe.Recipe
}

type TemplatePrinter struct {
}

// Print will return a rendered template given a summary recipe and timeline
// The template will be read from the tmplPath
// If this is empty, the default template will be used
func (m *TemplatePrinter) Print(s *summary.Summary, r *recipe.Recipe, timeline []string, tmplPath string) (string, error) {
	var n, tRaw string
	if tmplPath == "" {
		n = "md_default.tmpl"
		tRaw = tmpl
	} else {
		n = filepath.Base(tmplPath)
		raw, err := os.ReadFile(tmplPath)
		if err != nil {
			return "", err
		}
		tRaw = string(raw)
	}
	t, err := template.New(n).Funcs(template.FuncMap{
		"SplitString": func(st, sep string) []string {
			return strings.Split(st, sep)
		},
	}).Parse(tRaw)
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
