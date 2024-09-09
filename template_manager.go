package tmpltree

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

// TemplateManager manages the template tree and base templates
type TemplateManager struct {
	Root          *TemplateNode
	BaseTemplates map[string]string
}

// RenderTemplate renders a template with the given path and base template
func (tm *TemplateManager) RenderTemplate(tmplPath string, baseTemplateName string, w io.Writer, data interface{}) error {
	parts := strings.Split(tmplPath, "/")
	node, ok := tm.Root.GetNode(parts[:len(parts)-1]...)
	if !ok {
		return fmt.Errorf("template node not found for path: %s", tmplPath)
	}

	fileName := parts[len(parts)-1] + ".html"
	found := false
	for _, file := range node.Files {
		if file == fileName {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("template file not found: %s", fileName)
	}

	fullPath := filepath.Join(node.Path, fileName)

	baseTemplatePath, ok := tm.BaseTemplates[baseTemplateName]
	if !ok {
		return fmt.Errorf("base template not found: %s", baseTemplateName)
	}

	tmpl, err := template.ParseFiles(baseTemplatePath, fullPath)
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	if err := tmpl.Execute(w, data); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	return nil
}
