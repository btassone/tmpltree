package tmpltree

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// TemplateNode represents a node in the template tree structure
type TemplateNode struct {
	Name     string
	Path     string
	Children map[string]*TemplateNode
	Files    []string
}

// TemplateManager manages the template tree and base templates
type TemplateManager struct {
	Root          *TemplateNode
	BaseTemplates map[string]string
}

// NewTemplateNode creates a new TemplateNode
func NewTemplateNode(name string, path string) *TemplateNode {
	return &TemplateNode{
		Name:     name,
		Path:     path,
		Children: make(map[string]*TemplateNode),
		Files:    []string{},
	}
}

// NewTemplateManagerFunc is a function type for creating a new TemplateManager
type NewTemplateManagerFunc func(rootDir string, baseTemplates map[string]string) (*TemplateManager, error)

// NewTemplateManagerImpl is the actual implementation of NewTemplateManager
var NewTemplateManagerImpl NewTemplateManagerFunc = func(rootDir string, baseTemplates map[string]string) (*TemplateManager, error) {
	root, err := BuildTemplateTree(rootDir)
	if err != nil {
		return nil, err
	}

	return &TemplateManager{
		Root:          root,
		BaseTemplates: baseTemplates,
	}, nil
}

// NewTemplateManager creates a new TemplateManager
func NewTemplateManager(rootDir string, baseTemplates map[string]string) (*TemplateManager, error) {
	return NewTemplateManagerImpl(rootDir, baseTemplates)
}

// BuildTemplateTree constructs a tree structure of the template directory
func BuildTemplateTree(rootDir string) (*TemplateNode, error) {
	root := NewTemplateNode("templates", rootDir)
	requiredFolders := []string{"layouts", "pages", "partials"}

	for _, folder := range requiredFolders {
		root.Children[folder] = NewTemplateNode(folder, filepath.Join(rootDir, folder))
	}

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}

		parts := strings.Split(relPath, string(filepath.Separator))
		if len(parts) == 1 {
			return nil // Skip the root directory
		}

		current := root.Children[parts[0]]
		if current == nil {
			return nil // Skip if not in a required folder
		}

		for i, part := range parts[1:] {
			if i == len(parts)-2 { // If this is the last part and it's a file
				if !info.IsDir() {
					current.Files = append(current.Files, part)
				}
				break
			}
			if _, exists := current.Children[part]; !exists {
				current.Children[part] = NewTemplateNode(part, filepath.Join(current.Path, part))
			}
			current = current.Children[part]
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return root, nil
}

// Print prints the template tree structure to the provided writer
func (n *TemplateNode) Print(w io.Writer, indent string) {
	_, err := fmt.Fprintf(w, "%s%s/\n", indent, n.Name)
	if err != nil {
		log.Print(err)
		return
	}
	for _, file := range n.Files {
		_, err := fmt.Fprintf(w, "%s  %s\n", indent, file)
		if err != nil {
			log.Print(err)
			return
		}
	}
	for _, child := range n.Children {
		child.Print(w, indent+"  ")
	}
}

// GetNode retrieves a specific node from the template tree
func (n *TemplateNode) GetNode(path ...string) (*TemplateNode, bool) {
	current := n
	for _, part := range path {
		if child, ok := current.Children[part]; ok {
			current = child
		} else {
			return nil, false
		}
	}
	return current, true
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
