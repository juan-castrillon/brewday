package markdown

import (
	"brewday/internal/summary"
	"bytes"
	_ "embed"
	"strings"
	"text/template"
	"time"
)

//go:embed md.tmpl
var tmpl string

type MarkdownPrinter struct {
}

func (m *MarkdownPrinter) Print(s *summary.Summary, timeline []string) (string, error) {
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
	err = t.Execute(&buf, s)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
