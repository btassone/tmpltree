package tmpltree

import (
	"fmt"
	"io"
	"log"
)

// TemplateNode represents a node in the template tree structure
type TemplateNode struct {
	Name     string
	Path     string
	Children map[string]*TemplateNode
	Files    []string
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
