package renderer

import (
	"bytes"
	"html/template"
)

type TemplateRenderer struct{}

func NewTemplateRenderer() *TemplateRenderer {
	return &TemplateRenderer{}
}

func (r *TemplateRenderer) Render(
	templateText string,
	variables map[string]interface{},
) (string, error) {

	tmpl, err := template.New("email").
		Option("missingkey=error").
		Parse(templateText)

	if err != nil {
		return "", err
	}

	var output bytes.Buffer

	err = tmpl.Execute(&output, variables)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}