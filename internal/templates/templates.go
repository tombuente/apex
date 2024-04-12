package templates

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"strings"
)

//go:embed templates/*
var fs embed.FS

func Load(service string) (map[string]*template.Template, error) {
	entries, err := fs.ReadDir(fmt.Sprintf("templates/%v/views", service))
	if err != nil {
		return nil, err
	}

	var viewNames []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		viewNames = append(viewNames, entry.Name())
	}

	layout := "templates/layout.html"
	views := fmt.Sprintf("templates/%v/views", service)
	components := fmt.Sprintf("templates/%v/components/*.html", service)

	compiled := make(map[string]*template.Template)
	for _, viewName := range viewNames {
		compiled[viewName[:strings.LastIndex(viewName, ".")]], err = template.ParseFS(fs, layout, fmt.Sprintf("%v/%v", views, viewName), components)
		if err != nil {
			return nil, err
		}
	}

	slog.Info("Parsed templates", "service", service, "amount", len(compiled))

	return compiled, nil
}
