package templates

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
)

//go:embed layouts/*.html
var templateFiles embed.FS

type Engine struct {
	templates map[string]*template.Template
}

func NewEngine() (*Engine, error) {
	e := &Engine{
		templates: make(map[string]*template.Template),
	}
	if err := e.load(); err != nil {
		return nil, err
	}
	return e, nil
}

func (e *Engine) load() error {
	entries, err := templateFiles.ReadDir("layouts")
	if err != nil {
		return fmt.Errorf("error leyendo templates: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		content, err := templateFiles.ReadFile("layouts/" + entry.Name())
		if err != nil {
			return fmt.Errorf("error leyendo template %s: %w", entry.Name(), err)
		}

		name := entry.Name()[:len(entry.Name())-5]

		tmpl, err := template.New(name).Funcs(template.FuncMap{
			"mod": func(a, b int) int { return a % b },
		}).Parse(string(content))
		if err != nil {
			return fmt.Errorf("error parseando template %s: %w", name, err)
		}

		e.templates[name] = tmpl
	}

	return nil
}

func (e *Engine) Render(templateID string, data map[string]interface{}) (string, error) {
	tmpl, ok := e.templates[templateID]
	if !ok {
		return "", fmt.Errorf("template no encontrado: %s — disponibles: %v", templateID, e.availableTemplates())
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("error renderizando template %s: %w", templateID, err)
	}

	return buf.String(), nil
}

func (e *Engine) availableTemplates() []string {
	keys := make([]string, 0, len(e.templates))
	for k := range e.templates {
		keys = append(keys, k)
	}
	return keys
}
