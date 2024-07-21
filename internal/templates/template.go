package templates

import (
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"strings"
)

func Load(templateFS fs.FS, service string) (map[string]*template.Template, error) {
	entries, err := fs.ReadDir(templateFS, fmt.Sprintf("%v/views", service))
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

	layout := "layout.html"
	views := fmt.Sprintf("%v/views", service)
	components := fmt.Sprintf("%v/components/*.html", service)

	compiled := make(map[string]*template.Template)
	for _, viewName := range viewNames {
		compiled[viewName[:strings.LastIndex(viewName, ".")]], err = template.ParseFS(templateFS, layout, fmt.Sprintf("%v/%v", views, viewName), components)
		if err != nil {
			return nil, err
		}
	}

	slog.Info("Parsed templates", "service", service, "amount", len(compiled))

	return compiled, nil
}
