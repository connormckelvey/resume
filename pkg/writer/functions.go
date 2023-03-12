package writer

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"

	"github.com/connormckelvey/resume/pkg/schema"
)

var Utils = template.FuncMap{
	"link":      renderLink,
	"dates":     renderDates,
	"join":      renderJoin,
	"embed_css": renderCSS,
}

func renderTemplateString[T any](templateName string, templateString string, data T) (template.HTML, error) {
	t, err := template.
		New(templateName).
		Parse(templateString)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return template.HTML(buf.String()), nil
}

func renderLink(link *schema.Link) (template.HTML, error) {
	return renderTemplateString("link", `<a href="{{ .Href }}">{{ .Text }}</a>`, link)
}

func renderDates(dates *schema.DateRange) (template.HTML, error) {
	t := `<span class="dates">
{{ .Start }} -
{{if .End}}
	{{ .End }}
{{else}}
	Present
{{end}}
</span>`
	return renderTemplateString("dates", t, dates)
}

func renderJoin(s []string, sep string) (template.HTML, error) {
	return template.HTML(strings.Join(s, sep)), nil
}

func renderCSS(path string) (template.HTML, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	css := fmt.Sprintf(`<style type="text/css">%v</style>`, string(contents))
	return template.HTML(css), nil
}
